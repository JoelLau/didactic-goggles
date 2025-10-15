package customtypes_test

import (
	"didactic-goggles/internal/customtypes"
	"testing"

	gocsv "github.com/JoelLau/go-csv"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestDecimalUnmarshalCSV(t *testing.T) {
	t.Parallel()

	type Model struct {
		D customtypes.Decimal
	}

	given := []byte(`decimal
123.123
`)

	d := decimal.NewFromFloat(123.123)
	got := []Model{}
	want := []Model{{D: customtypes.Decimal{Decimal: d}}}
	want2 := []Model{{D: customtypes.Decimal{Decimal: d}}}

	err := gocsv.Unmarshal(given, &got)
	require.NoError(t, err)

	require.Len(t, got, len(want))
	require.True(t, want2[0].D.Decimal.Equal(want[0].D.Decimal))
}
