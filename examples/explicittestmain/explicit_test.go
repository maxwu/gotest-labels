package explicittestmain

import (
	"fmt"
	"os"
	"testing"

	"github.com/maxwu/gotest-labels"
)

func TestMain(m *testing.M) {
	tests := gotest_labels.MutateTestFilterByLabels()
	fmt.Printf("Filtered tests: %v\n", tests)
	code := m.Run()
	os.Exit(code)
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
