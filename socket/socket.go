package socket

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"io"
	"log"
	"net"

	"github.com/frozenkp/gopwn"
	"golang.org/x/crypto/chacha20"
)

type Socket struct {
	conn   net.Conn
	key    []byte
	nonce  []byte
	cipher *chacha20.Cipher
}

func Init(conn net.Conn, asymmetricKey []byte, client bool) (Socket, error) {
	// key exchange
	keyInfo, err := keyExchange(conn, asymmetricKey, client)
	if err != nil {
		return Socket{}, err
	}
	key := keyInfo[:chacha20.KeySize]
	nonce := keyInfo[chacha20.KeySize:]

	// build ChaCha20 cipher
	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return Socket{}, err
	}

	// build Socket
	s := Socket{
		conn:   conn,
		key:    key,
		nonce:  nonce,
		cipher: cipher,
	}

	return s, nil
}

func keyExchange(conn net.Conn, asymmetricKey []byte, client bool) ([]byte, error) {
	if client { // client
		// parse key
		publicKey, err := x509.ParsePKCS1PublicKey(asymmetricKey)
		if err != nil {
			return nil, err
		}

		// generate key
		keyInfo := make([]byte, chacha20.KeySize+chacha20.NonceSize)
		_, err = rand.Read(keyInfo)
		if err != nil {
			return nil, err
		}

		// encrypt key
		cipherKeyInfo, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, keyInfo, nil)
		if err != nil {
			return nil, err
		}

		// send key
		_, err = conn.Write(cipherKeyInfo)

		return keyInfo, err

	} else { // server
		// parse key
		privateKey, err := x509.ParsePKCS1PrivateKey(asymmetricKey)
		if err != nil {
			return nil, err
		}

		// receive key
		cipherKeyInfo := make([]byte, 256)
		_, err = io.ReadFull(conn, cipherKeyInfo)
		if err != nil {
			return nil, err
		}

		// decrypt key
		keyInfo, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cipherKeyInfo, nil)

		return keyInfo, err
	}
}

func (s Socket) Read() ([]byte, error) {
	// read size (4 bytes)
	sizeB := make([]byte, 4)
	_, err := io.ReadFull(s.conn, sizeB)
	if err != nil {
		return nil, err
	}
	size := gopwn.U32(string(sizeB))

	// read data ($size bytes)
	cipherMsg := make([]byte, size)
	_, err = io.ReadFull(s.conn, cipherMsg)
	if err != nil {
		return nil, err
	}

	// decrypt with ChaCha20
	plainMsg := make([]byte, size)
	s.cipher.XORKeyStream(plainMsg, cipherMsg)

	return plainMsg, nil
}

func (s Socket) Write(msg []byte) {
	// write size (4 bytes)
	s.conn.Write([]byte(gopwn.P32(len(msg))))

	// encrypt data
	cipherMsg := make([]byte, len(msg))
	s.cipher.XORKeyStream(cipherMsg, msg)

	// write data ($size bytes)
	s.conn.Write(cipherMsg)
}

func (s Socket) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s Socket) Close() {
	s.conn.Close()
}
