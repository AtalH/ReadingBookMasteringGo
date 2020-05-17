package main

import (
	"fmt"
	"time"
)

//BasicUse shows use basic of date time
func BasicUse() {
	fmt.Println("epoch time:", time.Now().Unix())
	t := time.Now()
	fmt.Println(t, t.Format(time.RFC3339))
	fmt.Println(t.Weekday, t.Date, t.Month, t.Year)

	// sleep 2 seconds
	time.Sleep(time.Second * 2)

	t1 := time.Now()
	fmt.Println("time diff:", t1.Sub(t))
}

func formating() {
	t := time.Now()
	cnDateFormatPattern := "2006/01/02 15:04:05 -0700"
	fmt.Println("CN Date Formating:", t.Format(cnDateFormatPattern))
}

func main() {
	BasicUse()
	formating()
}
