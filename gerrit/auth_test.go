package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"github.com/stretchr/testify/mock"
	"golang.org/x/build/gerrit"
)

func TestMissingCookieAndCredentialsMeansAnonymousLogin(t *testing.T) {
	//given source without authentication
	sourceWithOutAuthentication := Source{}

	//when authentication manager is created
	authManager := NewAuthManager(sourceWithOutAuthentication)

	//then anonymous login is used
	assert.True(t, authManager.Anonymous())
}

func TestProvidedCredentialsAreSufficientForAuthentication(t *testing.T) {
	//given valid credentials
	validCredentials := [...]Source {{Username: "filled"}, {Username: "filled", Password:"filled"}}

	for _, credentials := range validCredentials {
		//and source with filled credentials
		sourceWithCredentials := credentials

		//and missing cookie information
		sourceWithCredentials.Cookies = ""

		//when authentication manager is created
		authManager := NewAuthManager(sourceWithCredentials)

		//then uses credentials to login
		assert.False(t, authManager.Anonymous(), fmt.Sprintf("Because credentials are valid: %+v", credentials))
	}
}

func TestProvidingOnlyPasswordIsntSufficientForAuthentication(t *testing.T) {
	//given invalid credentials
	invalidCredentials := Source{Password:"filled"}

	//and missing cookie information
	invalidCredentials.Cookies = ""

	//when authentication manager is created
	authManager := NewAuthManager(invalidCredentials)

	//then uses credentials to login
	assert.True(t, authManager.Anonymous(), fmt.Sprintf("Because credentials are valid: %+v", invalidCredentials))
}

func TestProvidingCookieIsSufficientForAuthentication(t *testing.T) {
	//given valid authentication cookies
	validAuthentications := [...]Source {
		{Cookies: "valid", Username: "filled"},
		{Cookies: "valid", Username: "filled", Password:"filled"},
		{Cookies: "valid", Password:"filled"},
		{Cookies: "valid"},
	}

	for _, validAuthentication := range validAuthentications {
		//when validAuthentication manager is created
		authManager := NewAuthManager(validAuthentication)

		//then uses cookies to login
		assert.False(t, authManager.Anonymous(), fmt.Sprintf("Because authentication is valid: %+v", validAuthentication))
	}
}

type authManagerMock struct {
	mock.Mock
}

func (m *authManagerMock) GerritAuth() (gerrit.Auth, error) { return nil, nil }
func (m *authManagerMock) GitConfigArgs() ([]string, error) { return nil, nil }
func (m *authManagerMock) Anonymous() bool { return true }
func (m *authManagerMock) Cleanup() {}

func TestReturnsFirstAuthManagerThatHandlesGivenAuthenticationData(t *testing.T) {
	//given valid authentication cookies
	validCookiesAndCredentials := Source{Cookies: "filled"}

	//and valid credentials
	validCookiesAndCredentials.Username = "filled"

	//and
	credentialAuthManager := new(authManagerMock)
	credentialAuthManager.On("Anonymous").Return(false)

	cookieAuthManager := new(authManagerMock)
	cookieAuthManager.On("Anonymous").Return(false)

	//when creates new authentication manager
	authManager := newAuthManager(validCookiesAndCredentials,
		func(source *Source) AuthManager { return credentialAuthManager },
		func(source *Source) AuthManager { return cookieAuthManager })

	//then chooses credential manager which was declared first
	assert.Equal(t, credentialAuthManager, authManager)
}

func TestReturnsLastAuthenticationManagerInCaseNoneHandlesGivenAuthenticationData(t *testing.T) {
	//given missing authentication cookies
	missingCookiesAndCredentials := Source{}

	//and
	credentialAuthManager := new(authManagerMock)
	credentialAuthManager.On("Anonymous").Return(true)

	cookieAuthManager := new(authManagerMock)
	cookieAuthManager.On("Anonymous").Return(true)

	anonymousAuthManager := newAnonymousAuth(&Source{})

	//when creates new authentication manager
	authManager := newAuthManager(missingCookiesAndCredentials,
		func(source *Source) AuthManager { return credentialAuthManager },
		func(source *Source) AuthManager { return cookieAuthManager },
		func(source *Source) AuthManager { return anonymousAuthManager })

	//then chooses anonymous authentication manager which was declared as last
	assert.Equal(t, anonymousAuthManager, authManager)
}
