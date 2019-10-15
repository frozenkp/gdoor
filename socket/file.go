package socket

import(
  "fmt"
  "io"
  "os"
  "strconv"

  "../debug"
)

func (sock Socket) SendFile(fileName string)error{
  file, err := os.Open(fileName)
  if err != nil {
    debug.Println(err)
    sock.Write("<ERROR>")
    return err
  }
  defer file.Close()

  fileInfo, err := file.Stat()
  if err != nil {
    debug.Println(err)
    sock.Write("<ERROR>")
    return err
  }

  // send size
  sock.Write(fmt.Sprintf("%x", fileInfo.Size()))

  // wait for client create
  resp, _ := sock.Read()
  if resp != "<OK>" {
    return nil
  }

  // send file content
  buf := make([]byte, 1024)
  for {
    cnt, err := file.Read(buf)
    if err == io.EOF {
      break
    }
    sock.Write(string(buf[:cnt]))
  }

  return nil
}

func (sock Socket) RecvFile(fileName string) error {
  fileSize, _ := sock.Read()
  if fileSize == "<ERROR>" {
    return nil
  }
  fileSize_i, _ := strconv.ParseInt(fileSize, 16, 64)

  file, err := os.Create(fileName)
  if err != nil {
    debug.Println(err)
    sock.Write("<ERROR>")
    return err
  }
  sock.Write("<OK>")
  defer file.Close()

  var receivedBytes int64 = 0
  for receivedBytes != fileSize_i {
    recv, _ := sock.Read()
    fmt.Fprintf(file, "%s", recv)
    receivedBytes += int64(len(recv))
  }

  return nil
}
