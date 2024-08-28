package domonda

import (
	"context"
	"fmt"
	"net/url"

	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
)

const BankAccountsCSVHeader = `IBAN;BIC;Currency;Holder;AccountNumber;Name;Description`

// BankAccount represents a checking account
type BankAccount struct {
	IBAN     bank.IBAN
	BIC      bank.BIC
	Currency money.Currency
	Holder   notnull.TrimmedString

	// Optional
	AccountNumber nullable.TrimmedString
	Name          nullable.TrimmedString
	Description   nullable.TrimmedString
}

func PostBankAccounts(ctx context.Context, apiKey string, bankAccounts []*BankAccount, source string) error {
	vals := make(url.Values)
	vals.Set("source", source)
	response, err := postJSON(ctx, apiKey, "/masterdata/bank-accounts", vals, bankAccounts)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return nil
}
