package main

import (
	"fmt"
	"regexp"
)

func regexDemo() {
	regex := regexp.MustCompile("[1-5a-c]")
	test := "15isbd4"
	if regex.MatchString(test) {
		match := regex.FindStringSubmatch(test)
		fmt.Println("match:", match)
	}
	match := regex.FindAllString(test, -1)
	fmt.Println("match:", match)
}

func main() {
	regexDemo()
}
