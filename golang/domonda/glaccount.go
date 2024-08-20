package domonda

import (
	"context"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/nullable"
)

// GLAccount represents a general ledger account
type GLAccount struct {
	Number   account.Number
	Name     nullable.TrimmedString
	Category nullable.TrimmedString
}

func PostGLAccounts(ctx context.Context, apiKey string, accounts []GLAccount) error {
	panic("TODO: implement me")
}
