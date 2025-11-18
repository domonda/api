package domonda

//go:generate go tool go-enum $GOFILE

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

// RealEstateObject represents a real estate property managed in the system.
// Objects are used for property management, accounting segregation, and organizing
// related transactions and documents. Each object is identified by its Number.
type RealEstateObject struct {
	// Type specifies the kind of real estate object (WEG, HI, SUB, KREIS, MANDANT, MRG, MHV, SEV)
	Type RealEstateObjectType

	// Number is the unique identifier for this object (alphanumeric)
	Number account.Number

	// AccountingArea is an optional accounting segregation identifier
	AccountingArea account.NullableNumber

	// UserAccount is an optional user account number associated with this object
	UserAccount account.NullableNumber

	// Description provides additional details about the property
	Description nullable.TrimmedString

	// StreetAddress is the primary street address (required)
	StreetAddress notnull.TrimmedString

	// AlternativeAddresses contains additional addresses for the same property
	AlternativeAddresses nullable.StringArray

	// ZipCode is the postal/ZIP code
	ZipCode nullable.TrimmedString

	// City is the city name
	City nullable.TrimmedString

	// Country is the ISO 3166-1 alpha-2 country code (e.g., "DE", "AT")
	Country country.Code

	// BankAccounts are payment bank accounts associated with this object
	BankAccounts []bank.Account

	// Active indicates if this object is currently active
	Active bool
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
	if o.Country, err = o.Country.Normalized(); err != nil {
		errs = append(errs, fmt.Errorf("RealEstateObject.Country: %w", err))
	}
	if o.StreetAddress.IsEmpty() {
		errs = append(errs, errors.New("empty RealEstateObject.StreetAddress"))
	}
	for i := range o.BankAccounts {
		if err = o.BankAccounts[i].Validate(); err != nil {
			errs = append(errs, fmt.Errorf("RealEstateObject.BankAccounts[%d]: %w", i, err))
		}
	}
	return errors.Join(errs...)
}

// RealEstateObjectType categorizes real estate objects by their legal and management structure.
// Different types have different requirements and business rules in the system.
type RealEstateObjectType string //#enum

const (
	// RealEstateObjectTypeWEG represents a condominium owners' association
	// (Wohnungseigentümergemeinschaft in German/Austrian law)
	RealEstateObjectTypeWEG RealEstateObjectType = "WEG"

	// RealEstateObjectTypeHI represents a house/building (Hausverwaltung)
	RealEstateObjectTypeHI RealEstateObjectType = "HI"

	// RealEstateObjectTypeSUB represents a sub-object or unit within a larger property
	RealEstateObjectTypeSUB RealEstateObjectType = "SUB"

	// RealEstateObjectTypeKREIS represents a virtual grouping object (accounting circle)
	RealEstateObjectTypeKREIS RealEstateObjectType = "KREIS"

	// RealEstateObjectTypeMANDANT represents a virtual client-level object
	RealEstateObjectTypeMANDANT RealEstateObjectType = "MANDANT"

	// RealEstateObjectTypeMRG represents a property subject to Austrian rent control law
	// (Mietrechtsgesetz - MRG)
	RealEstateObjectTypeMRG RealEstateObjectType = "MRG"

	// RealEstateObjectTypeMHV represents a rental property management object
	// (Miethausverwaltung)
	RealEstateObjectTypeMHV RealEstateObjectType = "MHV"

	// RealEstateObjectTypeSEV represents a separate property management object
	// (Sondereigentumsverwaltung)
	RealEstateObjectTypeSEV RealEstateObjectType = "SEV"
)

// Valid indicates if r is any of the valid values for RealEstateObjectType
func (r RealEstateObjectType) Valid() bool {
	switch r {
	case
		RealEstateObjectTypeWEG,
		RealEstateObjectTypeHI,
		RealEstateObjectTypeSUB,
		RealEstateObjectTypeKREIS,
		RealEstateObjectTypeMANDANT,
		RealEstateObjectTypeMRG,
		RealEstateObjectTypeMHV,
		RealEstateObjectTypeSEV:
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

// Enums returns all valid values for RealEstateObjectType
func (RealEstateObjectType) Enums() []RealEstateObjectType {
	return []RealEstateObjectType{
		RealEstateObjectTypeWEG,
		RealEstateObjectTypeHI,
		RealEstateObjectTypeSUB,
		RealEstateObjectTypeKREIS,
		RealEstateObjectTypeMANDANT,
		RealEstateObjectTypeMRG,
		RealEstateObjectTypeMHV,
		RealEstateObjectTypeSEV,
	}
}

// EnumStrings returns all valid values for RealEstateObjectType as strings
func (RealEstateObjectType) EnumStrings() []string {
	return []string{
		"WEG",
		"HI",
		"SUB",
		"KREIS",
		"MANDANT",
		"MRG",
		"MHV",
		"SEV",
	}
}

// String implements the fmt.Stringer interface for RealEstateObjectType
func (r RealEstateObjectType) String() string {
	return string(r)
}

// IsVirtual returns true if this object type represents a virtual grouping
// rather than a physical property. Virtual objects (KREIS, MANDANT) are used
// for organizing and aggregating data but don't represent actual real estate.
func (r RealEstateObjectType) IsVirtual() bool {
	return r == RealEstateObjectTypeKREIS || r == RealEstateObjectTypeMANDANT
}

// PostRealEstateObjects upserts (inserts or updates) real estate objects via the Domonda API.
// Objects are identified by their Number field - if an object with the same number exists,
// it will be updated; otherwise, a new object is created.
//
// Arguments:
//   - ctx:     Context for the HTTP request (for cancellation and timeouts)
//   - apiKey:  API key (bearer token) for authentication
//   - objects: Slice of real estate objects to import
//   - source:  Optional identifier for the data source (e.g., your company name)
//
// Returns an error if validation fails or the API request fails.
// The function validates all objects before sending the request.
//
// API endpoint: https://domonda.app/api/public/masterdata/real-estate-objects
func PostRealEstateObjects(ctx context.Context, apiKey string, objects []*RealEstateObject, source string) error {
	var err error
	for i, obj := range objects {
		if e := obj.Validate(); e != nil {
			err = errors.Join(err, fmt.Errorf("RealEstateObject at index %d has error: %w", i, e))
		}
	}
	if err != nil {
		return err
	}

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
