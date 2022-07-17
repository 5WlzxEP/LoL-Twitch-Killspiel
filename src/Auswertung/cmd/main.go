package main

import (
	"encoding/json"
	"flag"
	"github.com/5WlzxEP/LoL-Twitch-Killspiel/src/Auswertung"
	"io"
	"log"
	"os"
)

func main() {
	confFilePtr := flag.String("conf", "", "Set custom config file")
	flag.Parse()
	var confFile string
	if confFile = "config.json"; *confFilePtr != "" {
		confFile = *confFilePtr
	}
	conf := readConf(confFile)

	if conf.LogFile != "" {

		f, err := os.OpenFile(conf.LogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Fatalf("Cloud not open log file: %v\n", err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Println("Error closing log file")
			}
		}(f)
		log.SetOutput(f)
	}
	Auswertung.SetGlobal(conf)

	Auswertung.Auswerten()

}

func readConf(file string) *Auswertung.Global {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Could not open config")
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Could not read config: %v\n", err)
	}
	var global Auswertung.Global
	err = json.Unmarshal(bytes, &global)
	if err != nil {
		log.Fatalf("Error while paring config: %v\n", err)

	}
	return &global
}
