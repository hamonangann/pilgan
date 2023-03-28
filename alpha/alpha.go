package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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
	setOption(string)
	getOption() string
	getDescription() string
}

// CorrectAnswer provides Answer interface
type CorrectAnswer struct {
	Option      string
	Description string
}

// WrongAnswer provides Answer interface
type WrongAnswer struct {
	Option      string
	Description string
}

func (c *CorrectAnswer) isCorrect() bool {
	return true
}

func (c *CorrectAnswer) setOption(opt string) {
	c.Option = opt
}

func (c *CorrectAnswer) getOption() string {
	return c.Option
}

func (c *CorrectAnswer) getDescription() string {
	return c.Description
}

func (w *WrongAnswer) isCorrect() bool {
	return false // wrong answer is NOT correct
}

func (w *WrongAnswer) setOption(opt string) {
	w.Option = opt
}

func (w *WrongAnswer) getOption() string {
	return w.Option
}

func (w *WrongAnswer) getDescription() string {
	return w.Description
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
	question.answers = append(question.answers, &CorrectAnswer{Description: correctDesc})

	for i := 1; i <= 3; i++ {
		wrongDesc, ok := rawQuestion["wrong"+strconv.Itoa(i)]
		if !ok {
			return question, errors.New(fmt.Sprintf("no wrong%s answer provided", strconv.Itoa(i)))
		}
		question.answers = append(question.answers, &WrongAnswer{Description: wrongDesc})
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
	fmt.Println()
}

// launchQuestionAlpha adds options to every answer description
func launchQuestionAlpha(question *Question) {
	fmt.Println(`Pertanyaan:`, question.description)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(question.answers), func(i, j int) {
		question.answers[i], question.answers[j] = question.answers[j], question.answers[i]
	})
	options := []string{"A", "B", "C", "D"}
	for idx, option := range options {
		question.answers[idx].setOption(option)
		fmt.Println(option, ".", question.answers[idx].getDescription())
	}
}

// formatQuizAlpha returns map of flags of every option, according to user answer
func formatAnswerAlpha(answerRawString string) (map[string]bool, bool) {
	opts := strings.Split(strings.TrimSuffix(strings.ToUpper(answerRawString), "\n"), "/")
	formattedAnswer := map[string]bool{"A": false, "B": false, "C": false, "D": false}
	for _, opt := range opts {
		if opt != "A" && opt != "B" && opt != "C" && opt != "D" { // only ABCD valid option
			return formattedAnswer, false
		}
		formattedAnswer[opt] = true
	}
	return formattedAnswer, true
}

// getCorrectOption returns options in correct answer. It is guaranteed that only 1 answer is correct
func getCorrectOption(available []Answer) string {
	var correctOption string
	for _, ans := range available {
		if ans.isCorrect() {
			correctOption = ans.getOption()
		}
	}
	return correctOption
}

// getNumberOfSelectedAnswer return number of selected answer by user. Guaranteed minimal of 1, max of 4.
func getNumberOfSelectedAnswer(selected map[string]bool) int {
	res := 0
	for _, value := range selected {
		if value {
			res += 1
		}
	}
	return res
}

// markAnswerAlpha gives marks message and call answer checker
func markAnswerAlpha(availableAnswers []Answer, selectedAnswer map[string]bool, score int) int {
	correctOption := getCorrectOption(availableAnswers)
	correctPoint := 12 / getNumberOfSelectedAnswer(selectedAnswer)

	for option, selected := range selectedAnswer {
		if selected && option == correctOption { // user guess the right answer
			score += correctPoint
			fmt.Println(`Yes! Jawaban`, option, `benar`)
			return score
		}
	}

	// none of the user guess are true
	fmt.Println(`Yah... Jawabanmu salah`)

	return score
}

// launchQuizAlpha prompts user feedback
func launchQuizAlpha(quiz *Quiz) int {
	reader := bufio.NewReader(os.Stdin)
	score := 0

	fmt.Println(`Tekan ENTER jika sudah siap!`)
	_, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Terjadi kesalahan. Silakan ulangi permainan.")
	}

	for idx := range quiz.questions {
		launchQuestionAlpha(&quiz.questions[idx])
		isAnswerValid := false
		answer := make(map[string]bool, 4)
		for !isAnswerValid {
			fmt.Println("Jawaban:")
			answerOpt, err := reader.ReadString('\n')
			if err == nil {
				answer, isAnswerValid = formatAnswerAlpha(answerOpt)
			}
		}
		score = markAnswerAlpha(quiz.questions[idx].answers, answer, score)

		fmt.Println(`Skormu adalah:`, score)
		fmt.Println()
	}

	return score
}

func main() {
	rawQuiz := readQuizAlpha()
	quiz, err := generateQuiz(rawQuiz)
	if err != nil {
		log.Fatal(err)
	}

	intro()

	finalScore := launchQuizAlpha(&quiz)

	fmt.Println()
	fmt.Println(`Permainan berakhir! Skormu adalah:`, finalScore)
}
