package oauth

import (
	"flag"
	"github.com/innovandalism/shodan"
	"errors"
)

type Module struct {
	shodan       *shodan.Shodan
	clientid     *string
	clientsecret *string
	returnuri    *string
	jwtSecret   *string
}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

func (_ *Module) GetIdentifier() string {
	return "oauth"
}

func (m *Module) FlagHook() {
	m.clientid = flag.String("oauth_clientid", "", "OAuth Client ID")
	m.clientsecret = flag.String("oauth_clientsecret", "", "OAuth Client Secret")
	m.returnuri = flag.String("oauth_returnuri", "", "OAuth Client Secret")
	m.jwtSecret = flag.String("oauth_jwt_secret", "", "JWT Secret")
}

func (m *Module) Attach(session *shodan.Shodan) {
	if *m.jwtSecret == "" {
		panic(errors.New("JWT secret not set. Refusing operation."))
	}
	m.shodan = session
	session.GetMux().HandleFunc("/oauth/authenticate/", handleAuthenticate)
	session.GetMux().HandleFunc("/oauth/callback/", handleExchangeToken)
	session.GetMux().HandleFunc("/oauth/exchange/", handleGetUserProfile)
	session.GetMux().HandleFunc("/user/profile/", handleGetUserProfile)
}