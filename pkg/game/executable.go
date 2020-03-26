package game

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func doOrDie(err error) {
	if err != nil {
		log.Fatalf("unable to continue: %+v", err)
	}
}

func Run(configPath string) {
	config, err := GetConfig(configPath)
	doOrDie(err)

	logLevel, err := config.GetLogLevel()
	doOrDie(err)
	log.SetLevel(logLevel)

	prometheus.Unregister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	prometheus.Unregister(prometheus.NewGoCollector())

	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf(":%d", config.Port)
	log.Infof("serving on %s", addr)
	go func() {
		http.ListenAndServe(addr, nil)
	}()

	stop := make(chan struct{})
	game := NewGame()
	gcw := NewGameConcurrencyWrapper(game, stop)

	log.Infof("instantiated game with concurrency wrapper: \n%s\n", gcw.GetJsonModel())

	<-stop
}
