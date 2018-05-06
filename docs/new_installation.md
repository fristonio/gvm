### Compiling and installing a new version of golang alongside with others

* Go 1.5+ removed the C compilers from the toolchain and replaced them with one written in Go. This creates a bootstrapping problem if we don't already have a working Go install. In order to compile Go 1.5+, make sure Go 1.4 is installed first.

* While installing go from source there is a catch in setting the environment variables so that it gets compiled at the required place and not somewhere else. These environment variables defines the installation behaviour for go.

* GOROOT_BOOTSTRAP
	* To install from source we need to bootstrap the go installation which requries either `gccgo` or an existing go installation. To make bootstrapping use the go version available on the system we need to set this environment variable to the GO_ROOT.

* Now we need to unset the existing environment variables associated with the parent installation. For this 
`unset GOARCH && unset GOOS && unset GOPATH && unset GOBIN && unset GOROOT`

* After this we set these environment variables for the new installation

```bash
export GOBIN=INSTALL_ROOT_PATH/bin &&
export PATH=$GOBIN:$PATH &&
export GOROOT=INSTALL_ROOT_PATH
```

Once this environment variable is set we are good to compile the package from source for this just cd to the source directory and run `./make.bash`

Alternatively one can directly use the platform dependent binary provided by of golang organization itself and use it setting the proper environemnt variables for it.
