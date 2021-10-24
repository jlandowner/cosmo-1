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

type GetWorkspaceResponse struct {
	Workspace *Workspace `json:"workspace"`
}

// AssertGetWorkspaceResponseRequired checks if the required fields are not zero-ed
func AssertGetWorkspaceResponseRequired(obj GetWorkspaceResponse) error {
	elements := map[string]interface{}{
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

// AssertRecurseGetWorkspaceResponseRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of GetWorkspaceResponse (e.g. [][]GetWorkspaceResponse), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseGetWorkspaceResponseRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aGetWorkspaceResponse, ok := obj.(GetWorkspaceResponse)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertGetWorkspaceResponseRequired(aGetWorkspaceResponse)
	})
}