package gotest_labels

// exp_parser.go provides a simple expression parser for logical expressions
// that can handle conditions and logical operators. It supports the following syntax:
// - Conditions in the form of "key=value"
// - Logical AND operator "&&"
// - Logical OR operator "||"
// - Parentheses for grouping expressions
//
// Example:	//   (key1=value1 && key2=value2) || (key3=value3 && key4=value4)
// The parser generates an abstract syntax tree (AST) representation of the expression.
// The AST consists of two types of nodes:
// - Condition nodes representing key-value pairs
// - LogicalOp nodes representing logical operations (AND/OR) with child nodes
// The parser can be used to evaluate expressions, validate syntax, and generate
// error messages for invalid input.

import (
	"fmt"
	"strings"
)

type Node any

type Condition struct {
	Key   string
	Value string
}

type LogicalOp struct {
	Operator string
	Children []Node
}

func tokenize(input string) ([]string, error) {
	var tokens []string
	runes := []rune(input)
	n := len(runes)
	i := 0
	buffer := make([]rune, 0, n)

	for i < n {
		if i+1 < n && runes[i] == '&' && runes[i+1] == '&' {
			if len(buffer) > 0 {
				tokens = append(tokens, string(buffer))
				buffer = buffer[:0]
			}
			tokens = append(tokens, "&&")
			i += 2
		} else if i+1 < n && runes[i] == '|' && runes[i+1] == '|' {
			if len(buffer) > 0 {
				tokens = append(tokens, string(buffer))
				buffer = buffer[:0]
			}
			tokens = append(tokens, "||")
			i += 2
		} else if runes[i] == '(' || runes[i] == ')' {
			if len(buffer) > 0 {
				tokens = append(tokens, string(buffer))
				buffer = buffer[:0]
			}
			tokens = append(tokens, string(runes[i]))
			i++
		} else {
			buffer = append(buffer, runes[i])
			i++
		}
	}

	if len(buffer) > 0 {
		tokens = append(tokens, string(buffer))
	}
	return tokens, nil
}

func parseFactor(tokens []string, pos int) (Node, int, error) {
	if pos >= len(tokens) {
		return nil, pos, fmt.Errorf("unexpected end of input")
	}
	token := tokens[pos]
	if token == "(" {
		node, newPos, err := parseExpr(tokens, pos+1)
		if err != nil {
			return nil, pos, err
		}
		if newPos >= len(tokens) || tokens[newPos] != ")" {
			return nil, pos, fmt.Errorf("expected closing bracket")
		}
		return node, newPos + 1, nil
	} else if strings.Contains(token, "=") {
		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 {
			return nil, pos, fmt.Errorf("invalid key=value pair: %s", token)
		}
		return Condition{Key: parts[0], Value: parts[1]}, pos + 1, nil
	}
	return nil, pos, fmt.Errorf("unexpected token: %s", token)
}

func parseTerm(tokens []string, pos int) (Node, int, error) {
	left, newPos, err := parseFactor(tokens, pos)
	if err != nil {
		return nil, pos, err
	}
	pos = newPos

	for pos < len(tokens) && tokens[pos] == "&&" {
		pos++
		right, newPos, err := parseFactor(tokens, pos)
		if err != nil {
			return left, pos, err
		}
		pos = newPos
		left = LogicalOp{
			Operator: "AND",
			Children: []Node{left, right},
		}
	}
	return left, pos, nil
}

func parseExpr(tokens []string, pos int) (Node, int, error) {
	left, newPos, err := parseTerm(tokens, pos)
	if err != nil {
		return nil, pos, err
	}
	pos = newPos

	for pos < len(tokens) && tokens[pos] == "||" {
		pos++
		right, newPos, err := parseTerm(tokens, pos)
		if err != nil {
			return left, pos, err
		}
		pos = newPos
		left = LogicalOp{
			Operator: "OR",
			Children: []Node{left, right},
		}
	}
	return left, pos, nil
}

// Evaluate traverses the AST and evaluates the expression
// against the provided labels. It returns true if the expression is satisfied.
func Evaluate(node Node, labels TestLabels) bool {
    switch n := node.(type) {
    case Condition:
        val, ok := labels[n.Key]
        return ok && val == n.Value
    case LogicalOp:
        switch n.Operator {
        case "AND":
            return Evaluate(n.Children[0], labels) && Evaluate(n.Children[1], labels)
        case "OR":
            return Evaluate(n.Children[0], labels) || Evaluate(n.Children[1], labels)
        default:
            return false
        }
    default:
        return false
    }
}

// The entry point for parsing label expressions
// It takes a string input and returns an AST representation of the expression
// or an error if the input is invalid.
func ParseLabelExp(input string) (Node, error) {
	tokens, err := tokenize(input)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	node, pos, err := parseExpr(tokens, 0)
	if err != nil {
		return nil, err
	}
	if pos < len(tokens) {
		return nil, fmt.Errorf("unexpected token %s at position %d", tokens[pos], pos)
	}
	return node, nil
}
