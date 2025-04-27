package gotest_labels

import (
	"os"
	"testing"
)

func TestMutateTestFilterByLabels(t *testing.T) {
	t.Run("Labels enabled in go test list", func(t *testing.T) {
		t.Parallel()

		origArgs := os.Args
		defer func() { os.Args = origArgs }()
		os.Args = []string{"theBinDoesntMatter", "-test.v", "-labels", "group=demo", "-test.list=."}
		origDefaultPkg := defaultPkg
		defer func() { defaultPkg = origDefaultPkg }()

		defaultPkg = "./examples/simple"

		tests := MutateTestFilterByLabels()

		t.Logf("Filtered tests: %v\n", tests)

		if len(tests) != 2 {
			t.Errorf("Expected 2 tests, got %v", len(tests))
			t.Fail()
		}

		if tests["TestSimpleAlpha"] == nil {
			t.Errorf("Expected TestSimpleAlpha")
			t.Fail()
		}

		if tests["TestSimpleGamma"] == nil {
			t.Errorf("Expected TestSimpleGamma")
			t.Fail()
		}

		if len(os.Args) != 5 {
			t.Errorf("Expected 5 args after mutation, got %#v", os.Args)
			t.Fail()
		}

		if os.Args[3] != "-test.list" {
			t.Errorf("Expected test.list, got %v", os.Args[3])
			t.Fail()
		}

		if os.Args[4] != "^TestSimpleAlpha|TestSimpleGamma$" && os.Args[4] != "^TestSimpleGamma|TestSimpleAlpha$" {
			// The order of the tests in the regex may vary
			t.Errorf("Expected ^TestSimpleAlpha|TestSimpleGamma$ or ^TestSimpleAlpha|TestSimpleGamma$, got %v", os.Args[4])
			t.Fail()
		}
	})

	t.Run("Labels enabled in go test run", func(t *testing.T) {
		origArgs := os.Args
		defer func() { os.Args = origArgs }()
		os.Args = []string{"theBinDoesntMatter", "-test.v", "-labels", "group=demo", "-test.run=Alpha"}
		origDefaultPkg := defaultPkg
		defer func() { defaultPkg = origDefaultPkg }()

		defaultPkg = "./examples/simple"

		tests := MutateTestFilterByLabels()

		t.Logf("Filtered tests: %v\n", tests)

		if len(tests) != 1 {
			t.Errorf("Expected 1 test, got %v", len(tests))
			t.Fail()
		}

		if tests["TestSimpleAlpha"] == nil {
			t.Errorf("Expected TestSimpleAlpha")
			t.Fail()
		}

		if len(os.Args) != 5 {
			t.Errorf("Expected 5 args after mutation, got %#v", os.Args)
			t.Fail()
		}

		if os.Args[3] != "-test.run" {
			t.Errorf("Expected test.list, got %v", os.Args[3])
			t.Fail()
		}

		if os.Args[4] != "^TestSimpleAlpha$" {
			t.Errorf("Expected ^TestSimpleAlpha$, got %v", os.Args[4])
			t.Fail()
		}
	})

	t.Run("Labels disabled", func(t *testing.T) {
		origArgs := os.Args
		defer func() { os.Args = origArgs }()
		os.Args = []string{"theBinDoesntMatter", "-test.v", "-test.list=."}
		origDefaultPkg := defaultPkg
		defer func() { defaultPkg = origDefaultPkg }()
		defaultPkg = "./examples/simple"
		tests := MutateTestFilterByLabels()
		t.Logf("Filtered tests: %v\n", tests)
		if len(tests) != 3 {
			t.Errorf("Expected 3 tests, got %v", len(tests))
			t.Fail()
		}
		if tests["TestSimpleAlpha"] == nil {
			t.Errorf("Expected TestSimpleAlpha")
			t.Fail()
		}
		if tests["TestSimpleBeta"] == nil {
			t.Errorf("Expected TestSimpleBeta")
			t.Fail()
		}
		if tests["TestSimpleGamma"] == nil {
			t.Errorf("Expected TestSimpleGamma")
			t.Fail()
		}
		if tests["TestUnderSimpleDelta"] != nil {
			t.Errorf("Unexpected TestUnderSimpleDelta")
			t.Fail()
		}
		if len(os.Args) != 3 {
			t.Errorf("Expected 3 args after mutation, got %#v", os.Args)
			t.Fail()
		}
	})
}
