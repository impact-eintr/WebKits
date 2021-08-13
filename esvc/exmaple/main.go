package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/impact-eintr/WebKits/esvc"
)

type program struct {
	LogFile *os.File
	svr     *server
	ctx     context.Context
}

func (p *program) Context() context.Context {
	return p.ctx
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	prg := program{
		svr: &server{},
		ctx: ctx,
	}

	defer func() {
		if prg.LogFile != nil {
			if closeErr := prg.LogFile.Close(); closeErr != nil {
				log.Printf("error closing '%s': %v\n", prg.LogFile.Name(), closeErr)
			}
		}
	}()

	// call svc.Run to start your program/service
	// svc.Run will call Init, Start, and Stop
	if err := esvc.Run(&prg); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init() error {
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
