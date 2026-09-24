package questions

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

// uses a set of predfefiend [options] to display the [question]
// if the user answer is not part of the options or the index is out of bounds
// it will prompt the user again for an answer
func ChooseOption(question string, options []string) (index int) {
	index, len_options := -1, len(options)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Q. %s:\n", question)
		for i, option := range options {
			fmt.Printf("%d: %s\n", i+1, option)
		}
		scanner.Scan()
		answer := strings.TrimSpace(scanner.Text())
		if result, err := strconv.Atoi(answer); err == nil { // first try to convert the answer directly to index
			index = result - 1
		} else {
			index = slices.Index(options, answer) // try to find the answer in options if it was found
		}
		if index >= 0 && index < len_options { // make sure the index is correct
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("this is probably not important but please report this: %s\n", err.Error())
	}
	return index
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
type Float interface {
	~float32 | ~float64
}
type Number interface {
	Integer | Float
}

type prompt_condition[T Number] func(T) bool

// displays the [question] and waits for the user to answer
// [condition] is a function that takes in an [Integer] and
// returns true if the constraints (such as answer >= 0) were met
// a description on what constraints should be followed by the answer
// can be sent to the function as [condition_description] argument
func PromptInteger[I Integer](question string, condition_description string, condition prompt_condition[I]) I {
	var answer I
	for {
		fmt.Printf("Q. %s\n", question)
		fmt.Printf("Condition. %s\n", condition_description)
		fmt.Scanf("%d", &answer)
		if condition(answer) {
			break
		}
		fmt.Println("the answer did not satisfy the condition")
	}
	return answer
}

// displays the [question] and waits for the user to answer
// [condition] is a function that takes in an [Float] and
// returns true if the constraints (such as answer >= 0) were met
// a description on what constraints should be followed by the answer
// can be sent to the function as [condition_description] argument
func PromptFloat[F Float](question string, condition_description string, condition prompt_condition[F]) F {
	var answer F
	for {
		fmt.Printf("Q. %s\n", question)
		fmt.Printf("Condition. %s\n", condition_description)
		fmt.Scanf("%f", &answer)
		if condition(answer) {
			break
		}
		fmt.Println("the answer did not satisfy the condition")
	}
	return answer
}
