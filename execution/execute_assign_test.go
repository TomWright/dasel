package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
)

func TestAssignVariable(t *testing.T) {
	t.Run("single assign", testCase{
		s: `$x=1`,
		outFn: func() *model.Value {
			r := model.NewIntValue(1)
			return r
		},
		opts: []execution.ExecuteOptionFn{
			execution.WithUnstable(),
		},
	}.run)
	t.Run("double assign", testCase{
		s: `$x=1;$y=$x+1`,
		outFn: func() *model.Value {
			r := model.NewIntValue(2)
			return r
		},
		opts: []execution.ExecuteOptionFn{
			execution.WithUnstable(),
		},
	}.run)
	t.Run("multiple assign with final statement", testCase{
		s: `$first = 'Tom';
$last = 'Wright';
$full = $first + ' ' + $last;
{first: $first, last: $last, full: $full}`,
		outFn: func() *model.Value {
			r := model.NewMapValue()
			if err := r.SetMapKey("first", model.NewStringValue("Tom")); err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if err := r.SetMapKey("last", model.NewStringValue("Wright")); err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if err := r.SetMapKey("full", model.NewStringValue("Tom Wright")); err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			return r
		},
		opts: []execution.ExecuteOptionFn{
			execution.WithUnstable(),
		},
	}.run)
	t.Run("multiple assign with final statement and mixed case variables", testCase{
		s: `$firstName = 'Tom';
$lastName = 'Wright';
$fullName = $firstName + ' ' + $lastName;
{firstName: $firstName, lastName: $lastName, fullName: $fullName}`,
		outFn: func() *model.Value {
			r := model.NewMapValue()
			if err := r.SetMapKey("firstName", model.NewStringValue("Tom")); err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if err := r.SetMapKey("lastName", model.NewStringValue("Wright")); err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if err := r.SetMapKey("fullName", model.NewStringValue("Tom Wright")); err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			return r
		},
		opts: []execution.ExecuteOptionFn{
			execution.WithUnstable(),
		},
	}.run)
	t.Run("self referencing variable", testCase{
		s: `$x=1;$x=$x*2`,
		outFn: func() *model.Value {
			r := model.NewIntValue(2)
			return r
		},
		opts: []execution.ExecuteOptionFn{
			execution.WithUnstable(),
		},
	}.run)
}
