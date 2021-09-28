package lib

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var differ = diffmatchpatch.New()

const twoFuncsTabs = `outline: twoFuncs
	path: twoFuncs
	functions:
		difference(a,b int) int
		sum(a,b int) int
			add two things together`

const twoFuncsSpaces = `outline: twoFuncs
	path: twoFuncs
  functions:
    difference(a,b int) int
    sum(a,b int) int
			add two things together`

var twoFuncs = &Doc{
	Name: "twoFuncs",
	Path: "twoFuncs",
	Functions: []*Function{
		{FuncName: "difference",
			Signature: "difference(a,b int) int",
			Receiver:  "twoFuncs",
		},
		{
			FuncName:    "sum",
			Signature:   "sum(a,b int) int",
			Description: "add two things together",
			Receiver:    "twoFuncs",
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
      a period of time
      methods:
        add(d duration) int
          params:
             d duration
      fields:
        hours float
          number of hours starting at zero
        minutes float
        nanoseconds int
        seconds
          number of seconds starting at zero
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
		{FuncName: "duration",
			Signature:   "duration(string) duration",
			Description: "parse a duration",
			Receiver:    "time"},
		{FuncName: "time",
			Signature:   "time(string, format=..., location=...) time",
			Description: "parse a time",
			Receiver:    "time"},
		{FuncName: "now",
			Signature:   "now() time",
			Description: "new time instance set to current time implementations are able to make this a constant",
			Receiver:    "time"},
		{FuncName: "zero",
			Signature:   "zero() time",
			Description: "a constant",
			Receiver:    "time"},
	},
	Types: []*Type{
		{Name: "duration",
			Description: "a period of time",
			Methods: []*Function{
				{FuncName: "add",
					Receiver:  "duration",
					Signature: "add(d duration) int",
					Params: []*Param{
						{Name: "d", Type: "duration"},
					},
				},
			},
			Fields: []*Field{
				{Name: "hours", Type: "float", Description: "number of hours starting at zero"},
				{Name: "minutes", Type: "float"},
				{Name: "nanoseconds", Type: "int"},
				{Name: "seconds", Description: "number of seconds starting at zero"},
			},
			Operators: []*Operator{
				{Opr: "duration - time = duration"},
				{Opr: "duration + time = time"},
				{Opr: "duration == duration = boolean"},
				{Opr: "duration < duration = booleans"},
			},
		},
		{Name: "time",
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
		{FuncName: "sum", Signature: "sum(a int, b int) int", Receiver: "doc"},
	},
}

const huhSpaces = `  outline: huh
  huh is a package that has no meaning or purpose
  functions:
    foo(bar string) int
      foo a bar, which is to to a bar and remove 'd' from 'food'
      params:
        bar string
          the name of a bar
    date() date
      make a date`

var huh = &Doc{
	Name:        "huh",
	Description: "huh is a package that has no meaning or purpose",
	Functions: []*Function{
		{FuncName: "foo",
			Receiver:    "huh",
			Signature:   "foo(bar string) int",
			Description: "foo a bar, which is to to a bar and remove 'd' from 'food'",
			Params: []*Param{
				{Name: "bar", Type: "string", Description: "the name of a bar"},
			},
		},
		{FuncName: "date",
			Signature:   "date() date",
			Description: "make a date",
			Receiver:    "huh"},
	},
}

const dataframeTabs = `
outline: dataframe
dataframe is a 2d columnar data structure that provides analysis and manipulation tools
path: dataframe
functions:
	DataFrame(data, index, columns, dtype) DataFrame
		constructs a DataFrame
		params:
			data any
				data for the content of the DataFrame
			index
				index for the rows of the DataFrame
`

var dataframe = &Doc{
	Name:        "dataframe",
	Description: "dataframe is a 2d columnar data structure that provides analysis and manipulation tools",
	Path:        "dataframe",
	Functions: []*Function{
		{FuncName: "DataFrame",
			Signature:   "DataFrame(data, index, columns, dtype) DataFrame",
			Description: "constructs a DataFrame",
			Receiver:    "dataframe",
			Params: []*Param{
				{Name: "data", Type: "any", Description: "data for the content of the DataFrame"},
				{Name: "index", Type: "", Description: "index for the rows of the DataFrame"},
			},
		},
	},
}

func TestParse(t *testing.T) {
	cases := []struct {
		name string
		in   string
		exp  *Doc
		err  string
	}{
		{"basic", "outline: foo\n", &Doc{Name: "foo"}, ""},
		{"two_funcs_tabs", twoFuncsTabs, twoFuncs, ""},
		{"two_funcs_spaces", twoFuncsSpaces, twoFuncs, ""},
		{"time", timeSpaces, time, ""},
		{"doc_with_description", docWithDescriptionTabs, docWithDescription, ""},
		{"huh", huhSpaces, huh, ""},
		{"dataframe", dataframeTabs, dataframe, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := bytes.NewBufferString(c.in)
			got, err := ParseFirst(b)
			if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
				t.Fatalf("error mismatch. expected: %s, got: %s", c.err, err)
			}

			if got == nil {
				t.Fatal("doc returned nil")
			}

			if diff := cmp.Diff(c.exp, got, cmpopts.IgnoreUnexported(Doc{})); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}

			gotB, _ := got.MarshalIndent(0, "  ")
			expB, _ := c.exp.MarshalIndent(0, "  ")
			if diff := cmp.Diff(string(expB), string(gotB)); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
