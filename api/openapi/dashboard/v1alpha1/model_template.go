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

type Template struct {
	Name string `json:"name"`

	Description string `json:"description,omitempty"`

	RequiredVars []TemplateRequiredVars `json:"requiredVars,omitempty"`

	IsDefaultUserAddon *bool `json:"isDefaultUserAddon,omitempty"`

	IsClusterScope bool `json:"isClusterScope,omitempty"`
}

// AssertTemplateRequired checks if the required fields are not zero-ed
func AssertTemplateRequired(obj Template) error {
	elements := map[string]interface{}{
		"name": obj.Name,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.RequiredVars {
		if err := AssertTemplateRequiredVarsRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseTemplateRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Template (e.g. [][]Template), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseTemplateRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aTemplate, ok := obj.(Template)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertTemplateRequired(aTemplate)
	})
}
