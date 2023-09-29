package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

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
			i, err := strconv.Atoi(args[idx+1])
			if err != nil {
				return options, err
			}

			if i < 0 {
				return options, errors.New("Введёное число не может быть отрицательным.\n")
			}

			options.F = i
		} else if argValue == "-s" {
			i, err := strconv.Atoi(args[idx+1])
			if err != nil {
				return options, err
			}

			if i < 0 {
				return options, errors.New("Введёное число не может быть отрицательным.\n")
			}

			options.S = i
		} else if _, err := strconv.Atoi(argValue); err == nil && (args[idx-1] == "-f" || args[idx-1] == "-s") {
			if err != nil {
				return options, err
			}
		} else if string(argValue[0]) == "-" {
			return options, errors.New("Такого параметра не существует.\n")
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

func main() {
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

	options, err := argsParser(os.Args[1:], optionsInitial)
	if err == nil {
		fmt.Println(inputManager(options))
	} else {
		fmt.Println(err)
	}
}
