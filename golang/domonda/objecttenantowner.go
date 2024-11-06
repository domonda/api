package domonda

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/notnull"
)

type ObjectTenantOwner struct {
	ObjectNo      account.Number
	TenantOwnerNo int64
	UnitNo        int64
	OwnerLinkNo   int64
	Owner         notnull.TrimmedString
}

func (o *ObjectTenantOwner) Validate() error {
	var (
		err  error
		errs []error
	)
	if err = o.ObjectNo.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("ObjectTenantOwner.ObjectNo: %w", err))
	}
	if o.Owner.IsEmpty() {
		errs = append(errs, errors.New("empty ObjectTenantOwner.Owner"))
	}
	return errors.Join(errs...)
}

func PostObjectTenantOwners(ctx context.Context, apiKey string, tenantOwners []*ObjectTenantOwner, source string) error {
	var err error
	for i, obj := range tenantOwners {
		if e := obj.Validate(); e != nil {
			err = errors.Join(err, fmt.Errorf("ObjectTenantOwner at index %d has error: %w", i, e))
		}
	}
	if err != nil {
		return err
	}

	vals := make(url.Values)
	if source != "" {
		vals.Set("source", source)
	}
	response, err := postJSON(ctx, apiKey, "/masterdata/real-estate-object-tenant-owners", vals, tenantOwners)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return nil
}
