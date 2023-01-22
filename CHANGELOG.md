# Changelog

This changelog goes through all the changes that have been made in each release.

## v0.1.0-alpha

First release of `r2`! This release includes all the commands of the AWS CLI's `s3` subcommand, but
not all the options.

- ADDED
  - [`configure` command](cmd/configure.go) — configure `r2` to use your R2 credentials
  - [`cp` command](cmd/cp.go) — copy objects and directories
  - [`ls` command](cmd/ls.go) — list objects and directories
  - [`mb` command](cmd/mb.go) — make buckets
  - [`mv` command](cmd/mv.go) — move objects and directories
  - [`presign` command](cmd/presign.go) — generate pre-signed URLs
  - [`rb` command](cmd/rb.go) — remove buckets
  - [`rm` command](cmd/rm.go) — remove objects and directories
  - [`sync` command](cmd/sync.go) — synchronize objects and directories
