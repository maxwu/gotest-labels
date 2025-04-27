package simple

import (
	"os"
	"testing"

	// The user testing packages only need to import the apply package to selectively
	// run the tests with labels
	_ "github.com/maxwu/gotest-labels/apply"
)

// @group=demo
// @regression
func TestSimpleAlpha(t *testing.T) {
	t.Log("Testing examples.simple.TestSimpleAlpha")
	t.Log("The args are: ", os.Args[1:])
	t.Log("The env is: ", os.Getenv("TEST_LABELS"))
}

// A test case with two labels: group=integration and env=dev
// @group=integration
/* @env=dev */
func TestSimpleBeta(t *testing.T) {
	t.Log("Testing examples.simple.TestSimpleBeta")
}

// @group=demo
// @env=prod
func TestSimpleGamma(t *testing.T) {
	t.Log("Testing examples.simple.TestSimpleGamma")
}
