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
	logEntry := log.WithField("function", funcName)

	flag.Parse()

	logEntry.Infof("Processing Cost Basis CSV File at path: %s", *csvFilePath)

	err := processor.ProcessCSV(*csvFilePath)
	if err != nil {

	}

}

func init() {
	csvFilePath = flag.String("csvfile", "mycsv.csv", "Betterment CostBasis CSV")
}
