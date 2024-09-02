package main

import (
	"fmt"
	"os"
	"strings"
)

var subj_prefix = `/C=CN/ST=GD/L=SZ/O=TwoWaySsl, Org.`

func getCertPath(domain string) (certPath, certKeyPath string) {
	home := os.Getenv("HOME")
	certPath = fmt.Sprintf(home+"/.2wayssl/%s.crt", domain)
	certKeyPath = fmt.Sprintf(home+"/.2wayssl/%s.key", domain)
	return certPath, certKeyPath
}

func createCA() {
	// 0. check ca exists
	confpath := os.Getenv("HOME") + "/.2wayssl"
	if IsFileExists(confpath + "/ca.crt") {
		return
	}
	// 1. generate ca key
	cmd := fmt.Sprintf(`openssl genrsa -out %s/ca.key 1024`, confpath)
	out, err := RunCommand("sh", "-c", cmd)
	if err != nil {
		fmt.Printf("failed to execute cmd(\033[31m %s \033[0m), err: %v, stdout: %s\n\n", cmd, err, out)
		os.Exit(0)
	}
	// 2. generate ca cert
	cmd = fmt.Sprintf(`openssl req -new -x509 -days 3650 -key ca.key -out ca.crt -subj "%s"`, subj_prefix)
	out, err = RunCommand("sh", "-c", cmd)
	if err != nil {
		fmt.Printf("failed to execute cmd(\033[31m %s \033[0m), err: %v, stdout: %s\n\n", cmd, err, out)
		os.Exit(0)
	}
}

func createServerCert(domain string) {
	// 1. check cert exists
	if IsFileExists(domain + ".server.crt") {
		return
	}
	// 2. generate domain's server key
	cmd := fmt.Sprintf(`openssl genrsa -out %s.server.key 1024`, domain)
	out, err := RunCommand("sh", "-c", cmd)
	if err != nil {
		fmt.Printf("failed to execute cmd(\033[31m %s \033[0m), err: %v, stdout: %s\n\n", cmd, err, out)
		os.Exit(0)
	}

	// 3. generate domain's server csr
	cmd = fmt.Sprintf(`openssl req -new -key %s.server.key -out %s.server.csr -subj "%s/CN=%s" -addext "subjectAltName = DNS:%s"`, domain, domain,subj_prefix, domain, domain)
	println(cmd)
	out, err = RunCommand("sh", "-c", cmd)
	if err != nil {
		fmt.Printf("failed to execute cmd(\033[31m %s \033[0m), err: %v, stdout: %s\n\n", cmd, err, out)
		os.Exit(0)
	}
	// 4. generate domain's server cert
	// cmd = fmt.Sprintf(`openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -in %s.server.csr -out %s.server.crt -days 3650`, domain, domain)
	cmd = fmt.Sprintf(`openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -in %s.server.csr -out %s.server.crt -days 3650 -extensions SAN -extfile <(printf "\n[SAN]\nsubjectAltName=DNS:%s")`, domain, domain, domain)
	println(cmd)
	out, err = RunCommand("bash", "-c", cmd)
	if err != nil {
		fmt.Printf("failed to execute cmd(\033[31m %s \033[0m), err: %v, stdout: %s\n\n", cmd, err, out)
		os.Exit(0)
	}
}

func createClientCert() {
	// 1. check cert exists
	if IsFileExists("client.crt") {
		return
	}
	// 2. generate domain's client key
	cmd := `openssl genrsa -out client.key 1024`
	out, err := RunCommand("sh", "-c", cmd)
	if err != nil {
		fmt.Printf("failed to execute cmd(\033[31m %s \033[0m), err: %v, stdout: %s\n\n", cmd, err, out)
		os.Exit(0)
	}

	// 3. generate domain's client csr
	cmd = fmt.Sprintf(`openssl req -new -key client.key -out client.csr -subj "%s/CN=client"`, subj_prefix)
	out, err = RunCommand("sh", "-c", cmd)
	if err != nil {
		fmt.Printf("failed to execute cmd(\033[31m %s \033[0m), err: %v, stdout: %s\n\n", cmd, err, out)
		os.Exit(0)
	}
	// 4. generate domain's client cert
	cmd = fmt.Sprintf(`openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -in client.csr -out client.crt -days %s`, "365")
	out, err = RunCommand("sh", "-c", cmd)
	if err != nil {
		fmt.Printf("failed to execute cmd(\033[31m %s \033[0m), err: %v, stdout: %s\n\n", cmd, err, out)
		os.Exit(0)
	}
}

func createCert(conf *Config) {
	// 1. init config 
	domain :=conf.DomainProxys[0].Domain
	if len(conf.SubjPrefix)>0 && conf.SubjPrefix[0] == '/' {
		subj_prefix = conf.SubjPrefix
	}

	// 2. ca + server + client
	createCA()
	createServerCert(domain)
	createClientCert()

	// 4. create nginx conf
	tpl := getNginxConf(domain)
	nginxPath := os.Getenv("HOME") + "/.2wayssl/nginx.conf"
	err := os.WriteFile(nginxPath, []byte(tpl), 0644)
	if err != nil {
		fmt.Printf("failed to write nginx conf file(%s), err: %v\n",nginxPath, err)
		os.Exit(0)
	}
	fmt.Printf("Nginx config: \033[94m ~/.2wayssl/nginx.conf \033[0m \n")

	// 5. curl command
	fmt.Printf("Have a try:\n\033[94m curl --cacert ~/.2wayssl/ca.crt --cert  ~/.2wayssl/client.crt --key ~/.2wayssl/client.key --tlsv1.2 -v https://%s:%s \033[0m \n", domain, conf.Port)
}

func getNginxConf(domain string) string {
	confPath := os.Getenv("HOME") + "/.2wayssl/"
	tpl := NGX_TEMPLATE
	tpl = strings.Replace(tpl, "{DOMAIN}", domain, 1)
	tpl = strings.Replace(tpl, "{SERVER_CRT_PATH}", confPath+domain+".server.crt", 1)
	tpl = strings.Replace(tpl, "{SERVER_KEY_PATH}", confPath+domain+".server.key", 1)
	tpl = strings.Replace(tpl, "{CA_CRT_PATH}", confPath+"ca.crt", 1)
	return tpl
}
