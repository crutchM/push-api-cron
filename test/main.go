package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Date(2021, time.Month(2), 21, 10, 00, 00, 0, time.UTC)
	if time.Now().UTC().Hour() > t.Hour() {
		fmt.Println("уже больше 10 часов")
	}
	fmt.Println(t.Hour())
	fmt.Println(time.Now().UTC().Hour())

}
