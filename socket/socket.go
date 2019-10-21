package socket

import(
  "encoding/base64"
  "net"
  "strings"
  "bufio"
  "fmt"
  "strconv"
  "io"

  "../crypto"
)

type Socket struct{
  conn  net.Conn
  key   string
}

func Init(conn net.Conn, key string)Socket{
  s := Socket{
    conn: conn,
    key: key,
  }
  if key == "" {
    s.key = "\x00"
  }

  return s
}

func (s Socket) Read()(string, error){
  r := bufio.NewReader(s.conn)
  header, err := r.ReadString('>')
  if err != nil {
    return "", err
  }

  length, _ := strconv.ParseInt(strings.Trim(header, "<>"), 16, 32)

  buf := make([]byte, length)
  cnt, _ := io.ReadFull(r, buf)

  content, err := base64.StdEncoding.DecodeString(string(buf[:cnt]))
  if err != nil {
    return "", err
  }

  return crypto.XOR(string(content), s.key), nil
}

func (s Socket) Write(resp string){
  w := bufio.NewWriter(s.conn)
  cipher := crypto.XOR(resp, s.key)
  content := base64.StdEncoding.EncodeToString([]byte(cipher))
  w.WriteString(fmt.Sprintf("<%x>%s", len(content), content))
  w.Flush()
}

func (s Socket) RemoteAddr() net.Addr {
  return s.conn.RemoteAddr()
}

func (s Socket) Close(){
  s.conn.Close()
}
