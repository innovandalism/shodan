package webui

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/bindata"
	"net/http"
	"flag"
	"log"
)

type Module struct {
	shodan *shodan.Shodan
	webroot *string
}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

func (_ *Module) GetIdentifier() string {
	return "webui"
}

func (m *Module) FlagHook() {
	m.webroot = flag.String("webui_webroot", "", "Specify to serve webui from filesystem; serve from embeded binary data if omitted")
}

func (m *Module) Attach(session *shodan.Shodan) {
	m.shodan = session
	if *m.webroot != "" {
		log.Printf("Serving webui from %s\n", *m.webroot)
		session.GetMux().Handle("/", http.FileServer(http.Dir(*m.webroot)))
	} else {
		log.Print("Serving webui from bindata\n")
		session.GetMux().Handle("/", http.FileServer(GetPublicAssetFS()))
	}
	session.GetMux().HandleFunc("/app.json", handleGetAppInfo)
}

func GetPublicAssetFS() *assetfs.AssetFS {
	return &assetfs.AssetFS{Asset: bindata.Asset, AssetDir: bindata.AssetDir, AssetInfo: bindata.AssetInfo, Prefix: "assets/webui/public"}
}
