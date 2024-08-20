package domonda

import (
	"context"

	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/notnull"
)

type BankAccount struct {
	Holder notnull.TrimmedString
	IBAN   bank.IBAN
	BIC    bank.BIC
}

func PostBankAccounts(ctx context.Context, apiKey string, bankAccounts []*BankAccount) error {
	panic("TODO: implement me")
}
