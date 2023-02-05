package costBasisProcessor

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

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
	colMinSize          = unrealizedPerCol + 1
)

type asset struct {
	Symbol            string
	NumAssets         float64
	CostBasis         float64
	MarketValue       float64
	CostBasisPerShare float64
}

type account struct {
	Name          string
	AccountNumber string
	Assets        map[string]asset
}

func cleanNumberOfCommas(number string) string {
	return strings.Replace(number, ",", "", -1)
}

func printAccountInformation(account account) {
	logger := log.WithField("function", "costBasisProcessor.printAccountInformation()")

	for _, asset := range account.Assets {
		logger.Infof("Account: %s, Symbol: %s, Shares: %f, CostBasis: %f, MarketValue %f, Cost Basis Per Share %f",
			account.Name, asset.Symbol, asset.NumAssets, asset.CostBasis, asset.MarketValue, asset.CostBasisPerShare)
	}
}

func processAccountAsset(account *account, asset asset) {
	logger := log.WithField("function", "costBasisProcessor.processAccountAsset()")

	// Check to see if this account already has this asset information
	// We need to append the data to the existing asset
	a, found := account.Assets[asset.Symbol]
	if !found {
		account.Assets[asset.Symbol] = asset

		logger.Infof("Asset %s not found, adding it", asset.Symbol)
		return
	}

	// Update/Append values
	a.CostBasis += asset.CostBasis
	a.MarketValue += asset.MarketValue
	a.NumAssets += asset.NumAssets
	a.CostBasisPerShare = (a.CostBasis / a.NumAssets)

	// Reassign value
	account.Assets[asset.Symbol] = a
}

func assetFromRow(row []string) (asset, error) {
	logger := log.WithField("function", "costBasisProcessor.assetFromRow()")

	symbol := row[symbolCol]

	numAssets, err := strconv.ParseFloat(cleanNumberOfCommas(row[sharesCol]), 64)
	if err != nil {
		logText := fmt.Sprintf("error while parsing the number of shares into a float: %s", row[sharesCol])
		logger.WithError(err).Error(logText)

		return asset{}, fmt.Errorf("%s: %w", logText, err)
	}

	costBasis, err := strconv.ParseFloat(cleanNumberOfCommas(row[costBasisCol]), 64)
	if err != nil {
		logText := fmt.Sprintf("error while parsing the cost basis into a float: %s", row[costBasisCol])
		logger.WithError(err).Error(logText)

		return asset{}, fmt.Errorf("%s: %w", logText, err)
	}

	marketValue, err := strconv.ParseFloat(cleanNumberOfCommas(row[marketValueCol]), 64)
	if err != nil {
		logText := fmt.Sprintf("error while parsing the market value into a float: %s", row[marketValueCol])
		logger.WithError(err).Error(logText)

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
	logger := log.WithField("function", "costBasisProcessor.processCSV()")

	accounts := make(map[string]account)

	reader := csv.NewReader(csvFile)

	firstRow := true

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		rowLength := len(row)

		if rowLength < colMinSize {
			return fmt.Errorf("the provided row length is below the minimum required %d", rowLength)
		}

		if firstRow {
			firstRow = false
			continue
		}

		logger.Infof("The row values %s", row)

		if err != nil {
			const logText = "error while reading csv file"
			logger.WithError(err).Error(logText)

			return fmt.Errorf("%s: %w", logText, err)
		}

		// Get the account, if it doesn't exist, we create it
		account, found := accounts[row[accountNumberCol]]

		// Initialize the asset map memory
		if !found {
			account.Assets = make(map[string]asset)
			account.Name = row[accountNameCol]
			account.AccountNumber = row[accountNumberCol]
		}

		// Process the row converting it to an asset
		asset, err := assetFromRow(row)
		if err != nil {
			const logText = "error while converting asset from row!"
			logger.WithError(err).Error(logText)

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
	logger := log.WithField("function", "costBasisProcessor.ReadCSV()")

	// Open the CSV file to read
	file, err := os.Open(csvFilePath)
	if err != nil {
		const logText = "failed to open the csv file"
		logger.WithError(err).Error(logText)

		return fmt.Errorf("%s: %w", logText, err)
	}

	defer file.Close()

	err = processCSV(file)
	if err != nil {
		const logText = "unable to process the provided csv"
		logger.WithError(err).Error(logText)

		return fmt.Errorf("%s: %w", logText, err)
	}

	return nil
}
