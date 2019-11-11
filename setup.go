package main

import(
  "flag"
  "os"
  "log"
  "strings"
  "bufio"
  "io"
)

const (
  IP_PATTERN    string = "XXXXXXXXXXXXXXX"
  PORT_PATTERN  string = "OOOOO"
)

func main(){
  // flag
  bin := flag.String("b", "", "the target binary name.")
  ip := flag.String("i", "", "Specified IP.")
  port := flag.String("p", "", "Specified port.")

  flag.Parse()

  // Subtituted IP / Port
  SubIP := *ip + strings.Repeat("\x00", len(IP_PATTERN) - len(*ip))
  SubPort := *port + strings.Repeat("\x00", len(PORT_PATTERN) - len(*port) - 1)

  // target binary
  file, err := os.Open(*bin)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  // create patched file
  PatchedFile, err := os.Create(*bin + ".patched")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  // Read, search, and write
  w := bufio.NewWriter(PatchedFile)
  r := bufio.NewReader(file)
  for true {
    b, err := r.ReadByte()
    if err == io.EOF {
      break
    }

    if b == 'X' {
      bs, err := r.Peek(len(IP_PATTERN)-1)
      if err == nil {
        verified := true
        for _, v := range bs {
          if v != 'X' {
            verified = false
            break
          }
        }
        if verified {
          w.WriteString(SubIP)
          r.Discard(len(IP_PATTERN)-1)
          log.Println("IP pattern found.")
          continue
        }
      }
    } else if b == 'O' {
      bs, err := r.Peek(len(PORT_PATTERN)-1)
      if err == nil {
        verified := true
        for _, v := range bs {
          if v != 'O' {
            verified = false
            break
          }
        }
        if verified {
          w.WriteString(SubPort)
          r.Discard(len(PORT_PATTERN)-1)
          log.Println("Port pattern found.")
          continue
        }
      }
    }

    w.WriteByte(b)
  }
}
