package explicittestmain

import (
	"testing"

	"github.com/maxwu/gotest-labels"
)

func init() {
	_ = gotest_labels.MutateTestFilterByLabels()
}

// @group=demo
func TestExplicitAlpha(t *testing.T) {
	t.Log("Testing examples.simple.TestSimpleAlpha")
}

func TestExplicitBeta(t *testing.T) {
	t.Log("Testing examples.simple.TestSimpleBeta")
}

// @group=demo
func TestExplicitGamma(t *testing.T) {
	t.Log("Testing examples.simple.TestSimpleGamma")
}
