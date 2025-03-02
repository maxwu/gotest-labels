package apply

import "gotestlabels"

// The apply package is a shortcurt to automatically apply the test filter by labels via anonymous import.
func init() {
	_ = gotestlabels.MutateTestFilterByLabels()
}
