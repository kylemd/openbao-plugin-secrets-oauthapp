Merge `upstream/main` into the checked-out fork while preserving the narrow
neutral default User-Agent compatibility behavior in
`pkg/oauth2ext/clientctx` and its backend call sites.

Requirements:

- Resolve only conflicts caused by the upstream update.
- Preserve explicitly supplied User-Agent values.
- Never add provider credentials, tokens, callback data, or retailer-specific
  secrets.
- Keep the fork patch minimal and generic.
- Update `fork-release.json` with the new upstream commit and increment the
  patch revision and release tag.
- Run formatting, the full race-enabled test suite, and `go vet ./...`.
- Leave the validated merged result as uncommitted workspace changes; do not
  commit, push, or publish.
