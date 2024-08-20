package domonda

import (
	"context"
	"fmt"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
)

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
}

func PostRealEstateObjects(ctx context.Context, apiKey string, objects []*RealEstateObject) error {
	panic("TODO: implement me")
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
		return fmt.Errorf("invalid value %#v for type idwell.RealEstateObjectType", r)
	}
	return nil
}

// String implements the fmt.Stringer interface for RealEstateObjectType
func (r RealEstateObjectType) String() string {
	return string(r)
}
