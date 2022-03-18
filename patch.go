package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
)

const (
	IP_PATTERN   string = "XXXXXXXXXXXXXXX"
	PORT_PATTERN string = "OOOOO"
)

func main() {
	// flag
	bin := flag.String("b", "", "the target binary name.")
	ip := flag.String("i", "", "Specified IP.")
	port := flag.String("p", "", "Specified port.")

	flag.Parse()

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
			bs, err := r.Peek(len(IP_PATTERN) - 1)
			if err == nil {
				verified := true
				for _, v := range bs {
					if v != 'X' {
						verified = false
						break
					}
				}
				if verified {
					w.WriteString(*ip)
					for i := len(*ip); i < len(IP_PATTERN); i++ {
						w.WriteByte('\x00')
					}
					r.Discard(len(IP_PATTERN) - 1)
					log.Println("IP pattern found.")
					continue
				}
			}
		} else if b == 'O' {
			bs, err := r.Peek(len(PORT_PATTERN) - 1)
			if err == nil {
				verified := true
				for _, v := range bs {
					if v != 'O' {
						verified = false
						break
					}
				}
				if verified {
					w.WriteString(*port)
					for i := len(*port); i < len(PORT_PATTERN); i++ {
						w.WriteByte('\x00')
					}
					r.Discard(len(PORT_PATTERN) - 1)
					log.Println("Port pattern found.")
					continue
				}
			}
		}

		w.WriteByte(b)
	}
}
