package domonda

import (
	"context"
	"errors"
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

func (a *BankAccount) Validate() error {
	var err error
	if e := a.IBAN.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid BankAccount.IBAN %q: %w", a.IBAN, e))
	}
	if e := a.BIC.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid BankAccount.BIC %q: %w", a.BIC, e))
	}
	if !a.Currency.Valid() {
		err = errors.Join(err, fmt.Errorf("invalid BankAccount.Currency %q", a.Currency))
	}
	if a.Holder.IsEmpty() {
		err = errors.Join(err, errors.New("empty BankAccount.Holder"))
	}
	return err
}

// PostBankAccounts posts the given bankAccounts to the domonda API.
//
// Usage example:
//
//	curl -X POST \
//	  -H "Authorization: Bearer ${DOMONDA_API_KEY}" \
//	  -H "Content-Type: application/json" \
//	  --data "[]"" \
//	  --include \
//	  https://domonda.app/api/public/masterdata/bank-accounts
func PostBankAccounts(ctx context.Context, apiKey string, bankAccounts []*BankAccount, source string) error {
	var err error
	for i, acc := range bankAccounts {
		if e := acc.Validate(); e != nil {
			err = errors.Join(err, fmt.Errorf("BankAccount at index %d has error: %w", i, e))
		}
	}
	if err != nil {
		return err
	}

	vals := make(url.Values)
	if source != "" {
		vals.Set("source", source)
	}
	response, err := postJSON(ctx, apiKey, "/masterdata/bank-accounts", vals, bankAccounts)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return nil
}
