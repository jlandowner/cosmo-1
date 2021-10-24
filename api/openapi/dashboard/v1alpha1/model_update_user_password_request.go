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

type UpdateUserPasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`

	NewPassword string `json:"newPassword"`
}

// AssertUpdateUserPasswordRequestRequired checks if the required fields are not zero-ed
func AssertUpdateUserPasswordRequestRequired(obj UpdateUserPasswordRequest) error {
	elements := map[string]interface{}{
		"currentPassword": obj.CurrentPassword,
		"newPassword":     obj.NewPassword,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseUpdateUserPasswordRequestRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of UpdateUserPasswordRequest (e.g. [][]UpdateUserPasswordRequest), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseUpdateUserPasswordRequestRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aUpdateUserPasswordRequest, ok := obj.(UpdateUserPasswordRequest)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertUpdateUserPasswordRequestRequired(aUpdateUserPasswordRequest)
	})
}
