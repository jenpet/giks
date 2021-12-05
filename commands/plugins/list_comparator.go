package plugins

import (
	"fmt"
	"strings"
)

const (
	varListA                     = "LIST_COMPARATOR_LIST_A"
	varListB                     = "LIST_COMPARATOR_LIST_B"
	varListOperator              = "LIST_COMPARATOR_OPERATION"
	varListComparatorFailOnMatch = "LIST_COMPARATOR_FAIL_ON_MATCH"
)

var operations = []string{"intersect"}

type ListComparator struct{}

func (lc ListComparator) ID() string {
	return "list-comparator"
}

func (lc ListComparator) Run(workingDir string, hook string, vars map[string]string, args []string) (bool, error) {
	listAStr, err := extractStringVar(varListA, vars, false)
	if err != nil {
		return false, err
	}
	listA := strings.Split(listAStr, " ")

	listBStr, err := extractStringVar(varListB, vars, false)
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
	if diff := compare(listA, listB, operation); len(diff) > 0 {
		if failOnMatch {
			return false, fmt.Errorf("elements which matched the comparison '%s'", strings.Join(diff, ","))
		}
	}
	return true, nil
}

func compare(a, b []string, operation string) []string {
	var diff []string
	switch operation {
	case "intersect":
		for _, el := range a {
			if contains(b, el) {
				diff = append(diff, el)
			}
		}
	}
	return diff
}

func contains(hs []string, n string) bool {
	for _, el := range hs {
		if el == n {
			return true
		}
	}
	return false
}
