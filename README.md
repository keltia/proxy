netrc
=====

[![GitHub release](https://img.shields.io/github/release/keltia/proxy.svg)](https://github.com/keltia/proxy/releases)
[![GitHub issues](https://img.shields.io/github/issues/keltia/proxy.svg)](https://github.com/keltia/proxy/issues)
[![Go Version](https://img.shields.io/badge/go-1.10-blue.svg)](https://golang.org/dl/)
[![Build Status](https://travis-ci.org/keltia/proxy.svg?branch=master)](https://travis-ci.org/keltia/proxy)
[![GoDoc](http://godoc.org/github.com/keltia/proxy?status.svg)](http://godoc.org/github.com/keltia/proxy)
[![SemVer](http://img.shields.io/SemVer/2.0.0.png)](https://semver.org/spec/v2.0.0.html)
[![License](https://img.shields.io/pypi/l/Django.svg)](https://opensource.org/licenses/BSD-2-Clause)
[![Go Report Card](https://goreportcard.com/badge/github.com/keltia/proxy)](https://goreportcard.com/report/github.com/keltia/proxy)

Go library to load and prepare proxy authentication by parsing the `netrc` file as defined in `ftp(1)`.  Its purpose is now mainly to have a standard way of specifying the credentials for proxy authentication.

## Requirements

* Go >= 1.10

## Installation

This is a pure library, there is no associated command (like I do in some of my other packages such as [RIPE Atlas](https://github.com/keltia/ripe-atlas/) or [Cryptcheck](https://github.com/keltia/cryptcheck/)).

Installation is like many Go libraries with a simple

    go get github.com/keltia/proxy

`Proxy` also has `vgo` support & metadata (see the articles on [vgo](https://research.swtch.com/vgo-intro)).  It respects the [Semantic Versioning](https://research.swtch.com/vgo-import) principle with tagged releases.

## API Usage

The API is very simple in `net/http`.  You have to create a custom transport and look for credentials.

The main work is done by `SetupAuthProxy()` which looks at the standard file `.netrc` file.  This file was defined a long time by the `ftp(1)` command to store FTP sites' credentials.  We (ab)use it with a special **site** called "proxy".

The goal is to avoid polluting (and leaking) your credentials in the environment variable.

    import "github.com/keltia/proxy"
    
    authstr, err := proxy.SetupProxyAuth()
    
This looks for proxy credentials and store that internally.  If you need the credentials later, you can still call `GetAuth()`:

    authstr := proxy.GetAuth()

`autstr` is suitable for inclusing in a `Proxy-Authorization` standard HTTP header like this (this only support Basic Authentication):

    req.Header.Add("Proxy-Authorization", authstr)

To create the tailored HTTP `Transport`, you can use `SetupTransport()`.

    req, transport := proxy.SetupTransport(URL)
    
URL is there to trigger the search for the various proxy definitions (through the environment variables like `HTTP_PROXY` or other means that are supported by `net/http`.

There are also two functions dealing with logging, log levels and stuff:

    proxy.SetLevel(N)         // 0 (default), 1 (verbose), 2 (debug)

    proxy.SetLog(logger)      // logger is a *log.Logger object

By default, nothing is logged but if you set to 1 or more, the default is to log to Stderr in a fairly classic way.

## `netrc` file

On UNIX systems like FreeBSD, macOS or Linux, the `.netrc`file is located in the user's home directory (aka `$HOME`).  On Windows, I have decided to emulate this by looking for a `netrc` file (no ".") located in the AFAIK traditional location, designed by the `%LOCALAPPDATA%` variable.

Format:

    machine HOST username USER password PASS

in our case, `HOST` **must** be **`proxy`** or **`default`**.

If there is no `netrc` file or no `proxy` entry with credentials, the HTTP proxy can still be used but without authentication.

## License

This is under the 2-Clause BSD license, see `LICENSE.md`.

## History

I originally wrote this code for the [erc-cimbl](https://github.com/keltia/erc-cimbl/) project and have re-used it enough time to think about putting it into its own module.

## Contributing

Please see CONTRIBUTING.md for some simple rules.
