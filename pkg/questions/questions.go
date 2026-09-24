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
func GetAnswer(question string, options []string) (index int) {
	index, len_options := -1, len(options)
	for {
		fmt.Printf("Q. %s:\n", question)
		for i, option := range options {
			fmt.Printf("%d: %s\n", i+1, option)
		}
		scanner := bufio.NewScanner(os.Stdin)
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
	return index
}
