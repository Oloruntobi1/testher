package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
// "golang.org/x/net/http2"

)
 
func main() {
	client := &http.Client{Transport: transport1()}
 
	res, err := client.Get("https://localhost:8443")
	if err != nil {
		log.Fatal(err)
	}
 
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
 
	res.Body.Close()
 
	fmt.Printf("Code: %d\n", res.StatusCode)
	fmt.Printf("Body: %s\n", body)
}
 
// func transport2() *http2.Transport {
// 	return &http2.Transport{
// 		TLSClientConfig:     tlsConfig(),
// 		DisableCompression:  true,
// 		AllowHTTP:           false,
// 	}
// }
 
func transport1() *http.Transport {
	return &http.Transport{
		// Original configurations from `http.DefaultTransport` variable.
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true, // Set it to false to enforce HTTP/1
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		// Our custom configurations.
		ResponseHeaderTimeout: 10 * time.Second,
		DisableCompression:    true,
		// Set DisableKeepAlives to true when using HTTP/1 otherwise it will cause error: dial tcp [::1]:8090: socket: too many open files
		DisableKeepAlives:     false,
		TLSClientConfig:       tlsConfig(),
	}
}
 
func tlsConfig() *tls.Config {
	crt, err := ioutil.ReadFile("../server.crt")
	if err != nil {
		log.Fatal(err)
	}
 
	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(crt)
 
	return &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: false,
		ServerName:         "localhost",
	}
}