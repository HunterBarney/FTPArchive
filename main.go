package main

import (
	"flag"
	"fmt"
)

func main() {
	profilePath := flag.String("profile", "", "Path to profile file")

	flag.Parse()

	fmt.Println("profilePath:", *profilePath)
}
