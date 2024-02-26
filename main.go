package main // Package declaration for the main package

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	// Define command-line flags
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 10, "the time limit for the quiz in seconds")
	flag.Parse() // Parse the command-line flags

	// Open and read the CSV file
	file, err := os.Open(*csvFilename) // Open the CSV file
	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file: %s\n", *csvFilename)) // Handle error if failed to open file
	}
	r := csv.NewReader(file)  // Create a new CSV reader
	lines, err := r.ReadAll() // Read all lines from the CSV file
	if err != nil {
		exit("Failed to parse the provided csv file.")
	}

	// Parse CSV lines into problems
	problems := parseLines(lines)                                   // Parse CSV lines into problems struct
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second) // Create a new timer

	correct := 0 // Initialize correct answers counter

problemloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		answerCh := make(chan string) // Create a channel for receiving answers

		go func() { // Start a goroutine to handle user input
			var answer string
			fmt.Scanf("%s\n", &answer) // Read user input
			answerCh <- answer         // Send user input to channel
		}()

		select {
		case <-timer.C: // If time limit is reached
			fmt.Println()
			break problemloop // Exit loop
		case answer := <-answerCh: // Receive answer from channel
			if answer == p.a { // Check if answer is correct
				correct++ // Increment correct answers counter
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems)) // Print final score
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines)) // Initialize slice of problems
	for i, line := range lines {
		ret[i] = problem{ // Create new problem struct and assign values from CSV line
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret // Return slice of problems
}

type problem struct {
	q string // Struct definition for a problem with question and answer fields
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
