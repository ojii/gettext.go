package pluralforms

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type match struct {
	OpenPos  int
	ClosePos int
}

var pat = regexp.MustCompile(`(\?|:|\|\||&&|==|!=|>=|>|<=|<|%|\d+|n)`)

type expr_token interface {
	Compile(tokens []string) (expr Expression, err error)
}

type test_token interface {
	Compile(tokens []string) (test Test, err error)
}

type cmp_test_builder func(val uint32, flipped bool) Test
type logic_test_build func(left Test, right Test) Test

var ternary ternary_

type ternary_ struct{}

func (ternary_) Compile(tokens []string) (expr Expression, err error) {
	main, err := split_tokens(tokens, "?")
	if err != nil {
		return expr, err
	}
	test, err := compile_test(strings.Join(main.Left, ""))
	if err != nil {
		return expr, err
	}
	actions, err := split_tokens(main.Right, ":")
	if err != nil {
		return expr, err
	}
	true_action, err := compile_expression(strings.Join(actions.Left, ""))
	if err != nil {
		return expr, err
	}
	false_action, err := compile_expression(strings.Join(actions.Right, ""))
	if err != nil {
		return expr, nil
	}
	return Ternary{
		Test:  test,
		True:  true_action,
		False: false_action,
	}, nil
}

var const_val const_val_

type const_val_ struct{}

func (const_val_) Compile(tokens []string) (expr Expression, err error) {
	if len(tokens) == 0 {
		return expr, errors.New("Got nothing instead of constant")
	}
	if len(tokens) != 1 {
		return expr, errors.New(fmt.Sprintf("Invalid constant: %s", strings.Join(tokens, "")))
	}
	i, err := strconv.Atoi(tokens[0])
	if err != nil {
		return expr, err
	}
	return Const{Value: i}, nil
}

func compile_logic_test(tokens []string, sep string, builder logic_test_build) (test Test, err error) {
	split, err := split_tokens(tokens, sep)
	if err != nil {
		return test, err
	}
	left, err := compile_test(strings.Join(split.Left, ""))
	if err != nil {
		return test, err
	}
	right, err := compile_test(strings.Join(split.Right, ""))
	if err != nil {
		return test, err
	}
	return builder(left, right), nil
}

var or or_

type or_ struct{}

func (or_) Compile(tokens []string) (test Test, err error) {
	return compile_logic_test(tokens, "||", build_or)
}
func build_or(left Test, right Test) Test {
	return Or{Left: left, Right: right}
}

var and and_

type and_ struct{}

func (and_) Compile(tokens []string) (test Test, err error) {
	return compile_logic_test(tokens, "&&", build_and)
}
func build_and(left Test, right Test) Test {
	return And{Left: left, Right: right}
}

func compile_mod(tokens []string) (math Math, err error) {
	split, err := split_tokens(tokens, "%")
	if err != nil {
		return math, err
	}
	if len(split.Left) != 1 || split.Left[0] != "n" {
		return math, errors.New("Modulus operation requires 'n' as left operand")
	}
	if len(split.Right) != 1 {
		return math, errors.New("Modulus operation requires simple integer as right operand")
	}
	i, err := parse_uint32(split.Right[0])
	if err != nil {
		return math, err
	}
	return Mod{Value: uint32(i)}, nil
}

func _pipe(mod_tokens []string, action_tokens []string, builder cmp_test_builder, flipped bool) (test Test, err error) {
	modifier, err := compile_mod(mod_tokens)
	if err != nil {
		return test, err
	}
	if len(action_tokens) != 1 {
		return test, errors.New("Can only get modulus of integer")
	}
	i, err := parse_uint32(action_tokens[0])
	if err != nil {
		return test, err
	}
	action := builder(uint32(i), flipped)
	return Pipe{
		Modifier: modifier,
		Action:   action,
	}, nil
}

func compile_equality(tokens []string, sep string, builder cmp_test_builder) (test Test, err error) {
	split, err := split_tokens(tokens, sep)
	if err != nil {
		return test, err
	}
	if len(split.Left) == 1 && split.Left[0] == "n" {
		if len(split.Right) != 1 {
			return test, errors.New("test can only compare n to integers")
		}
		i, err := parse_uint32(split.Right[0])
		if err != nil {
			return test, err
		}
		return builder(i, false), nil
	} else if len(split.Right) == 1 && split.Right[0] == "n" {
		if len(split.Left) != 1 {
			return test, errors.New("test can only compare n to integers")
		}
		i, err := parse_uint32(split.Left[0])
		if err != nil {
			return test, err
		}
		return builder(i, true), nil
	} else if contains(split.Left, "n") && contains(split.Left, "%") {
		return _pipe(split.Left, split.Right, builder, false)
	} else {
		return test, errors.New("equality test must have 'n' as one of the two tests")
	}
}

var eq eq_

type eq_ struct{}

func (eq_) Compile(tokens []string) (test Test, err error) {
	return compile_equality(tokens, "==", build_eq)
}
func build_eq(val uint32, flipped bool) Test {
	return Equal{Value: val}
}

var neq neq_

type neq_ struct{}

func (neq_) Compile(tokens []string) (test Test, err error) {
	return compile_equality(tokens, "!=", build_neq)
}
func build_neq(val uint32, flipped bool) Test {
	return NotEqual{Value: val}
}

var gt gt_

type gt_ struct{}

func (gt_) Compile(tokens []string) (test Test, err error) {
	return compile_equality(tokens, ">", build_gt)
}
func build_gt(val uint32, flipped bool) Test {
	return Gt{Value: val, Flipped: flipped}
}

var gte gte_

type gte_ struct{}

func (gte_) Compile(tokens []string) (test Test, err error) {
	return compile_equality(tokens, ">=", build_gte)
}
func build_gte(val uint32, flipped bool) Test {
	return GtE{Value: val, Flipped: flipped}
}

var lt lt_

type lt_ struct{}

func (lt_) Compile(tokens []string) (test Test, err error) {
	return compile_equality(tokens, "<", build_lt)
}
func build_lt(val uint32, flipped bool) Test {
	return Lt{Value: val, Flipped: flipped}
}

var lte lte_

type lte_ struct{}

func (lte_) Compile(tokens []string) (test Test, err error) {
	return compile_equality(tokens, "<=", build_lte)
}
func build_lte(val uint32, flipped bool) Test {
	return LtE{Value: val, Flipped: flipped}
}

type test_token_def struct {
	Op    string
	Token test_token
}

var precedence = []test_token_def{
	test_token_def{Op: "||", Token: or},
	test_token_def{Op: "&&", Token: and},
	test_token_def{Op: "==", Token: eq},
	test_token_def{Op: "!=", Token: neq},
	test_token_def{Op: ">=", Token: gte},
	test_token_def{Op: ">", Token: gt},
	test_token_def{Op: "<=", Token: lte},
	test_token_def{Op: "<", Token: lt},
}

type splitted struct {
	Left  []string
	Right []string
}

// Find index of token in list of tokens
func index(tokens []string, sep string) int {
	for index, token := range tokens {
		if token == sep {
			return index
		}
	}
	return -1
}

// Split a list of tokens by a token into a splitted struct holding the tokens
// before and after the token to be split by.
func split_tokens(tokens []string, sep string) (s splitted, err error) {
	index := index(tokens, sep)
	if index == -1 {
		return s, errors.New(fmt.Sprintf("'%s' not found in ['%s']", sep, strings.Join(tokens, "','")))
	}
	return splitted{
		Left:  tokens[:index],
		Right: tokens[index+1:],
	}, nil
}

// Scan a string for parenthesis
func scan(s string) []match {
	ret := []match{}
	depth := 0
	opener := 0
	for index, char := range s {
		switch char {
		case '(':
			if depth == 0 {
				opener = index
			}
			depth++
		case ')':
			depth--
			if depth == 0 {
				ret = append(ret, match{OpenPos: opener, ClosePos: index + 1})
			}
		}

	}
	return ret
}

// Split the string into tokens
func split(s string) []string {
	s = strings.Replace(s, " ", "", -1)
	if !strings.Contains(s, "(") {
		return []string{s}
	}
	last := 0
	end := len(s)
	ret := []string{}
	for _, info := range scan(s) {
		if last != info.OpenPos {
			ret = append(ret, s[last:info.OpenPos])
		}
		ret = append(ret, s[info.OpenPos:info.ClosePos])
		last = info.ClosePos
	}
	if last != end {
		ret = append(ret, s[last:])
	}
	return ret
}

// Tokenizes a string into a list of strings, tokens grouped by parenthesis are
// not split! If the string starts with ( and ends in ), those are stripped.
func tokenize(s string) []string {
	/*
		TODO: Properly detect if the string starts with a ( and ends with a )
		and that those two form a matching pair.

		Eg: (foo) -> true; (foo)(bar) -> false;
	*/
	if s[0] == '(' && s[len(s)-1] == ')' {
		s = s[1 : len(s)-1]
	}
	ret := []string{}
	for _, chunk := range split(s) {
		if len(chunk) != 0 {
			if chunk[0] == '(' && chunk[len(chunk)-1] == ')' {
				ret = append(ret, chunk)
			} else {
				for _, token := range pat.FindAllStringSubmatch(chunk, -1) {
					ret = append(ret, token[0])
				}
			}
		} else {
			fmt.Printf("Empty chunk in string '%s'\n", s)
		}
	}
	return ret
}

// Compile a string containing a plural form expression to a Expression object.
func Compile(s string) (expr Expression, err error) {
	if s == "0" {
		return Const{Value: 0}, nil
	}
	if !strings.Contains(s, "?") {
		s += "?1:0"
	}
	return compile_expression(s)
}

// Check if a token is in a slice of strings
func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// Compiles an expression (ternary or constant)
func compile_expression(s string) (expr Expression, err error) {
	tokens := tokenize(s)
	if contains(tokens, "?") {
		return ternary.Compile(tokens)
	} else {
		return const_val.Compile(tokens)
	}
}

// Compiles a test (comparison)
func compile_test(s string) (test Test, err error) {
	tokens := tokenize(s)
	for _, token_def := range precedence {
		if contains(tokens, token_def.Op) {
			return token_def.Token.Compile(tokens)
		}
	}
	return test, errors.New("Cannot compile")
}

func parse_uint32(s string) (ui uint32, err error) {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return ui, err
	} else {
		return uint32(i), nil
	}
}
