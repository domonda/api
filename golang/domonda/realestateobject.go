package domonda

//go:generate go-enum $GOFILE

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
)

const RealEstateObjectsCSVHeader = `Type;Number;AccountingArea;UserAccount;Description;StreetAddress;AlternativeAddresses;ZipCode;City;Country;IBAN;BIC;Active`

type RealEstateObject struct {
	Type                 RealEstateObjectType
	Number               account.Number
	AccountingArea       account.NullableNumber
	UserAccount          account.NullableNumber
	Description          nullable.TrimmedString
	StreetAddress        notnull.TrimmedString
	AlternativeAddresses nullable.StringArray
	ZipCode              nullable.TrimmedString
	City                 nullable.TrimmedString
	Country              country.Code
	IBAN                 bank.NullableIBAN
	BIC                  bank.NullableBIC
	Active               bool
}

func (o *RealEstateObject) Validate() error {
	var (
		err  error
		errs []error
	)
	if err = o.Type.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("RealEstateObject.Type: %w", err))
	}
	if err = o.Number.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("RealEstateObject.Number: %w", err))
	}
	if err = o.AccountingArea.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("RealEstateObject.AccountingArea: %w", err))
	}
	if err = o.UserAccount.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("RealEstateObject.UserAccount: %w", err))
	}
	if o.Country, err = o.Country.NormalizedWithAltCodes(); err != nil {
		errs = append(errs, fmt.Errorf("RealEstateObject.Country: %w", err))
	}
	if o.IBAN, err = o.IBAN.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("RealEstateObject.IBAN: %w", err))
	}
	if o.BIC, err = o.BIC.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("RealEstateObject.BIC: %w", err))
	}
	return errors.Join(errs...)
}

type RealEstateObjectType string //#enum

const (
	RealEstateObjectTypeHI      RealEstateObjectType = "HI"
	RealEstateObjectTypeWEG     RealEstateObjectType = "WEG"
	RealEstateObjectTypeSUB     RealEstateObjectType = "SUB"
	RealEstateObjectTypeKREIS   RealEstateObjectType = "KREIS"
	RealEstateObjectTypeMANDANT RealEstateObjectType = "MANDANT"
)

// Valid indicates if r is any of the valid values for RealEstateObjectType
func (r RealEstateObjectType) Valid() bool {
	switch r {
	case
		RealEstateObjectTypeHI,
		RealEstateObjectTypeWEG,
		RealEstateObjectTypeSUB,
		RealEstateObjectTypeKREIS,
		RealEstateObjectTypeMANDANT:
		return true
	}
	return false
}

// Validate returns an error if r is none of the valid values for RealEstateObjectType
func (r RealEstateObjectType) Validate() error {
	if !r.Valid() {
		return fmt.Errorf("invalid value %#v for type domonda.RealEstateObjectType", r)
	}
	return nil
}

// String implements the fmt.Stringer interface for RealEstateObjectType
func (r RealEstateObjectType) String() string {
	return string(r)
}

func PostRealEstateObjects(ctx context.Context, apiKey string, objects []*RealEstateObject, source string) error {
	vals := make(url.Values)
	if source != "" {
		vals.Set("source", source)
	}
	response, err := postJSON(ctx, apiKey, "/masterdata/real-estate-objects", vals, objects)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return nil
}
