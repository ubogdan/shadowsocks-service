package main

import (
	"log"

	"github.com/judwhite/go-svc/svc"
)

// implements svc.Service
type program struct {
	svr *server
}

func main() {
	prg := program{
		svr: &server{},
	}

	// call svc.Run to start your program/service
	// svc.Run will call Init, Start, and Stop
	if err := svc.Run(&prg); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	log.Printf("is win service? %v\n", env.IsWindowsService())

	return nil
}

func (p *program) Start() error {
	log.Printf("Starting...\n")
	go p.svr.start()
	return nil
}

func (p *program) Stop() error {
	log.Printf("Stopping...\n")
	if err := p.svr.stop(); err != nil {
		return err
	}
	log.Printf("Stopped.\n")
	return nil
}
