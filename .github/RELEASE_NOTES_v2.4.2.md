# Release v2.4.2 - Surface real API error messages on 4xx

Patch release that fixes a long-standing transport bug where 4xx responses from SailPoint were masked by the misleading `resty: content decoder not found` error, hiding the real API message and forcing users to re-run failed calls manually via `sail api` to discover the actual cause.

## Bug Fixes

- **Error handling**: 4xx responses from SailPoint now surface the real API error message instead of `resty: content decoder not found`. Root cause: the SailPoint edge (Cloudflare) returns a non-standard `Content-Encoding: UTF-8` header on some 4xx responses (`UTF-8` is a charset, not an encoding — the body is plain JSON). Resty v3 only registers `gzip` / `deflate` decompressers by default and bailed before the JSON body could be read. The provider now registers a no-op decompresser keyed on `UTF-8` so the response body reaches the per-resource error formatters. Affects every resource — most visibly transform create with rejected names, identity profile update during background tasks, and any other endpoint hitting a 4xx with this header. Closes #81.

## Notes for upgraders

- The previous mitigation (`Accept-Encoding: identity`) added in v2.4.0 (PR #82) is preserved as a microbenefit on 2xx responses, but it was not effective on 4xx because Cloudflare ignores the header for error responses. v2.4.2 is the first version where the underlying transport behavior is actually fixed.
- After upgrading, an `apply` that previously failed with `resty: content decoder not found` will fail with the real SailPoint message instead — typically a `400.x.y` `detailCode` plus a human-readable `text` field. This is intentional: you'll now see what to fix.

## Full Changelog

See [CHANGELOG.md](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/blob/main/CHANGELOG.md) for complete details.

---

**Questions or Issues?** Please open an issue on [GitHub](https://github.com/AnasSahel/terraform-provider-sailpoint-isc-community/issues).
