package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	addrFlag = flag.String("addr", ":80", "server address:port")
)

func main() {
	flag.Parse()

	http.Handle("/", appHandler(rootHandler))
	http.Handle("/solid-srv/", appHandler(solidHandler))
	err := http.ListenAndServe(*addrFlag, nil)
	if err != nil {
		log.Fatal(err)
	}
}

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err != nil {
		switch err := err.(type) {
		case appError:
			http.Error(w, err.Error(), err.Code)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type appError struct {
	Code int
	Msg  string
	Err  error
}

func (err appError) Error() string {
	return fmt.Sprintf("app-error: %s: %v", err.Msg, err.Err)
}

func rootHandler(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Welcome to clrbinetsrv.in2p3.fr\n")
	return nil
}
