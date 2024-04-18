server:
	go build -tags debug -tags server server.go

client:
	go build -tags release -tags client client.go

client_debug:
	go build -tags debug -tags client client.go

