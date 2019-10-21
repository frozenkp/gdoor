package socket

import(
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

  return crypto.XOR(string(buf[:cnt]), s.key), nil
}

func (s Socket) Write(resp string){
  w := bufio.NewWriter(s.conn)
  cipher := crypto.XOR(resp, s.key)
  w.WriteString(fmt.Sprintf("<%x>%s", len(cipher), cipher))
  w.Flush()
}

func (s Socket) RemoteAddr() net.Addr {
  return s.conn.RemoteAddr()
}

func (s Socket) Close(){
  s.conn.Close()
}
