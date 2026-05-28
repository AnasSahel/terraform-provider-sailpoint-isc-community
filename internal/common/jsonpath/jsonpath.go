// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

// Package jsonpath implements a minimal subset of JSONPath for preserving
// server-managed fields in connector_attributes arrays.
//
// Supported syntax:
//
//	$.key
//	$.key.nested
//	$.key[N].field   (N is a non-negative integer index)
//	$.key[*].field   (wildcard: all array elements)
package jsonpath

import (
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
			} else {
				idx, err := strconv.Atoi(inner)
				if err != nil || idx < 0 {
					return nil, fmt.Errorf("invalid index %q in path: %q (expected non-negative integer or *)", inner, path)
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
