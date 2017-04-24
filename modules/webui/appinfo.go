package webui

import (
	"encoding/json"
	"fmt"
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/config"
	"net/http"
)

// ShodanAppInfo holds meta information requested by the UI on load.
type ShodanAppInfo struct {
	Name            string `json:"name"`
	Description     string `json:"des"`
	Icon            string `json:"icon"`
	ID              string `json:"id"`
	VersionMajor    string `json:"version_major"`
	VersionMinor    string `json:"version_minor"`
	VersionRevision string `json:"version_revision"`
	VersionGit      string `json:"version_git"`
	SingleGuildMode string `json:"single_guild_mode"`
}

// CurrentApplicationInfo gets the app information from Discord
func CurrentApplicationInfo() (*ShodanAppInfo, error) {
	info := &ShodanAppInfo{}
	app, err := mod.shodan.GetDiscord().Application("@me")
	if err != nil {
		return nil, shodan.WrapError(err)
	}
	info.Name = app.Name
	info.Icon = app.Icon
	info.ID = app.ID
	info.Description = app.Description

	info.VersionMajor = fmt.Sprint(config.VersionMajor)
	info.VersionMinor = fmt.Sprint(config.VersionMinor)
	info.VersionRevision = fmt.Sprint(config.VersionRevision)
	info.VersionGit = config.VersionGitHash
	info.SingleGuildMode = config.SingleGuildMode

	return info, nil
}

func handleGetAppInfo(w http.ResponseWriter, _ *http.Request) {
	appInfo, err := CurrentApplicationInfo()
	if err != nil {
		shodan.ReportThreadError(false, err)
	}
	responseBytes, err := json.Marshal(appInfo)
	if err != nil {
		shodan.ReportThreadError(false, err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBytes)
}
