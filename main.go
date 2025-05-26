package main

import (
	"fmt"
	"merge/constants"
	"merge/logics"
	"time"
)

func main() {
	startTime := time.Now()
	fmt.Println(constants.Message)
	fmt.Println(constants.Empty)
	logics.MainLogic()
	endTime := time.Now()
	fmt.Printf("We took: %v ms for the job\n", endTime.Sub(startTime).Milliseconds())
}
