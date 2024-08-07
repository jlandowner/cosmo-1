// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: dashboard/v1alpha1/workspace.proto

package dashboardv1alpha1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on NetworkRule with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *NetworkRule) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on NetworkRule with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in NetworkRuleMultiError, or
// nil if none found.
func (m *NetworkRule) ValidateAll() error {
	return m.validate(true)
}

func (m *NetworkRule) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if val := m.GetPortNumber(); val <= 0 || val >= 65536 {
		err := NetworkRuleValidationError{
			field:  "PortNumber",
			reason: "value must be inside range (0, 65536)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for CustomHostPrefix

	// no validation rules for HttpPath

	// no validation rules for Url

	// no validation rules for Public

	if len(errors) > 0 {
		return NetworkRuleMultiError(errors)
	}

	return nil
}

// NetworkRuleMultiError is an error wrapping multiple validation errors
// returned by NetworkRule.ValidateAll() if the designated constraints aren't met.
type NetworkRuleMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m NetworkRuleMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m NetworkRuleMultiError) AllErrors() []error { return m }

// NetworkRuleValidationError is the validation error returned by
// NetworkRule.Validate if the designated constraints aren't met.
type NetworkRuleValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e NetworkRuleValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e NetworkRuleValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e NetworkRuleValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e NetworkRuleValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e NetworkRuleValidationError) ErrorName() string { return "NetworkRuleValidationError" }

// Error satisfies the builtin error interface
func (e NetworkRuleValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sNetworkRule.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = NetworkRuleValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = NetworkRuleValidationError{}

// Validate checks the field values on WorkspaceSpec with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *WorkspaceSpec) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on WorkspaceSpec with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in WorkspaceSpecMultiError, or
// nil if none found.
func (m *WorkspaceSpec) ValidateAll() error {
	return m.validate(true)
}

func (m *WorkspaceSpec) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Template

	// no validation rules for Replicas

	// no validation rules for Vars

	for idx, item := range m.GetNetwork() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, WorkspaceSpecValidationError{
						field:  fmt.Sprintf("Network[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, WorkspaceSpecValidationError{
						field:  fmt.Sprintf("Network[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return WorkspaceSpecValidationError{
					field:  fmt.Sprintf("Network[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return WorkspaceSpecMultiError(errors)
	}

	return nil
}

// WorkspaceSpecMultiError is an error wrapping multiple validation errors
// returned by WorkspaceSpec.ValidateAll() if the designated constraints
// aren't met.
type WorkspaceSpecMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m WorkspaceSpecMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m WorkspaceSpecMultiError) AllErrors() []error { return m }

// WorkspaceSpecValidationError is the validation error returned by
// WorkspaceSpec.Validate if the designated constraints aren't met.
type WorkspaceSpecValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e WorkspaceSpecValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e WorkspaceSpecValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e WorkspaceSpecValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e WorkspaceSpecValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e WorkspaceSpecValidationError) ErrorName() string { return "WorkspaceSpecValidationError" }

// Error satisfies the builtin error interface
func (e WorkspaceSpecValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sWorkspaceSpec.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = WorkspaceSpecValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = WorkspaceSpecValidationError{}

// Validate checks the field values on WorkspaceStatus with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *WorkspaceStatus) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on WorkspaceStatus with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// WorkspaceStatusMultiError, or nil if none found.
func (m *WorkspaceStatus) ValidateAll() error {
	return m.validate(true)
}

func (m *WorkspaceStatus) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Phase

	// no validation rules for MainUrl

	if all {
		switch v := interface{}(m.GetLastStartedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, WorkspaceStatusValidationError{
					field:  "LastStartedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, WorkspaceStatusValidationError{
					field:  "LastStartedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetLastStartedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return WorkspaceStatusValidationError{
				field:  "LastStartedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return WorkspaceStatusMultiError(errors)
	}

	return nil
}

// WorkspaceStatusMultiError is an error wrapping multiple validation errors
// returned by WorkspaceStatus.ValidateAll() if the designated constraints
// aren't met.
type WorkspaceStatusMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m WorkspaceStatusMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m WorkspaceStatusMultiError) AllErrors() []error { return m }

// WorkspaceStatusValidationError is the validation error returned by
// WorkspaceStatus.Validate if the designated constraints aren't met.
type WorkspaceStatusValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e WorkspaceStatusValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e WorkspaceStatusValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e WorkspaceStatusValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e WorkspaceStatusValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e WorkspaceStatusValidationError) ErrorName() string { return "WorkspaceStatusValidationError" }

// Error satisfies the builtin error interface
func (e WorkspaceStatusValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sWorkspaceStatus.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = WorkspaceStatusValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = WorkspaceStatusValidationError{}

// Validate checks the field values on Workspace with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Workspace) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Workspace with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in WorkspaceMultiError, or nil
// if none found.
func (m *Workspace) ValidateAll() error {
	return m.validate(true)
}

func (m *Workspace) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Name

	// no validation rules for OwnerName

	if all {
		switch v := interface{}(m.GetSpec()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, WorkspaceValidationError{
					field:  "Spec",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, WorkspaceValidationError{
					field:  "Spec",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetSpec()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return WorkspaceValidationError{
				field:  "Spec",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetStatus()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, WorkspaceValidationError{
					field:  "Status",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, WorkspaceValidationError{
					field:  "Status",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetStatus()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return WorkspaceValidationError{
				field:  "Status",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.Raw != nil {
		// no validation rules for Raw
	}

	if m.RawInstance != nil {
		// no validation rules for RawInstance
	}

	if m.RawIngressRoute != nil {
		// no validation rules for RawIngressRoute
	}

	if m.DeletePolicy != nil {

		if _, ok := DeletePolicy_name[int32(m.GetDeletePolicy())]; !ok {
			err := WorkspaceValidationError{
				field:  "DeletePolicy",
				reason: "value must be one of the defined enum values",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return WorkspaceMultiError(errors)
	}

	return nil
}

// WorkspaceMultiError is an error wrapping multiple validation errors returned
// by Workspace.ValidateAll() if the designated constraints aren't met.
type WorkspaceMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m WorkspaceMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m WorkspaceMultiError) AllErrors() []error { return m }

// WorkspaceValidationError is the validation error returned by
// Workspace.Validate if the designated constraints aren't met.
type WorkspaceValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e WorkspaceValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e WorkspaceValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e WorkspaceValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e WorkspaceValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e WorkspaceValidationError) ErrorName() string { return "WorkspaceValidationError" }

// Error satisfies the builtin error interface
func (e WorkspaceValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sWorkspace.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = WorkspaceValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = WorkspaceValidationError{}
