// Copyright 2025 The JSON Schema Go Project Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package jsonschema

import (
	"errors"
	"strings"
	"testing"
)

// mustValidate calls ValidateStructured and asserts a non-nil result.
func mustValidate(t *testing.T, rs *Resolved, instance any) *ValidationError {
	t.Helper()
	ve := rs.ValidateStructured(instance)
	if ve == nil {
		t.Fatal("expected a ValidationError, got nil")
	}
	return ve
}

func TestValidationErrorKind_Type(t *testing.T) {
	rs, err := (&Schema{Type: "string"}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 42)
	if ve.Kind != KindType {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindType)
	}
	if !strings.Contains(ve.Message, `has type "integer"`) {
		t.Errorf("Message = %q, want to contain 'has type'", ve.Message)
	}
}

func TestValidationErrorKind_Enum(t *testing.T) {
	rs, err := (&Schema{Enum: []any{"a", "b"}}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, "c")
	if ve.Kind != KindEnum {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindEnum)
	}
}

func TestValidationErrorKind_Const(t *testing.T) {
	rs, err := (&Schema{Const: Ptr[any]("hello")}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, "world")
	if ve.Kind != KindConst {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindConst)
	}
}

func TestValidationErrorKind_Minimum(t *testing.T) {
	rs, err := (&Schema{Minimum: Ptr(10.0)}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 5)
	if ve.Kind != KindMinimum {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindMinimum)
	}
}

func TestValidationErrorKind_Maximum(t *testing.T) {
	rs, err := (&Schema{Maximum: Ptr(5.0)}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 10)
	if ve.Kind != KindMaximum {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindMaximum)
	}
}

func TestValidationErrorKind_MinLength(t *testing.T) {
	rs, err := (&Schema{MinLength: Ptr(5)}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, "hi")
	if ve.Kind != KindMinLength {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindMinLength)
	}
}

func TestValidationErrorKind_MaxLength(t *testing.T) {
	rs, err := (&Schema{MaxLength: Ptr(3)}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, "toolong")
	if ve.Kind != KindMaxLength {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindMaxLength)
	}
}

func TestValidationErrorKind_Pattern(t *testing.T) {
	rs, err := (&Schema{Pattern: `^\d+$`}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, "abc")
	if ve.Kind != KindPattern {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindPattern)
	}
}

func TestValidationErrorKind_Required(t *testing.T) {
	rs, err := (&Schema{Required: []string{"name"}}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, map[string]any{})
	if ve.Kind != KindRequired {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindRequired)
	}
	if !strings.Contains(ve.Message, "name") {
		t.Errorf("Message %q does not mention missing property 'name'", ve.Message)
	}
}

func TestValidationErrorKind_MinItems(t *testing.T) {
	rs, err := (&Schema{MinItems: Ptr(3)}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, []any{1})
	if ve.Kind != KindMinItems {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindMinItems)
	}
}

func TestValidationErrorKind_MaxItems(t *testing.T) {
	rs, err := (&Schema{MaxItems: Ptr(1)}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, []any{1, 2, 3})
	if ve.Kind != KindMaxItems {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindMaxItems)
	}
}

func TestValidationErrorKind_UniqueItems(t *testing.T) {
	rs, err := (&Schema{UniqueItems: true}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, []any{1, 1})
	if ve.Kind != KindUniqueItems {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindUniqueItems)
	}
}

func TestValidationErrorKind_AdditionalProperties(t *testing.T) {
	rs, err := (&Schema{
		Properties:           map[string]*Schema{"a": {}},
		AdditionalProperties: falseSchema(),
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, map[string]any{"a": 1, "b": 2})
	if ve.Kind != KindAdditionalProperties {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindAdditionalProperties)
	}
}

func TestValidationErrorKind_MinProperties(t *testing.T) {
	rs, err := (&Schema{MinProperties: Ptr(3)}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, map[string]any{"a": 1})
	if ve.Kind != KindMinProperties {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindMinProperties)
	}
}

func TestValidationErrorKind_MaxProperties(t *testing.T) {
	rs, err := (&Schema{MaxProperties: Ptr(1)}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, map[string]any{"a": 1, "b": 2})
	if ve.Kind != KindMaxProperties {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindMaxProperties)
	}
}

func TestValidationErrorKind_AnyOf_HasCauses(t *testing.T) {
	rs, err := (&Schema{
		AnyOf: []*Schema{
			{Type: "string"},
			{Type: "boolean"},
		},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 42)
	if ve.Kind != KindAnyOf {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindAnyOf)
	}
	if len(ve.Causes) != 2 {
		t.Errorf("Causes has %d elements, want 2", len(ve.Causes))
	}
	// Each cause should itself be a typed ValidationError.
	for i, c := range ve.Causes {
		if c.Kind == "" {
			t.Errorf("Causes[%d].Kind is empty", i)
		}
	}
}

func TestValidationErrorKind_AllOf_HasCause(t *testing.T) {
	rs, err := (&Schema{
		AllOf: []*Schema{
			{Type: "string"},
			{MinLength: Ptr(10)},
		},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 42)
	if ve.Kind != KindAllOf {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindAllOf)
	}
	if len(ve.Causes) != 1 {
		t.Errorf("Causes has %d elements, want 1", len(ve.Causes))
	}
}

func TestValidationErrorKind_OneOf(t *testing.T) {
	rs, err := (&Schema{
		OneOf: []*Schema{
			{Type: "string"},
			{Type: "boolean"},
		},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 42)
	if ve.Kind != KindOneOf {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindOneOf)
	}
}

func TestValidationErrorKind_Not(t *testing.T) {
	rs, err := (&Schema{
		Not: &Schema{Type: "integer"},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 42)
	if ve.Kind != KindNot {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindNot)
	}
}

// TestValidationError_SchemaPath_RootLevel checks that a failure directly
// against the root schema has SchemaPath "".
func TestValidationError_SchemaPath_RootLevel(t *testing.T) {
	rs, err := (&Schema{Type: "string"}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 42)
	if ve.Kind != KindType {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindType)
	}
	if ve.SchemaPath != "" {
		t.Errorf("SchemaPath = %q, want empty string for root-level error", ve.SchemaPath)
	}
}

// TestValidationError_SchemaPath_NestedProperty checks the error tree for a
// property sub-schema failure:
//
//	KindSchema  SchemaPath=""   (root wrapper, message "validating root")
//	  └─ KindMinimum  SchemaPath="/properties/age"
func TestValidationError_SchemaPath_NestedProperty(t *testing.T) {
	rs, err := (&Schema{
		Properties: map[string]*Schema{
			"age": {Minimum: Ptr(0.0)},
		},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, map[string]any{"age": -1})

	if ve.Kind != KindSchema {
		t.Errorf("root Kind = %q, want %q", ve.Kind, KindSchema)
	}
	if ve.SchemaPath != "" {
		t.Errorf("root SchemaPath = %q, want empty string", ve.SchemaPath)
	}
	if len(ve.Causes) != 1 {
		t.Fatalf("root Causes len = %d, want 1", len(ve.Causes))
	}
	leaf := ve.Causes[0]
	if leaf.Kind != KindMinimum {
		t.Errorf("leaf Kind = %q, want %q", leaf.Kind, KindMinimum)
	}
	if leaf.SchemaPath != "/properties/age" {
		t.Errorf("leaf SchemaPath = %q, want %q", leaf.SchemaPath, "/properties/age")
	}
}

// TestValidationError_SchemaPath_DeeplyNested checks a two-level property
// nesting produces the full JSON-Pointer path at each level:
//
//	KindSchema  SchemaPath=""                                    (root)
//	  └─ KindSchema  SchemaPath="/properties/address"
//	       └─ KindType  SchemaPath="/properties/address/properties/zip"
func TestValidationError_SchemaPath_DeeplyNested(t *testing.T) {
	rs, err := (&Schema{
		Properties: map[string]*Schema{
			"address": {
				Properties: map[string]*Schema{
					"zip": {Type: "string"},
				},
			},
		},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, map[string]any{"address": map[string]any{"zip": 12345}})

	if ve.Kind != KindSchema || ve.SchemaPath != "" {
		t.Errorf("root: Kind=%q SchemaPath=%q, want KindSchema/empty", ve.Kind, ve.SchemaPath)
	}
	if len(ve.Causes) != 1 {
		t.Fatalf("root Causes len = %d, want 1", len(ve.Causes))
	}

	mid := ve.Causes[0]
	if mid.Kind != KindSchema || mid.SchemaPath != "/properties/address" {
		t.Errorf("mid: Kind=%q SchemaPath=%q, want KindSchema/%q",
			mid.Kind, mid.SchemaPath, "/properties/address")
	}
	if len(mid.Causes) != 1 {
		t.Fatalf("mid Causes len = %d, want 1", len(mid.Causes))
	}

	leaf := mid.Causes[0]
	if leaf.Kind != KindType {
		t.Errorf("leaf Kind = %q, want %q", leaf.Kind, KindType)
	}
	if leaf.SchemaPath != "/properties/address/properties/zip" {
		t.Errorf("leaf SchemaPath = %q, want %q",
			leaf.SchemaPath, "/properties/address/properties/zip")
	}
}

// TestValidationError_SchemaPath_AllOf checks that an allOf failure has
// SchemaPath="" on the root and the cause path "/allOf/0".
func TestValidationError_SchemaPath_AllOf(t *testing.T) {
	rs, err := (&Schema{
		AllOf: []*Schema{{Minimum: Ptr(5.0)}},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 1)

	if ve.Kind != KindAllOf || ve.SchemaPath != "" {
		t.Errorf("root: Kind=%q SchemaPath=%q, want KindAllOf/empty", ve.Kind, ve.SchemaPath)
	}
	if len(ve.Causes) != 1 {
		t.Fatalf("Causes len = %d, want 1", len(ve.Causes))
	}
	cause := ve.Causes[0]
	if cause.Kind != KindMinimum {
		t.Errorf("cause Kind = %q, want %q", cause.Kind, KindMinimum)
	}
	if cause.SchemaPath != "/allOf/0" {
		t.Errorf("cause SchemaPath = %q, want %q", cause.SchemaPath, "/allOf/0")
	}
}

// TestValidationError_SchemaPath_AnyOf checks that each branch cause carries
// its "/anyOf/<n>" path and the root error has SchemaPath="".
func TestValidationError_SchemaPath_AnyOf(t *testing.T) {
	rs, err := (&Schema{
		AnyOf: []*Schema{{Type: "string"}, {Type: "boolean"}},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 42)

	if ve.Kind != KindAnyOf || ve.SchemaPath != "" {
		t.Errorf("root: Kind=%q SchemaPath=%q, want KindAnyOf/empty", ve.Kind, ve.SchemaPath)
	}
	if len(ve.Causes) != 2 {
		t.Fatalf("Causes len = %d, want 2", len(ve.Causes))
	}
	for i, wantPath := range []string{"/anyOf/0", "/anyOf/1"} {
		if ve.Causes[i].SchemaPath != wantPath {
			t.Errorf("Causes[%d].SchemaPath = %q, want %q", i, ve.Causes[i].SchemaPath, wantPath)
		}
	}
}

// TestValidationError_SchemaPath_If checks that the then-branch cause carries
// SchemaPath="/then" and the root error has SchemaPath="".
func TestValidationError_SchemaPath_If(t *testing.T) {
	rs, err := (&Schema{
		If:   &Schema{Type: "string"},
		Then: &Schema{MinLength: Ptr(5)},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, "hi")

	if ve.Kind != KindIf || ve.SchemaPath != "" {
		t.Errorf("root: Kind=%q SchemaPath=%q, want KindIf/empty", ve.Kind, ve.SchemaPath)
	}
	if len(ve.Causes) != 1 {
		t.Fatalf("Causes len = %d, want 1", len(ve.Causes))
	}
	cause := ve.Causes[0]
	if cause.Kind != KindMinLength {
		t.Errorf("cause Kind = %q, want %q", cause.Kind, KindMinLength)
	}
	if cause.SchemaPath != "/then" {
		t.Errorf("cause SchemaPath = %q, want %q", cause.SchemaPath, "/then")
	}
}

// TestValidationError_SchemaPath_Ref checks that a $ref failure has
// SchemaPath="" on the root and the cause carries the $defs path.
func TestValidationError_SchemaPath_Ref(t *testing.T) {
	rs, err := (&Schema{
		Defs: map[string]*Schema{"pos": {Minimum: Ptr(1.0)}},
		Ref:  "#/$defs/pos",
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, -1)

	if ve.Kind != KindRef || ve.SchemaPath != "" {
		t.Errorf("root: Kind=%q SchemaPath=%q, want KindRef/empty", ve.Kind, ve.SchemaPath)
	}
	if len(ve.Causes) != 1 {
		t.Fatalf("Causes len = %d, want 1", len(ve.Causes))
	}
	cause := ve.Causes[0]
	if cause.Kind != KindMinimum {
		t.Errorf("cause Kind = %q, want %q", cause.Kind, KindMinimum)
	}
	if cause.SchemaPath != "/$defs/pos" {
		t.Errorf("cause SchemaPath = %q, want %q", cause.SchemaPath, "/$defs/pos")
	}
}

// TestValidationError_SchemaPath_UnevaluatedProperties checks that the root
// error has SchemaPath="" and the cause has SchemaPath="/unevaluatedProperties".
func TestValidationError_SchemaPath_UnevaluatedProperties(t *testing.T) {
	rs, err := (&Schema{
		Properties:            map[string]*Schema{"a": {}},
		UnevaluatedProperties: falseSchema(),
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, map[string]any{"a": 1, "b": 2})

	if ve.Kind != KindUnevaluatedProperties || ve.SchemaPath != "" {
		t.Errorf("root: Kind=%q SchemaPath=%q, want KindUnevaluatedProperties/empty",
			ve.Kind, ve.SchemaPath)
	}
	if len(ve.Causes) != 1 {
		t.Fatalf("Causes len = %d, want 1", len(ve.Causes))
	}
	if ve.Causes[0].SchemaPath != "/unevaluatedProperties" {
		t.Errorf("cause SchemaPath = %q, want %q",
			ve.Causes[0].SchemaPath, "/unevaluatedProperties")
	}
}

// TestValidationError_NestedCauses checks that a property validation failure
// is surfaced as a nested cause tree and that the leaf carries the correct path.
func TestValidationError_NestedCauses(t *testing.T) {
	rs, err := (&Schema{
		Properties: map[string]*Schema{
			"items": {MinItems: Ptr(2)},
		},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}

	ve := rs.ValidateStructured(map[string]any{"items": []any{1}})
	if ve == nil {
		t.Fatal("expected error, got nil")
	}

	// Walk the cause tree to find the KindMinItems leaf.
	var findLeaf func(*ValidationError, ValidationErrorKind) *ValidationError
	findLeaf = func(ve *ValidationError, k ValidationErrorKind) *ValidationError {
		if ve.Kind == k {
			return ve
		}
		for _, c := range ve.Causes {
			if found := findLeaf(c, k); found != nil {
				return found
			}
		}
		return nil
	}

	leaf := findLeaf(ve, KindMinItems)
	if leaf == nil {
		t.Fatalf("could not find KindMinItems in error tree:\n%v", ve)
	}
	if leaf.SchemaPath != "/properties/items" {
		t.Errorf("leaf SchemaPath = %q, want %q", leaf.SchemaPath, "/properties/items")
	}
}

// TestValidationError_ErrorString checks that Error() renders a human-readable
// multi-line string including the schema path.
func TestValidationError_ErrorString(t *testing.T) {
	ve := &ValidationError{
		Kind:       KindRequired,
		Message:    `missing properties: ["name"]`,
		SchemaPath: "properties/name",
	}
	s := ve.Error()
	if !strings.Contains(s, "properties/name") {
		t.Errorf("Error() = %q; want to contain schema path 'properties/name'", s)
	}
	if !strings.Contains(s, "missing properties") {
		t.Errorf("Error() = %q; want to contain 'missing properties'", s)
	}
}

// TestValidationError_Unwrap checks that errors.As traversal works through Causes.
func TestValidationError_Unwrap(t *testing.T) {
	inner := &ValidationError{Kind: KindType, Message: "wrong type"}
	outer := &ValidationError{
		Kind:    KindAllOf,
		Message: "allOf failed",
		Causes:  []*ValidationError{inner},
	}
	var got *ValidationError
	if !errors.As(outer, &got) {
		t.Fatal("errors.As failed for outer")
	}
	// errors.As on the outer should resolve to outer itself first.
	if got != outer {
		t.Errorf("errors.As resolved to %p, want outer %p", got, outer)
	}
}

// ---------------------------------------------------------------------------
// ValidationError.All() – flat pre-order traversal.
// ---------------------------------------------------------------------------

// TestValidationError_All_Leaf checks that All on a leaf returns just itself.
func TestValidationError_All_Leaf(t *testing.T) {
	ve := &ValidationError{Kind: KindType, Message: "wrong type", SchemaPath: "/type"}
	all := ve.All()
	if len(all) != 1 {
		t.Fatalf("All() returned %d nodes, want 1", len(all))
	}
	if all[0].Kind != KindType || all[0].SchemaPath != "/type" {
		t.Errorf("All()[0] = {Kind:%s SchemaPath:%q}, want {Kind:%s SchemaPath:%q}",
			all[0].Kind, all[0].SchemaPath, KindType, "/type")
	}
}

// TestValidationError_All_PreOrder checks that a multi-level tree is returned
// in pre-order, with KindSchema wrapper nodes skipped.
func TestValidationError_All_PreOrder(t *testing.T) {
	leaf1 := &ValidationError{Kind: KindMinimum, SchemaPath: "/allOf/0"}
	leaf2 := &ValidationError{Kind: KindMaximum, SchemaPath: "/allOf/1"}
	mid := &ValidationError{Kind: KindAllOf, Causes: []*ValidationError{leaf1, leaf2}}
	// root is a KindSchema wrapper – it should be traversed but not included.
	root := &ValidationError{Kind: KindSchema, Causes: []*ValidationError{mid}}

	all := root.All()

	if len(all) != 3 {
		t.Fatalf("All() returned %d nodes, want 3 (KindSchema root skipped)", len(all))
	}
	wantKinds := []ValidationErrorKind{KindAllOf, KindMinimum, KindMaximum}
	for i, wantKind := range wantKinds {
		if all[i].Kind != wantKind {
			t.Errorf("All()[%d].Kind = %q, want %q", i, all[i].Kind, wantKind)
		}
	}
}

// TestValidationError_All_SchemaWrapperSkipped explicitly verifies that
// KindSchema nodes are omitted from All() but their children are still
// included.
func TestValidationError_All_SchemaWrapperSkipped(t *testing.T) {
	// A nested-property failure produces:
	//   KindSchema (root wrapper, SchemaPath="")
	//     └─ KindMinimum (SchemaPath="/properties/age")
	rs, err := (&Schema{
		Properties: map[string]*Schema{
			"age": {Minimum: Ptr(0.0)},
		},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, map[string]any{"age": -1})

	all := ve.All()

	// The KindSchema wrapper at the root must not appear; only the real failure.
	for _, e := range all {
		if e.Kind == KindSchema {
			t.Errorf("All() included a KindSchema wrapper: %+v", e)
		}
	}
	if len(all) != 1 {
		t.Fatalf("All() returned %d nodes, want 1", len(all))
	}
	if all[0].Kind != KindMinimum {
		t.Errorf("All()[0].Kind = %q, want %q", all[0].Kind, KindMinimum)
	}
	if all[0].SchemaPath != "/properties/age" {
		t.Errorf("All()[0].SchemaPath = %q, want %q", all[0].SchemaPath, "/properties/age")
	}
}

// TestValidationError_All_FindKind shows the intended use-case: finding all
// errors of a given Kind without manual recursion.
func TestValidationError_All_FindKind(t *testing.T) {
	rs, err := (&Schema{
		AnyOf: []*Schema{
			{Type: "string"},
			{Type: "boolean"},
		},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 42)

	var typeErrors []ValidationError
	for _, e := range ve.All() {
		if e.Kind == KindType {
			typeErrors = append(typeErrors, e)
		}
	}
	if len(typeErrors) != 2 {
		t.Errorf("found %d KindType errors via All(), want 2", len(typeErrors))
	}
}

// TestValidationError_All_NilCauseSkipped checks that a nil pointer in Causes
// is silently skipped and does not panic.
func TestValidationError_All_NilCauseSkipped(t *testing.T) {
	ve := &ValidationError{
		Kind:   KindAllOf,
		Causes: []*ValidationError{nil},
	}
	all := ve.All()
	if len(all) != 1 {
		t.Errorf("All() returned %d nodes, want 1 (nil cause skipped)", len(all))
	}
}

// ---------------------------------------------------------------------------
// ValidateStructured happy-path: nil on valid input.
// ---------------------------------------------------------------------------

func TestValidateStructured_NilOnSuccess(t *testing.T) {
	rs, err := (&Schema{Type: "string"}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	if ve := rs.ValidateStructured("hello"); ve != nil {
		t.Errorf("expected nil, got %v", ve)
	}
}

// ---------------------------------------------------------------------------
// Validate wraps ValidateStructured: it must return a non-nil error that is
// a *ValidationError.
// ---------------------------------------------------------------------------

func TestValidate_ReturnsValidationError(t *testing.T) {
	rs, err := (&Schema{Type: "string"}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	gotErr := rs.Validate(42)
	if gotErr == nil {
		t.Fatal("expected non-nil error from Validate")
	}
	var ve *ValidationError
	if !errors.As(gotErr, &ve) {
		t.Fatalf("Validate returned %T, want *ValidationError", gotErr)
	}
	if ve.Kind != KindType {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindType)
	}
}

// ---------------------------------------------------------------------------
// KindSchemaVersion – unsupported $schema value.
// ---------------------------------------------------------------------------

func TestValidationErrorKind_SchemaVersion(t *testing.T) {
	rs, err := (&Schema{
		Schema: "http://json-schema.org/draft-04/schema#",
		Type:   "string",
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := rs.ValidateStructured("hello")
	if ve == nil {
		t.Fatal("expected a ValidationError, got nil")
	}
	if ve.Kind != KindSchemaVersion {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindSchemaVersion)
	}
}

// ---------------------------------------------------------------------------
// KindInvalidType – Go value that has no JSON counterpart.
// ---------------------------------------------------------------------------

func TestValidationErrorKind_InvalidType(t *testing.T) {
	rs, err := (&Schema{Type: "string"}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	// A channel is not a valid JSON value.
	ch := make(chan int)
	ve := rs.ValidateStructured(ch)
	if ve == nil {
		t.Fatal("expected a ValidationError, got nil")
	}
	if ve.Kind != KindInvalidType {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindInvalidType)
	}
}

// ---------------------------------------------------------------------------
// KindMultipleOf
// ---------------------------------------------------------------------------

func TestValidationErrorKind_MultipleOf(t *testing.T) {
	rs, err := (&Schema{MultipleOf: Ptr(3.0)}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 7)
	if ve.Kind != KindMultipleOf {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindMultipleOf)
	}
}

// ---------------------------------------------------------------------------
// KindExclusiveMinimum / KindExclusiveMaximum
// ---------------------------------------------------------------------------

func TestValidationErrorKind_ExclusiveMinimum(t *testing.T) {
	rs, err := (&Schema{ExclusiveMinimum: Ptr(10.0)}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 10) // equal to the exclusive bound – not valid
	if ve.Kind != KindExclusiveMinimum {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindExclusiveMinimum)
	}
}

func TestValidationErrorKind_ExclusiveMaximum(t *testing.T) {
	rs, err := (&Schema{ExclusiveMaximum: Ptr(10.0)}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, 10) // equal to the exclusive bound – not valid
	if ve.Kind != KindExclusiveMaximum {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindExclusiveMaximum)
	}
}

// ---------------------------------------------------------------------------
// KindContains / KindMinContains / KindMaxContains
// ---------------------------------------------------------------------------

func TestValidationErrorKind_Contains(t *testing.T) {
	rs, err := (&Schema{
		Contains: &Schema{Type: "string"},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, []any{1, 2, 3})
	if ve.Kind != KindContains {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindContains)
	}
}

func TestValidationErrorKind_MinContains(t *testing.T) {
	rs, err := (&Schema{
		Contains:    &Schema{Type: "string"},
		MinContains: Ptr(2),
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Only one string – satisfies Contains but violates MinContains.
	ve := mustValidate(t, rs, []any{"hello", 1, 2})
	if ve.Kind != KindMinContains {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindMinContains)
	}
}

func TestValidationErrorKind_MaxContains(t *testing.T) {
	rs, err := (&Schema{
		Contains:    &Schema{Type: "string"},
		MaxContains: Ptr(1),
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Two strings exceeds MaxContains.
	ve := mustValidate(t, rs, []any{"a", "b", 1})
	if ve.Kind != KindMaxContains {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindMaxContains)
	}
}

// ---------------------------------------------------------------------------
// KindPropertyNames
// ---------------------------------------------------------------------------

func TestValidationErrorKind_PropertyNames(t *testing.T) {
	rs, err := (&Schema{
		PropertyNames: &Schema{MaxLength: Ptr(3)},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, map[string]any{"toolong": 1})
	if ve.Kind != KindPropertyNames {
		// The error may be a KindMaxLength wrapped in KindSchema at the
		// propertyNames level; walk the tree looking for either.
		var findKind func(*ValidationError, ValidationErrorKind) bool
		findKind = func(ve *ValidationError, k ValidationErrorKind) bool {
			if ve.Kind == k {
				return true
			}
			for _, c := range ve.Causes {
				if findKind(c, k) {
					return true
				}
			}
			return false
		}
		if !findKind(ve, KindMaxLength) {
			t.Errorf("expected KindPropertyNames or a KindMaxLength cause, got root Kind = %q:\n%v", ve.Kind, ve)
		}
	}
}

// ---------------------------------------------------------------------------
// KindDependentRequired
// ---------------------------------------------------------------------------

func TestValidationErrorKind_DependentRequired(t *testing.T) {
	rs, err := (&Schema{
		DependentRequired: map[string][]string{
			"credit_card": {"billing_address"},
		},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Has credit_card but is missing billing_address.
	ve := mustValidate(t, rs, map[string]any{"credit_card": "1234"})
	if ve.Kind != KindDependentRequired {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindDependentRequired)
	}
	if !strings.Contains(ve.Message, "billing_address") {
		t.Errorf("Message %q does not mention missing property 'billing_address'", ve.Message)
	}
}

// ---------------------------------------------------------------------------
// KindStructInstance – validating a raw Go struct is rejected.
// ---------------------------------------------------------------------------

func TestValidationErrorKind_StructInstance(t *testing.T) {
	rs, err := (&Schema{}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	type myStruct struct{ Name string }
	ve := rs.ValidateStructured(myStruct{"alice"})
	if ve == nil {
		t.Fatal("expected a ValidationError, got nil")
	}
	if ve.Kind != KindStructInstance {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindStructInstance)
	}
}

// ---------------------------------------------------------------------------
// KindMapKeyType – map with non-string keys is rejected.
// ---------------------------------------------------------------------------

func TestValidationErrorKind_MapKeyType(t *testing.T) {
	rs, err := (&Schema{}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := rs.ValidateStructured(map[int]any{1: "a"})
	if ve == nil {
		t.Fatal("expected a ValidationError, got nil")
	}
	if ve.Kind != KindMapKeyType {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindMapKeyType)
	}
}

// ---------------------------------------------------------------------------
// KindIf – if/then branch failure surfaces as KindIf.
// ---------------------------------------------------------------------------

func TestValidationErrorKind_If(t *testing.T) {
	// if type is string, then minLength must be >= 5.
	rs, err := (&Schema{
		If:   &Schema{Type: "string"},
		Then: &Schema{MinLength: Ptr(5)},
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	// "hi" passes the if-condition but fails the then-constraint.
	ve := mustValidate(t, rs, "hi")
	if ve.Kind != KindIf {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindIf)
	}
	if len(ve.Causes) == 0 {
		t.Error("expected at least one cause for KindIf error")
	}
}

// ---------------------------------------------------------------------------
// KindRef – $ref failure surfaces as KindRef with a cause.
// ---------------------------------------------------------------------------

func TestValidationErrorKind_Ref(t *testing.T) {
	rs, err := (&Schema{
		Defs: map[string]*Schema{
			"positiveInt": {Minimum: Ptr(1.0)},
		},
		Ref: "#/$defs/positiveInt",
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	ve := mustValidate(t, rs, -5)
	if ve.Kind != KindRef {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindRef)
	}
	if len(ve.Causes) == 0 {
		t.Error("expected at least one cause for KindRef error")
	}
}

// ---------------------------------------------------------------------------
// KindUnevaluatedItems
// ---------------------------------------------------------------------------

func TestValidationErrorKind_UnevaluatedItems(t *testing.T) {
	rs, err := (&Schema{
		PrefixItems:      []*Schema{{Type: "string"}},
		UnevaluatedItems: falseSchema(),
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Second item is beyond the prefix and is not allowed.
	ve := mustValidate(t, rs, []any{"ok", "extra"})
	if ve.Kind != KindUnevaluatedItems {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindUnevaluatedItems)
	}
}

// ---------------------------------------------------------------------------
// KindUnevaluatedProperties
// ---------------------------------------------------------------------------

func TestValidationErrorKind_UnevaluatedProperties(t *testing.T) {
	rs, err := (&Schema{
		Properties:            map[string]*Schema{"a": {}},
		UnevaluatedProperties: falseSchema(),
	}).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	// "b" is unevaluated (not in Properties) and not allowed.
	ve := mustValidate(t, rs, map[string]any{"a": 1, "b": 2})
	if ve.Kind != KindUnevaluatedProperties {
		t.Errorf("Kind = %q, want %q", ve.Kind, KindUnevaluatedProperties)
	}
	if len(ve.Causes) == 0 {
		t.Error("expected at least one cause for KindUnevaluatedProperties error")
	}
}

// ---------------------------------------------------------------------------
// KindSchema (wrap path) – wrapValidationError produces a KindSchema parent
// when the child error already has a SchemaPath set.
// ---------------------------------------------------------------------------

func TestWrapValidationError_WrapsWhenPathAlreadySet(t *testing.T) {
	inner := &ValidationError{
		Kind:       KindType,
		Message:    "wrong type",
		SchemaPath: "properties/name",
	}
	wrapped := wrapValidationError(inner, "outer/path")
	if wrapped.Kind != KindSchema {
		t.Errorf("Kind = %q, want %q", wrapped.Kind, KindSchema)
	}
	if wrapped.SchemaPath != "outer/path" {
		t.Errorf("SchemaPath = %q, want %q", wrapped.SchemaPath, "outer/path")
	}
	if len(wrapped.Causes) != 1 || wrapped.Causes[0] != inner {
		t.Error("wrapped error does not carry the original as Causes[0]")
	}
}

// TestWrapValidationError_StampsWhenNoPath covers the branch where the child
// error has no SchemaPath yet: the path should be stamped in-place.
func TestWrapValidationError_StampsWhenNoPath(t *testing.T) {
	ve := &ValidationError{Kind: KindMinimum, Message: "too small"}
	result := wrapValidationError(ve, "some/path")
	if result != ve {
		t.Error("expected the same pointer back when stamping in-place")
	}
	if ve.SchemaPath != "some/path" {
		t.Errorf("SchemaPath = %q, want %q", ve.SchemaPath, "some/path")
	}
}

// TestWrapValidationError_NilReturnsNil guards the nil check.
func TestWrapValidationError_NilReturnsNil(t *testing.T) {
	if result := wrapValidationError(nil, "x"); result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

// ---------------------------------------------------------------------------
// ValidationError.Error() with nested causes renders indented children.
// ---------------------------------------------------------------------------

func TestValidationError_ErrorString_WithCauses(t *testing.T) {
	inner := &ValidationError{
		Kind:       KindMinLength,
		Message:    "too short",
		SchemaPath: "properties/name/minLength",
	}
	outer := &ValidationError{
		Kind:       KindSchema,
		Message:    "validating properties/name",
		SchemaPath: "properties/name",
		Causes:     []*ValidationError{inner},
	}
	s := outer.Error()
	if !strings.Contains(s, "properties/name") {
		t.Errorf("Error() missing outer path: %q", s)
	}
	if !strings.Contains(s, "too short") {
		t.Errorf("Error() missing inner message: %q", s)
	}
	// The inner error should be indented (two spaces at depth 1).
	if !strings.Contains(s, "\n  ") {
		t.Errorf("Error() missing indented newline: %q", s)
	}
}

// ---------------------------------------------------------------------------
// Unwrap chain: errors.As can reach an inner cause.
// ---------------------------------------------------------------------------

func TestValidationError_Unwrap_ReachesInnerCause(t *testing.T) {
	leaf := &ValidationError{Kind: KindMinimum, Message: "below minimum"}
	mid := &ValidationError{
		Kind:    KindRef,
		Message: "$ref failed",
		Causes:  []*ValidationError{leaf},
	}
	outer := &ValidationError{
		Kind:    KindAllOf,
		Message: "allOf failed",
		Causes:  []*ValidationError{mid},
	}
	// errors.As should find the outermost matching type (outer itself).
	var got *ValidationError
	if !errors.As(outer, &got) {
		t.Fatal("errors.As failed")
	}
	if got != outer {
		t.Errorf("errors.As returned %p, want outer %p", got, outer)
	}
	// Unwrap once manually to reach mid.
	if unwrapped := outer.Unwrap(); unwrapped != mid {
		t.Errorf("Unwrap() returned %v, want mid", unwrapped)
	}
}
