package socket

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/frozenkp/gopwn"

	"gdoor/debug"
)

func (sock Socket) SendFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		debug.Println(err)
		sock.Write([]byte("<ERROR>"))
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		debug.Println(err)
		sock.Write([]byte("<ERROR>"))
		return err
	}

	// send size
	sock.Write([]byte(gopwn.P64(int(fileInfo.Size()))))

	// wait for client create
	resp, _ := sock.Read()
	if string(resp) != "<OK>" {
		return errors.New("Received side open file failed.")
	}

	// send file content
	buf := make([]byte, 1024)
	for {
		cnt, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		sock.Write(buf[:cnt])
	}

	return nil
}

func (sock Socket) RecvFile(fileName string) error {
	fileName = filepath.Base(fileName)
	fileSizeB, _ := sock.Read()
	if string(fileSizeB) == "<ERROR>" {
		return errors.New("Sent side open file failed.")
	}
	fileSize := int64(gopwn.U64(string(fileSizeB)))

	file, err := os.Create(fileName)
	if err != nil {
		debug.Println(err)
		sock.Write([]byte("<ERROR>"))
		return err
	}
	sock.Write([]byte("<OK>"))
	defer file.Close()

	var receivedBytes int64 = 0
	for receivedBytes != fileSize {
		recv, _ := sock.Read()
		file.Write(recv)
		receivedBytes += int64(len(recv))
	}

	return nil
}
