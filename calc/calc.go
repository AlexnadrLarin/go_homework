package main

import (
	"os"
	"fmt"
	"calc_utils/include/utils"
)

func main() {
	if len(os.Args) == 2 {
		result, err := calc_utils.EvaluateExpression(os.Args[1])
 		if err != nil {
			fmt.Println()
  			fmt.Println("Ошибка при вычислении выражения:", err)
  			return
 		}

 		fmt.Println("Результат:", result)
	}

}
