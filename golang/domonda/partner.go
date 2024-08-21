package domonda

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/email"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
	"github.com/domonda/go-types/vat"
)

type Partner struct {
	Name             string
	AlternativeNames notnull.StringArray // used when merging

	// main location
	Street    nullable.TrimmedString
	City      nullable.TrimmedString
	ZIP       nullable.TrimmedString
	Country   country.NullableCode
	Phone     nullable.TrimmedString
	Email     email.NullableAddress
	Website   nullable.TrimmedString
	CompRegNo nullable.TrimmedString
	TaxIDNo   nullable.TrimmedString
	VATIDNo   vat.NullableID

	// partner accounts
	VendorAccountNumber account.NullableNumber // "" means not set -> will not create a partner account
	ClientAccountNumber account.NullableNumber // "" means not set -> will not create a partner account

	IBAN  bank.NullableIBAN
	BIC   bank.NullableBIC
	IBAN2 bank.NullableIBAN
	BIC2  bank.NullableBIC
}

func (p *Partner) String() string {
	var b strings.Builder
	b.WriteString(p.Name)
	if p.VATIDNo.IsNotNull() {
		b.WriteString(" ")
		b.WriteString(string(p.VATIDNo))
	}
	return b.String()
}

func (p *Partner) VendorAccountNumberUint() uint64 {
	u, _ := strconv.ParseUint(p.VendorAccountNumber.String(), 10, 64)
	return u
}

func (p *Partner) ClientAccountNumberUint() uint64 {
	u, _ := strconv.ParseUint(p.ClientAccountNumber.String(), 10, 64)
	return u
}

func PostPartners(ctx context.Context, apiKey string, partners []*Partner) error {
	response, err := postJSON(ctx, apiKey, "/masterdata/partner-companies", partners)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return nil
}
