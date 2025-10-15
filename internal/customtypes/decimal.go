package customtypes

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Decimal struct {
	decimal.Decimal
}

func (d *Decimal) UnmarshalCSV(b []byte) error {
	var dec decimal.Decimal
	var err error
	s := string(b)

	if strings.TrimSpace(s) == "" {
		dec = decimal.Zero
	} else {
		dec, err = decimal.NewFromString(s)
		if err != nil {
			return fmt.Errorf("failed to unmarshal CSV into Decimal: %w", err)
		}
	}

	d.Decimal = dec
	return nil
}
