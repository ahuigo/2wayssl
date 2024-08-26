package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ahuigo/gofnext"
)

var loadCertificate = loadCertificateRaw

func init() {
	loadCertificate = gofnext.CacheFn2Err(loadCertificateRaw)
}

func loadCertificateRaw(certFile, keyFile string) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Printf("load certificate failed, err: %v\n", err)
		return nil, err
	}
	return &cert, nil
}

// dual https
func startProxyServer(cleanup func()) *http.Server {
	config := GetConfig()
	// handler := getHander(config.DomainProxys)
	// ssl 双向检验
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

// self https
func CreateProxyServer(cleanup func()) *http.Server {
	config := GetConfig()
	handler := getHander(config.DomainProxys)
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: handler,
		TLSConfig: &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				certPath, certKeyPath := getCertPath(info.ServerName)
				return loadCertificate(certPath, certKeyPath)
			},
		},
	}
	go func() {
		// ts = httptest.NewUnstartedServer(http.HandlerFunc(fn))
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Println(err)
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
		cleanup()
	}()
	return server
}
