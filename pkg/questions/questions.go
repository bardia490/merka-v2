package questions

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// precompiled regex patterns for getting the time
var (
	hoursRegex   = regexp.MustCompile(`(?i)^\s*(\d+(?:\.\d+)?)\s*h\s*$`)               // Matches: "2.5h", "3 h", "0.5H" (with (?i) case-insensitivity)
	minutesRegex = regexp.MustCompile(`(?i)^\s*(\d+(?:\.\d+)?)\s*m\s*$`)               // Matches: "45m", "30.5 m", "10M"
	clockRegex   = regexp.MustCompile(`^\s*(\d+(?:\.\d+)?)\s*:\s*(\d+(?:\.\d+)?)\s*$`) // Matches: "1:30", "2.5:15", "01 : 45.5"
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
		fmt.Print("> ")
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

type prompt_condition[N Number] func(N) bool
type prompt_string_condition func(string) bool

// displays the [question] and waits for the user to answer
// [condition] is a function that takes in an [Integer] and
// returns true if the constraints (such as answer >= 0) were met
// a description on what constraints should be followed by the answer
// can be sent to the function as [condition_description] argument
func PromptInteger[I Integer](question string, condition_description string, condition prompt_condition[I]) I {
	var answer I
	for {
		fmt.Printf("Q. %s\n", question)
		fmt.Printf("Condition. %s\n> ", condition_description)
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
		fmt.Printf("Condition. %s\n> ", condition_description)
		fmt.Scanf("%f", &answer)
		if condition(answer) {
			break
		}
		fmt.Println("the answer did not satisfy the condition")
	}
	return answer
}

// displays the [question] and waits for the user to answer
// [condition] is a function that takes in an [String] and
// returns true if the constraints (such as len(answer) >= 0) were met
// a description on what constraints should be followed by the answer
// can be sent to the function as [condition_description] argument
func PromptString(question string, condition_description string, condition prompt_string_condition) string {
	var answer string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Q. %s\n", question)
		fmt.Printf("Condition. %s\n> ", condition_description)
		scanner.Scan()
		answer = strings.TrimSpace(scanner.Text())
		if condition(answer) {
			break
		}
		fmt.Println("the answer did not satisfy the condition")
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("this is probably not important but please report this: %s\n", err.Error())
	}
	return answer
}

// it can recieve time in on of the following formats:
// hours h or hours H(e.g. 2.5h, 3 h, 0.5H)
// minutes m or minutes M (e.g. 45m, 30.5 m, 10M)
// hours:minutes or hours : minutes (e.g. 1:30, 2.5:15, 01 : 45.5)
// returns the time in minutes
func PromptTime() (minutes float64) {
	minutes = 0
	scanner := bufio.NewScanner(os.Stdin)
	var answer string

	for {
		fmt.Printf("Q. Please Enter the time\n")
		fmt.Printf("Condition. The following formats are allowed:\n hours h or hours H(e.g. 2.5h, 3 h, 0.5H)\n minutes m or minutes M (e.g. 45m, 30.5 m, 10M)\n hours:minutes or hours : minutes (e.g. 1:30, 2.5:15, 01 : 45.5)\n> ")
		scanner.Scan()
		answer = strings.TrimSpace(scanner.Text())
		// 1. Check "number h" (Hours)
		if matches := hoursRegex.FindStringSubmatch(answer); len(matches) > 1 {
			hours, err := strconv.ParseFloat(matches[1], 64)
			if err == nil {
				minutes = hours * 60
				break
			}
		}

		// 2. Check "number m" (Minutes)
		if matches := minutesRegex.FindStringSubmatch(answer); len(matches) > 1 {
			mins, err := strconv.ParseFloat(matches[1], 64)
			if err == nil {
				minutes = mins
				break
			}
		}

		// 3. Check "number1 : number2" (Hours : Minutes)
		if matches := clockRegex.FindStringSubmatch(answer); len(matches) > 2 {
			hours, err := strconv.ParseFloat(matches[1], 64)
			if err != nil {
				continue
			}
			mins, err := strconv.ParseFloat(matches[2], 64)
			if err != nil {
				continue
			}

			minutes = hours*60 + mins
		}
		fmt.Println("the time was not in correct format, please try again")
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("this is probably not important but please report this: %s\n", err.Error())
	}
	return minutes
}
