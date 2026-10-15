// Copyright 2025 The JSON Schema Go Project Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package jsonschema

import (
	"fmt"
	"strings"
)

// ValidationErrorKind identifies the schema keyword (or meta-error) that caused
// a validation failure.
type ValidationErrorKind string

const (
	// Meta / structural kinds
	KindSchemaVersion ValidationErrorKind = "schemaVersion"
	KindSchema        ValidationErrorKind = "schema"      // wraps a child error with schema-path context
	KindInvalidType   ValidationErrorKind = "invalidType" // Go value is not a valid JSON value at all

	// Generic applicator kinds
	KindRef        ValidationErrorKind = "$ref"
	KindDynamicRef ValidationErrorKind = "$dynamicRef"
	KindAllOf      ValidationErrorKind = "allOf"
	KindAnyOf      ValidationErrorKind = "anyOf"
	KindOneOf      ValidationErrorKind = "oneOf"
	KindNot        ValidationErrorKind = "not"
	KindIf         ValidationErrorKind = "if"

	// Scalar validation kinds
	KindType             ValidationErrorKind = "type"
	KindEnum             ValidationErrorKind = "enum"
	KindConst            ValidationErrorKind = "const"
	KindMultipleOf       ValidationErrorKind = "multipleOf"
	KindMinimum          ValidationErrorKind = "minimum"
	KindMaximum          ValidationErrorKind = "maximum"
	KindExclusiveMinimum ValidationErrorKind = "exclusiveMinimum"
	KindExclusiveMaximum ValidationErrorKind = "exclusiveMaximum"
	KindMinLength        ValidationErrorKind = "minLength"
	KindMaxLength        ValidationErrorKind = "maxLength"
	KindPattern          ValidationErrorKind = "pattern"

	// Array validation kinds
	KindMinItems    ValidationErrorKind = "minItems"
	KindMaxItems    ValidationErrorKind = "maxItems"
	KindUniqueItems ValidationErrorKind = "uniqueItems"
	KindContains    ValidationErrorKind = "contains"
	KindMinContains ValidationErrorKind = "minContains"
	KindMaxContains ValidationErrorKind = "maxContains"

	// Object validation kinds
	KindRequired              ValidationErrorKind = "required"
	KindMinProperties         ValidationErrorKind = "minProperties"
	KindMaxProperties         ValidationErrorKind = "maxProperties"
	KindAdditionalProperties  ValidationErrorKind = "additionalProperties"
	KindPropertyNames         ValidationErrorKind = "propertyNames"
	KindDependentRequired     ValidationErrorKind = "dependentRequired"
	KindDependentSchemas      ValidationErrorKind = "dependentSchemas"
	KindUnevaluatedItems      ValidationErrorKind = "unevaluatedItems"
	KindUnevaluatedProperties ValidationErrorKind = "unevaluatedProperties"

	// Unsupported instance types
	KindStructInstance ValidationErrorKind = "structInstance"
	KindMapKeyType     ValidationErrorKind = "mapKeyType"
)

// ValidationError is a structured, typed error returned by [Resolved.Validate].
//
// Each error carries a [ValidationErrorKind] identifying the schema keyword that
// failed, a human-readable Message, and an optional Causes slice for nested or
// aggregate failures (e.g. each branch of anyOf, or a child property that failed).
type ValidationError struct {
	// Kind is the schema keyword (or meta category) that caused this error.
	Kind ValidationErrorKind

	// Message is a human-readable explanation of the failure.
	Message string

	// SchemaPath is the path to the schema node that rejected the instance
	// (e.g. "properties/address/minLength").
	SchemaPath string

	// Causes holds nested errors.  For allOf/anyOf/oneOf/if each branch's
	// error is stored here; for wrapper errors (KindSchema) the single wrapped
	// error is stored as Causes[0].
	Causes []*ValidationError
}

// Ensure ValidationError satisfies the error interface at compile time.
var _ error = (*ValidationError)(nil)

// Error implements the error interface.
func (e *ValidationError) Error() string {
	var b strings.Builder
	e.writeError(&b, 0)
	return b.String()
}

func (e *ValidationError) writeError(b *strings.Builder, depth int) {
	indent := strings.Repeat("  ", depth)
	if e.SchemaPath != "" {
		fmt.Fprintf(b, "%s%s: %s", indent, e.SchemaPath, e.Message)
	} else {
		fmt.Fprintf(b, "%s%s", indent, e.Message)
	}
	for _, c := range e.Causes {
		b.WriteByte('\n')
		c.writeError(b, depth+1)
	}
}

// All returns a flat slice of every ValidationError in the tree rooted at e,
// including e itself, in pre-order (parent before its children).
// This is useful when callers want to iterate over all failures without
// writing their own recursive traversal.
func (e *ValidationError) All() []ValidationError {
	var out []ValidationError
	e.collect(&out)
	return out
}

func (e *ValidationError) collect(out *[]ValidationError) {
	if e == nil {
		return
	}
	// KindSchema nodes are structural wrappers that provide schema-path
	// context; they are not independent failures, so skip appending them
	// while still recursing into their children.
	if e.Kind != KindSchema {
		*out = append(*out, *e)
	}
	for i := range e.Causes {
		e.Causes[i].collect(out)
	}
}

// Unwrap returns the first cause, enabling errors.Is / errors.As traversal.
func (e *ValidationError) Unwrap() error {
	if len(e.Causes) == 0 {
		return nil
	}
	return e.Causes[0]
}

// validationErrorf creates a leaf ValidationError.
func validationErrorf(kind ValidationErrorKind, format string, args ...any) *ValidationError {
	return &ValidationError{
		Kind:    kind,
		Message: fmt.Sprintf(format, args...),
	}
}

// wrapValidationError stamps schemaPath onto ve if it doesn't already have one,
// or wraps it in a KindSchema parent if it does.
func wrapValidationError(ve *ValidationError, schemaPath string) *ValidationError {
	if ve == nil {
		return nil
	}
	if ve.SchemaPath == "" {
		ve.SchemaPath = schemaPath
		return ve
	}
	// Already has a path, wrap to preserve both
	label := schemaPath
	if label == "" {
		label = "root"
	}
	return &ValidationError{
		Kind:       KindSchema,
		Message:    fmt.Sprintf("validating %s", label),
		SchemaPath: schemaPath,
		Causes:     []*ValidationError{ve},
	}
}
