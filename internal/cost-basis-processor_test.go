package costBasisProcessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCleanNumberOfCommas(t *testing.T) {
	myStr := "1,304"
	expected := "1304"

	retval := cleanNumberOfCommas(myStr)

	assert.Equal(t, expected, retval)
}

func TestAssetFromRow(t *testing.T) {
	// Invalid number of shares
	row := [...]string{"My Betterment Account", "test-account-number", "VWO", "invalid-shares", "1/20/2023", "100.79", "99.28", "1.50", "1.51"}
	_, err := assetFromRow(row[:])
	require.Error(t, err)

	// Invalid cost basis
	row = [...]string{"My Betterment Account", "test-account-number", "VWO", "50", "1/20/2023", "100.79", "invalid-cost-basis", "1.50", "1.51"}
	_, err = assetFromRow(row[:])
	require.Error(t, err)

	// Invalid market value
	row = [...]string{"My Betterment Account", "test-account-number", "VWO", "50", "1/20/2023", "invalid-market-value", "99.28", "1.50", "1.51"}
	_, err = assetFromRow(row[:])
	require.Error(t, err)
}

func TestProcessCSV(t *testing.T) {
	err := ReadCSV("../test/good-test.csv")
	require.NoError(t, err)

	err = ReadCSV("../test/bad-test.csv")
	require.Error(t, err)

	err = ReadCSV("../test/nofile.csv")
	require.Error(t, err)

	err = ReadCSV("../test/invalid-col-len.csv")
	require.Error(t, err)
}
