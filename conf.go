package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func GetConfig() *Config {
	if conf == nil {
		conf = &Config{}
		home := os.Getenv("HOME")
		if !IsFileExists(home + "/.2wayssl") {
			err := os.Mkdir(home+"/.2wayssl", 0755)
			if err != nil {
				log.Fatalf("Failed to create directory: %v", err)
			}
		}
		err:=os.Chdir(home + "/.2wayssl")
		if err != nil {
			println(err.Error())
			os.Exit(0)
		}
	}
	return conf
}

type DomainProxy struct {
	Domain string
	// ProxyPass string
}

type Config struct {
	Port         string
	Silent       bool
	DomainProxys []DomainProxy
	// Domain      string
	// ProxyPass   string
	// CertPath    string
	// CertKeyPath string
}

var conf *Config

func initConf(conf *Config, domainProxys []string) {
	if len(domainProxys) == 0 {
		fmt.Println("Usage: 2wayssl -d local.com=http://localhost:5000")
		os.Exit(0)
	}
	domainProxy := domainProxys[0]
	domain := strings.TrimSpace(domainProxy)
	domain = strings.ToLower(domain)
	if !regexp.MustCompile(`^\w[\w\-\.]*$`).MatchString(domain) {
		fmt.Println("Invalid domain: " + domain)
		os.Exit(1)
	}
	conf.DomainProxys = append(conf.DomainProxys, DomainProxy{
		Domain: domain,
		// ProxyPass: proxyPass,
	})
	createCert(conf)
}
