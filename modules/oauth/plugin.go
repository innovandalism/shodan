package oauth

import (
	"flag"
	"github.com/innovandalism/shodan"
	"errors"
	"github.com/innovandalism/shodan/api"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"github.com/innovandalism/shodan/util"
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

type ShodanClaims struct{
	Id string
	Discord_uid string
}

func (sc *ShodanClaims) Valid() error {
	_, err := sc.GetID()
	return err
}

func (sc *ShodanClaims) GetID() (int64, error) {
	return strconv.ParseInt(sc.Id, 10, 32)
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

	api.VerifyJWTFunc = verifyJwt
}

func verifyJwt(t string) (string, error) {
	var sc ShodanClaims
	jwt, err := jwt.ParseWithClaims(t, &sc, func(t *jwt.Token) (interface{}, error) {
		return []byte(*mod.jwtSecret), nil
	})
	if err != nil {
		return "", util.WrapErrorHttp(err, 400)
	}
	if !jwt.Valid {
		return "", util.WrapErrorHttp(err, 403)
	}
	id, _ := sc.GetID()
	token, err := mod.shodan.GetDatabase().GetToken(int(id))
	if err != nil {
		return "", err
	}
	return token, nil
}