package recurrency

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChangeRecurrencyEnd(t *testing.T) {
	var recurrency string = "FREQ=MONTHLY;UNTIL=20000301"

	results, err := ChangeRecurrencyEnd(recurrency, "2000-02-01")
	require.NoError(t, err)

	expected := "FREQ=MONTHLY;UNTIL=20000131T000000Z"

	require.Equal(t, expected, results)
}
