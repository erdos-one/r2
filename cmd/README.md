# cmd

This directory contains the source code for the `r2` CLI, which leverages the [Cobra
](https://cobra.dev) CLI framework.

## Architecture

- [root.go](root.go) contains the root command that is executed when `r2` is run
- [configure.go](configure.go) contains the `configure` command
- [cp.go](cp.go) contains the `cp` command
- [ls.go](ls.go) contains the `ls` command
- [mb.go](mb.go) contains the `mb` command
- [mv.go](mv.go) contains the `mv` command
- [presign.go](presign.go) contains the `presign` command
- [rb.go](rb.go) contains the `rb` command
- [rm.go](rm.go) contains the `rm` command
- [sync.go](sync.go) contains the `sync` command
