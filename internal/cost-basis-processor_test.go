package costBasisProcessor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessCSV(t *testing.T) {
	err := ProcessCSV("test")
	require.NoError(t, err)
}
