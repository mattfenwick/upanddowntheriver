package game

import log "github.com/sirupsen/logrus"

func doOrDie(err error) {
	if err != nil {
		log.Fatalf("unable to continue: %+v", err)
	}
}

func Run(configPath string) {
	config, err := GetConfig(configPath)
	doOrDie(err)
}
