package undersimple

import (
	"testing"

	// Limited by go package loading mechanism, if sub packages are involved, the apply package
	// is also needed to be imported in the sub packages to activate the label filtering.
	// Otherwise, only applied testing packages will run tests filtered by labels, the others
	// will run all tests as defined in CLI args.
	_ "gotestlabels/apply"
)

// func TestMain(m *testing.M) {
// 	gotestlabels.MutateTestFuncsByLabels()
// 	code := m.Run()
// 	os.Exit(code)
// }

func TestUnderSimpleDelta(t *testing.T) {
	t.Log("Testing examples.simple.undersimple.TestUnderSimpleDelta")
}

// @group=demo
func TestUnderSimpleEpsilon(t *testing.T) {
	t.Log("Testing examples.simple.undersimple.TestUnderSimpleEpsilon")
}