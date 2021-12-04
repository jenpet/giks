package plugins

import (
	"fmt"
	"strings"
)

const (
	varListA                     = "LIST_COMPARATOR_LIST_A"
	varListB                     = "LIST_COMPARATOR_LIST_B"
	varListOperator              = "LIST_COMPARATOR_OPERATOR"
	varListComparatorFailOnMatch = "LIST_COMPARATOR_FAIL_ON_SUCCESS"
)

var operations = []string{"intersect"}

type ListComparator struct{}

func (lc ListComparator) ID() string {
	return "list-comparator"
}

func (lc ListComparator) Run(workingDir string, hook string, vars map[string]string, args []string) (bool, error) {
	listAStr, err := extractStringVar(varListA, vars, true)
	if err != nil {
		return false, err
	}
	listA := strings.Split(listAStr, " ")

	listBStr, err := extractStringVar(varListB, vars, true)
	if err != nil {
		return false, err
	}
	listB := strings.Split(listBStr, " ")

	operation, err := extractStringVar(varListOperator, vars, true)
	if err != nil {
		return false, err
	}
	if !contains(operations, operation) {
		return false, fmt.Errorf("list-comparator does not support operation '%s'", operation)
	}

	failOnMatch, err := extractBoolVar(varListComparatorFailOnMatch, vars, false)
	if err != nil {
		return false, err
	}
	if compare(listA, listB, operation) {
		if failOnMatch {
			return true, fmt.Errorf("matched lists with operation '%s'", operation)
		}
	}
	return true, nil
}

func compare(a, b []string, operation string) bool {
	switch operation {
	case "intersect":
		for _, el := range a {
			if contains(b, el) {
				return true
			}
		}
	}
	return false
}

func contains(hs []string, n string) bool {
	for _, el := range hs {
		if el == n {
			return true
		}
	}
	return false
}
