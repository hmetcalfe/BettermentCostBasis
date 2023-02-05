package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	processor "github.com/hmetcalfe/betterment-cost-basis/internal"
)

var (
	csvFilePath *string
)

func main() {
	const funcName = "main.main()"
	logger := log.WithField("function", funcName)

	flag.Parse()

	logger.Infof("Processing Cost Basis CSV File at path: %s", *csvFilePath)

	err := processor.ReadCSV(*csvFilePath)
	if err != nil {
		const logText = "Error while reading the Betterment Cost Basis CSV!"
		logger.WithError(err).Error(logText)
	}
}

func init() {
	csvFilePath = flag.String("csvfile", "mycsv.csv", "Betterment CostBasis CSV")
}
