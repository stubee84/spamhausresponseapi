package helpers

import (
	"testing"
)

type expectedReverse struct {
	source   string
	expected string
}

func initActualExpected() []expectedReverse {
	return []expectedReverse{
		{
			source:   "37.156.173.11",
			expected: "11.173.156.37",
		},
		{
			source:   "185.77.248.1",
			expected: "1.248.77.185",
		},
	}
}

func TestReverse(t *testing.T) {
	for _, expected := range initActualExpected() {
		actual, _ := ReverseUsingSeperator(expected.source, ".")

		if actual != expected.expected {
			t.Logf("actual: %s. is not equal to expected: %s. for source: %s\n", actual, expected.expected, expected.source)
			t.FailNow()
		}
	}
}
