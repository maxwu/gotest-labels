package explicittestmain

import (
	"testing"

	"gotestlabels"
)

func init() {
	_ = gotestlabels.MutateTestFilterByLabels()
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
