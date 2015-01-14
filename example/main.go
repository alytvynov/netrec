package main

import (
	"bytes"
	"log"
	"net"
	"net/http"

	"github.com/alytvynov/netrec"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	if err = http.Serve(netrec.NewRecordListener(l, printCb), nil); err != nil {
		log.Println(err)
	}
}

func printCb(in *bytes.Buffer, out *bytes.Buffer) {
	log.Printf("request:\n%s", in)
	log.Printf("response:\n%s", out)
}
