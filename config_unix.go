// config_unix.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>
// +build !windows

/*
File location: $HOME/.netrc
*/
package proxy

import (
	"os"
	"path/filepath"
)

/*
File location: $HOME/.netrc
*/
var (
	netrcFile = filepath.Join(os.Getenv("HOME"), ".netrc")
)
