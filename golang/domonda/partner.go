package domonda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"slices"
	"sort"
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

type ImportPartnerResult struct {
	// Shows how the input was normalized
	NormalizedInput *Partner

	// Warnings from normalizing and validating the input
	InputWarnings []string

	// Data of the partner after import
	// TODO replace json.RawMessage with struct types
	PartnerCompany   json.RawMessage `json:",omitempty"`
	PartnerLocations json.RawMessage `json:",omitempty"` // Main location first
	VendorAccount    json.RawMessage `json:",omitempty"`
	ClientAccount    json.RawMessage `json:",omitempty"`
	PaymentPresets   json.RawMessage `json:",omitempty"`

	// State of the partner after import
	State ImportState

	// Error message from the import in case of State "ERROR"
	Error string `json:",omitempty"`
}

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

// NormalizedAlternativeNames returns the alternative names
// with all whitespace trimmed and empty strings removed.
// The names are sorted alphabetically.
func (p *Partner) NormalizedAlternativeNames() []string {
	names := make([]string, 0, len(p.AlternativeNames))
	for _, name := range p.AlternativeNames {
		if n := strutil.TrimSpace(name); n != "" {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	return names
}

func (p *Partner) Normalize(resetInvalid bool) []error {
	var errs []error
	if p.Name.IsEmpty() {
		errs = append(errs, errors.New("Name is empty"))
	}

	p.AlternativeNames = p.NormalizedAlternativeNames()

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
		// Prepend IBAN/BIC as first bank account
		p.BankAccounts = append(
			[]bank.Account{{IBAN: p.IBAN.Get(), BIC: p.BIC}},
			p.BankAccounts...,
		)
		// After prepending, set IBAN and BIC to null
		p.IBAN.SetNull()
		p.BIC.SetNull()
	}
	// Check for duplicate bank accounts
	for i := 0; i < len(p.BankAccounts); i++ {
		seenBefore := slices.ContainsFunc(p.BankAccounts[:i], func(b bank.Account) bool {
			return b.IBAN == p.BankAccounts[i].IBAN
		})
		if seenBefore {
			errs = append(errs, fmt.Errorf("duplicate bank account: %s", p.BankAccounts[i]))
			if resetInvalid {
				p.BankAccounts = slices.Delete(p.BankAccounts, i, i+1)
				i--
			}
		}
	}
	return errs
}

func (p *Partner) String() string {
	var b strings.Builder
	b.WriteString(p.Name.String())
	if p.Country.IsNotNull() {
		fmt.Fprintf(&b, "|%s", p.Country)
	}
	if p.VATIDNo.IsNotNull() {
		fmt.Fprintf(&b, "|%s", p.VATIDNo)
	}
	if p.VendorAccountNumber.IsNotNull() {
		fmt.Fprintf(&b, "|Vendor:%s", p.VendorAccountNumber)
	}
	if p.ClientAccountNumber.IsNotNull() {
		fmt.Fprintf(&b, "|Client:%s", p.ClientAccountNumber)
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
		name = strings.TrimSpace(name)
		if !slices.ContainsFunc(p.AlternativeNames, func(altName string) bool {
			return name == strings.TrimSpace(altName)
		}) {
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
func PostPartners(ctx context.Context, apiKey string, partners []*Partner, failOnInvalid, useCleanedInvalid, allOrNone bool, source string) (results []ImportPartnerResult, err error) {
	vals := make(url.Values)
	if failOnInvalid {
		vals.Set("failOnInvalid", "true")
	}
	if useCleanedInvalid {
		vals.Set("useCleanedInvalid", "true")
	}
	if allOrNone {
		vals.Set("allOrNone", "true")
	}
	if source != "" {
		vals.Set("source", source)
	}
	response, err := postJSON(ctx, apiKey, "/masterdata/partner-companies", vals, partners)
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
