package simple

import (
	"testing"

	// The user testing packages only need to import the apply package to selectively
	// run the tests with labels
	_ "github.com/maxwu/gotest-labels/apply"
)

// @group=demo
func TestSimpleAlpha(t *testing.T) {
	t.Log("Testing examples.simple.TestSimpleAlpha")
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
