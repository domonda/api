package domonda

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/email"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
	"github.com/domonda/go-types/strutil"
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

func (p *Partner) Normalize(resetInvalid bool) []error {
	var errs []error
	if p.Name.IsEmpty() {
		errs = append(errs, errors.New("Name is empty"))
	}
	// Trim whitespace and remove empty alternative names
	for i := 0; i < len(p.AlternativeNames); i++ {
		p.AlternativeNames[i] = strutil.TrimSpace(p.AlternativeNames[i])
		if p.AlternativeNames[i] == "" {
			p.AlternativeNames = slices.Delete(p.AlternativeNames, i, i+1)
			i--
		}
	}
	var err error
	p.Country, err = p.Country.Normalized()
	if err != nil {
		errs = append(errs, fmt.Errorf("Country '%s' has error: %w", p.Country, err))
		if resetInvalid {
			p.Country.SetNull()
		}
	}
	p.VATIDNo, err = p.VATIDNo.Normalized()
	if err != nil {
		errs = append(errs, fmt.Errorf("VATIDNo '%s' has error: %w", p.VATIDNo, err))
		if resetInvalid {
			p.VATIDNo.SetNull()
		}
	}
	// if p.VATIDNo.ValidAndNotNull() && p.Country.ValidAndNotNull() {
	// 	vatCountry := p.VATIDNo.Get().CountryCode()
	// 	if vatCountry != vat.MOSSSchemaVATCountryCode && vatCountry != p.Country.Get() {
	// 		errs = append(errs, fmt.Errorf("Country '%s' is different from VATIDNo '%s' country code", p.Country, p.VATIDNo))
	// 		if resetInvalid {
	// 			if p.Street.IsNotNull() && p.ZIP.IsNotNull() && p.City.IsNotNull() {
	// 				// If there is a complete address, don't set the country to null
	// 				p.VATIDNo.SetNull()
	// 			} else {
	// 				// If there is no address, keep the VAT ID
	// 				p.Country.SetNull()
	// 			}
	// 		}
	// 	}
	// }
	p.Email, err = p.Email.Normalized()
	if err != nil {
		errs = append(errs, fmt.Errorf("Email '%s' has error: %w", p.Email, err))
		if resetInvalid {
			p.Email.SetNull()
		}
	}
	if err = p.VendorAccountNumber.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("VendorAccountNumber '%s' has error: %w", p.VendorAccountNumber, err))
		if resetInvalid {
			p.VendorAccountNumber.SetNull()
		}
	}
	if err = p.ClientAccountNumber.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("ClientAccountNumber '%s' has error: %w", p.ClientAccountNumber, err))
		if resetInvalid {
			p.ClientAccountNumber.SetNull()
		}
	}
	p.IBAN, err = bank.NullableIBAN(strings.ToUpper(string(p.IBAN))).Normalized()
	if err != nil {
		errs = append(errs, fmt.Errorf("IBAN '%s' has error: %w", p.IBAN, err))
		if resetInvalid {
			p.IBAN.SetNull()
		}
	}
	p.BIC, err = p.BIC.Normalized()
	if err != nil {
		errs = append(errs, fmt.Errorf("BIC '%s' has error: %w", p.BIC, err))
		if resetInvalid {
			p.BIC.SetNull()
		}
	}
	for i := 0; i < len(p.BankAccounts); i++ {
		if err = p.BankAccounts[i].Normalize(); err != nil {
			errs = append(errs, fmt.Errorf("BankAccounts[%d] has error: %w", i, err))
			if resetInvalid {
				p.BankAccounts = slices.Delete(p.BankAccounts, i, i+1)
				i--
			}
		}
	}
	if p.IBAN.IsNotNull() {
		// Use IBAN/BIC as first bank account
		p.BankAccounts = append(
			[]bank.Account{{IBAN: p.IBAN.Get(), BIC: p.BIC}},
			p.BankAccounts...,
		)
		p.IBAN.SetNull()
		p.BIC.SetNull()
	}
	return errs
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

func (p *Partner) HasLocation() bool {
	return p.Street.IsNotNull() ||
		p.City.IsNotNull() ||
		p.ZIP.IsNotNull() ||
		p.Country.IsNotNull() ||
		p.Phone.IsNotNull() ||
		p.Email.IsNotNull() ||
		p.Website.IsNotNull() ||
		p.CompRegNo.IsNotNull() ||
		p.TaxIDNo.IsNotNull() ||
		p.VATIDNo.IsNotNull()
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
