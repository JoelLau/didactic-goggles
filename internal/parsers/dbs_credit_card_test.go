package parsers_test

import (
	"didactic-goggles/internal/customtypes"
	"didactic-goggles/internal/parsers"
	"testing"

	gocsv "github.com/JoelLau/go-csv"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEmpty(t *testing.T) {
	// arrange
	p := parsers.DbsCreditCardParser{}

	// act
	have, err := p.Parse([][]string{})

	// assert
	require.NoError(t, err)
	require.NotNil(t, have)
}

func TestParseString(t *testing.T) {
	// arrange
	given := [][]string{
		{"Card Transaction Details For:", "DBS yuu Visa Card 1234-1234-1234-1234", "", "", "", "", "", ""},
		{"Transactions as at:", "12 Sep 2025", "", "", "", "", ""},
		{"", "", "", "", "", "", ""},
		{"Credit Limit:", "SGD undefined", "", "", "", "", ""},
		{"Available Limit:", "SGD undefined", "", "", "", "", ""},
		{"", "", "", "", "", "", ""},
		{"Transaction Date", "Transaction Posting Date", "Transaction Description", "Transaction Type", "Payment Type", "Transaction Status", "Debit Amount", "Credit Amount"},
		{"22 Aug 2025", "23 Aug 2025", "BOON LAY BI            SINGAPORE     SG", "PURCHASE", "Contactless", "Settled", "2.3", ""},
		{"22 Aug 2025", "23 Aug 2025", "FAVEPAY - HOMIES BAKER SINGAPORE     SG", "PURCHASE", "Online/In-App Payment", "Settled", "1.75", ""},
		{"21 Aug 2025", "23 Aug 2025", "MALAYSIA BOLEH - JURON SINGAPORE     SG", "PURCHASE", "Contactless", "Settled", "2.5", ""},
		{"21 Aug 2025", "23 Aug 2025", "MALAYSIA BOLEH - JURON SINGAPORE     SG", "PURCHASE", "Contactless", "Settled", "3.5", ""},
	}

	parser := parsers.NewDbsCreditCardParser()

	// act
	got, err := parser.Parse(given)
	require.NoError(t, err)

	want := []parsers.DbsCreditCardStatementItem{
		{TransactionDate: newDbsDate(t, "22 Aug 2025"), TransactionPostingDate: newDbsDate(t, "23 Aug 2025"), TransactionDescription: "BOON LAY BI            SINGAPORE     SG", TransactionType: "PURCHASE", PaymentType: "Contactless", TransactionStatus: "Settled", CreditAmount: newDecimal(t, "2.3")},
		{TransactionDate: newDbsDate(t, "22 Aug 2025"), TransactionPostingDate: newDbsDate(t, "23 Aug 2025"), TransactionDescription: "FAVEPAY - HOMIES BAKER SINGAPORE     SG", TransactionType: "PURCHASE", PaymentType: "Online/In-App Payment", TransactionStatus: "Settled", CreditAmount: newDecimal(t, "1.75")},
		{TransactionDate: newDbsDate(t, "21 Aug 2025"), TransactionPostingDate: newDbsDate(t, "23 Aug 2025"), TransactionDescription: "MALAYSIA BOLEH - JURON SINGAPORE     SG", TransactionType: "PURCHASE", PaymentType: "Contactless", TransactionStatus: "Settled", CreditAmount: newDecimal(t, "2.5")},
		{TransactionDate: newDbsDate(t, "21 Aug 2025"), TransactionPostingDate: newDbsDate(t, "23 Aug 2025"), TransactionDescription: "MALAYSIA BOLEH - JURON SINGAPORE     SG", TransactionType: "PURCHASE", PaymentType: "Contactless", TransactionStatus: "Settled", CreditAmount: newDecimal(t, "3.5")},
	}

	// assert
	require.Equalf(t, "DBS yuu Visa Card 1234-1234-1234-1234", got.CardTransactionDetailsFor, "'Transaction Details For' does not match")
	require.Equalf(t, "12 Sep 2025", got.TransactionsAsAt, "'Transactions At' does not match")

	require.Len(t, got.LineItems, len(want))
	for i := 0; i < len(got.LineItems); i++ {
		require.True(t, got.LineItems[i].Equals(want[i]))
	}
}

func newDbsDate(t *testing.T, s string) parsers.DbsDate {
	t.Helper()

	d := parsers.DbsDate{}
	err := d.UnmarshalCSV([]byte(s))
	assert.NoError(t, err)

	return d
}

func newDecimal(t *testing.T, s string) customtypes.Decimal {
	t.Helper()

	d, err := decimal.NewFromString(s)
	require.NoError(t, err)

	return customtypes.Decimal{Decimal: d}
}

func TestDbsDate(t *testing.T) {
	given := []byte(`Transaction Date
22 Aug 2025`)

	type Model struct {
		TransactionDate parsers.DbsDate `csv:"Transaction Date"`
	}

	got := []Model{}
	want := []Model{{TransactionDate: newDbsDate(t, "22 Aug 2025")}}

	err := gocsv.Unmarshal(given, &got)
	require.NoError(t, err)

	require.Len(t, got, len(want))
	require.True(t, got[0].TransactionDate.Time.Equal(want[0].TransactionDate.Time))
}
