package uniq_utils

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
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
			fmt.Errorf("Неправильный формат ввода!\n" + 
								   "Параметры c, d, u взаимозаменяемы, поэтому их использование вместе не имеет никакого смысла.\n" +
								   "Использование утилиты uniq:\nuniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]"),
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
		res, err := parseArgs(e.args, e.options)
		if res != e.exp || (err != nil && err.Error() != e.err.Error()){
			t.Errorf("Test parseArgs error %s expected %s",
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

func TestEditBuf(t *testing.T) {
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
		res := editBuf(e.options, e.buf)
		assert.Equal(t, res, e.exp, "TestEditBuf")
	}
}

func TestFindRepeatedLines(t *testing.T) {

	var tests = []struct {
		buf []string
		edittedBuf []string
		exp []string
	}{
		{
			[]string{"aaa", "aaa", "b"}, 
			[]string{"aaa", "aaa", "b"}, 
			[]string{"aaa"},
		},
		{
			[]string{"baa", "baa"}, 
			[]string{"aa", "aa"}, 
			[]string{"baa"},
		},
		{
			[]string{"aaa", "baa", "aaa"}, 
			[]string{"aaa", "baa", "aaa"}, 
			nil,
		},

	}

	for _, e := range tests {
		res := findRepeatedLines(e.buf, e.edittedBuf)
		assert.Equal(t, res, e.exp, "TestFindRepeatedLines")
	}
}

func TestFindUniqLines(t *testing.T) {

	var tests = []struct {
		buf []string
		edittedBuf []string
		exp []string
	}{
		{
			[]string{"aaa", "aaa", "b"}, 
			[]string{"aaa", "aaa", "b"}, 
			[]string{"aaa", "b"},
		},
		{
			[]string{"baa", "baa"}, 
			[]string{"aa", "aa"}, 
			[]string{"baa"},
		},
		{
			[]string{"aaa", "baa", "aaa"}, 
			[]string{"aaa", "baa", "aaa"}, 
			[]string{"aaa", "baa", "aaa"},
		},

	}

	for _, e := range tests {
		res := findUniqLines(e.buf, e.edittedBuf)
		assert.Equal(t, res, e.exp, "TestFindUniqLines")
	}
}

func TestFindNotRepeatedLines(t *testing.T) {

	var tests = []struct {
		buf []string
		edittedBuf []string
		exp []string
	}{
		{
			[]string{"aaa", "aaa", "b"}, 
			[]string{"aaa", "aaa", "b"}, 
			[]string{"b"},
		},
		{
			[]string{"baa", "baa"}, 
			[]string{"aa", "aa"}, 
			nil,
		},
		{
			[]string{"aaa", "baa", "aaa"}, 
			[]string{"aaa", "baa", "aaa"}, 
			[]string{"aaa", "baa", "aaa"},
		},

	}

	for _, e := range tests {
		res := findNotRepeatedLines(e.buf, e.edittedBuf)
		assert.Equal(t, res, e.exp, "TestFindNotRepeatedLines")
	}
}

func TestCountOccurrencesNumber(t *testing.T) {

	var tests = []struct {
		buf []string
		edittedBuf []string
		exp []string
	}{
		{
			[]string{"aaa", "aaa"}, 
			[]string{"aaa", "aaa"}, 
			[]string{"2 aaa"},
		},
		{
			[]string{"baa", "baa"}, 
			[]string{"aa", "aa"}, 
			[]string{"2 baa"},
		},
		{
			[]string{"baa", "baa", "aab"}, 
			[]string{"aa", "aa", "ab"}, 
			[]string{"2 baa", "1 aab"},
		},

	}

	for _, e := range tests {
		res := countOccurrencesNumber(e.buf, e.edittedBuf)
		assert.Equal(t, res, e.exp, "TestCountOccurrencesNumber")
	}
}

func TestCheckUniqLines(t *testing.T) {
	var tests = []struct {
		options Options
		buf []string
		edittedBuf []string
		exp []string
	}{
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
		res := checkUniqLines(e.options, e.buf, e.edittedBuf)
		assert.Equal(t, res, e.exp, "TestCheckUniqLines")
	}
}
