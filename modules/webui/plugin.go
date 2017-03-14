package webui

import (
	"flag"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/bindata"
	"log"
	"net/http"
)

// Module holds data for this module and implements the shodan.Module interface
type Module struct {
	shodan  shodan.Shodan
	webroot *string
}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

// GetIdentifier returns the identifier for this module
func (_ *Module) GetIdentifier() string {
	return "webui"
}

// FlagHook triggers before flags are parsed to allow this module to add options
func (m *Module) FlagHook() {
	m.webroot = flag.String("webui_webroot", "", "Specify to serve webui from filesystem; serve from embedded binary data if omitted")
}

// Attach attaches this module to a Shodan session
func (m *Module) Attach(session shodan.Shodan) {
	m.shodan = session

	session.GetMux().HandleFunc("/app.json", handleGetAppInfo)
	session.GetMux().HandleFunc("/user/profile/", handleGetUserProfile)
	session.GetMux().HandleFunc("/guilds/{guild}/roles/", handleGetRolesForGuild)
	if *m.webroot != "" {
		log.Printf("Serving webui from %s\n", *m.webroot)
		session.GetMux().PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir(*m.webroot))))
	} else {
		log.Print("Serving webui from bindata\n")
		session.GetMux().PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(getPublicAssetFS())))
	}
	session.GetMux().Handle("/", http.RedirectHandler("/ui/", 301))
}

func getPublicAssetFS() *assetfs.AssetFS {
	return &assetfs.AssetFS{Asset: bindata.Asset, AssetDir: bindata.AssetDir, AssetInfo: bindata.AssetInfo, Prefix: "assets/webui/public"}
}
