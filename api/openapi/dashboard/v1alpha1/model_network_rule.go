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

type NetworkRule struct {
	PortName string `json:"portName"`

	PortNumber int32 `json:"portNumber"`

	Group string `json:"group,omitempty"`

	HttpPath string `json:"httpPath,omitempty"`

	Url string `json:"url,omitempty"`
}

// AssertNetworkRuleRequired checks if the required fields are not zero-ed
func AssertNetworkRuleRequired(obj NetworkRule) error {
	elements := map[string]interface{}{
		"portName":   obj.PortName,
		"portNumber": obj.PortNumber,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseNetworkRuleRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of NetworkRule (e.g. [][]NetworkRule), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseNetworkRuleRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aNetworkRule, ok := obj.(NetworkRule)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertNetworkRuleRequired(aNetworkRule)
	})
}