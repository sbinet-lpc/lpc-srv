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
	http.Handle("/snfusion", snfusionProxy)
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
		log.Printf("error: %v\n", err)
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
	fmt.Fprintf(w, rootIndex)
	return nil
}

var rootIndex = `<!doctype html>
<html>
	<head>
		<meta charset="utf-8"/>
		<title>LPC-binet-srv</title>
		<meta name="viewport" content="width=device-width, minimum-scale=1.0, initial-scale=1.0, user-scalable=yes" />
	</head>

<body>
	<h1>Welcome to clrbinetsrv.in2p3.fr.</hi>
	<ul>
		<li><a href="http://clrbinetsrv.in2p3.fr:5555"><code>/fcs-motor-ctl</code></a></li>
		<li><a href="/snfusion"><code>/snfusion</code></a></li>
		<li><a href="http://clrbinetsrv.in2p3.fr:7073"><code>/fouracc</code></a></li>
		<li><a href="/solid-srv"><code>/solid-srv</code></a></li>
		<li><a href="http://clrbinetsrv.in2p3.fr:8080"><code>/solid-runctl-srv</code></a></li>
		<li><a href="http://clreco.in2p3.fr"><code>/eco-lpc</code></a></li>
		<li><a href="https://clrbinetsrv.in2p3.fr:2000"><code>UCA-print via PaperCut</code></a></li>
		<li><a href="https://clrbinetsrv.in2p3.fr:2001"><code>Melies Server</code></a></li>
	</ul>
</body>
</html>
`
