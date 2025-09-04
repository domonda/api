package domonda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
	"github.com/domonda/go-types/uu"
)

// BankAccount represents a checking account
type BankAccount struct {
	IBAN     bank.IBAN
	BIC      bank.BIC
	Currency money.Currency
	Holder   notnull.TrimmedString

	// Optional
	AccountNumber nullable.TrimmedString `json:",omitempty"`
	Name          nullable.TrimmedString `json:",omitempty"`
	Description   nullable.TrimmedString `json:",omitempty"`
}

func (a *BankAccount) Validate() (err error) {
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

func (a *BankAccount) Normalize() (err error) {
	var e error
	if a.IBAN, e = a.IBAN.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid BankAccount.IBAN %q: %w", a.IBAN, e))
	}
	if a.BIC, e = a.BIC.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid BankAccount.BIC %q: %w", a.BIC, e))
	}
	if a.Currency, e = a.Currency.Normalized(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid BankAccount.Currency %q: %w", a.Currency, e))
	}
	if a.Holder.IsEmpty() {
		err = errors.Join(err, errors.New("empty BankAccount.Holder"))
	}
	return err
}

type ImportBankAccountResult struct {
	// ID of the bank account that was created or updated
	ID uu.NullableID `json:",omitzero"`

	BankAccount

	// State of the account after import
	State ImportState

	// Error message from the import in case of State "ERROR"
	Error string `json:",omitempty"`
}

// PostBankAccounts posts the given bankAccounts to the domonda API.
//
// Arguments:
//   - apiKey:          API key (bearer token) for the domonda API
//   - accounts:        Bank accounts to insert or update
//   - failOnInvalid:   Fail if any account data is invalid
//   - allOrNone:       Import either all accounts or none in case of any error
//   - source:          Optional name or ID of who did the import
//
// Usage example:
//
//	curl -X POST \
//	  -H "Authorization: Bearer ${DOMONDA_API_KEY}" \
//	  -H "Content-Type: application/json" \
//	  --data "[]"" \
//	  --include \
//	  https://domonda.app/api/public/masterdata/bank-accounts?failOnInvalid=true&source=MY_SERVICE
func PostBankAccounts(ctx context.Context, apiKey string, accounts []*BankAccount, failOnInvalid, allOrNone bool, source string) (results []*ImportBankAccountResult, err error) {
	for i, acc := range accounts {
		if e := acc.Normalize(); e != nil {
			err = errors.Join(err, fmt.Errorf("BankAccount at index %d has error: %w", i, e))
		}
	}
	if err != nil {
		return nil, err
	}

	vals := make(url.Values)
	if failOnInvalid {
		vals.Set("failOnInvalid", "true")
	}
	if allOrNone {
		vals.Set("allOrNone", "true")
	}
	if source != "" {
		vals.Set("source", source)
	}
	response, err := postJSON(ctx, apiKey, "/masterdata/bank-accounts", vals, accounts)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return results, nil
}
