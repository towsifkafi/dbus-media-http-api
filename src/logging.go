package main

import (
	"fmt"
	"log"
	"time"
)

type logWriter struct{}

func (l logWriter) Write(bytes []byte) (int, error) {
	return fmt.Printf("[%.19s] %s", time.Now().Format("2006/01/02 15:04:05"), string(bytes))
}

func setupLogging() {
	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(logWriter{})
}
