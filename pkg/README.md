# pkg

This directory contains the source code for `r2`'s main functionality. Functions defined here are
to be used by the [CLI](../cmd) and are publicly exported for API use.

## Architecture

- [client.go](client.go) contains all R2 client-level operations (e.g. configuration, bucket
  creation, etc.)
- [bucket.go](bucket.go) contains all bucket-level operations (e.g. listing objects, fetching
  objects, etc.)
- [helpers.go](helpers.go) contains miscellaneous helper functions used throughout the CLI
