package main

import (
	"log"
	"net/http"
	"restapi/chapter_09/encryptService/helpers"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	svc := helpers.EncryptServiceInstance{}
	encryptHandler := httptransport.NewServer(helpers.MakeEncryptEndpoint(svc), helpers.DecodeEncryptRequest, helpers.EncodeResponse)
	decryptHandler := httptransport.NewServer(helpers.MakeDecryptEndpoint(svc), helpers.DecodeDecryptRequest, helpers.EncodeResponse)

	http.Handle("/encrypt", encryptHandler)
	http.Handle("/decrypt", decryptHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}