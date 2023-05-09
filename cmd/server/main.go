package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/loicsikidi/wif-go/cmd/server/static"
)

const PORT = "8080"

var httpRoot http.FileSystem

func main() {
	if static.IsEmbedded {
		contentStatic, err := fs.Sub(static.Box, "kodata")
		if err != nil {
			panic(err)
		}
		httpRoot = http.FS(contentStatic)
	} else {
		httpRoot = http.Dir(os.Getenv("KO_DATA_PATH"))
	}

	http.Handle("/", http.FileServer(httpRoot))
	log.Printf("server is listening on :%s", PORT)

	server := &http.Server{
		Addr:              ":" + PORT,
		ReadHeaderTimeout: 15 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
