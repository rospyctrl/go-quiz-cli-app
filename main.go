package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func probPuller(fileName string) ([]problem, error) {
	// read all the problems from the quiz.csv

	// open the file
	if fObj, err := os.Open(fileName); err == nil {
		//create a new reader instance
		csvR := csv.NewReader(fObj)
		//it will need to read the file
		if cLines, err := csvR.ReadAll(); err == nil {
			//call the parser function
			return probParser(cLines), nil
		} else {
			return nil, fmt.Errorf("error in reading the csv"+"format from %s file; %s", fileName, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening %s file; %s", fileName, err.Error())
	}
}

func main() {
	//1. Input the name of the file
	fName := flag.String("f", "quiz.csv", "path of csv file1")
	//2. Set the duration of the timer
	timer := flag.Int("t", 30, "timer for the quiz")
	flag.Parse()
	//3. Pull the file (problem puller function)
	problems, error := probPuller(*fName)
	//4. Handle the error
	if error != nil {
		exit(fmt.Sprintf("Something went wrong:%s", error.Error()))
	}
	//5. Create a variable to count our correct answer
	correctAns := 0
	//6. Using the duration of the timer, we want to initialize the timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)
	//7. Loop through the problems, print the questions, we'll accept the answer
problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s=", i+1, p.q)
		//process the answer
		go func() {
			fmt.Scanf("%s", &answer)
			ansC <- answer
		}()
		select {
		case <-tObj.C:
			fmt.Printf("case 1")
			fmt.Println()
			break problemLoop
		case iAns := <-ansC:
			if iAns == p.a {
				fmt.Printf("case 2")
				correctAns++
			}
			if i == len(problems)-1 {
				fmt.Printf("case 3")
				close(ansC)
			}
		}
	}
	//8. We'll calculate and print out the result
	fmt.Printf("Your Result is %d out of %d\n", correctAns, len(problems))
	fmt.Printf("Press enter to exit")
	<-ansC
}

func probParser(lines [][]string) []problem {
	// go over the lines and parse the data with problem struct
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		// [ [70+22, 92], [0+2,2]]
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}
	return r
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

type problem struct {
	q string
	a string
}
