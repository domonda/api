package domonda

import "fmt"

//go:generate go tool go-enum $GOFILE

type ImportState string //#enum

const (
	ImportStateUnchanged ImportState = "UNCHANGED"
	ImportStateUpdated   ImportState = "UPDATED"
	ImportStateCreated   ImportState = "CREATED"
	ImportStateError     ImportState = "ERROR"
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

// String implements the fmt.Stringer interface for ImportState
func (i ImportState) String() string {
	return string(i)
}
