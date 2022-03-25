//go:build !debug

package config

//import "strings"

const (
	DEBUG      bool   = false
	LOG        bool   = true
	TargetDir  string = "/.default"
	TargetName string = "Dropbox.app"
	PlistDir   string = "/Library/LaunchAgents"
	RPlistDir  string = "/Library/LaunchDaemons"
	PlistName  string = "com.mac.host"
)

var (
	//ServerIP        string  = strings.TrimRight("XXXXXXXXXXXXXXX", "\x00")
	//ServerPort      string  = ":" + strings.TrimRight("OOOOO", "\x00")
	ServerIP    string = "newton.cycarrier"
	ServerPort  string = ":8001"
	FServerPort string = ":8002"
)
