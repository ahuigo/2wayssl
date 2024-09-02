package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	conf := GetConfig()
	flags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "d",
			Usage:    "domain",
			Required: true,
			// Destination: &conf.Domain,
		},
		&cli.StringFlag{
			Name:        "p",
			Value:       "443",
			Usage:       "server `PORT`",
			Destination: &conf.Port,
		},
		&cli.StringFlag{
			Name:        "subj",
			Value:       "/C=CN/ST=GD/L=SZ/O=TwoWaySsl, Org.",
			Usage:       "subj prefix",
			Destination: &conf.Port,
		},
		&cli.BoolFlag{
			Name: "silent",
			Aliases:     []string{"s"},
			Value:       false,
			Usage:       "Silent mode, no prompt",
			Destination: &conf.Silent,
		},
	}
	app := &cli.App{
        Name:        "2wayssl",
        Description: fmt.Sprintf("simple 2wayssl generator(version:%s,go:%s)",BuildVersion, GoVersion),
		UsageText:   "2wayssl [-p 443] [-s] [-h] -d your-domain",
		Usage: `Automatically generate 2way-ssl certificates, configuration, and test server.

# 1. generate certificate and start a test server
$ 2wayssl -p 444 -d 2wayssl.local

# 2. test certificate via test server
$ echo "127.0.0.1 2wayssl.local" | sudo tee -a /etc/hosts
$ cd ~/.2wayssl && curl --cacert ca.crt --cert  client.crt --key client.key --tlsv1.2  https://2wayssl.local:444

# 3. view certificate and nginx.conf
$ ls ~/.2wayssl 
2wayssl.local.server.crt 2wayssl.local.server.key ca.key                   client.crt               client.key
2wayssl.local.server.csr ca.crt                   ca.srl                   client.csr               nginx.conf
		`,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("load")),
		Flags:  flags,
		Action: func(cCtx *cli.Context) error {
			domainProxys := cCtx.StringSlice("d")
			initConf(conf, domainProxys)
			quit := make(chan os.Signal, 1)
			proxyServer := startHttpsServer(func() {
				log.Println("cleanup")
			})
			domains := lo.Map(conf.DomainProxys, func(dp DomainProxy, i int) string {
				return dp.Domain
			})

			cmd:=fmt.Sprintf("echo '127.0.0.1 %s' | sudo tee -a /etc/hosts", strings.Join(domains, " "))
			runCmdWithConfirm("Config /etc/hosts", cmd, conf.Silent)
			fmt.Printf("Press Ctrl+C to shutdown\n")
			signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
			sig := <-quit
			err := proxyServer.Close()
			log.Printf("gracefully shutdown Server ...(sig=%v, err=%v)\n", sig, err)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
