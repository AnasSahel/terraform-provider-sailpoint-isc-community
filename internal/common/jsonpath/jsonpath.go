// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

// Package jsonpath implements a minimal subset of JSONPath for preserving
// server-managed fields in JSON attributes.
//
// Supported syntax:
//
//	$.key
//	$.key.nested
//	$.key[N].field          (N is a non-negative integer index)
//	$.key[*].field          (wildcard: all array elements)
//	$.key['obj key']        (bracket-quoted, single quotes: for keys containing spaces or dots)
//	$.key["obj key"]        (bracket-quoted, double quotes: equivalent double-quote variant)
//	$['obj key']            (bracket-quoted key directly after root '$')
//	$['outer']['inner']     (chained bracket-quoted keys)
package jsonpath

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type segKind int

const (
	segKey      segKind = iota
	segWildcard         // [*]
	segIndex            // [N]
)

type segment struct {
	kind segKind
	key  string
	idx  int
}

// Parse parses a minimal JSONPath expression into segments.
// Returns an error for any syntax not in the supported subset.
func Parse(path string) ([]segment, error) {
	if !strings.HasPrefix(path, "$") {
		return nil, fmt.Errorf("path must start with '$', got: %q", path)
	}
	rest := path[1:]
	var segs []segment

	for len(rest) > 0 {
		switch rest[0] {
		case '.':
			rest = rest[1:]
			end := strings.IndexAny(rest, ".[")
			var key string
			if end < 0 {
				key, rest = rest, ""
			} else {
				key, rest = rest[:end], rest[end:]
			}
			if key == "" {
				return nil, fmt.Errorf("empty key segment in path: %q", path)
			}
			segs = append(segs, segment{kind: segKey, key: key})

		case '[':
			end := strings.Index(rest, "]")
			if end < 0 {
				return nil, fmt.Errorf("unclosed '[' in path: %q", path)
			}
			inner := rest[1:end]
			rest = rest[end+1:]
			if inner == "*" {
				segs = append(segs, segment{kind: segWildcard})
			} else if len(inner) >= 2 && (inner[0] == '\'' || inner[0] == '"') {
				quote := inner[0]
				if inner[len(inner)-1] != quote {
					return nil, fmt.Errorf("mismatched quotes in bracket key in path: %q", path)
				}
				key := inner[1 : len(inner)-1]
				if key == "" {
					return nil, fmt.Errorf("empty bracket-quoted key in path: %q", path)
				}
				segs = append(segs, segment{kind: segKey, key: key})
			} else {
				idx, err := strconv.Atoi(inner)
				if err != nil || idx < 0 {
					return nil, fmt.Errorf("invalid index %q in path: %q (expected non-negative integer, *, 'key', or \"key\")", inner, path)
				}
				segs = append(segs, segment{kind: segIndex, idx: idx})
			}

		default:
			return nil, fmt.Errorf("unexpected character %q in path: %q", string(rest[0]), path)
		}
	}

	if len(segs) == 0 {
		return nil, fmt.Errorf("path %q has no segments", path)
	}
	return segs, nil
}

// Validate returns an error if path is not a valid minimal JSONPath expression.
func Validate(path string) error {
	_, err := Parse(path)
	return err
}

// PreservePaths re-injects values from serverAttrs into merged at each path.
// Paths that don't exist in serverAttrs are silently skipped.
// Returns an error only if a path cannot be parsed.
func PreservePaths(merged, serverAttrs map[string]interface{}, paths []string) error {
	for _, path := range paths {
		segs, err := Parse(path)
		if err != nil {
			return fmt.Errorf("preserve paths: %w", err)
		}
		preserveNode(merged, serverAttrs, segs)
	}
	return nil
}

// PreservePathsInJSON unmarshals mergedJSON and priorJSON, re-injects the values
// from priorJSON into mergedJSON at each path (using PreservePaths), and
// re-marshals the result. Paths absent from priorJSON are silently skipped.
// Returns the merged JSON string, or an error if either input is invalid JSON or
// a path cannot be parsed.
func PreservePathsInJSON(mergedJSON, priorJSON string, paths []string) (string, error) {
	var merged map[string]interface{}
	if err := json.Unmarshal([]byte(mergedJSON), &merged); err != nil {
		return "", fmt.Errorf("unmarshal merged JSON: %w", err)
	}
	var prior map[string]interface{}
	if err := json.Unmarshal([]byte(priorJSON), &prior); err != nil {
		return "", fmt.Errorf("unmarshal prior JSON: %w", err)
	}
	if err := PreservePaths(merged, prior, paths); err != nil {
		return "", err
	}
	out, err := json.Marshal(merged)
	if err != nil {
		return "", fmt.Errorf("marshal merged JSON: %w", err)
	}
	return string(out), nil
}

// RemovePaths deletes the leaf at each path from m. Paths that don't exist in m
// are silently skipped. Returns an error only if a path cannot be parsed.
//
// Unlike PreservePaths (which copies a value in), this prunes the path entirely.
// It is used to honour an "ignore this nested field" contract: a path the
// practitioner declared as ignored is removed from the projected state so it can
// never produce a diff (the config omits it, so the state must omit it too).
func RemovePaths(m map[string]interface{}, paths []string) error {
	for _, path := range paths {
		segs, err := Parse(path)
		if err != nil {
			return fmt.Errorf("remove paths: %w", err)
		}
		removeNode(m, segs)
	}
	return nil
}

// RemovePathsInJSON unmarshals jsonStr, removes the leaf at each path (using
// RemovePaths), and re-marshals the result. Paths absent from the input are
// silently skipped. Returns an error if the input is invalid JSON or a path
// cannot be parsed.
func RemovePathsInJSON(jsonStr string, paths []string) (string, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		return "", fmt.Errorf("unmarshal JSON: %w", err)
	}
	if err := RemovePaths(m, paths); err != nil {
		return "", err
	}
	out, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("marshal JSON: %w", err)
	}
	return string(out), nil
}

// removeNode walks node following segs and deletes the leaf key. Array segments
// (wildcard/index) recurse into the matching element(s); a leaf key inside an
// array element is deleted from each element map.
func removeNode(node interface{}, segs []segment) {
	if len(segs) == 0 {
		return
	}

	seg := segs[0]
	rest := segs[1:]

	switch seg.kind {
	case segKey:
		m, ok := node.(map[string]interface{})
		if !ok {
			return
		}
		if len(rest) == 0 {
			delete(m, seg.key)
			return
		}
		removeNode(m[seg.key], rest)

	case segWildcard:
		arr, ok := node.([]interface{})
		if !ok {
			return
		}
		for i := range arr {
			removeNode(arr[i], rest)
		}

	case segIndex:
		arr, ok := node.([]interface{})
		if !ok {
			return
		}
		if seg.idx >= len(arr) {
			return
		}
		removeNode(arr[seg.idx], rest)
	}
}

// preserveNode walks merged and server in lock-step following segs,
// and at the leaf key copies the server value into merged.
func preserveNode(mergedNode, serverNode interface{}, segs []segment) {
	if len(segs) == 0 {
		return
	}

	seg := segs[0]
	rest := segs[1:]

	switch seg.kind {
	case segKey:
		mergedMap, ok1 := mergedNode.(map[string]interface{})
		serverMap, ok2 := serverNode.(map[string]interface{})
		if !ok1 || !ok2 {
			return
		}
		serverVal, exists := serverMap[seg.key]
		if !exists {
			return
		}
		if len(rest) == 0 {
			// Leaf: copy server value into merged map.
			mergedMap[seg.key] = serverVal
			return
		}
		mergedVal := mergedMap[seg.key]
		preserveNode(mergedVal, serverVal, rest)

	case segWildcard:
		mergedArr, ok1 := mergedNode.([]interface{})
		serverArr, ok2 := serverNode.([]interface{})
		if !ok1 || !ok2 {
			return
		}
		for i := range mergedArr {
			if i >= len(serverArr) {
				break
			}
			preserveNode(mergedArr[i], serverArr[i], rest)
		}

	case segIndex:
		mergedArr, ok1 := mergedNode.([]interface{})
		serverArr, ok2 := serverNode.([]interface{})
		if !ok1 || !ok2 {
			return
		}
		i := seg.idx
		if i >= len(mergedArr) || i >= len(serverArr) {
			return
		}
		preserveNode(mergedArr[i], serverArr[i], rest)
	}
}
