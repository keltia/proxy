// utils.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

package proxy

// debug displays only if fDebug is set
func debug(str string, a ...interface{}) {
	if ctx.level >= 2 {
		ctx.Log.Printf(str, a...)
	}
}

// debug displays only if fVerbose is set
func verbose(str string, a ...interface{}) {
	if ctx.level >= 1 {
		ctx.Log.Printf(str, a...)
	}
}
