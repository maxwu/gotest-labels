package apply

import "github.com/maxwu/gotest-labels"

// The apply package is a shortcurt to automatically apply the test filter by labels via anonymous import.
func init() {
	_ = gotest_labels.MutateTestFilterByLabels()
}
