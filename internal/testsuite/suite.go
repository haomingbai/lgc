// Package testsuite provides the small suite harness shared by LGC test
// categories.
package testsuite

import (
	"fmt"
	"testing"
)

// Case is one named test body in a suite.
type Case struct {
	Name string
	Run  func(t *testing.T)
}

// Suite groups related unit-test cases.
type Suite struct {
	Name  string
	Cases []Case
}

// Run executes a suite and its named cases as Go subtests.
func Run(t *testing.T, suite Suite) {
	t.Helper()

	t.Run(suite.Name, func(t *testing.T) {
		for index, testCase := range suite.Cases {
			name := testCase.Name
			if name == "" {
				name = fmt.Sprintf("case_%d", index)
			}
			if testCase.Run == nil {
				t.Fatalf("%s: test case has no body", name)
			}
			t.Run(name, testCase.Run)
		}
	})
}
