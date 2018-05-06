# Go Version Manager

> GVM is a simple and clean go version manager in GO itself. 
It provides command line interface to manage go versions.

## Background

It is a hobby project which I made while learning Go. Go unlike Node.js and other languages does not have a
standard version manager, there are some open source alternatives but they are too buggy to be used.
I developed this keeping simplicity in mind, it solves all the painful tasks involved when someone is dealing
with multiple version of go. This is not a full fledged tool to meet all the requirements, but it performs
well to serve the purpose it was made for.

Do give a read to [Installing a new go version](/docs/new_installation.md).

## Install

Clone the github repository, or get the latest release from github.

```bash
$ cd gvm
$ dep ensure
$ make build
```

Add **gvm** location to $PATH

## Usage

```
Usage:
  gvm [flags]
  gvm [command]

Available Commands:
  help        Help about any command
  install     Installs the version of go mentioned against this flag
  list        List local version of go available for use
  list-remote List remote version of go available
  uninstall   Uninstall the specified version of go
  version     Displays the version of the current build of gvm

Flags:
  -h, --help   help for gvm

Use "gvm [command] --help" for more information about a command.
```

#### Installing a go version

To install a go version run `gvm install go1.8`

GO installation using gvm are source first. Installation procedure fetches the release from official golang releases
and downloads it, caching it for later use. It then sets the necessery environment variable for compilation and compile the source
generating an environment file for each go version.

#### Uninstalling a go version

To uninstall a perviously installed go version run `gvm uninstall go1.8`

#### List available versions

GVM provides a command to list go version released by google from their official release website. To view this list use
`gvm list-remote` 

#### List installed versions

To view locally available and installed version of GO run `gvm list`

#### Activating a GO version

Installing a Go version creates an environment file for it, which can then be used by the shell to set up proper environment variables to run that version.
To use or activate a version of go just source the environment file.

```bash
source ~/.gvm/environments/go1.8
```

## License

This project is licensed under MIT license. View [License](/LICENSE.md)

## Links

* [moovweb/gvm](https://github.com/moovweb/gvm)
