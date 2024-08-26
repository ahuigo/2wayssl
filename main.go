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
		UsageText:   "2wayssl [-p 443] [-s] -d your-domain",
		Usage: ` $ 2wayssl -p 444 -d 2wayssl.local

echo "127.0.0.1 2wayssl.local" | sudo tee -a /etc/hosts
cd ~/.2wayssl && curl --cacert ca.crt --cert  client.crt --key client.key --tlsv1.2  https://2wayssl.local:444

				+---------------------------+
				|curl -k https://local1.com |
				+------+--------------------+
						  |
						  v 
				  +-------+------+
				  | nginx gateway| 
				  | (port:444)   |  
				  ++-----+-------+  
					 |         | 
					 v         v
		   +-------+---+        +-----------+  
		   | upstream1 |        | upstream2 |  
		   |(port:4500)|        |(port:4501)|  
		   +-----------+        +-----------+  

		`,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("load")),
		Flags:  flags,
		Action: func(cCtx *cli.Context) error {
			domainProxys := cCtx.StringSlice("d")
			initConf(conf, domainProxys)
			quit := make(chan os.Signal, 1)
			proxyServer := startProxyServer(func() {
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
