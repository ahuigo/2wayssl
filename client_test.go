package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

// cd ~/.2wayssl && curl --cacert ca.crt --cert  client.crt --key client.key --tlsv1.2  https://2wayssl.local:444
func TestClient(t *testing.T) {
		pool := x509.NewCertPool()
		confpath := os.Getenv("HOME")+ "/.2wayssl"
		caPath := confpath + "/ca.crt"
		clientCertPath := confpath + "/client.crt"
		clientKeyPath := confpath + "/client.key"
		caCrt, err := os.ReadFile(caPath)
		if err != nil {
			log.Fatal("read ca.crt file error:", err.Error())
		}
		pool.AppendCertsFromPEM(caCrt)
		clientCrt, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
		if err != nil {
			log.Fatalln("LoadX509KeyPair error:", err.Error())
		}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{clientCrt},
				MinVersion:   tls.VersionTLS12,
				MaxVersion:   tls.VersionTLS12,
			},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get("https://2wayssl.local:444/")
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(string(body))
}