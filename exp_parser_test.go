package gotest_labels

import (
	"fmt"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := map[string]struct {
		exp  string
		want []string
	}{
		"empty": {
			exp:  "",
			want: []string{},
		},
		"single condition": {
			exp:  "key=value",
			want: []string{"key=value"},
		},
		"AND condition": {
			exp:  "key=value&&key2=value2",
			want: []string{"key=value", "&&", "key2=value2"},
		},
		"OR condition": {
			exp:  "key=value||key2=value2",
			want: []string{"key=value", "||", "key2=value2"},
		},
		"NOT condition": {
			exp:  "!key=value",
			want: []string{"!", "key=value"},
		},
		"mixed condition": {
			exp:  "key=value&&key2=value2||key3=value3",
			want: []string{"key=value", "&&", "key2=value2", "||", "key3=value3"},
		},
		"parentheses": {
			exp:  "(key=value&&key2=value2)||!key3=value3",
			want: []string{"(", "key=value", "&&", "key2=value2", ")", "||", "!", "key3=value3"},
		},
		"nested parentheses": {
			exp:  "((key=value&&key2=value2)||key3=value3)&&key4=value4",
			want: []string{"(", "(", "key=value", "&&", "key2=value2", ")", "||", "key3=value3", ")", "&&", "key4=value4"},
		},
		"complex expression": {
			exp:  "((key=value&&key2=value2)||key3=value3)&&((key4=value4||key5=value5)&&key6=value6)",
			want: []string{"(", "(", "key=value", "&&", "key2=value2", ")", "||", "key3=value3", ")", "&&", "(", "(", "key4=value4", "||", "key5=value5", ")", "&&", "key6=value6", ")"},
		},
		"spaces around operators": {
			exp:  "key=value && key2=value2 || key3=value3",
			want: []string{"key=value", "&&", "key2=value2", "||", "key3=value3"},
		},
		"spaces around parentheses": {
			exp: " ( key=value && key2=value2 ) || key3=value3 ",
			want: []string{"(", "key=value", "&&", "key2=value2",
				")", "||", "key3=value3"},
		},
		"spaces around conditions": {
			exp:  " key = value && key2 = value2 || key3 = value3 ",
			want: []string{"key", "=", "value", "&&", "key2", "=", "value2", "||", "key3", "=", "value3"},
		},
		"spaces around operators and parentheses": {
			exp:  " ( key = value && key2 = value2 ) || key3 = value3 ",
			want: []string{"(", "key", "=", "value", "&&", "key2", "=", "value2", ")", "||", "key3", "=", "value3"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tokenize(test.exp)
			if len(got) != len(test.want) {
				t.Errorf("tokenize(%q) returned %#v, want %#v", test.exp, got, test.want)
				return
			}
			for i := range got {
				if got[i] != test.want[i] {
					t.Errorf("tokenize(%q)[%d] returned %#v, want %#v", test.exp, i, got[i], test.want[i])
				}
			}
		})
	}
}

func TestParseLabelExp(t *testing.T) {
	tests := map[string]struct {
		exp  string
		want string
		err  string
	}{
		"empty": {
			exp: "",
			err: "empty input",
		},
		"single condition": {
			exp:  "key=value",
			want: `gotest_labels.Condition{Key:"key", Value:"value"}`,
			err:  "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseLabelExp(test.exp)
			if err != nil && test.err == "" {
				t.Errorf("ParseLabelExp(%q) generated \"%v\", want no error", test.exp, err)
			}
			if err != nil && test.err != "" {
				if err.Error() != test.err {
					t.Errorf("ParseLabelExp(%q) generated \"%v\", want %v", test.exp, err, test.err)
				}
				return
			}
			output := fmt.Sprintf("%#v", got)
			if output != test.want {
				t.Errorf("ParseLabelExp(%q) = %v, want %v", test.exp, output, test.want)
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	tests := map[string]struct {
		exp    string
		labels TestLabels
		want   bool
		err    string
	}{
		"single condition": {
			exp:    "env=dev",
			labels: TestLabels{"env": "dev", "group": "demo"},
			want:   true,
		},
		"AND condition - positive": {
			exp:    "env=dev&&group=demo",
			labels: TestLabels{"env": "dev", "group": "demo", "integration": "true"},
			want:   true,
		},
		"AND condition - negative": {
			exp:    "env=dev&&group=demo",
			labels: TestLabels{"env": "dev", "group": "prod", "integration": "true"},
			want:   false,
		},
		"OR condition": {
			exp:    "env=dev||group=demo",
			labels: TestLabels{"env": "dev", "group": "prod"},
			want:   true,
		},
		"OR condition with Parentheses - 1": {
			exp:    "(env=dev||env=int)&&group=demo",
			labels: TestLabels{"env": "dev", "group": "demo"},
			want:   true,
		},
		"OR condition with Parentheses - 2": {
			exp:    "(env=dev||env=int)&&group=demo",
			labels: TestLabels{"env": "int", "group": "demo"},
			want:   true,
		},
		"NOT condition - positive - 1 missing": {
			exp:    "!env=dev",
			labels: TestLabels{"group": "prod"},
			want:   true,
		},
		"NOT condition - positive - 2 mismatching": {
			exp:    "!env=dev",
			labels: TestLabels{"dev": "prod"},
			want:   true,
		},
		"NOT condition - positive - 1 brackets": {
			exp:    "!(env=dev)",
			labels: TestLabels{"env": "somthing-else", "key": "dev"},
			want:   true,
		},
		"Combined conditions - 1": {
			exp:    "!env=dev&&group=demo",
			labels: TestLabels{"group": "demo", "env": "prod"},
			want:   true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			node, err := ParseLabelExp(test.exp)
			if err != nil && test.err == "" {
				t.Errorf("Evaluate(%q) generated \"%v\", want no error", test.exp, err)
			}
			if err != nil && test.err != "" {
				if err.Error() != test.err {
					t.Errorf("Evaluate(%q) generated \"%v\", want %v", test.exp, err, test.err)
				}
				return
			}

			got := Evaluate(node, test.labels)
			if got != test.want {
				t.Errorf("Evaluate(%q) = %v, want %v", test.exp, got, test.want)
			}
		})
	}
}
