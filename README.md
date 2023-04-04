<p align="center">
<img src="docs/auditory.png" alt="yatas-logo" width="30%">
<p align="center">

# YATAS
[![codecov](https://codecov.io/gh/StanGirard/yatas-template/branch/main/graph/badge.svg?token=OFGny8Za4x)](https://codecov.io/gh/StanGirard/YATAS) [![goreport](https://goreportcard.com/badge/github.com/stangirard/yatas-template)](https://goreportcard.com/badge/github.com/stangirard/yatas)

Yet Another Testing &amp; Auditing Solution 

The goal of YATAS is to help you create a secure AWS environment without too much hassle. It won't check for all best practices but only for the ones that are important for you based on my experience. Please feel free to tell me if you find something that is not covered.

## Plugins implementation

Simply use this repository as a template and implement your own plugin.

Don't forget to change the name of the plugin in the `Makefile` and the `go.mod` file.


Add you code in `main.go` and in the function `RunPlugin` 
## Usage

Implement your own plugin and build it with the `make build` command. 


## How to tests ? 

Simply run `make install` and in your Yatas config instead of the version of the plugin use `local` and it will use the version you just installed.

## How to deploy ?

You can use the provided workflow to run your plugin in a CI environment. 
The plugin needs to have a binary that starts with `yatas-` and ends with the name of the plugin in the releases.