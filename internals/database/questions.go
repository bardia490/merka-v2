package database

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
// NOTE: this is not complete yet (add looping for wrong answers)
func GetAnswer(question string, options []string) (index int) {
	index = 0
	var answer string
	fmt.Printf("Q. %s:\n", question)
	for i, option := range options {
		fmt.Printf("%d: %s\n", i, option)
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	if result, err := strconv.Atoi(input); err == nil {
		index = result
	} else {
		answer = input
	}
	in_slice := slices.Index(options, answer)
	if in_slice != -1 {
		//fmt.Println(options[in_slice])
		return in_slice
	} else {
		fmt.Println(options[index])
	}
	return index
}
