server:
	go build -tags debug server.go

client_linux_debug:
	go build -tags debug client.go

client_mac:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s" client.go

client_mac_debug:
	GOOS=darwin GOARCH=amd64 go build -tags debug client.go

build: server client_mac

build_linux: server client_linux_debug

build_debug: server client_mac_debug
