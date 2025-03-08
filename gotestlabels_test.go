package gotest_labels

import (
	"os"
	"testing"
)

func TestMutateTestFilterByLabels(t *testing.T) {
	t.Run("Labels enabled", func(t *testing.T) {
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

		if tests[0] != "TestSimpleAlpha" {
			t.Errorf("Expected TestSimpleAlpha, got %v", tests[0])
			t.Fail()
		}

		if tests[1] != "TestSimpleGamma" {
			t.Errorf("Expected TestSimpleGamma, got %v", tests[1])
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

		if os.Args[4] != "^TestSimpleAlpha|TestSimpleGamma$" {
			t.Errorf("Expected ^TestSimpleAlpha|TestSimpleGamma$, got %v", os.Args[4])
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
		if len(tests)!= 3 {
			t.Errorf("Expected 3 tests, got %v", len(tests))
			t.Fail()
		}
		if tests[0] != "TestSimpleAlpha" {
			t.Errorf("Expected TestSimpleAlpha, got %v", tests[0])
			t.Fail()
		}
		if tests[1] != "TestSimpleBeta" {
			t.Errorf("Expected TestSimpleBeta, got %v", tests[1])
			t.Fail()
		}
		if tests[2] != "TestSimpleGamma" {
			t.Errorf("Expected TestSimpleGamma, got %v", tests[2])
			t.Fail()
		}
		if len(os.Args)!= 3 {
			t.Errorf("Expected 3 args after mutation, got %#v", os.Args)
			t.Fail()	
		}
	})
}
