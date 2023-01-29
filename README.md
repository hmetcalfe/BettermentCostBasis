# BettermentCostBasis
A small golang application to calculate your Betterment portfolio's cost basis per share from the Betterment's exportable CSV.

## Pre-requesites

Linting of this application requires `golangci-lint` 

## Installation

To build the application, you can simply run `make build`.

## Installation Output

Upon compilation, you'll get several executables in the `bin` directory
1. costBasisCalculator.mac
2. costBasisCalculator.exe
3. costBasisCalculator 

## Application Execution

To run `costBasisCalculator -csvfile /path/to/test.csv`
