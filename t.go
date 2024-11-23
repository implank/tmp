package main

import "time"

func main() {
	print("t")
	time.Sleep(time.Second / time.Duration(3))
	print("a")
}
