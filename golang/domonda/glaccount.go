package domonda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/nullable"
)

// GLAccount represents a general ledger account
type GLAccount struct {
	Number   account.Number         // Alphanumeric account number
	Name     nullable.TrimmedString // Name of the account
	Category nullable.TrimmedString // Higher level description of the account
	ObjectNo account.NullableNumber // Real estate object number
}

func (a *GLAccount) Validate() error {
	var err error
	if e := a.Number.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid GLAccount.Number %q: %w", a.Number, e))
	}
	if e := a.ObjectNo.Validate(); e != nil {
		err = errors.Join(err, fmt.Errorf("invalid GLAccount.ObjectNo %q: %w", a.ObjectNo, e))
	}
	return err
}

type ImportGLAccountResult struct {
	// General ledger account number
	Number account.Number

	// State of the partner after import
	State ImportState

	// Error message from the import in case of State "ERROR"
	Error string `json:",omitempty"`
}

// PostGLAccounts upserts general ledger accounts
// using the API endpoint https://domonda.app/api/public/masterdata/gl-accounts.
//
// nil for objectSpecificAccountNos means we do not care about this
func PostGLAccounts(ctx context.Context, apiKey string, accounts []*GLAccount, objectSpecificAccountNos *bool, source string) (results []*ImportGLAccountResult, err error) {
	for i, acc := range accounts {
		if e := acc.Validate(); e != nil {
			err = errors.Join(err, fmt.Errorf("GLAccount at index %d has error: %w", i, e))
		}
	}
	if err != nil {
		return nil, err
	}

	vals := make(url.Values)
	if objectSpecificAccountNos != nil {
		vals.Set("objectSpecificAccountNos", fmt.Sprint(*objectSpecificAccountNos))
	}
	if source != "" {
		vals.Set("source", source)
	}
	response, err := postJSON(ctx, apiKey, "/masterdata/gl-accounts", vals, accounts)
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
