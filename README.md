<p align="center"><img src="https://raw.githubusercontent.com/gobuffalo/buffalo/master/logo.svg" width="360"></p>

<p align="center">
<a href="https://pkg.go.dev/github.com/gobuffalo/cli"><img src="https://pkg.go.dev/badge/github.com/gobuffalo/cli" alt="PkgGoDev"></a>
<img src="https://github.com/gobuffalo/cli/workflows/Tests/badge.svg" alt="Tests Status" />
<a href="https://goreportcard.com/report/github.com/gobuffalo/cli"><img src="https://goreportcard.com/badge/github.com/gobuffalo/cli" alt="Go Report Card" /></a>
</p>

# Buffalo CLI

This is the repo for the Buffalo CLI. The Buffalo CLI is a tool to develop, test and deploy your Buffalo applications.

## Installation

Installing the Buffalo CLI requires Go 1.16 or older as it depends heavily on the embed package. Once you have ensured you installed Go 1.16 or older, you can install the Buffalo CLI by running the following command:

```
go install github.com/gobuffalo/cli/cmd/buffalo@latest
```
## Usage
Once installed, the Buffalo CLI can be used by invoking the `buffalo` command. To know more about the available commands, run the `buffalo help` command. or you can also get details on a specific command by running `buffalo help <command>`.
