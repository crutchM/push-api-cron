package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().UTC().Hour() + 5)
}
