package webui

import (
	"github.com/innovandalism/shodan/util"
	"net/http"
	"encoding/json"
	"github.com/innovandalism/shodan/config"
	"fmt"
)

type ShodanAppInfo struct{
	Name string `json:"name"`
	Description string `json:"des"`
	Icon string `json:"icon"`
	ID string `json:"id"`
	VersionMajor string `json:"version_major"`
	VersionMinor string  `json:"version_minor"`
	VersionRevision string `json:"version_revision"`
	VersionGit string `json:"version_git"`
}

func CurrentApplicationInfo() (*ShodanAppInfo, error) {
	info := &ShodanAppInfo{}
	app, err := mod.shodan.GetDiscord().Application("@me")
	if err != nil {
		return nil, util.WrapError(err)
	}
	info.Name = app.Name
	info.Icon = app.Icon
	info.ID = app.ID
	info.Description = app.Description

	info.VersionMajor = fmt.Sprint(config.VersionMajor)
	info.VersionMinor = fmt.Sprint(config.VersionMinor)
	info.VersionRevision = fmt.Sprint(config.VersionRevision)
	info.VersionGit = config.VersionGitHash

	return info, nil
}

func handleGetAppInfo(w http.ResponseWriter, _ *http.Request) {
	appInfo, err := CurrentApplicationInfo()
	if err != nil {
		util.ReportThreadError(false, err)
	}
	responseBytes, err := json.Marshal(appInfo)
	if err != nil {
		util.ReportThreadError(false, err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(responseBytes)
}