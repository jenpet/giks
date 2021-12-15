package plugins

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListComparator(t *testing.T) {
	listComparatorTests := []struct {
		name        string
		listA       string
		listB       string
		operation   string
		failOnMatch  string
		exitExpected bool
		errExpected  bool
	}{
		{
			"lists intersect fail",
			"a b c",
			"c d e",
			"intersect",
			"true",
			true,
			true,
		},
		{
			"lists intersect pass",
			"a b c",
			"c d e",
			"intersect",
			"false",
			false,
			true,
		},
		{
			"lists no intersection",
			"a b",
			"d e",
			"intersect",
			"true",
			false,
			false,
		},
		{
			"invalid operation",
			"a b",
			"d e",
			"foo",
			"true",
			true,
			true,
		},
		{
			"empty lists",
			" ",
			"    ",
			"intersect",
			"true",
			false,
			false,
		},
	}
	lc, _ := Get("list-comparator")
	for _, tt := range listComparatorTests {
		t.Run(tt.name, func(t *testing.T) {
			vars := map[string]string{
				"LIST_COMPARATOR_LIST_A":        tt.listA,
				"LIST_COMPARATOR_LIST_B":        tt.listB,
				"LIST_COMPARATOR_OPERATION":     tt.operation,
				"LIST_COMPARATOR_FAIL_ON_MATCH": tt.failOnMatch,
			}
			ok, err := lc.Run("", "pre-commit", vars, nil)
			assert.Equal(t, tt.exitExpected, ok, "expected bool does not match executed plugin result")
			assert.Equal(t, tt.errExpected, err != nil, "error expectation and result does not match")
		})
	}
}
