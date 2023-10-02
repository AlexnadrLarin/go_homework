package uniq_utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Структура входных параметров
type Options struct {
	C              bool
	D              bool
	U              bool
	F              int
	S              int
	I              bool
	inputFileName  string
	outputFileName string
}

// Функция для проверки на вхождение элемента в slice
func entryNumber(slice []string, element string) int {
	counter := 0
    for _, a := range slice {
        if a == element {
            counter++
        }
    }
    return counter
}

// Функция для чтения входных данных
func scanner(input *os.File) []string {
	var buf []string

	inputScanner := bufio.NewScanner(input)

	if inputScanner.Err() != nil {
		return nil
	}

	for inputScanner.Scan() {
		buf = append(buf, inputScanner.Text())
	}

	return buf
}

// Функция для записи в файл
func writer(output *os.File, buf []string) {
	datawriter := bufio.NewWriter(output)

	for _, line := range buf {
		datawriter.WriteString(line + "\n")
	}

	datawriter.Flush()  
}

// Парсер входных данных
func argsParser(args []string, options Options) (Options, error) {
	for idx, argValue := range args {
		if argValue == "-c" {
			options.C = true
		} else if argValue == "-d" {
			options.D = true
		} else if argValue == "-u" {
			options.U = true
		} else if argValue == "-i" {
			options.I = true
		} else if argValue == "-f" {
			i, err := strconv.Atoi(args[idx + 1])
			if err != nil {
				return options, err
			}

			if i < 0 {
				return options, fmt.Errorf("Введёное число не может быть отрицательным.")
			}

			options.F = i
		} else if argValue == "-s" {
			i, err := strconv.Atoi(args[idx+1])
			if err != nil {
				return options, err
			}

			if i < 0 {
				return options, fmt.Errorf("Введёное число не может быть отрицательным.")
			}

			options.S = i
		} else if _, err := strconv.Atoi(argValue); err == nil && (args[idx-1] == "-f" || args[idx-1] == "-s") {
			if err != nil {
				return options, err
			}
		} else if string(argValue[0]) == "-" {
			return options, fmt.Errorf("Такого параметра не существует.")
		} else {
			file, err := os.Open(argValue)
			if err == nil {
				if options.inputFileName == "" {
					options.inputFileName = argValue
				} else if options.outputFileName == "" {
					options.outputFileName = argValue
				}

				file.Close()
			} else {
				return options, err
			}
		}
	}

	return options, nil
}

// Проверка на наличие файла со входными данными
func inputManager(options Options) []string {
	if options.inputFileName != "" {
		inputFile, err := os.Open(options.inputFileName)
		defer inputFile.Close()

		if err == nil {
			return scanner(inputFile)
		}

	} else {
		return scanner(os.Stdin)
	}

	return nil
}

// Подготовка и изменение массива входных данных в зависимости от условий
func uniqManager(options Options, buf []string) []string {
	if options.F >= 0 {
		for idx := range buf {
			splitBuf := strings.Split(buf[idx], " ")
			if len(splitBuf) > options.F {
				buf[idx] = strings.Join(splitBuf[options.F:], " ")
			}
		}
	} 

	if options.S >= 0 {
		for idx := range buf {
			if len(buf[idx]) > options.S {
				buf[idx] = buf[idx][options.S:]
			}
			
		}

	}

	if options.I {
		for idx := range buf {
			buf[idx] = strings.ToLower(buf[idx])
		}
	}

	return buf
}

// Составление конечного массива по алгоритму в зависимости от условий
func uniqStringsChecker(options Options, buf []string, edittedBuf []string) []string {
	if (options.C && options.D) ||
	   (options.D && options.U) ||
	   (options.C && options.U) {
		return nil
	}

	if options.C {
		var uniqStringsBuf []string

		counter := 1

		edittedBuf = append(edittedBuf, "")
		prevoiusLine := edittedBuf[0]
		for idx, value := range edittedBuf[1:len(edittedBuf) + 1] {
			if value != prevoiusLine {
				uniqStringsBuf = append(uniqStringsBuf, strconv.Itoa(counter) + " " + buf[idx])
				counter = 1
			} else {
				counter++
			}
			prevoiusLine = value
		}

		return uniqStringsBuf
	} else if options.D {
		var repeatStringsBuf []string

		prevoiusLine := edittedBuf[0]
		counter := 1
		edittedBuf = append(edittedBuf, "")
		for idx, value := range edittedBuf[1:len(edittedBuf) + 1] {
			if value != prevoiusLine  {
				if counter > 0 {
					repeatStringsBuf = append(repeatStringsBuf, buf[idx])
				}
				counter = 0
			} else {
				counter++
			}
			prevoiusLine = value
		}

		return repeatStringsBuf
	} else if options.U {
		var uniqStringsBuf []string
		var repeatStringsBuf []string
		var notRepeatStringsBuf []string

		prevoiusLine := ""
		for idx, value := range edittedBuf {
			if value != prevoiusLine {
				uniqStringsBuf = append(uniqStringsBuf, buf[idx])
			} else {
				if entryNumber(repeatStringsBuf, buf[idx]) == 0 {
					repeatStringsBuf = append(repeatStringsBuf, buf[idx])
				}
			}
			prevoiusLine = value
		}
		
		for idx, value := range uniqStringsBuf {
			if entryNumber(repeatStringsBuf, value) == 0 {
				notRepeatStringsBuf = append(notRepeatStringsBuf, uniqStringsBuf[idx])
			}
		}

		return notRepeatStringsBuf
	}

	var uniqStringsBuf []string

	prevoiusLine := ""
	for idx, value := range edittedBuf {
		if value != prevoiusLine {
			uniqStringsBuf = append(uniqStringsBuf, buf[idx])
		}
		prevoiusLine = value
	}

	return uniqStringsBuf
}

// Основная функция
func Uniq() {
	var optionsInitial Options = Options{
		C:              false,
		D:              false,
		U:              false,
		F:              -1,
		S:              -1,
		I:              false,
		inputFileName:  "",
		outputFileName: "",
	}

	var buf []string
	var edittedBuf []string

	options, err := argsParser(os.Args[1:], optionsInitial)
	if err != nil {
		fmt.Println("Ошибка:\n", err)
		return
	}

	buf = inputManager(options)
	if buf == nil {
		fmt.Println("Вы ввели пустые данные!")
		return
	}

	edittedBuf = append(edittedBuf, buf...)
	uniqManager(options, edittedBuf)
	resultBuf := uniqStringsChecker(options, buf, edittedBuf)

	if resultBuf == nil {
		fmt.Println("Неправильный формат ввода!\n", 
					"Параметры c, d, u взаимозаменяемы, поэтому их использование вместе не имеет никакого смысла.\n", 
					"Использование утилиты uniq:\nuniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]")
		return
	}

	if options.outputFileName != "" {
		// Вывод результата в файл
		outputFile, err := os.OpenFile(options.outputFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		defer outputFile.Close()

		if err != nil {
			fmt.Println("Ошибка:\n", err)
			return
		} 
		
		writer(outputFile, resultBuf)
	} else {
		// Вывод результата в StdOut
		for _, value := range resultBuf {
			fmt.Println(value)
		}
	}

}