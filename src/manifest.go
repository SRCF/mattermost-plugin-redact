package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

var manifest *model.Manifest

const manifestStr = `
{
  "id": "srcf.redact",
  "name": "Message Redaction Plugin",
  "description": "This plugin allows team administrators to delete old messages",
  "version": "1.0.0",
  "min_server_version": "5.6.0",
  "server": {
    "executables": {
      "linux-amd64": "plugin-linux-amd64",
      "darwin-amd64": "plugin-darwin-amd64",
      "windows-amd64": "plugin-windows-amd64"
    },
    "executable": ""
  }
}
`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}
