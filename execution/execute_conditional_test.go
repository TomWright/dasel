package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestConditional(t *testing.T) {
	t.Run("true", testCase{
		s:   `if (true) { "yes" } else { "no" }`,
		out: model.NewStringValue("yes"),
	}.run)
	t.Run("false", testCase{
		s:   `if (false) { "yes" } else { "no" }`,
		out: model.NewStringValue("no"),
	}.run)
	t.Run("nested", testCase{
		s: `
				if (true) {
					if (true) { "yes" }
					else { "no" }
				} else { "no" }`,
		out: model.NewStringValue("yes"),
	}.run)
	t.Run("nested false", testCase{
		s: `
				if (true) {
					if (false) { "yes" }
					else { "no" }
				} else { "no" }`,
		out: model.NewStringValue("no"),
	}.run)
	t.Run("else if", testCase{
		s: `
				if (false) { "yes" }
				elseif (true) { "no" }
				else { "maybe" }`,
		out: model.NewStringValue("no"),
	}.run)
	t.Run("else if else", testCase{
		s: `
				if (false) { "yes" }
				elseif (false) { "no" }
				else { "maybe" }`,
		out: model.NewStringValue("maybe"),
	}.run)
	t.Run("if elseif elseif else", testCase{
		s: `
				if (false) { "yes" }
				elseif (false) { "no" }
				elseif (false) { "maybe" }
				else { "nope" }`,
		out: model.NewStringValue("nope"),
	}.run)
}
