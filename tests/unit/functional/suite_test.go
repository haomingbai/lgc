package functional_test

import (
	"testing"

	"github.com/haomingbai/lgc/internal/testsuite"
)

// TestFunctionalSuite is the intentionally empty functional unit-test shell.
func TestFunctionalSuite(t *testing.T) {
	testsuite.Run(t, testsuite.Suite{
		Name: "functional",
	})
}
