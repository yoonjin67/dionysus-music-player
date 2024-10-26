package main

import (
	"log"
	"strings"
)

func errorHandler(e error, logMsg string, logLevel uint8, opt string) {
	if e == nil {
		return
	}
	if strings.Compare(opt, "warn") == 0 && logLevel > 1 {
		log.Println(
			"Waring: error while "+logMsg,
			"\nError Code: ", e)
	}
	if strings.Compare(opt, "info") == 0 && logLevel > 2 {
		log.Println("INFO: " + logMsg)
	}
	if strings.Compare(opt, "error") == 0 {
		log.Println("Panic while " + logMsg)
		panic(e)
	}
}
