package domonda

import (
	"context"
	"fmt"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/nullable"
)

// GLAccount represents a general ledger account
type GLAccount struct {
	Number   account.Number
	Name     nullable.TrimmedString
	Category nullable.TrimmedString // Higher level description of the account
}

func PostGLAccounts(ctx context.Context, apiKey string, accounts []*GLAccount) error {
	response, err := postJSON(ctx, apiKey, "/masterdata/gl-accounts", accounts)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return nil
}
