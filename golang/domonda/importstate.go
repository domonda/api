package domonda

import "fmt"

//go:generate go tool go-enum $GOFILE

// ImportState represents the result state of importing a single data item
// (partner, GL account, bank account, etc.) via the Domonda API.
type ImportState string //#enum

const (
	// ImportStateUnchanged indicates the item already exists with identical data
	ImportStateUnchanged ImportState = "UNCHANGED"

	// ImportStateUpdated indicates an existing item was updated with new data
	ImportStateUpdated ImportState = "UPDATED"

	// ImportStateCreated indicates a new item was created
	ImportStateCreated ImportState = "CREATED"

	// ImportStateError indicates the import failed for this item
	// Check the Error field in the result for details
	ImportStateError ImportState = "ERROR"
)

// Valid indicates if i is any of the valid values for ImportState
func (i ImportState) Valid() bool {
	switch i {
	case
		ImportStateUnchanged,
		ImportStateUpdated,
		ImportStateCreated,
		ImportStateError:
		return true
	}
	return false
}

// Validate returns an error if i is none of the valid values for ImportState
func (i ImportState) Validate() error {
	if !i.Valid() {
		return fmt.Errorf("invalid value %#v for type domonda.ImportState", i)
	}
	return nil
}

// Enums returns all valid values for ImportState
func (ImportState) Enums() []ImportState {
	return []ImportState{
		ImportStateUnchanged,
		ImportStateUpdated,
		ImportStateCreated,
		ImportStateError,
	}
}

// EnumStrings returns all valid values for ImportState as strings
func (ImportState) EnumStrings() []string {
	return []string{
		"UNCHANGED",
		"UPDATED",
		"CREATED",
		"ERROR",
	}
}

// String implements the fmt.Stringer interface for ImportState
func (i ImportState) String() string {
	return string(i)
}
