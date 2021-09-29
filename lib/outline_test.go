package lib

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const unsorted = `
	outline: twoFuncs
		path: twoFuncs
		functions:
			difference(a,b int) int
			sum(a,b int) int
				add two things together

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

func TestSort(t *testing.T) {
	docs, err := Parse(strings.NewReader(unsorted), AlphaSortFuncs(), AlphaSortTypes())
	if err != nil {
		t.Fatal(err)
	}

	docs.Sort()
	data, err := docs[0].MarshalIndent(0, "\t")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(expectA, string(data)); diff != "" {
		t.Errorf("serialized output mismatch (-want +got):\n%s", diff)
	}

	data, err = docs[1].MarshalIndent(0, "\t")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
	if diff := cmp.Diff(expectB, string(data)); diff != "" {
		t.Errorf("serialized output mismatch (-want +got):\n%s", diff)
	}
}

const expectA = `outline: time
	functions:
		duration(string) duration
			parse a duration
		now() time
			new time instance set to current time implementations are able to make this a constant
		time(string, format=..., location=...) time
			parse a time
		zero() time
			a constant
	types:
		duration
				a period of time
			fields:
				hours float
				minutes float
				nanoseconds int
				seconds float
		time
`

const expectB = `outline: twoFuncs
	path: twoFuncs
	functions:
		difference(a,b int) int
		sum(a,b int) int
			add two things together
`
