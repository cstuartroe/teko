package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cstuartroe/teko/src/checker"
	"github.com/cstuartroe/teko/src/interpreter"
	"github.com/cstuartroe/teko/src/lexparse"
	"github.com/cstuartroe/teko/src/shared"
)

const TEST_DIR string = "./tests/"

type TestOutcome int

const (
	test_errored = iota
	test_passed
	test_failed
	test_skipped
)

const colorReset string = "\033[0m"
const colorRed string = "\033[31m"
const colorGreen string = "\033[32m"
const colorYellow string = "\033[33m"
const colorBlue string = "\033[34m"

func printRed(s string) {
	fmt.Println(colorRed + s + colorReset)
}

func printYellow(s string) {
	fmt.Println(colorYellow + s + colorReset)
}

func printGreen(s string) {
	fmt.Println(colorGreen + s + colorReset)
}

func printBlue(s string) {
	fmt.Println(colorBlue + s + colorReset)
}

func RunTest(filepath string, test_outcome *TestOutcome) {
	program_log := new(bytes.Buffer)
	error_log := new(bytes.Buffer)

	shared.SetErrorDest(error_log)
	shared.SetPrintDest(program_log)

	var expected_output string
	var expected_error string

	defer func() {
		if e := recover(); e == shared.TekoErrorMessage {
			actual_error := error_log.String()

			if expected_error == actual_error {
				printGreen(fmt.Sprintf("%s passed\n", filepath))
				*test_outcome = test_passed
			} else {
				printRed(fmt.Sprintf("An error occurred while executing %s:\n", filepath))
				fmt.Println(actual_error)

				if expected_error != "" {
					printRed("Expected error:\n")
					fmt.Println(expected_error)
				}

				*test_outcome = test_errored
			}
		} else if e != nil {
			panic(e)
		}
	}()

	lexer := lexparse.FromFile(filepath)
	tokens := []lexparse.Token{}

	for lexer.HasMore() {
		t := lexer.GrabToken()

		if t.TType == lexparse.BlockCommentT {
			comment_body := string(t.Value)
			comment_body = strings.TrimSpace(comment_body[2 : len(comment_body)-2])

			if len(comment_body) >= 12 && comment_body[0:12] == "EXPECT ERROR" {
				expected_error = strings.TrimSpace(comment_body[12:]) + "\n"
			} else if len(comment_body) >= 6 && comment_body[0:6] == "EXPECT" {
				expected_output = strings.TrimSpace(comment_body[6:]) + "\n"
			} else if len(comment_body) >= 4 && comment_body[0:4] == "SKIP" {
				*test_outcome = test_skipped
				printBlue(fmt.Sprintf("%s skipped\n", filepath))
				return
			}
		}

		tokens = append(tokens, t)
	}

	p := lexparse.Parser{}
	p.Parse(tokens, true)

	c := checker.NewBaseChecker()
	c.CheckTree(p.Codeblock)

	i := interpreter.New(&interpreter.StdLibModule)
	i.Execute(p.Codeblock)

	actual_output := program_log.String()

	if actual_output == expected_output {
		printGreen(fmt.Sprintf("%s passed\n", filepath))
		*test_outcome = test_passed
	} else {
		printYellow(fmt.Sprintf("%s failed\n", filepath))

		if expected_error != "" {
			printYellow("Expected error:")
			fmt.Println(expected_error)
			printYellow("No error occurred.")
		} else if expected_output != "" {
			printYellow("Expected output:")
			fmt.Println(expected_output)
			printYellow("Actual output:")
			fmt.Println(actual_output)
		} else {
			printYellow("No output expected\n")
			printYellow("Got output:")
			fmt.Println(actual_output)
		}

		*test_outcome = test_failed
	}
}

func RunTests() {
	files, err := ioutil.ReadDir(TEST_DIR)

	if err != nil {
		fmt.Println("Error reading " + TEST_DIR + " - make sure there is a directory with that name and it is accessible")
		os.Exit(1)
	}

	total_tests := len(files)
	passed_tests := 0
	failed_tests := 0
	errored_tests := 0
	skipped_tests := 0

	for _, file := range files {
		full_path := filepath.Join(TEST_DIR, file.Name())
		var outcome *TestOutcome = new(TestOutcome)
		RunTest(full_path, outcome)

		switch *outcome {
		case test_errored:
			errored_tests++
		case test_passed:
			passed_tests++
		case test_failed:
			failed_tests++
		case test_skipped:
			skipped_tests++
		}
	}

	fmt.Printf("\n%2d tests run\n", total_tests-skipped_tests)
	printBlue(fmt.Sprintf("%2d tests skipped ", skipped_tests) + strings.Repeat("•", skipped_tests))
	printGreen(fmt.Sprintf("%2d tests passed  ", passed_tests) + strings.Repeat("•", passed_tests))
	printYellow(fmt.Sprintf("%2d tests failed  ", failed_tests) + strings.Repeat("•", failed_tests))
	printRed(fmt.Sprintf("%2d tests errored ", errored_tests) + strings.Repeat("•", errored_tests))
}
