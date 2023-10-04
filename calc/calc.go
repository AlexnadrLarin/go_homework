package main

import (
	"os"
	"fmt"
	"calc_utils/include/utils"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Неправильный формат ввода!\n" +
					`Пример использования программы: go run calc.go "(1+2)-3"`)
		return
	}

	result, err := calc_utils.EvaluateExpression(os.Args[1])
 		if err != nil {
  			fmt.Println("Ошибка при вычислении выражения:", err)
  			return
 		}

 		fmt.Println("Результат:", result)

}
