package recurrency

import (
	"cashflow/models"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitTransactionsWithRecurrency(t *testing.T) {
	expectedWith := []models.Transaction{
		{
			Date:        "2000/01/01",
			Amount:      2000,
			Id:          "1",
			Description: "",
			Recurrency:  "FREQ=MONTHLY;BYMONTHDAY=15",
		},
	}

	expectedWithout := []models.Transaction{
		{
			Date:        "2000/01/01",
			Amount:      2000,
			Id:          "2",
			Description: "",
			Recurrency:  "",
		},
	}

	resultWith, resultWithout := SplitTransactionsWithRecurrency(append(expectedWith, expectedWithout...))

	require.Equal(t, expectedWith, resultWith)
	require.Equal(t, expectedWithout, resultWithout)
}
