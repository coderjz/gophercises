package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

func main() {
	//TODO: Flag
	var filePath string
	var quizTimeSec int
	flag.StringVar(&filePath, "filepath", "questions.csv", "location for the CSV file with the questions")
	flag.IntVar(&quizTimeSec, "quiztime", 30, "time to take the quiz (in seconds)")
	flag.Parse()

	records, err := getRecordsFromCSVFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, r := range records {
		if len(r) != 2 {
			fmt.Printf("Expected all rows to have two fields.  Record has %d fields %+v\n: ", len(r), r)
			return
		}
	}

	quiz := newQuiz(records, quizTimeSec)
	quiz.askAllQuestions()
}

func getRecordsFromCSVFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error opening file %s: %v", filePath, err)
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error reading file %s in CSV format: %v", filePath, err)
	}
	return records, nil
}

type quiz struct {
	questions       [][]string
	numCorrect      int32
	allowedQuizTime time.Duration
}

func newQuiz(questions [][]string, quizTimeSec int) *quiz {
	return &quiz{
		questions:       questions,
		allowedQuizTime: time.Duration(quizTimeSec) * time.Second,
		numCorrect:      0,
	}
}

func (q *quiz) askAllQuestions() error {
	doneChan := make(chan struct{})
	numQuestions := len(q.questions)
	go func() {
		for index := range q.questions {
			err := q.askQuestion(index)
			if err != nil {
				fmt.Println("Error asking questions: ", err)
				doneChan <- struct{}{}
			}
		}
		doneChan <- struct{}{}
	}()

	select {
	case <-doneChan:
		break
	case <-time.After(q.allowedQuizTime):
		break
	}

	numCorrectFinal := atomic.LoadInt32(&q.numCorrect)
	fmt.Printf("\n......\n......\nResults: You answered %d out of %d questions correctly.", numCorrectFinal, numQuestions)
	return nil
}

func (q *quiz) askQuestion(index int) error {
	question := q.questions[index][0]
	realAnswer := q.questions[index][1]

	fmt.Printf("Problem %d: %s ", index+1, question)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %+v", err)
	}
	userAnswer := scanner.Text()
	if stringsEqualForUserInput(userAnswer, realAnswer) {
		atomic.AddInt32(&q.numCorrect, 1)
	}
	return nil
}

func stringsEqualForUserInput(str1, str2 string) bool {
	return strings.ToLower(strings.TrimSpace(str1)) == strings.ToLower(strings.TrimSpace(str2))
}
