package oauth

import (
	"flag"
	"github.com/innovandalism/shodan"
)

type Module struct {
	clientid     *string
	clientsecret *string
	returnuri    *string
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
}

func (m *Module) Attach(session *shodan.Shodan) {
	session.GetMux().HandleFunc("/oauth/authenticate/", handleAuthenticate)
	session.GetMux().HandleFunc("/oauth/callback/", handleCallback)
	session.GetMux().HandleFunc("/oauth/testtoken/", handleTokenTest)
}
