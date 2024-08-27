package domonda

import (
	"context"
	"errors"
	"fmt"

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

func PostGLAccounts(ctx context.Context, apiKey string, accounts []*GLAccount, objectSpecificAccountNos *bool) error {
	query := ""
	if objectSpecificAccountNos != nil {
		query = fmt.Sprintf("?objectSpecificAccountNos=%t", *objectSpecificAccountNos)
	}
	response, err := postJSON(ctx, apiKey, "/masterdata/gl-accounts"+query, accounts)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return nil
}
