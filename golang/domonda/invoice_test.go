package domonda

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/domonda/go-types/money"
)

func ptr[T any](v T) *T { return &v }

func TestInvoice_Validate_costAmounts(t *testing.T) {
	for _, scenario := range []struct {
		name           string
		net            *money.Amount
		conversionRate *money.Rate
		costCenters    map[string]money.Amount
		costUnits      map[string]money.Amount
		wantErr        string
	}{
		{
			name: "no cost centers or cost units",
			net:  ptr(money.Amount(1000)),
		},
		{
			name:      "cost unit sum below net",
			net:       ptr(money.Amount(1000)),
			costUnits: map[string]money.Amount{"100": 400, "200": 500},
		},
		{
			name:      "cost unit sum equal to net",
			net:       ptr(money.Amount(1000)),
			costUnits: map[string]money.Amount{"100": 400, "200": 600},
		},
		{
			name:      "cost unit sum above net",
			net:       ptr(money.Amount(1000)),
			costUnits: map[string]money.Amount{"100": 400, "200": 700},
			wantErr:   "sum of cost unit amounts 1100.000000 greater than invoice net sum 1000.000000",
		},
		{
			name:      "cost unit sum above net without net",
			costUnits: map[string]money.Amount{"100": 999999},
		},
		{
			name:           "cost unit sum within converted net",
			net:            ptr(money.Amount(1000)),
			conversionRate: ptr(money.Rate(1.5)),
			costUnits:      map[string]money.Amount{"100": 1400},
		},
		{
			name:      "empty cost unit number",
			net:       ptr(money.Amount(1000)),
			costUnits: map[string]money.Amount{"": 100},
			wantErr:   "empty costUnit string",
		},
		{
			name:      "zero cost unit amount",
			net:       ptr(money.Amount(1000)),
			costUnits: map[string]money.Amount{"100": 0},
			wantErr:   "cost unit '100' amount must not be zero",
		},
		{
			name:      "negative cost unit amount",
			net:       ptr(money.Amount(1000)),
			costUnits: map[string]money.Amount{"100": -50},
			wantErr:   "cost unit '100' amount (-50.000000) must not be negative",
		},
		{
			name:        "cost center sum above net",
			net:         ptr(money.Amount(1000)),
			costCenters: map[string]money.Amount{"1000": 1100},
			wantErr:     "sum of cost center amounts 1100.000000 greater than invoice net sum 1000.000000",
		},
		{
			name:        "empty cost center number",
			net:         ptr(money.Amount(1000)),
			costCenters: map[string]money.Amount{"": 100},
			wantErr:     "empty costCenter string",
		},
		{
			name:        "zero cost center amount",
			net:         ptr(money.Amount(1000)),
			costCenters: map[string]money.Amount{"1000": 0},
			wantErr:     "cost center '1000' amount must not be zero",
		},
		{
			name:        "negative cost center amount",
			net:         ptr(money.Amount(1000)),
			costCenters: map[string]money.Amount{"1000": -50},
			wantErr:     "cost center '1000' amount (-50.000000) must not be negative",
		},
		{
			// Cost centers and cost units are independent splits of the net
			// amount, so they are validated against the net separately
			// instead of being summed up together.
			name:        "cost centers and cost units both equal to net",
			net:         ptr(money.Amount(1000)),
			costCenters: map[string]money.Amount{"1000": 1000},
			costUnits:   map[string]money.Amount{"100": 1000},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			// given
			inv := &Invoice{
				Net:            scenario.net,
				ConversionRate: scenario.conversionRate,
				CostCenters:    scenario.costCenters,
				CostUnits:      scenario.costUnits,
			}

			// when
			err := inv.Validate()

			// then
			if scenario.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.EqualError(t, err, scenario.wantErr)
		})
	}
}
