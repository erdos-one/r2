# Architecture

Knowing where to find things in the repo can be difficult. This document aims to help you find your
way.

## [cmd](cmd)

The cmd directory contains the source code for the CLI. The CLI is built using the [cobra
](https://cobra.dev) framework.

The CLI is split into several files:

- [cmd/root.go](cmd/root.go) contains the root command that is executed when `r2` is run
- [cmd/configure.go](cmd/configure.go) contains the `configure` command
- [cmd/cp.go](cmd/cp.go) contains the `cp` command
- [cmd/ls.go](cmd/ls.go) contains the `ls` command
- [cmd/mb.go](cmd/mb.go) contains the `mb` command
- [cmd/mv.go](cmd/mv.go) contains the `mv` command
- [cmd/presign.go](cmd/presign.go) contains the `presign` command
- [cmd/rb.go](cmd/rb.go) contains the `rb` command
- [cmd/rm.go](cmd/rm.go) contains the `rm` command
- [cmd/sync.go](cmd/sync.go) contains the `sync` command

## [pkg](pkg)

The pkg directory contains the source code for the R2 package — the library enabling the CLI to
communicate with R2.

- [pkg/client.go](pkg/client.go) contains all R2 client-level operations (e.g. configuration, bucket
  creation, etc.)
- [pkg/bucket.go](pkg/bucket.go) contains all bucket-level operations (e.g. listing objects, fetching
  objects, etc.)
- [pkg/helpers.go](pkg/helpers.go) contains miscellaneous helper functions used throughout the CLI

## [workflows](.github/workflows)

The workflows directory contains the GitHub Actions workflows used for this repo.

- [.github/workflow/assembly.yml](.github/workflows/assembly.yml)
  - Bumps the version of the CLI in the [install script](install.sh) to fetch the latest release
- [.github/workflow/release.yml](.github/workflows/release.yml)
  - Builds and deploys the CLI to GitHub Releases

## [assets](assets)

The assets directory contains the assets used for the repo.

- [assets/bucket.svg](assets/bucket.svg) is the R2 bucket icon used in the [README](README.md)

## [install.sh](install.sh)

The install script is used to install the latest release of the CLI.

## Thanks

Thanks to [Alex Kladov](https://matklad.github.io/) for his [blog post
](https://matklad.github.io//2021/02/06/ARCHITECTURE.md.html) on the importance of having an
ARCHITECTURE.md file.
