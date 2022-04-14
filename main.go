// HTTPs server with Letsencrypt certificate
//
// https://pkg.go.dev/golang.org/x/crypto/acme/autocert#example-Manager
//
// package main
//
// import (
// 	"net/http"
//
// 	"golang.org/x/crypto/acme/autocert"
// )
//
// func main() {
// 	m := &autocert.Manager{
// 		Cache:      autocert.DirCache("secret-dir"),
// 		Prompt:     autocert.AcceptTOS,
// 		Email:      "kirill@scherba.ru",
// 		HostPolicy: autocert.HostWhitelist("ex.myteo.net", "ex.myteo.net"),
// 	}
// 	s := &http.Server{
// 		Addr:      ":https",
// 		TLSConfig: m.TLSConfig(),
// 	}
// 	s.ListenAndServeTLS("", "")
// }
//
// https://pkg.go.dev/golang.org/x/crypto/acme/autocert#example-NewListener
//
//
// This is sample HTTPS server with Letsencrypt automatic SSL sertificat and
// redirection from HTTP to HTTPS.
//
// By default server uses 80 and 443 ports.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

var domain string

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+domain+":443"+r.RequestURI, http.StatusMovedPermanently)
}

func main() {

	flag.StringVar(&domain, "domain", "", "domain name to process HTTP/s server")
	flag.Parse()
	if len(domain) == 0 {
		fmt.Println("The domain parameter is required")
		flag.Usage()
		os.Exit(0)
	}

	// Redirect HTTP to HTTPS
	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// HTTPS server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, TLS user! Your config: %+v", r.TLS)
	})
	log.Fatal(http.Serve(autocert.NewListener(domain), mux))
}
