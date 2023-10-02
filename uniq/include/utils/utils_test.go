package uniq_utils

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestArgsParser(t *testing.T) {
	var OptionsOriginalState = Options{
		C:              false,
		D:              false,
		U:              false,
		F:              -1,
		S:              -1,
		I:              false,
		inputFileName:  "",
		outputFileName: "",
	}

	var tests = []struct {
		args []string
		options Options
		exp Options
		err error
	}{
		{
			[]string{"-c", "-d", "-u", "-f", "1", "-s", "1", "input.txt", "output.txt", "-i"}, 
			OptionsOriginalState, 
			Options {
				C:              true,
				D:              true,
				U:              true,
				F:              1,
				S:              1,
				I:              true,
				inputFileName:  "input.txt",
				outputFileName: "output.txt",
			}, 
			nil,
		},
		{
			[]string{"-f", "safffsa"}, 
			OptionsOriginalState, 
			OptionsOriginalState,
			fmt.Errorf(`strconv.Atoi: parsing "safffsa": invalid syntax`),
		},
		{
			[]string{"-s", "safffsa"}, 
			OptionsOriginalState, 
			OptionsOriginalState,
			fmt.Errorf(`strconv.Atoi: parsing "safffsa": invalid syntax`),
		},
		{
			[]string{"-s", "-1"}, 
			OptionsOriginalState, 
			OptionsOriginalState,
			fmt.Errorf("Введёное число не может быть отрицательным."),
		},
		{
			[]string{"-f", "-1"}, 
			OptionsOriginalState, 
			OptionsOriginalState,
			fmt.Errorf("Введёное число не может быть отрицательным."),
		},
		{
			[]string{"-ffsdfsdf"}, 
			OptionsOriginalState, 
			OptionsOriginalState,
			fmt.Errorf("Такого параметра не существует."),
		},
		{
			[]string{"int.txt"}, 
			OptionsOriginalState, 
			OptionsOriginalState,
			fmt.Errorf("open int.txt: no such file or directory"),
		},
		
	}

	for _, e := range tests {
		res, err := argsParser(e.args, e.options)
		if res != e.exp || (err != nil && err.Error() != e.err.Error()){
			t.Errorf("Test argsParser error %s expected %s",
					err.Error(), e.err.Error())
		}
	}
}

func TestEntryNumber(t *testing.T) {
	var tests = []struct {
		slice []string
		element string
		exp int
	}{
		{[]string{"1", "2", "3"}, "1", 1},
		{[]string{"1", "1", "3"}, "1", 2},
		{[]string{"1", "2", "3"}, "0", 0},
	}

	for _, e := range tests {
		res := entryNumber(e.slice, e.element)
		if res != e.exp {
			t.Errorf("Test entryNumber = %s expected %d",
					e.element, e.exp)
		}
	}
}

func TestUniqManager(t *testing.T) {
	var tests = []struct {
		options Options
		buf []string
		exp []string
	}{
		{
			Options {
				C:              false,
				D:              false,
				U:              false,
				F:              -1,
				S:              1,
				I:              false,
				inputFileName:  "",
				outputFileName: "",
			}, 
			[]string{"aaa", "aaa"}, 
			[]string{"aa", "aa"}, 
		},
		{
			Options {
				C:              false,
				D:              false,
				U:              false,
				F:              1,
				S:              -1,
				I:              false,
				inputFileName:  "",
				outputFileName: "",
			}, 
			[]string{"aaa aaa", "aaa aaa"}, 
			[]string{"aaa", "aaa"}, 
		},
		{
			Options {
				C:              false,
				D:              false,
				U:              false,
				F:              1,
				S:              1,
				I:              false,
				inputFileName:  "",
				outputFileName: "",
			}, 
			[]string{"aaa aaa", "aaa aaa"}, 
			[]string{"aa", "aa"}, 
		},
		{
			Options {
				C:              false,
				D:              false,
				U:              false,
				F:              1,
				S:              1,
				I:              true,
				inputFileName:  "",
				outputFileName: "",
			}, 
			[]string{"BAA BAA", "BAA BAA"}, 
			[]string{"aa", "aa"}, 
		},
	}

	for _, e := range tests {
		res := uniqManager(e.options, e.buf)
		assert.Equal(t, res, e.exp, "TestUniqManager")
	}
}

func TestUniqStringsChecker(t *testing.T) {

	var tests = []struct {
		options Options
		buf []string
		edittedBuf []string
		exp []string
	}{
		{
			Options {
				C:              true,
				D:              true,
				U:              true,
				F:              -1,
				S:              -1,
				I:              false,
				inputFileName:  "",
				outputFileName: "",
			}, 
			[]string{}, 
			[]string{}, 
			nil,
		},
		{
			Options {
				C:              false,
				D:              false,
				U:              false,
				F:              -1,
				S:              -1,
				I:              false,
				inputFileName:  "",
				outputFileName: "",
			}, 
			[]string{"aaa", "aaa"}, 
			[]string{"aaa", "aaa"}, 
			[]string{"aaa"},
		},
		{
			Options {
				C:              false,
				D:              false,
				U:              false,
				F:              -1,
				S:              -1,
				I:              false,
				inputFileName:  "",
				outputFileName: "",
			}, 
			[]string{"baa", "baa"}, 
			[]string{"aa", "aa"}, 
			[]string{"baa"},
		},
		{
			Options {
				C:              true,
				D:              false,
				U:              false,
				F:              -1,
				S:              -1,
				I:              false,
				inputFileName:  "",
				outputFileName: "",
			}, 
			[]string{"aaa", "aaa"}, 
			[]string{"aaa", "aaa"}, 
			[]string{"2 aaa"},
		},
		{
			Options {
				C:              false,
				D:              true,
				U:              false,
				F:              -1,
				S:              -1,
				I:              false,
				inputFileName:  "",
				outputFileName: "",
			}, 
			[]string{"aaa", "aaa", "bbb"}, 
			[]string{"aaa", "aaa", "bbb"}, 
			[]string{"aaa"},
		},
		{
			Options {
				C:              false,
				D:              false,
				U:              true,
				F:              -1,
				S:              -1,
				I:              false,
				inputFileName:  "",
				outputFileName: "",
			}, 
			[]string{"aaa", "aaa", "bbb"}, 
			[]string{"aaa", "aaa", "bbb"}, 
			[]string{"bbb"},
		},
	}

	for _, e := range tests {
		res := uniqStringsChecker(e.options, e.buf, e.edittedBuf)
		assert.Equal(t, res, e.exp, "TestUniqStringsChecker")
	}
}
