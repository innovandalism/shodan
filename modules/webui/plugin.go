package webui

import (
	"log"
	"net/http"
	"os"

	"github.com/innovandalism/shodan"
)

// Module holds data for this module and implements the shodan.Module interface
type Module struct {
	shodan  shodan.Shodan
	webroot string
}

var mod = Module{}

func init() {
	mod.webroot = os.Getenv("WEBUI_WEBROOT")
	shodan.Loader.LoadModule(&mod)
}

// GetIdentifier returns the identifier for this module
func (_ *Module) GetIdentifier() string {
	return "webui"
}

// Attach attaches this module to a Shodan session
func (m *Module) Attach(session shodan.Shodan) error {
	m.shodan = session

	session.GetMux().HandleFunc("/app.json", handleGetAppInfo)
	session.GetMux().HandleFunc("/user/profile/", handleGetUserProfile)
	session.GetMux().HandleFunc("/guilds/{guild}/roles/", handleGetRolesForGuild)
	if m.webroot != "" {
		log.Printf("Serving webui from %s\n", m.webroot)
		session.GetMux().PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir(m.webroot))))
	} else {
		shodan.ReportThreadError(true, shodan.Error("No webui path provided"))
	}
	session.GetMux().Handle("/", http.RedirectHandler("/ui/", 301))
	return nil
}
