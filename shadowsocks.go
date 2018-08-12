package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/shadowsocks/go-shadowsocks2/core"
)

var settings struct {
	Server     string `json:"server_addr"`
	Local      string `json:"local_addr"`
	ServerPort int    `json:"server_port"`
	SocksPort  int    `json:"socks_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
}

type server struct {
	data chan int

	exit chan struct{}
	wg   sync.WaitGroup
}

func (s *server) start() {
	s.exit = make(chan struct{})

	s.wg.Add(1)
	go s.startServer()
}

func (s *server) stop() error {
	close(s.exit)
	s.wg.Wait()
	return nil
}

func (s *server) startServer() {
	defer s.wg.Done()
	logger := log.New(os.Stderr, "cow: ", log.LstdFlags|log.Lshortfile)

	configFile, err := os.Open("config.json")

	defer configFile.Close()

	if err != nil {
		logger.Fatal(err)
	}

	// parse config file
	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(&settings); err != nil {
		logger.Fatal("Parsing config file failed", err.Error())
	}

	// setup shadowsocks client
	ciph, err := core.PickCipher(settings.Method, nil, settings.Password)

	if err != nil {
		logger.Fatal("No such cipher", err.Error())
	}

	// setup local route
	server := fmt.Sprintf("%s:%d", settings.Server, settings.ServerPort)
	localSOCKS := fmt.Sprintf("%s:%d", settings.Local, settings.SocksPort)

	go socksLocal(localSOCKS, server, ciph.StreamConn)

	<-s.exit
}
