package domonda

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/nullable"
)

const GLAccountsCSVHeader = `Number;Name;Category`

// GLAccount represents a general ledger account
type GLAccount struct {
	Number   account.Number
	Name     nullable.TrimmedString
	Category nullable.TrimmedString // Higher level description of the account
	ObjectNo account.NullableNumber // Real estate object number
}

func (a *GLAccount) Validate() error {
	var err error
	if e := a.Number.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid GLAccount.Number %q: %w", a.Number, e))
	}
	if e := a.ObjectNo.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid GLAccount.ObjectNo %q: %w", a.ObjectNo, e))
	}
	return err
}

// nil for objectSpecificAccountNos means we do not care about this
func PostGLAccounts(ctx context.Context, apiKey string, accounts []*GLAccount, objectSpecificAccountNos *bool, source string) error {
	var err error
	for i, acc := range accounts {
		if e := acc.Validate(); e != nil {
			err = errors.Join(err, fmt.Errorf("GLAccount at index %d has error: %w", i, e))
		}
	}
	if err != nil {
		return err
	}

	vals := make(url.Values)
	if objectSpecificAccountNos != nil {
		vals.Set("objectSpecificAccountNos", fmt.Sprint(*objectSpecificAccountNos))
	}
	if source != "" {
		vals.Set("source", source)
	}
	response, err := postJSON(ctx, apiKey, "/masterdata/gl-accounts", vals, accounts)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return nil
}
