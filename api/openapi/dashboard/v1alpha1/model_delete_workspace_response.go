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

type DeleteWorkspaceResponse struct {
	Message string `json:"message"`

	Workspace *Workspace `json:"workspace"`
}

// AssertDeleteWorkspaceResponseRequired checks if the required fields are not zero-ed
func AssertDeleteWorkspaceResponseRequired(obj DeleteWorkspaceResponse) error {
	elements := map[string]interface{}{
		"message":   obj.Message,
		"workspace": obj.Workspace,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if obj.Workspace != nil {
		if err := AssertWorkspaceRequired(*obj.Workspace); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseDeleteWorkspaceResponseRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of DeleteWorkspaceResponse (e.g. [][]DeleteWorkspaceResponse), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseDeleteWorkspaceResponseRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aDeleteWorkspaceResponse, ok := obj.(DeleteWorkspaceResponse)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertDeleteWorkspaceResponseRequired(aDeleteWorkspaceResponse)
	})
}
