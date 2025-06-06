package domonda

//go:generate go tool go-enum $GOFILE

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// PostObjectInstancesWithIDProp updates or inserts instances of the class "className"
// using the prop idPropName as the identifier for the objects.
// The objectsProps is a slice of maps, where each map represents the properties of an object.
// The ID prop with idPropName must be present in each object.
// The source argument is used to identify the source of the request.
func PostObjectInstancesWithIDProp(ctx context.Context, apiKey string, className, idPropName string, objectsProps []map[string]any, source string) (err error) {
	if className == "" {
		return errors.New("className is required")
	}
	if idPropName == "" {
		return errors.New("idPropName is required")
	}
	for i, props := range objectsProps {
		if props[idPropName] == nil {
			err = errors.Join(err, fmt.Errorf("object at index %d has no ID prop %q", i, idPropName))
		}
	}
	if err != nil {
		return err
	}

	vals := make(url.Values)
	if source != "" {
		vals.Set("source", source)
	}
	endpoint := fmt.Sprintf("/masterdata/upsert-objects/%s/id-prop/%s", className, idPropName)
	response, err := postJSON(ctx, apiKey, endpoint, vals, objectsProps)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return nil
}
