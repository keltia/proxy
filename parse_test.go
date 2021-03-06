// proxy_test.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

package proxy

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// GoodAuth is test:test
	GoodAuth = "Basic dGVzdDp0ZXN0"
)

func TestVersion(t *testing.T) {
	str := Version()
	require.Equal(t, MyVersion, str)
}

func setvars(t *testing.T) {
	// Insert our values
	require.NoError(t, os.Setenv("HTTP_PROXY", "http://proxy:8080/"))
	require.NoError(t, os.Setenv("HTTPS_PROXY", "http://proxy:8080/"))
	require.NoError(t, os.Setenv("http_proxy", "http://proxy:8080/"))
	require.NoError(t, os.Setenv("https_proxy", "http://proxy:8080/"))
}

func unsetvars(t *testing.T) {
	// Remove our values
	require.NoError(t, os.Unsetenv("HTTP_PROXY"))
	require.NoError(t, os.Unsetenv("HTTPS_PROXY"))
	require.NoError(t, os.Unsetenv("http_proxy"))
	require.NoError(t, os.Unsetenv("https_proxy"))
}

// --- setupProxyAuth
func TestSetupProxyAuthNoNetrc(t *testing.T) {
	f := filepath.Join(".", "test/no-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	auth, err := SetupProxyAuth()
	assert.Error(t, err, "should be an error")
	assert.Equal(t, ErrNoAuth, err)
	assert.Empty(t, auth)

	os.Unsetenv("NETRC")
}

func TestSetupProxyAuthVerboseNoNetrc(t *testing.T) {
	SetLevel(1)

	f := filepath.Join(".", "test/no-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	auth, err := SetupProxyAuth()
	assert.Error(t, err, "should be an error")
	assert.Equal(t, ErrNoAuth, err)
	assert.Empty(t, auth)
	SetLevel(0)

	os.Unsetenv("NETRC")
}

func TestSetupProxyAuth(t *testing.T) {
	f := filepath.Join(".", "test/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	// We must ensure propre perms
	err = os.Chmod(f, 0600)
	require.NoError(t, err)

	auth, err := SetupProxyAuth()
	assert.NoError(t, err, "no error")
	assert.Equal(t, GoodAuth, auth)

	os.Unsetenv("NETRC")
}

func TestSetupProxyAuthVerbose(t *testing.T) {
	SetLevel(1)

	f := filepath.Join(".", "test/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	// We must ensure propre perms
	err = os.Chmod(f, 0600)
	require.NoError(t, err)

	auth, err := SetupProxyAuth()
	assert.NoError(t, err, "no error")
	assert.Equal(t, GoodAuth, auth)
	SetLevel(0)

	os.Unsetenv("NETRC")
}

// -- loadNetrc
func TestLoadNetrcNoFile(t *testing.T) {
	f := filepath.Join(".", "test/no-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	user, password := loadNetrc()
	assert.EqualValues(t, "", user, "null user")
	assert.EqualValues(t, "", password, "null password")

	os.Unsetenv("NETRC")
}

func TestLoadNetrcZero(t *testing.T) {
	err := os.Setenv("NETRC", filepath.Join(".", "test/zero-netrc"))
	require.NoError(t, err)

	user, password := loadNetrc()
	assert.EqualValues(t, "", user, "test user")
	assert.EqualValues(t, "", password, "test password")

	os.Unsetenv("NETRC")
}

func TestLoadNetrcVarEmpty(t *testing.T) {
	err := os.Setenv("NETRC", "ignore")
	require.NoError(t, err)

	user, password := loadNetrc()
	assert.EqualValues(t, "", user, "test user")
	assert.EqualValues(t, "", password, "test password")

	os.Unsetenv("NETRC")
}

func TestLoadNetrcPerms(t *testing.T) {
	f := filepath.Join(".", "test/perms-netrc")
	err := os.Setenv("NETRC", f)
	assert.NoError(t, err)

	err = os.Chmod(f, 0644)
	require.NoError(t, err)

	user, password := loadNetrc()
	err = os.Chmod(f, 0600)
	require.NoError(t, err)

	assert.EqualValues(t, "", user, "test user")
	assert.EqualValues(t, "", password, "test password")

	os.Unsetenv("NETRC")
}

func TestLoadNetrcGood(t *testing.T) {
	f := filepath.Join(".", "test/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	// We must ensure propre perms
	err = os.Chmod(f, 0600)
	require.NoError(t, err)

	user, password := loadNetrc()
	assert.EqualValues(t, "test", user, "test user")
	assert.EqualValues(t, "test", password, "test password")

	os.Unsetenv("NETRC")
}

func TestLoadNetrcGoodVerbose(t *testing.T) {
	SetLevel(1)

	f := filepath.Join(".", "test/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	// We must ensure propre perms
	err = os.Chmod(f, 0600)
	require.NoError(t, err)

	user, password := loadNetrc()
	assert.EqualValues(t, "test", user, "test user")
	assert.EqualValues(t, "test", password, "test password")
	SetLevel(0)

	os.Unsetenv("NETRC")
}

func TestLoadNetrcBad(t *testing.T) {
	f := filepath.Join(".", "test/bad-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	// We must ensure propre perms
	err = os.Chmod(f, 0600)
	require.NoError(t, err)

	user, password := loadNetrc()
	assert.EqualValues(t, "", user, "test user")
	assert.EqualValues(t, "", password, "test password")

	os.Unsetenv("NETRC")
}

func TestGetAuth(t *testing.T) {
	f := filepath.Join(".", "test/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	// We must ensure propre perms
	err = os.Chmod(f, 0600)
	require.NoError(t, err)
	auth, err := SetupProxyAuth()
	assert.NoError(t, err)

	str := GetAuth()
	assert.Equal(t, auth, str)

	os.Unsetenv("NETRC")
}

func TestSetLog(t *testing.T) {
	nl := log.New(os.Stderr, "", log.Lshortfile)
	SetLog(nl)
	assert.EqualValues(t, nl, ctx.Log)
}

func TestSetLevel(t *testing.T) {
	SetLevel(1)
	assert.Equal(t, 1, ctx.level)

	SetLevel(2)
	assert.Equal(t, 2, ctx.level)
	SetLevel(0)
}

func TestSetupTransport(t *testing.T) {
	req, trsp := SetupTransport("https://www.example.com/")
	assert.NotNil(t, req)
	assert.NotNil(t, trsp)
}

func TestSetupTransport2(t *testing.T) {
	req, trsp := SetupTransport(":foo")
	assert.Nil(t, req)
	assert.Nil(t, trsp)
}

func TestGetProxy(t *testing.T) {
	unsetvars(t)
	req, err := http.NewRequest("GET", "https://www.example.com/", nil)
	assert.NotNil(t, req)
	assert.NoError(t, err)
	uri := getProxy(req)
	assert.Nil(t, uri)
}

func TestGetProxySet(t *testing.T) {
	setvars(t)

	req, err := http.NewRequest("GET", "http://www.example.com/", nil)
	assert.NotNil(t, req)
	assert.NoError(t, err)

	u, err := http.ProxyFromEnvironment(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, u)

	uri := getProxy(req)
	assert.NotNil(t, uri)

	assert.EqualValues(t, uri, u)

	prx, _ := url.Parse("http://proxy:8080/")
	assert.EqualValues(t, prx, uri)

	unsetvars(t)
}
