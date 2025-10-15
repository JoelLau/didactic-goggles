package parsers

import (
	"didactic-goggles/internal/customtypes"
	"fmt"
	"strings"
	"time"

	gocsv "github.com/JoelLau/go-csv"
)

type DbsCreditCardParser struct{}

func NewDbsCreditCardParser() *DbsCreditCardParser {
	return &DbsCreditCardParser{}
}

func (*DbsCreditCardParser) Parse(contents [][]string) (DbsCreditCardStatement, error) {
	s := DbsCreditCardStatement{}
	var table [][]string = nil

	for y, row := range contents {
		if len(row) <= 0 || len(strings.TrimSpace(strings.Join(row, ""))) <= 0 {
			continue
		}

		if strings.Contains(strings.ToLower(row[0]), "card transaction details for") {
			s.CardTransactionDetailsFor = row[1]
			continue
		}

		if strings.Contains(strings.ToLower(row[0]), "transactions as at") {
			s.TransactionsAsAt = row[1]
			continue
		}

		if strings.Contains(strings.ToLower(row[0]), "credit limit") || strings.Contains(strings.ToLower(row[0]), "available limit") {
			continue
		}

		// if its none of the others, assume this is the start of the table
		table = contents[y:]
		break
	}

	// dbs comes with headers
	r := make([]string, len(table))
	for i, row := range table {
		r[i] = strings.Join(row, ",")
	}

	v := []byte(strings.Join(r, "\n"))
	m := make([]DbsCreditCardStatementItem, 0)

	if err := gocsv.Unmarshal(v, &m); err != nil {
		return s, err
	}

	s.LineItems = m
	return s, nil
}

type DbsCreditCardStatement struct {
	CardTransactionDetailsFor string `csv-h:"Card Transaction Details For:"` // e.g. "DBS yuu Visa Card  <card_num>"
	TransactionsAsAt          string `csv-h:"Transactions as at:"`           // e.g. "12 Sept 2025"

	LineItems []DbsCreditCardStatementItem
}

// remember to update `Equals` method when updating this struct
type DbsCreditCardStatementItem struct {
	TransactionDate        DbsDate             `csv:"Transaction Date"`         // e.g. "22 Aug 2025"
	TransactionPostingDate DbsDate             `csv:"Transaction Posting Date"` // e.g. "22 Aug 2025"
	TransactionDescription string              `csv:"Transaction Description"`  // e.g. "MALAYSIA BOLEH - JURON SINGAPORE     SG"
	TransactionType        string              `csv:"Transaction Type"`         // e.g. "REFUND & CREDITS", "PURCHASE"
	PaymentType            string              `csv:"Payment Type"`             // e.g. "CONTACTLESS", "Online/In-App Payment"
	TransactionStatus      string              `csv:"Transaction Status"`       // e.g. "Settled"
	DebitAmount            customtypes.Decimal `csv:"Debit Amount"`             // e.g. "Settled"
	CreditAmount           customtypes.Decimal `csv:"Credit Amount"`            // e.g. "Settled"
}

func (a DbsCreditCardStatementItem) Equals(b DbsCreditCardStatementItem) bool {
	if a.TransactionDate.UTC().Compare(b.TransactionDate.Time) != 0 {
		return false
	}

	if a.TransactionPostingDate.UTC().Compare(b.TransactionPostingDate.Time) != 0 {
		return false
	}

	if a.TransactionDescription != b.TransactionDescription {
		return false
	}

	if a.TransactionType != b.TransactionType {
		return false
	}

	if a.PaymentType != b.PaymentType {
		return false
	}

	if a.TransactionStatus != b.TransactionStatus {
		return false
	}

	if a.DebitAmount.Equal(b.DebitAmount.Decimal) {
		return false
	}

	if a.CreditAmount.Equal(b.CreditAmount.Decimal) {
		return false
	}

	return true
}

type DbsDate struct {
	time.Time
}

const DbsDateLayout = "02 Jan 2006"

func (d *DbsDate) UnmarshalCSV(data []byte) (err error) {
	d.Time, err = time.Parse(DbsDateLayout, string(data))
	if err != nil {
		err = fmt.Errorf("failed to parse dbs date: %v", err)
		return
	}

	return
}
