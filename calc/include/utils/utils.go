package calc_utils

import (
	"fmt"
	"strconv"
	"strings"
)

func expressionParse(expression string) ([]string, error) {
	elements := []string{}
	currentElement := ""
   
	for _, value := range expression {
	 	if isOperator(value) || value == '(' || value == ')' {
	  		if currentElement != "" {
	   			elements = append(elements, currentElement)
	  		}

			currentElement = ""
			elements = append(elements, string(value))
	 	} else {
			_, err := strconv.ParseFloat(string(value), 64)
			if err != nil {
				return nil, err
			}
	  		currentElement += string(value)
		}
	}
   
	if currentElement != "" {
		elements = append(elements, currentElement)
	}
   
	return elements, nil
}
   
func isOperator(value rune) bool {
	return value == '+' || value == '-' || value == '*' || value == '/' 
}

// Основная функция для анализа выражения
func EvaluateExpression(expression string) (float64, error) {
 	expression = strings.ReplaceAll(expression, " ", "") 

	var operandsStack []float64
	var operatorsStack []string

 	elements, err := expressionParse(expression)

	if err != nil {
		return 0, err
	}

 	for _, element := range elements {
  		switch element {
  		case "+", "-", "*", "/":
   			for len(operatorsStack) > 0 && hasHigherPrecedence(operatorsStack[len(operatorsStack)-1], element) {
    			err := calculateExpression(&operandsStack, &operatorsStack)
				if err != nil {
					return 0, err
				}
   			}
   			operatorsStack = append(operatorsStack, element)
  		case "(":
   			operatorsStack = append(operatorsStack, element)
  		case ")":
			for len(operatorsStack) > 0 && operatorsStack[len(operatorsStack)-1] != "(" {
				err := calculateExpression(&operandsStack, &operatorsStack)
				if err != nil {
					return 0, err
				}
			}

			if len(operatorsStack) > 0 && operatorsStack[len(operatorsStack)-1] == "(" {
				operatorsStack = operatorsStack[:len(operatorsStack)-1]
			}
		default:
			number, err := strconv.ParseFloat(element, 64)
			if err != nil {
				return 0, err
			}
			operandsStack = append(operandsStack, number)
		}
	}

	for len(operatorsStack) > 0 {
		err := calculateExpression(&operandsStack, &operatorsStack)
		if err != nil {
			return 0, err
		}
	}

	if len(operandsStack) > 0 {
		return operandsStack[0], nil
	}
	return 0, fmt.Errorf("Ввыражение не может быть вычислено")
}

// Функция для корректной передачи операторов и операндов из стека в функцию подсчёта
func calculateExpression(operands *[]float64, operators *[]string) error {
	if len(*operands) < 2 || len(*operators) == 0 {
	 	return fmt.Errorf("Некорректное выражение")
	}
   
	operator := (*operators)[len(*operators)-1]
	(*operators) = (*operators)[:len(*operators)-1]
   
	operand2 := (*operands)[len(*operands)-1]
	(*operands) = (*operands)[:len(*operands)-1]
   
	operand1 := (*operands)[len(*operands)-1]
	(*operands) = (*operands)[:len(*operands)-1]
   
	result, err := performOperation(operand1, operand2, operator)
	if err != nil {
	 	return err
	}
   
	(*operands) = append((*operands), result)

	return nil
}

// Функция подсчёта 
func performOperation(operand1, operand2 float64, operator string) (float64, error) {
 	switch operator {
 	case "+":
  		return operand1 + operand2, nil
 	case "-":
  		return operand1 - operand2, nil
 	case "*":
  		return operand1 * operand2, nil
 	case "/":
		if operand2 == 0 {
			return 0, fmt.Errorf("Деление на ноль")
		}
  		return operand1 / operand2, nil
 	default:
  		return 0, fmt.Errorf("Неподдерживаемая операция")
 	}
}

// Функция для определение приоритета операции
func hasHigherPrecedence(operand1, operand2 string) bool {
 	if operand1 == "*" || operand1 == "/" {
  		return true
 	}
 	if operand1 == "+" || operand1 == "-" {
  		return operand2 	== "+" || operand2 == "-"
 	}
 	return false
}