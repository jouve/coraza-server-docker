package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"sync"

	_ "github.com/jptosso/coraza-libinjection"
	_ "github.com/jptosso/coraza-pcre"
	"github.com/jptosso/coraza-server/config"
	"github.com/jptosso/coraza-server/protocols"
	"github.com/jptosso/coraza-waf/v2"
	"github.com/jptosso/coraza-waf/v2/seclang"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	f := flag.String("f", "/etc/coraza-server/config.yml", "Absolute path to configuration file (.yml)")
	flag.Parse()
	cfg, err := config.ReadFile(*f)
	if err != nil {
		logrus.Fatal(err)
	}
	wg := &sync.WaitGroup{}
	for _, a := range cfg.Agents {
		proto, err := protocols.GetProtocol(a.Protocol)
		if err != nil {
			logrus.Fatal(err)
		}
		wg.Add(1)
		logrus.Info("Initializing waf")
		waf := coraza.NewWaf()
		parser, _ := seclang.NewParser(waf)
		if len(a.Include) == 0 {
			logrus.Warn("No rules detected for agent")
		}
		for _, file := range a.Include {
			if err := parser.FromFile(file); err != nil {
				logrus.Fatal(err)
			}
		}
		logrus.Infof("Initializing protocol %s", a.Protocol)
		proto.Init(waf, a)
		logrus.Infof("Starting protocol %s on %s", a.Protocol, a.Bind)
		go func() {
			defer wg.Done()
			if err := proto.Start(); err != nil {
				logrus.Fatal(err)
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logrus.Info("%s", http.ListenAndServe(":6060", nil))
	}()
	wg.Wait()
	logrus.Info("Coraza server finished.")
}
