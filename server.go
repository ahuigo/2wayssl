package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
)


func startHttpsServer(cleanup func()) *http.Server {
	config := GetConfig()
	// handler := getHander(config.DomainProxys)
	pool := x509.NewCertPool()
	domain := config.DomainProxys[0].Domain
	confpath := os.Getenv("HOME") + "/.2wayssl"
	crt, err := os.ReadFile(confpath + "/ca.crt")
	if err != nil {
		log.Fatalf("failed to read cert path:%s, err:%s\n", confpath+"ca.crt", err.Error())
		os.Exit(0)
	}
	pool.AppendCertsFromPEM(crt)
	http.HandleFunc("/", handler)
	s := &http.Server{
		Addr: ":" + config.Port,
		// Handler: handler,
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert, // verify client cert

		},
	}
	go func() {
		serverCertPath := fmt.Sprintf(confpath+"/%s.server.crt", domain)
		serverKeyPath := fmt.Sprintf(confpath+"/%s.server.key", domain)
		log.Fatal(s.ListenAndServeTLS(serverCertPath, serverKeyPath))
		cleanup()
	}()
	return s
}
