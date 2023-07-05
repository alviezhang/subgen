package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/alviezhang/subgen/internal"
	"github.com/go-chi/chi"
)

func main() {
	id := flag.String("id", "", "Subscription id")
	filename := flag.String("config", "", "Source configuration file")
	port := flag.String("port", "", "API server port")

	flag.Parse()

	if *id == "" || *filename == "" {
		fmt.Println("Usage: genweb --id <id> --config <filename> --port <port>")
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Get("/sub", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()

		subType := queryParams.Get("type")
		inputID := queryParams.Get("id")

		if *id != inputID {
			http.NotFound(w, r)
			return
		}

		result, err := internal.Generate(*filename, subType)
		if err != nil {
			http.Error(w, "Generate failed", 500)
		}
		w.Write([]byte(result))
	})

	http.ListenAndServe(fmt.Sprintf(":%s", *port), r)
}
