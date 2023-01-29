package costBasisProcessor

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	accountNameCol      = 0
	accountNumberCol    = 1
	symbolCol           = 2
	sharesCol           = 3
	purchaseDateCol     = 4
	marketValueCol      = 5
	costBasisCol        = 6
	unrealizedDollarCol = 7
	unrealizedPerCol    = 8
)

type asset struct {
	Symbol      string
	NumAssets   float64
	CostBasis   float64
	MarketValue float64
}

type account struct {
	Name          string
	AccountNumber string
	Assets        map[string]asset
}

func printAccountInformation(account account) {
	const funcName = "costBasisProcessor.printAccountInformation()"
	logEntry := log.WithField("function", funcName)

	for _, asset := range account.Assets {
		logEntry.Infof("Account: %s, Symbol: %s, Shares: %f, CostBasis: %f, MarketValue %f", account.Name, asset.Symbol, asset.NumAssets, asset.CostBasis, asset.MarketValue)
	}
}

func processAccountAsset(account *account, asset asset) {
	const funcName = "costBasisProcessor.processAccountAsset()"
	logEntry := log.WithField("function", funcName)

	// Check to see if this account already has this asset information
	// We need to append the data to the existing asset
	a, found := account.Assets[asset.Symbol]
	if !found {
		account.Assets[asset.Symbol] = asset

		logEntry.Infof("Asset %s not found, adding it", asset.Symbol)
		return
	}

	// Update/Append values
	a.CostBasis += asset.CostBasis
	a.MarketValue += asset.MarketValue
	a.NumAssets += asset.NumAssets

	// Reassign value
	account.Assets[asset.Symbol] = a
}

func assetFromRow(row []string) (asset, error) {
	const funcName = "costBasisProcessor.assetFromRow()"
	logEntry := log.WithField("function", funcName)

	symbol := row[symbolCol]

	numAssets, err := strconv.ParseFloat(row[sharesCol], 64)
	if err != nil {
		logText := fmt.Sprintf("Error while parsing the number of shares into a float: %s", row[sharesCol])
		logEntry.WithError(err).Error(logText)

		return asset{}, fmt.Errorf("%s: %w", logText, err)
	}

	costBasis, err := strconv.ParseFloat(row[costBasisCol], 64)
	if err != nil {
		logText := fmt.Sprintf("Error while parsing the cost basis into a float: %s", row[costBasisCol])
		logEntry.WithError(err).Error(logText)

		return asset{}, fmt.Errorf("%s: %w", logText, err)
	}

	marketValue, err := strconv.ParseFloat(row[costBasisCol], 64)
	if err != nil {
		logText := fmt.Sprintf("Error while parsing the market value into a float: %s", row[marketValueCol])
		logEntry.WithError(err).Error(logText)

		return asset{}, fmt.Errorf("%s: %w", logText, err)
	}

	return asset{
		Symbol:      symbol,
		NumAssets:   numAssets,
		CostBasis:   costBasis,
		MarketValue: marketValue,
	}, nil
}

func processCSV(csvFile *os.File) error {
	const funcName = "costBasisProcessor.processCSV()"
	logEntry := log.WithField("function", funcName)

	accounts := make(map[string]account)

	reader := csv.NewReader(csvFile)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		logEntry.Infof("The row values %s", row)

		if err != nil {
			const logText = "Error while reading csv file"
			logEntry.WithError(err).Error(logText)

			return fmt.Errorf("%s: %w", logText, err)
		}

		// Get the account, if it doesn't exist, we create it
		account := accounts[row[accountNumberCol]]

		// Process the row converting it to an asset
		asset, err := assetFromRow(row)
		if err != nil {
			const logText = "Error while converting asset from row!"
			logEntry.WithError(err).Error(logText)

			return fmt.Errorf("%s: %w", logText, err)
		}

		// Add the asset information to the account
		processAccountAsset(&account, asset)

		accounts[row[accountNumberCol]] = account
	}

	for _, account := range accounts {
		printAccountInformation(account)
	}

	return nil
}

func ReadCSV(csvFilePath string) error {
	const funcName = "costBasisProcessor.ReadCSV()"
	logEntry := log.WithField("function", funcName)

	// Open the CSV file to read
	file, err := os.Open(csvFilePath)
	if err != nil {
		const logText = "Failed to open the CSV file!!"
		logEntry.WithError(err).Error(logText)

		return fmt.Errorf("%s: %w", logText, err)
	}

	defer file.Close()

	err = processCSV(file)
	if err != nil {
		const logText = "Unable to process the provided CSV!"
		logEntry.WithError(err).Error(logText)

		return fmt.Errorf("%s: %w", logText, err)
	}

	return nil
}
