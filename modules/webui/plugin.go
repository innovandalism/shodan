package webui

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/bindata"
	"net/http"
)

type Module struct {
	shodan *shodan.Shodan
}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

func (_ *Module) GetIdentifier() string {
	return "webui"
}

func (m *Module) FlagHook() {

}

func (m *Module) Attach(session *shodan.Shodan) {
	session.GetMux().Handle("/", http.FileServer(GetPublicAssetFS()))
	registerStatusPage(session)
}

func GetPublicAssetFS() *assetfs.AssetFS {
	return &assetfs.AssetFS{Asset: bindata.Asset, AssetDir: bindata.AssetDir, AssetInfo: bindata.AssetInfo, Prefix: "assets/webui/public"}
}
