// parse.go
//
// Copyright © 2018 by Ollivier Robert <roberto@keltia.net>

package proxy // import "github.com/keltia/proxy"

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	proxyTag = "proxy"

	// MyVersion is our API Version
	MyVersion = "0.9.2"

	// MyName is the library name
	MyName = "proxy"
)

var (
	ctx Context
)

type Context struct {
	proxyauth string
	level     int
	Log       *log.Logger
}

func init() {
	// Default is stderr
	ctx = Context{Log: log.New(os.Stderr, "", log.LstdFlags)}
}

// ErrNoAuth is just to say we do not use auth for proxy
var ErrNoAuth = fmt.Errorf("no proxy auth")

func SetupProxyAuth() (proxyauth string, err error) {
	// Try to load $HOME/.netrc or file pointed at by $NETRC
	user, password := loadNetrc()

	if user != "" {
		verbose("Proxy user %s found.", user)
	}

	err = ErrNoAuth

	// Do we have a proxy user/password?
	if user != "" && password != "" {
		auth := fmt.Sprintf("%s:%s", user, password)
		proxyauth = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		err = nil

		// Store copy in context
		ctx.proxyauth = proxyauth
	}
	return
}

// SetLevel change the log level (currently 0, 1 or 2)
func SetLevel(level int) {
	ctx.level = level
	debug("logging level set to %d", level)
}

// SetLog allow to change the default logger
func SetLog(l *log.Logger) {
	ctx.Log = l
	debug("logging logger set to %v", l)

}

// GetAuth returns the proxyauth encoded string
func GetAuth() string {
	debug("returns cached credentials")
	return ctx.proxyauth
}

// SetupTransport is the way to have a custom http client
func SetupTransport(str string) (*http.Request, *http.Transport) {
	/*
	   Proxy code taken from https://github.com/LeoCBS/poc-proxy-https/blob/master/main.go
	*/
	myurl, err := url.Parse(str)
	if err != nil {
		log.Printf("error parsing %s: %v", str, err)
		return nil, nil
	}

	req, err := http.NewRequest("GET", str, nil)
	if err != nil {
		debug("error: req is nil: %v", err)
		return nil, nil
	}
	req.Header.Set("Host", myurl.Host)
	req.Header.Add("User-Agent", fmt.Sprintf("%s/%s", MyName, MyVersion))

	// Get proxy URL
	proxyURL := getProxy(req)
	if ctx.proxyauth != "" {
		req.Header.Add("Proxy-Authorization", ctx.proxyauth)
	}

	transport := &http.Transport{
		Proxy:              http.ProxyURL(proxyURL),
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		ProxyConnectHeader: req.Header,
	}
	debug("transport=%#v", transport)
	return req, transport
}

// Private functions

func getProxy(req *http.Request) (uri *url.URL) {
	uri, err := http.ProxyFromEnvironment(req)
	if err != nil {
		verbose("no proxy in environment")
		uri = &url.URL{}
	} else if uri == nil {
		verbose("No proxy configured or url excluded")
	}
	return
}

// loadNetrc supports a subset of the original ftp(1) .netrc file.
/*
We support:

  machine
  default
  login
  password

Format:
  machine <host> login <user> password <pass>
*/
func loadNetrc() (user, password string) {
	var dnetrc string

	// is $NETRC defined?
	dnetVar := os.Getenv("NETRC")

	// Allow override
	if dnetVar == "" {
		dnetrc = netrcFile
	} else {
		dnetrc = dnetVar
	}

	if dnetrc == "ignore" {
		return "", ""
	}

	verbose("NETRC=%s", dnetrc)

	// First check for permissions
	fh, err := os.Open(dnetrc)
	if err != nil {
		verbose("warning: can not find/read %s: %v", dnetrc, err)
		return "", ""
	}
	defer fh.Close()

	// Now check permissions
	st, err := fh.Stat()
	if err != nil {
		verbose("unable to stat: %v", err)
		return "", ""
	}

	if (st.Mode() & 077) != 0 {
		verbose("invalid permissions, must be 0400/0600")
		return "", ""
	}

	verbose("now parsing")
	user, password = parseNetrc(fh)
	return
}

/*
   Format:
   machine proxy|default login <user> password <pass>
*/
// parseDbrc loads the file format historically defined by DBI::Dbrc
func parseNetrc(r io.Reader) (user, password string) {
	verbose("found netrc")

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			break
		}

		flds := strings.Split(line, " ")
		debug("%s: %d fields", line, len(flds))

		if flds[0] != "machine" {
			verbose("machine is not the first word")
			continue
		}

		// Check what we need
		if len(flds) != 6 {
			verbose("bad format")
			continue
		}

		if flds[1] == proxyTag || flds[1] == "default" {

			if flds[2] == "login" && flds[4] == "password" {
				user = flds[3]
				password = flds[5]
				verbose("got %s/default entry for user %s", proxyTag, user)
			}
			break
		}
	}
	if err := s.Err(); err != nil {
		verbose("error reading netrc: %v", err)
		return "", ""
	}

	debug("nothing found for %s", proxyTag)

	if user == "" {
		verbose("no user/password for %s/default in netrc", proxyTag)
	}

	return
}
