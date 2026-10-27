package regression_test

import (
	"testing"

	"github.com/haomingbai/lgc/internal/testsuite"
)

// TestRegressionSuite is the intentionally empty regression unit-test shell.
func TestRegressionSuite(t *testing.T) {
	testsuite.Run(t, testsuite.Suite{
		Name: "regression",
	})
}
