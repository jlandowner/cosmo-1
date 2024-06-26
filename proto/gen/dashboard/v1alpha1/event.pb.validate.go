// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: dashboard/v1alpha1/event.proto

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

// Validate checks the field values on Event with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Event) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Event with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in EventMultiError, or nil if none found.
func (m *Event) ValidateAll() error {
	return m.validate(true)
}

func (m *Event) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for User

	if all {
		switch v := interface{}(m.GetEventTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, EventValidationError{
					field:  "EventTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, EventValidationError{
					field:  "EventTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetEventTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return EventValidationError{
				field:  "EventTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Type

	// no validation rules for Note

	// no validation rules for Reason

	if all {
		switch v := interface{}(m.GetRegarding()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, EventValidationError{
					field:  "Regarding",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, EventValidationError{
					field:  "Regarding",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRegarding()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return EventValidationError{
				field:  "Regarding",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for ReportingController

	if all {
		switch v := interface{}(m.GetSeries()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, EventValidationError{
					field:  "Series",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, EventValidationError{
					field:  "Series",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetSeries()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return EventValidationError{
				field:  "Series",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.RegardingWorkspace != nil {
		// no validation rules for RegardingWorkspace
	}

	if len(errors) > 0 {
		return EventMultiError(errors)
	}

	return nil
}

// EventMultiError is an error wrapping multiple validation errors returned by
// Event.ValidateAll() if the designated constraints aren't met.
type EventMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EventMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EventMultiError) AllErrors() []error { return m }

// EventValidationError is the validation error returned by Event.Validate if
// the designated constraints aren't met.
type EventValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EventValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EventValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EventValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EventValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EventValidationError) ErrorName() string { return "EventValidationError" }

// Error satisfies the builtin error interface
func (e EventValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEvent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EventValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EventValidationError{}

// Validate checks the field values on EventSeries with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *EventSeries) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on EventSeries with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in EventSeriesMultiError, or
// nil if none found.
func (m *EventSeries) ValidateAll() error {
	return m.validate(true)
}

func (m *EventSeries) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Count

	if all {
		switch v := interface{}(m.GetLastObservedTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, EventSeriesValidationError{
					field:  "LastObservedTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, EventSeriesValidationError{
					field:  "LastObservedTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetLastObservedTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return EventSeriesValidationError{
				field:  "LastObservedTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return EventSeriesMultiError(errors)
	}

	return nil
}

// EventSeriesMultiError is an error wrapping multiple validation errors
// returned by EventSeries.ValidateAll() if the designated constraints aren't met.
type EventSeriesMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EventSeriesMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EventSeriesMultiError) AllErrors() []error { return m }

// EventSeriesValidationError is the validation error returned by
// EventSeries.Validate if the designated constraints aren't met.
type EventSeriesValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EventSeriesValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EventSeriesValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EventSeriesValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EventSeriesValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EventSeriesValidationError) ErrorName() string { return "EventSeriesValidationError" }

// Error satisfies the builtin error interface
func (e EventSeriesValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEventSeries.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EventSeriesValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EventSeriesValidationError{}

// Validate checks the field values on ObjectReference with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ObjectReference) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ObjectReference with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ObjectReferenceMultiError, or nil if none found.
func (m *ObjectReference) ValidateAll() error {
	return m.validate(true)
}

func (m *ObjectReference) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ApiVersion

	// no validation rules for Kind

	// no validation rules for Name

	// no validation rules for Namespace

	if len(errors) > 0 {
		return ObjectReferenceMultiError(errors)
	}

	return nil
}

// ObjectReferenceMultiError is an error wrapping multiple validation errors
// returned by ObjectReference.ValidateAll() if the designated constraints
// aren't met.
type ObjectReferenceMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ObjectReferenceMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ObjectReferenceMultiError) AllErrors() []error { return m }

// ObjectReferenceValidationError is the validation error returned by
// ObjectReference.Validate if the designated constraints aren't met.
type ObjectReferenceValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ObjectReferenceValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ObjectReferenceValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ObjectReferenceValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ObjectReferenceValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ObjectReferenceValidationError) ErrorName() string { return "ObjectReferenceValidationError" }

// Error satisfies the builtin error interface
func (e ObjectReferenceValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sObjectReference.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ObjectReferenceValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ObjectReferenceValidationError{}
