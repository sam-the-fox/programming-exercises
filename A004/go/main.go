package main

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
)

//go:embed template.html
var rawTmpl string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Println("meow")

	services, err := ReadServices()
	if err != nil {
		log.Printf("error creating services container: %v", err)
		os.Exit(1)
	}

	tmpl, err := template.New("main").Parse(rawTmpl)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		os.Exit(1)
	}

	go services.Query(ctx)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		services.mu.RLock()
		defer services.mu.RUnlock()

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		err := tmpl.Execute(w, services)
		if err != nil {
			log.Printf("error executing template: %v", err)
		}
	})

	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(":5001", nil)
	}()

	select {
	case <-ctx.Done():
		os.Exit(0)
	case <-errChan:
		fmt.Printf("error serving on :5001: %v", err)
		os.Exit(1)
	}
}
