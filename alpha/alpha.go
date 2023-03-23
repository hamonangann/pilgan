package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Quiz struct {
	questions []Question
}

type Question struct {
	description string
	answers     []Answer
}

type Answer interface {
	isCorrect() bool
}

type CorrectAnswer struct {
	Option      string
	Description string
}

type WrongAnswer struct {
	Option      string
	Description string
}

func (c CorrectAnswer) isCorrect() bool {
	return true
}

func (w WrongAnswer) isCorrect() bool {
	return false
}

// intro describes the game rule in Bahasa Indonesia
func intro() {
	fmt.Println(`Halo!`)
	fmt.Println(`	1. Hanya ada 1 jawaban yang benar`)
	fmt.Println(`	2. Pilihlah jawaban-jawaban yang menurutmu mungkin benar`)
	fmt.Println()
	fmt.Println(`Cara menjawab: disediakan opsi jawaban A/B/C/D, tuliskan tiap opsi Anda dengan dipisahkan garis miring`)
	fmt.Println(`Misalnya Anda yakin jawabannya adalah B, maka ketik "B"`)
	fmt.Println(`Tapi jika Anda ragu jawabannya antara antara A atau C, maka ketik "A/C"`)
	fmt.Println()
	fmt.Println(`Note: urutan menulis opsi tidak jadi masalah. A/C dengan C/A dianggap jawaban yang sama.`)
}

// generateQuestion provides a question for the quiz
func generateQuestion(rawQuestion map[string]string) (Question, error) {
	question := Question{}

	desc, ok := rawQuestion["description"]
	if !ok {
		return question, errors.New("no question description provided")
	}
	question.description = desc

	correctDesc, ok := rawQuestion["correct"]
	if !ok {
		return question, errors.New("no correct answer provided")
	}
	question.answers = append(question.answers, CorrectAnswer{Description: correctDesc})

	for i := 1; i <= 3; i++ {
		wrongDesc, ok := rawQuestion["wrong"+strconv.Itoa(i)]
		if !ok {
			return question, errors.New(fmt.Sprintf("no wrong%s answer provided", strconv.Itoa(i)))
		}
		question.answers = append(question.answers, WrongAnswer{Description: wrongDesc})
	}

	return question, nil
}

// generateQuiz provides quiz ready to be consumed by the program/launcher
func generateQuiz(jsonMap map[string]any) (Quiz, error) {
	quiz := Quiz{}

	for key, val := range jsonMap {
		valMapWithStringKey, ok := val.(map[string]any)
		if !ok {
			return quiz, errors.New(fmt.Sprintf("invalid question format on %s", key))
		}
		rawQuestion := make(map[string]string)
		for innerKey, innerVal := range valMapWithStringKey {
			rawQuestion[innerKey], ok = innerVal.(string)
			if !ok {
				return quiz, errors.New(fmt.Sprintf("invalid question format on %s", key))
			}
		}

		question, err := generateQuestion(rawQuestion)
		if err != nil {
			return quiz, errors.New(fmt.Sprintf("invalid in question %s: %s", key, err.Error()))
		}
		quiz.questions = append(quiz.questions, question)
	}

	return quiz, nil
}

// readQuizAlpha reads from raw JSON file
func readQuizAlpha() map[string]any {
	// Extract JSON file
	jsonString, err := os.ReadFile("question.json")
	if err != nil {
		log.Fatal(err)
	}
	var jsonMap map[string]any
	err = json.Unmarshal(jsonString, &jsonMap)
	if err != nil {
		log.Fatal(err)
	}

	return jsonMap
}

func main() {
	rawQuiz := readQuizAlpha()
	quiz, err := generateQuiz(rawQuiz)
	if err != nil {
		log.Fatal(err)
	}

	// reader := bufio.NewReader(os.Stdin)
	log.Println(quiz)
}
