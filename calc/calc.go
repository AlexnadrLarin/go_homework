package main

import (
	
	"os"
	"fmt"
	"utils/include/utils"
)

func main() {
	if len(os.Args) == 2 {
		result, err := utils.EvaluateExpression(os.Args[1])
 		if err != nil {
			fmt.Println()
  			fmt.Println("Ошибка при вычислении выражения:", err)
  			return
 		}

 		fmt.Println("Результат:", result)
	}

}
