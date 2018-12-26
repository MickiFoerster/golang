package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var tpl *template.Template
var pseudoRandomGenSeed = rand.NewSource(time.Now().UnixNano())
var pseudoRandomGen = rand.New(pseudoRandomGenSeed)
var taskCounter = 1

type task struct {
	Operation    string
	OperandLeft  int
	OperandRight int
}

type answer struct {
	Answered         bool
	GivenAnswer      int
	CorrectAnswer    int
	AnswerWasCorrect bool
}

var tasks = make(map[*task]answer)
var currentTask *task

func init() {
	tpl = template.Must(template.ParseFiles("tpl.gohtml"))
}

func main() {
	http.HandleFunc("/", serveMainRoute)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	serverAddress := ":1234"
	log.Println("Server is running on ", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

func serveMainRoute(w http.ResponseWriter, req *http.Request) {
	answerWasCorrect := false
	log.Printf("Serving URL %q", req.URL)
	err := req.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	if len(req.PostForm) > 0 && len(req.PostForm["answer"]) == 1 {
		a := req.PostForm["answer"][0]
		taskCounter++
		// show results
		fmt.Println(a)
		i, err := strconv.Atoi(a)
		if err == nil {
			if i == tasks[currentTask].CorrectAnswer {
				answerWasCorrect = true
			}
			currentTask = nil
		}
	}

	if currentTask == nil {
		currentTask = createTask()
		//fmt.Println(currentTask)
		//fmt.Println(tasks[currentTask])
	}
	challenge := fmt.Sprintf("Was ist %s + %s?", fmt.Sprint(currentTask.OperandLeft), fmt.Sprint(currentTask.OperandRight))
	data := struct {
		Challenge   string
		Response    bool
		Answerlabel string
		Counter     int
	}{
		Challenge:   challenge,
		Response:    answerWasCorrect,
		Answerlabel: "Antwort",
		Counter:     taskCounter,
	}
	err = tpl.ExecuteTemplate(w, "tpl.gohtml", data)
	if err != nil {
		log.Fatal(err)
	}
}

func createTask() *task {
	t := new(task)
	t.Operation = "+"
	t.OperandLeft = pseudoRandomGen.Intn(10) + 1
	t.OperandRight = pseudoRandomGen.Intn(10) + 1
	tasks[t] = answer{
		Answered:      false,
		CorrectAnswer: t.OperandLeft + t.OperandRight,
	}
	return t
}

func (a answer) String() string {
	s := fmt.Sprintln(a.Answered)
	s += fmt.Sprintln(a.GivenAnswer)
	s += fmt.Sprintln(a.CorrectAnswer)
	s += fmt.Sprintln(a.AnswerWasCorrect)
	return s
}

func (t task) String() string {
	return fmt.Sprint(t.OperandLeft) + t.Operation + fmt.Sprint(t.OperandRight)
}
