// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: candles.proto

package pb

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

// Validate checks the field values on CandlesRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CandlesRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CandlesRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CandlesRequestMultiError,
// or nil if none found.
func (m *CandlesRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CandlesRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Names

	if _, ok := _CandlesRequest_Mode_InLookup[m.GetMode()]; !ok {
		err := CandlesRequestValidationError{
			field:  "Mode",
			reason: "value must be in list [practice competition]",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for UserId

	if len(errors) > 0 {
		return CandlesRequestMultiError(errors)
	}

	return nil
}

// CandlesRequestMultiError is an error wrapping multiple validation errors
// returned by CandlesRequest.ValidateAll() if the designated constraints
// aren't met.
type CandlesRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CandlesRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CandlesRequestMultiError) AllErrors() []error { return m }

// CandlesRequestValidationError is the validation error returned by
// CandlesRequest.Validate if the designated constraints aren't met.
type CandlesRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CandlesRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CandlesRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CandlesRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CandlesRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CandlesRequestValidationError) ErrorName() string { return "CandlesRequestValidationError" }

// Error satisfies the builtin error interface
func (e CandlesRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCandlesRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CandlesRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CandlesRequestValidationError{}

var _CandlesRequest_Mode_InLookup = map[string]struct{}{
	"practice":    {},
	"competition": {},
}

// Validate checks the field values on CandlesResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *CandlesResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CandlesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CandlesResponseMultiError, or nil if none found.
func (m *CandlesResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CandlesResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Name

	if all {
		switch v := interface{}(m.GetOneChart()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CandlesResponseValidationError{
					field:  "OneChart",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CandlesResponseValidationError{
					field:  "OneChart",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetOneChart()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CandlesResponseValidationError{
				field:  "OneChart",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for BtcRatio

	// no validation rules for EntryTime

	// no validation rules for EntryPrice

	// no validation rules for Identifier

	if len(errors) > 0 {
		return CandlesResponseMultiError(errors)
	}

	return nil
}

// CandlesResponseMultiError is an error wrapping multiple validation errors
// returned by CandlesResponse.ValidateAll() if the designated constraints
// aren't met.
type CandlesResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CandlesResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CandlesResponseMultiError) AllErrors() []error { return m }

// CandlesResponseValidationError is the validation error returned by
// CandlesResponse.Validate if the designated constraints aren't met.
type CandlesResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CandlesResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CandlesResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CandlesResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CandlesResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CandlesResponseValidationError) ErrorName() string { return "CandlesResponseValidationError" }

// Error satisfies the builtin error interface
func (e CandlesResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCandlesResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CandlesResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CandlesResponseValidationError{}

// Validate checks the field values on CandleData with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CandleData) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CandleData with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CandleDataMultiError, or
// nil if none found.
func (m *CandleData) ValidateAll() error {
	return m.validate(true)
}

func (m *CandleData) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetPData() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CandleDataValidationError{
						field:  fmt.Sprintf("PData[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CandleDataValidationError{
						field:  fmt.Sprintf("PData[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CandleDataValidationError{
					field:  fmt.Sprintf("PData[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetVData() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CandleDataValidationError{
						field:  fmt.Sprintf("VData[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CandleDataValidationError{
						field:  fmt.Sprintf("VData[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CandleDataValidationError{
					field:  fmt.Sprintf("VData[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return CandleDataMultiError(errors)
	}

	return nil
}

// CandleDataMultiError is an error wrapping multiple validation errors
// returned by CandleData.ValidateAll() if the designated constraints aren't met.
type CandleDataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CandleDataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CandleDataMultiError) AllErrors() []error { return m }

// CandleDataValidationError is the validation error returned by
// CandleData.Validate if the designated constraints aren't met.
type CandleDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CandleDataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CandleDataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CandleDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CandleDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CandleDataValidationError) ErrorName() string { return "CandleDataValidationError" }

// Error satisfies the builtin error interface
func (e CandleDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCandleData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CandleDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CandleDataValidationError{}

// Validate checks the field values on PriceData with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PriceData) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PriceData with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PriceDataMultiError, or nil
// if none found.
func (m *PriceData) ValidateAll() error {
	return m.validate(true)
}

func (m *PriceData) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Open

	// no validation rules for Close

	// no validation rules for High

	// no validation rules for Low

	// no validation rules for Time

	if len(errors) > 0 {
		return PriceDataMultiError(errors)
	}

	return nil
}

// PriceDataMultiError is an error wrapping multiple validation errors returned
// by PriceData.ValidateAll() if the designated constraints aren't met.
type PriceDataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PriceDataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PriceDataMultiError) AllErrors() []error { return m }

// PriceDataValidationError is the validation error returned by
// PriceData.Validate if the designated constraints aren't met.
type PriceDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PriceDataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PriceDataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PriceDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PriceDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PriceDataValidationError) ErrorName() string { return "PriceDataValidationError" }

// Error satisfies the builtin error interface
func (e PriceDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPriceData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PriceDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PriceDataValidationError{}

// Validate checks the field values on VolumeData with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *VolumeData) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on VolumeData with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in VolumeDataMultiError, or
// nil if none found.
func (m *VolumeData) ValidateAll() error {
	return m.validate(true)
}

func (m *VolumeData) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Value

	// no validation rules for Time

	// no validation rules for Color

	if len(errors) > 0 {
		return VolumeDataMultiError(errors)
	}

	return nil
}

// VolumeDataMultiError is an error wrapping multiple validation errors
// returned by VolumeData.ValidateAll() if the designated constraints aren't met.
type VolumeDataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m VolumeDataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m VolumeDataMultiError) AllErrors() []error { return m }

// VolumeDataValidationError is the validation error returned by
// VolumeData.Validate if the designated constraints aren't met.
type VolumeDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e VolumeDataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e VolumeDataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e VolumeDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e VolumeDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e VolumeDataValidationError) ErrorName() string { return "VolumeDataValidationError" }

// Error satisfies the builtin error interface
func (e VolumeDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sVolumeData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = VolumeDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = VolumeDataValidationError{}
