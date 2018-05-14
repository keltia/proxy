// config_windows.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>
// +build windows

package proxy

import (
	"os"
	"path/filepath"
)

/*
File location: %LOCALAPPDATA%\netrc
*/
var (
	netrcFile = filepath.Join(os.Getenv("%LOCALAPPDATA%"), "netrc")
)
