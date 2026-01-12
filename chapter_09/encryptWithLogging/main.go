package main

import (
	"log"
	"net/http"
	"os"
	"restapi/chapter_09/encryptWithLogging/helpers"

	kitlog "github.com/go-kit/kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stderr)
	var svc helpers.EncryptService
	svc = helpers.LoggingMiddleware{Logger: logger, Next: svc}
	encryptHandler := httptransport.NewServer(helpers.MakeEncryptEndpoint(svc), helpers.DecodeEncryptRequest, helpers.EncodeResponse)
	decryptHandler := httptransport.NewServer(helpers.MakeDecryptEndpoint(svc), helpers.DecodeDecryptRequest, helpers.EncodeResponse)

	http.Handle("/encrypt", encryptHandler)
	http.Handle("/decrypt", decryptHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}