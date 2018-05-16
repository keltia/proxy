// proxy_test.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

package proxy

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	// GoodAuth is test:test
	GoodAuth = "Basic dGVzdDp0ZXN0"
)

// --- setupProxyAuth
func TestSetupProxyAuthNoNetrc(t *testing.T) {
	f := filepath.Join(".", "test/no-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	_, err = SetupProxyAuth()
	assert.Error(t, err, "should be an error")
	assert.Equal(t, ErrNoAuth, err)
}

func TestSetupProxyAuthVerboseNoNetrc(t *testing.T) {
	SetLevel(1)

	f := filepath.Join(".", "test/no-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	_, err = SetupProxyAuth()
	assert.Error(t, err, "should be an error")
	assert.Equal(t, ErrNoAuth, err)
	SetLevel(0)
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
}

// -- loadNetrc
func TestLoadNetrcNoFile(t *testing.T) {
	f := filepath.Join(".", "test/no-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	user, password := loadNetrc()
	assert.EqualValues(t, "", user, "null user")
	assert.EqualValues(t, "", password, "null password")
}

func TestLoadNetrcZero(t *testing.T) {
	err := os.Setenv("NETRC", filepath.Join(".", "test/zero-netrc"))
	require.NoError(t, err)

	user, password := loadNetrc()
	assert.EqualValues(t, "", user, "test user")
	assert.EqualValues(t, "", password, "test password")
}

func TestLoadNetrcVarEmpty(t *testing.T) {
	err := os.Setenv("NETRC", "")
	require.NoError(t, err)

	user, password := loadNetrc()
	assert.EqualValues(t, "", user, "test user")
	assert.EqualValues(t, "", password, "test password")
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
}

func TestSetLog(t *testing.T) {
	nl := log.New(os.Stderr, MyName, log.Lshortfile)
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
