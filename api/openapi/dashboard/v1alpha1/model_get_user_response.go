/*
 * Cosmo Dashboard API
 *
 * Manipulate cosmo dashboard resource API
 *
 * API version: v1alpha1
 * Contact: jlandowner8@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package v1alpha1

type GetUserResponse struct {
	User *User `json:"user"`
}

// AssertGetUserResponseRequired checks if the required fields are not zero-ed
func AssertGetUserResponseRequired(obj GetUserResponse) error {
	elements := map[string]interface{}{
		"user": obj.User,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if obj.User != nil {
		if err := AssertUserRequired(*obj.User); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseGetUserResponseRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of GetUserResponse (e.g. [][]GetUserResponse), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseGetUserResponseRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aGetUserResponse, ok := obj.(GetUserResponse)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertGetUserResponseRequired(aGetUserResponse)
	})
}