package calc_utils

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestExpressionParse(t *testing.T) {
	var tests = []struct {
		expression string
		exp []string
	}{
		{"1+3", []string{"1", "+", "3"}},
		{"11+3", []string{"11", "+", "3"}},
		{"aaa", nil},
		{"312414214/5", []string{"312414214", "/", "5"}},
		{"(312414214/5)-1", []string{"(", "312414214", "/", "5", ")", "-", "1"}},
	}

	for _, e := range tests {
		res := expressionParse(e.expression)
		assert.Equal(t, res, e.exp, "TestExpressionParse")
	}
}

func TestHasHigherPrecedence(t *testing.T) {
	var tests = []struct {
		op1 string
		op2 string
		exp bool
	}{
		{"*", "-", true},
		{"/", "-", true},
		{"*", "+", true},
		{"/", "+", true},
		{"/", "/", true},
		{"+", "+", true},
		{"+", "-", true},
		{"-", "+", true},
		{"-", "-", true},
		{"-", "*", false},
		{"-", "/", false},
		{"+", "*", false},
		{"+", "/", false},
	}

	for _, e := range tests {
		res := hasHigherPrecedence(e.op1, e.op2)
		if res != e.exp {
			t.Errorf("hasHigherPrecedence(%s, %s) = %t, expected %t",
				e.op1, e.op2, res, e.exp)
		}
	}
}

func TestPerformOperation(t *testing.T) {
	var tests = []struct {
		op1 float64
		op2 float64
		operator string
		exp float64
		err error
	}{
		{1, 1, "+", 2, nil}, 
		{1, 1, "-", 0, nil}, 
		{1, 2, "*", 2, nil}, 
		{2, 1, "/", 2, nil}, 
		{2, 0, "/", 0, fmt.Errorf("Деление на ноль")}, 
		{2, 1, "|", 0, fmt.Errorf("Неподдерживаемая операция")},
	}

	for _, e := range tests {
		res, err := performOperation(e.op1, e.op2, e.operator)
		if res != e.exp || (err != nil && err.Error() != e.err.Error()) {
			t.Errorf("performOperation(%e, %e) = %s, expected %e and error %s",
				e.op1, e.op2, e.operator, e.exp, e.err)
		}
	}
}

func TestCalculateExpression(t *testing.T) {
	var tests = []struct {
		operands *[]float64
		operators *[]string
		err error
	}{
		{&[]float64{1, 2}, &[]string{"+"}, nil}, 
		{&[]float64{1, 2}, &[]string{}, fmt.Errorf("Некорректное выражение")}, 
		{&[]float64{1}, &[]string{"+"}, fmt.Errorf("Некорректное выражение")}, 
		{&[]float64{2, 0}, &[]string{"/"}, fmt.Errorf("Деление на ноль")}, 
		{&[]float64{1, 2}, &[]string{"|"}, fmt.Errorf("Неподдерживаемая операция")}, 
	}

	for _, e := range tests {
		err := calculateExpression(e.operands, e.operators)
		if (err != nil && err.Error() != e.err.Error()) {
			t.Errorf("calculateExpression expected error %s",
				e.err)
		}
	}
}

func TestEvaluateExpression(t *testing.T) {
	var tests = []struct {
		expression string
		res float64
		err error
	}{
		{"1 + 2", 3, nil}, 
		{"(1 + 2)", 3, nil}, 
		{"1 + 2 * (3 - 1)", 5, nil}, 
		{"1 / 0", 0, fmt.Errorf("Деление на ноль")}, 
		{"asdasdsa", 0, fmt.Errorf("Неправильный формат ввода!")}, 
	}

	for _, e := range tests {
		res, err := EvaluateExpression(e.expression)
		if res != e.res || (err != nil && err.Error() != e.err.Error()) {
			t.Errorf("EvaluateExpression(%s) = %e error %s, expected %e error %s",
				e.expression, res, err.Error(), e.res, e.err)
		}
	}
}
