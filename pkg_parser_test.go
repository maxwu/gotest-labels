package gotest_labels

import (
	"go/ast"
	"regexp"
	"testing"
)

func TestFilterTestFuncs(t *testing.T) {
	tests := []struct {
		name     string
		funcs    map[string]TestLabels
		regex    *regexp.Regexp
		expected map[string]TestLabels
	}{
		{
			name: "Nil regex returns all funcs",
			funcs: map[string]TestLabels{
				"TestA": {},
				"TestB": {},
			},
			regex: nil,
			expected: map[string]TestLabels{
				"TestA": {},
				"TestB": {},
			},
		},
		{
			name: "Regex matches some funcs",
			funcs: map[string]TestLabels{
				"TestA":       {},
				"TestB":       {},
				"ExampleTest": {},
			},
			regex: regexp.MustCompile("^Test"),
			expected: map[string]TestLabels{
				"TestA": {},
				"TestB": {},
			},
		},
		{
			name: "Regex matches no funcs",
			funcs: map[string]TestLabels{
				"TestA": {},
				"TestB": {},
			},
			regex:    regexp.MustCompile("^Example"),
			expected: map[string]TestLabels{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := filterTestFuncs(tt.funcs, tt.regex)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
			for key, _ := range tt.expected {
				if result[key] == nil {
					t.Errorf("expected %v key", key)
				}
			}
		})
	}
}

func TestGetFuncLabels(t *testing.T) {
	tests := []struct {
		name     string
		fn       *ast.FuncDecl
		expected TestLabels
	}{
		{
			name: "No comments",
			fn: &ast.FuncDecl{
				Doc: nil,
			},
			expected: make(TestLabels),
		},
		{
			name: "Single valid comment with key-value",
			fn: &ast.FuncDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// @key=value",
						},
					},
				},
			},
			expected: TestLabels{"key": "value"},
		},
		{
			name: "Multiple comments with mixed styles",
			fn: &ast.FuncDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// @key1=value1",
						},
						{
							Text: "/* @key2=value2 */",
						},
						{
							Text: "// @key3",
						},
					},
				},
			},
			expected: TestLabels{
				"key1": "value1",
				"key2": "value2",
				"key3": DefaultLabelValue,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := getFuncLabels(tt.fn)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
			for key, value := range tt.expected {
				if result[key] != value {
					t.Errorf("for key %q, expected %q, got %q", key, value, result[key])
				}
			}
		})
	}
}
