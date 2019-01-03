package lib

import (
	"bytes"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var differ = diffmatchpatch.New()

const twoFuncsTabs = `outline: twoFuncs
	functions:
		difference(a,b int) int
		sum(a,b int) int
			add two things together`

const twoFuncsSpaces = `outline: twoFuncs
  functions:
    difference(a,b int) int
    sum(a,b int) int
			add two things together`

var twoFuncs = &Doc{
	Name: "twoFuncs",
	Functions: []*Function{
		{
			Signature: "difference(a,b int) int",
		},
		{
			Signature:   "sum(a,b int) int",
			Description: "add two things together",
		},
	},
}

const timeSpaces = `here's some leading gak that shouldn't get read into the doc

outline: time
  functions:
    duration(string) duration
      parse a duration
    time(string, format=..., location=...) time
      parse a time
    now() time
      new time instance set to current time
      implementations are able to make this a constant
    zero() time
      a constant

  types:
    duration
      fields:
        hours float
        minutes float
        nanoseconds int
        seconds float
      operators:
        duration - time = duration
        duration + time = time
        duration == duration = boolean
        duration < duration = booleans
    time
      fields:
      operators:
        time == time = boolean
        time < time = boolean`

var time = &Doc{
	Name: "time",
	Functions: []*Function{
		{Signature: "duration(string) duration",
			Description: "parse a duration"},
		{Signature: "time(string, format=..., location=...) time",
			Description: "parse a time"},
		{Signature: "now() time",
			Description: "new time instance set to current time implementations are able to make this a constant"},
		{Signature: "zero() time",
			Description: "a constant"},
	},
	Types: []*Type{
		{Name: "duration",
			Fields: []*Field{
				{Name: "hours float"},
				{Name: "minutes float"},
				{Name: "nanoseconds int"},
				{Name: "seconds float"},
			},
			Operators: []*Operator{
				{Opr: "duration - time = duration"},
				{Opr: "duration + time = time"},
				{Opr: "duration == duration = boolean"},
				{Opr: "duration < duration = booleans"},
			},
		},
		{Name: "time",
			Fields: []*Field{},
			Operators: []*Operator{
				{Opr: "time == time = boolean"},
				{Opr: "time < time = boolean"},
			},
		},
	},
}

const docWithDescriptionTabs = `outline: doc
	this is a document description.
	It's written across two lines
	functions:
		sum(a int, b int) int`

var docWithDescription = &Doc{
	Name:        "doc",
	Description: "this is a document description. It's written across two lines",
	Functions: []*Function{
		{Signature: "sum(a int, b int) int"},
	},
}

const huhSpaces = `  outline: time
  time defines time primitives for starlark
  this is a second line of description
  functions:
    duration(string) duration
			parse a duration
		date() date
			make a date`

var huh = &Doc{
	Name:        "time",
	Description: "time defines time primitives for starlark this is a second line of description",
	Functions: []*Function{
		{Signature: "duration(string) duration",
			Description: "parse a duration"},
		{Signature: "date() date",
			Description: "make a date"},
	},
}

func TestParse(t *testing.T) {
	cases := []struct {
		in  string
		exp *Doc
		err string
	}{
		{"outline: foo\n", &Doc{Name: "foo"}, ""},
		{twoFuncsTabs, twoFuncs, ""},
		{twoFuncsSpaces, twoFuncs, ""},
		{timeSpaces, time, ""},
		{docWithDescriptionTabs, docWithDescription, ""},
		{huhSpaces, huh, ""},
	}

	for i, c := range cases {
		b := bytes.NewBufferString(c.in)
		got, err := Parse(b)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
			continue
		}

		if got == nil {
			t.Errorf("case %d doc returned nil", i)
			continue
		}

		gotB, _ := got.MarshalIndent(0, "  ")
		expB, _ := c.exp.MarshalIndent(0, "  ")
		if string(expB) != string(gotB) {
			t.Errorf("case %d equality mismatch. expected:\n%s\ngot:\n%s\n", i, string(expB), string(gotB))
			t.Log("\n", gotB, "\n", expB)
			diffs := differ.DiffMain(string(gotB), string(expB), true)
			t.Log(differ.DiffPrettyText(diffs))
		}
	}
}
