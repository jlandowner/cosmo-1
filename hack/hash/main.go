package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"golang.org/x/crypto/argon2"
)

func main() {
	raw := os.Args[1]
	fmt.Println("raw   : ", raw)
	fmt.Println("hashed: ", string(WebAuthnID(raw)))
}

func WebAuthnID(x string) []byte {
	id := make([]byte, 64)
	hashed := argon2.IDKey([]byte(x), nil, 1, 2048, 4, 32)
	n := hex.Encode(id, hashed)
	if n != 64 {
		panic(fmt.Errorf("invalid hash length: n=%d", n))
	}
	return id
}
