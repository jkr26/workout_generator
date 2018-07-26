package main

import (
"bufio"
"io/ioutil"
"fmt"
"log"
"os"
"encoding/json"
)

type workoutClass struct{
	name string
	attrs []map[string][]string
}

type allWorkouts struct{
	workouts []workoutClass
}

func readFile(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
        fmt.Print(err)
    }
    return b, err
}

//writeLines writes the lines to stdout
func writeLines(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	text, err := readFile("Workouts.json")
	if err != nil {
		log.Fatalf("readFile: %s", err)
	}
//	fmt.Printf(string(text))
	var t allWorkouts
	err = json.Unmarshal(text, &t)
	if err != nil {
		log.Fatalf("Unmarshal: %s", err)
	}
	fmt.Printf("%+v\n", t)
}