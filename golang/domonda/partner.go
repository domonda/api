package domonda

import (
	"context"
	"errors"
	"fmt"
	"net/url"
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
	Name             notnull.TrimmedString
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

	// A single payment bank account for the partner.
	// IBAN and BIC are well suited as CSV columns.
	IBAN bank.NullableIBAN
	BIC  bank.NullableBIC
	// More payment bank accounts for the partner.
	// As struct better suited for JSON import.
	BankAccounts []bank.Account
}

func (p *Partner) Validate() error {
	var err error
	if p.Name.IsEmpty() {
		err = errors.Join(err, errors.New("empty Partner.Name"))
	}
	if e := p.Country.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid Partner.Country %q: %w", p.Country, e))
	}
	if e := p.Email.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid Partner.Email %q: %w", p.Email, e))
	}
	if e := p.VATIDNo.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid Partner.VATIDNo %q: %w", p.VATIDNo, e))
	}
	if e := p.VendorAccountNumber.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid Partner.VendorAccountNumber %q: %w", p.VendorAccountNumber, e))
	}
	if e := p.ClientAccountNumber.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid Partner.ClientAccountNumber %q: %w", p.ClientAccountNumber, e))
	}
	if e := p.IBAN.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid Partner.IBAN %q: %w", p.IBAN, e))
	}
	if e := p.BIC.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid Partner.BIC %q: %w", p.BIC, e))
	}
	for i := range p.BankAccounts {
		if e := p.BankAccounts[i].Validate(); e != nil {
			err = errors.Join(err, fmt.Errorf("invalid Partner.BankAccounts[%d] %s: %w", i, p.BankAccounts[i], e))
		}
	}
	return err
}

func (p *Partner) String() string {
	var b strings.Builder
	b.WriteString(p.Name.String())
	if p.VATIDNo.IsNotNull() {
		b.WriteString(" ")
		b.WriteString(string(p.VATIDNo))
	}
	return b.String()
}

// EqualAlternativeNames returns true if the partner
// has exactly all of the passed names in its
// AlternativeNames, but independent of order.
func (p *Partner) EqualAlternativeNames(names []string) bool {
	if len(p.AlternativeNames) != len(names) {
		return false
	}
	for _, name := range names {
		if !p.AlternativeNames.Contains(name) {
			return false
		}
	}
	return true
}

func (p *Partner) VendorAccountNumberUint() uint64 {
	u, _ := strconv.ParseUint(p.VendorAccountNumber.String(), 10, 64)
	return u
}

func (p *Partner) ClientAccountNumberUint() uint64 {
	u, _ := strconv.ParseUint(p.ClientAccountNumber.String(), 10, 64)
	return u
}

// PostPartners upserts partner companies.
// Endpoint: https://domonda.app/api/public/masterdata/partner-companies
func PostPartners(ctx context.Context, apiKey string, partners []*Partner, source string) error {
	var err error
	for i, partner := range partners {
		if e := partner.Validate(); e != nil {
			err = errors.Join(err, fmt.Errorf("Partner at index %d has error: %w", i, e))
		}
	}
	if err != nil {
		return err
	}

	vals := make(url.Values)
	if source != "" {
		vals.Set("source", source)
	}
	response, err := postJSON(ctx, apiKey, "/masterdata/partner-companies", vals, partners)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return nil
}
