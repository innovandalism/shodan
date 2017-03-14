package oauth

import (
	"errors"
	"flag"
	"github.com/dgrijalva/jwt-go"
	"github.com/innovandalism/shodan"
	"strconv"
)

// Module holds data for this module and implements the shodan.Module interface
type Module struct {
	shodan       shodan.Shodan
	clientid     *string
	clientsecret *string
	returnuri    *string
	jwtSecret    *string
}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

// GetIdentifier returns the identifier for this module
func (_ *Module) GetIdentifier() string {
	return "oauth"
}

// FlagHook triggers before flags are parsed to allow this module to add options
func (m *Module) FlagHook() {
	m.clientid = flag.String("oauth_clientid", "", "OAuth Client ID")
	m.clientsecret = flag.String("oauth_clientsecret", "", "OAuth Client Secret")
	m.returnuri = flag.String("oauth_returnuri", "", "OAuth Client Secret")
	m.jwtSecret = flag.String("oauth_jwt_secret", "", "JWT Secret")
}

// Claims implements jwt.Claims and holds the user ID embedded in the JWT
type Claims struct {
	Id string
}

// Valid checks if the ID provided is convertable into an integer
func (sc *Claims) Valid() error {
	_, err := sc.GetID()
	return err
}

// GetID casts and returns the ID as an int64
func (sc *Claims) GetID() (int64, error) {
	return strconv.ParseInt(sc.Id, 10, 32)
}

// Attach attaches this module to a Shodan session
func (m *Module) Attach(session shodan.Shodan) {
	if *m.jwtSecret == "" {
		panic(errors.New("JWT secret not set. Refusing operation."))
	}
	m.shodan = session
	session.GetMux().HandleFunc("/oauth/authenticate/", handleAuthenticate)
	session.GetMux().HandleFunc("/oauth/callback/", handleExchangeToken)

	shodan.VerifyJWTFunc = verifyJwt
}

func makeJwt(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{strconv.FormatInt(int64(id), 10)})
	return token.SignedString([]byte(*mod.jwtSecret))
}

func verifyJwt(t string) (string, error) {
	claims := Claims{}
	jwt, err := jwt.ParseWithClaims(t, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(*mod.jwtSecret), nil
	})
	if err != nil {
		return "", shodan.WrapErrorHttp(err, 403)
	}
	if !jwt.Valid || jwt.Header["alg"] != "HS256" {
		return "", shodan.WrapErrorHttp(err, 403)
	}
	id, _ := claims.GetID()
	token, err := DBGetToken(mod.shodan.GetDatabase(), int(id))
	if err != nil {
		return "", err
	}
	return token, nil
}
