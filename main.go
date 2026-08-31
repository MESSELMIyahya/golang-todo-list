package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	LIST_COMMAND           = "list"
	LIST_COMPLETED_COMMAND = "list-completed"
	CREATE_COMMAND         = "create"     // needs the content
	UPDATE_COMMAND         = "update"     // needs the task number
	DELETE_COMMAND         = "delete"     // needs the task number
	COMPLETE_COMMAND       = "complete"   // needs the task number
	UNCOMPLETE_COMMAND     = "uncomplete" // needs the task number
	CLEAR_COMMAND          = "clear"      // clears the screen
	CLEAR_COMMAND2         = "cls"        // clears the screen (for windows users)
	EXIST_COMMAND          = "exit"       // clears the screen
)

var (
	errored = false
)

// todo item Struct

type Task struct {
	content     string
	completed   bool
	createdAt   time.Time
	updatedAt   time.Time
	completedAt time.Time
}

func main() {

	// tasks slice
	tasksList := &[]Task{
		{content: "initial task", completed: true, createdAt: time.Now(), updatedAt: time.Now(), completedAt: time.Now()},
		{content: "Make mony", completed: false, createdAt: time.Now(), updatedAt: time.Now()},
		{content: "writer the newsletter", completed: false, createdAt: time.Now(), updatedAt: time.Now()},
		{content: "pack the car", completed: true, createdAt: time.Now(), updatedAt: time.Now(), completedAt: time.Now()},
		{content: "buying a domainname", completed: false, createdAt: time.Now(), updatedAt: time.Now()},
		{content: "go the gym", completed: true, createdAt: time.Now(), updatedAt: time.Now(), completedAt: time.Now()},
	}

	fmt.Println(tasksList)
	renderUI(tasksList)
}

// error handling func

func handleError(err error) {
	if err != nil {
		errored = true
		// clearing the UI
		clearUI()

		// printing error
		fmt.Println("========================================================")
		fmt.Println("=== Something When Wrong, Please Restart The Program ===")
		fmt.Println("========================================================")
	}
}

// clearing screen (cross-platform)
func clearUI() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout

	cmd.Run()
}

// UI render functions

func printUIHeader() {
	fmt.Println("============= Welcome To Task Tracker App =============")
	fmt.Println("== list: to list all tasks ")
	fmt.Println("== list-completed: to list all completed tasks ")
	fmt.Println("== create 'Your Task Info': to create a new task ")
	fmt.Println("== update 'task-number-id' 'Your Task Info': to update an existing task ")
	fmt.Println("== complete 'task-number-id' : to mark a task as completed")
	fmt.Println("== uncomplete 'task-number-id' : to mark a task as uncompleted")
	fmt.Println("=======================================================")
	fmt.Println("== clear : to clear the CLI")
	fmt.Println("== exist : to exist the program")
	fmt.Println("=======================================================")
}

func renderUI(tasks *[]Task) {

	running := true

	// first clearing the screen
	clearUI()

	if errored {

		// stopping the running
		running = false

		// printing error
		fmt.Println("========================================================")
		fmt.Println("=== Something When Wrong, Please Restart The Program ===")
		fmt.Println("========================================================")

		return
	}

	for running {
		// checking global errors
		if errored {
			clearUI()

			// stopping the running
			running = false

			// printing error
			fmt.Println("========================================================")
			fmt.Println("=== Something When Wrong, Please Restart The Program ===")
			fmt.Println("========================================================")

			break
		}

		fmt.Print("\n")
		// clearUI()
		// printing the app ui header
		printUIHeader()

		command := renderInputAndGetInput()

		if command == CLEAR_COMMAND || command == CLEAR_COMMAND2 {
			clearUI()
			fmt.Println("======================= Cleared =======================")

		} else if command == LIST_COMMAND { // Listing Tasks

			for idx, task := range *tasks {

				completedStr := ""

				if task.completed {
					completedStr = "Completed"
				} else {
					completedStr = "Pending"
				}

				fmt.Print("\n", idx+1, " - ", task.content, " (", completedStr, ")", " (", task.updatedAt.Format(time.RFC822), ")", "\n")
			}

		} else if command == LIST_COMPLETED_COMMAND { // Listing Completed Tasks

			for idx, task := range *tasks {

				if !task.completed {
					continue
				}

				completeDate := task.completedAt.Format(time.RFC822)

				fmt.Print("\n", idx+1, " - ", task.content, " (", completeDate, ")", "\n")
			}

		} else if isCreateCmd := strings.HasPrefix(command, CREATE_COMMAND+" "); isCreateCmd { // create task

			*tasks = append(*tasks, Task{
				content:   strings.TrimSpace(strings.TrimPrefix(command, CREATE_COMMAND)),
				createdAt: time.Now(),
				completed: false,
				updatedAt: time.Now(),
			})

			fmt.Println("================== New Task Created ===================")
		} else if isUpdateCmd := strings.HasPrefix(command, UPDATE_COMMAND+" "); isUpdateCmd { // update task by id with handling edge-cases

			resolvedStr := strings.TrimSpace(strings.TrimPrefix(command, UPDATE_COMMAND))

			selectedId := strings.Split(resolvedStr, " ")[0]

			if selectedId == "" {
				fmt.Println("======== The Task ID Should Be Valid Number Id ========")
				continue
			}

			parsedSelectedId, err := strconv.ParseInt(selectedId, 0, 32)

			if err != nil || parsedSelectedId <= 0 {
				fmt.Println("======== The Task ID Should Be Valid Number Id ========")
				continue
			}

			if int(parsedSelectedId) > len(*tasks) {
				fmt.Println("=========== Task With This ID Doesn't Exist ===========")
				continue
			}

			selectedTask := (*tasks)[int(parsedSelectedId)-1]

			newContent := strings.TrimSpace(strings.Join(strings.Split(resolvedStr, " ")[1:], " ")) // splitting excluding the ID and joining and trimming the space

			(*tasks)[int(parsedSelectedId)-1] = Task{
				content:     newContent,
				updatedAt:   time.Now(),
				createdAt:   selectedTask.createdAt,
				completedAt: selectedTask.completedAt,
				completed:   selectedTask.completed,
			}

			fmt.Printf("\n================ Task %v Was Updated =================\n", parsedSelectedId)
		} else if isCompleteCmd := strings.HasPrefix(command, COMPLETE_COMMAND+" "); isCompleteCmd { // set task as completed by id with handling edge-cases

			resolvedStr := strings.TrimSpace(strings.TrimPrefix(command, COMPLETE_COMMAND))

			selectedId := strings.Split(resolvedStr, " ")[0]

			if selectedId == "" {
				fmt.Println("======== The Task ID Should Be Valid Number Id ========")
				continue
			}

			parsedSelectedId, err := strconv.ParseInt(selectedId, 0, 32)

			if err != nil || parsedSelectedId <= 0 {
				fmt.Println("======== The Task ID Should Be Valid Number Id ========")
				continue
			}

			if int(parsedSelectedId) > len(*tasks) {
				fmt.Println("=========== Task With This ID Doesn't Exist ===========")
				continue
			}

			selectedTask := (*tasks)[int(parsedSelectedId)-1]

			(*tasks)[int(parsedSelectedId)-1] = Task{
				content:     selectedTask.content,
				updatedAt:   time.Now(),
				createdAt:   selectedTask.createdAt,
				completedAt: time.Now(),
				completed:   true,
			}

			fmt.Printf("\n========== Task %v Was Marked As Completed ===========\n", parsedSelectedId)

		} else if isUncompleteCmd := strings.HasPrefix(command, UNCOMPLETE_COMMAND+" "); isUncompleteCmd { // set task as uncompleted by id with handling edge-cases

			resolvedStr := strings.TrimSpace(strings.TrimPrefix(command, UNCOMPLETE_COMMAND))

			selectedId := strings.Split(resolvedStr, " ")[0]

			if selectedId == "" {
				fmt.Println("======== The Task ID Should Be Valid Number Id ========")
				continue
			}

			parsedSelectedId, err := strconv.ParseInt(selectedId, 0, 32)

			if err != nil || parsedSelectedId <= 0 {
				fmt.Println("======== The Task ID Should Be Valid Number Id ========")
				continue
			}

			if int(parsedSelectedId) > len(*tasks) {
				fmt.Println("=========== Task With This ID Doesn't Exist ===========")
				continue
			}

			selectedTask := (*tasks)[int(parsedSelectedId)-1]

			(*tasks)[int(parsedSelectedId)-1] = Task{

				content:     selectedTask.content,
				updatedAt:   time.Now(),
				createdAt:   selectedTask.createdAt,
				completedAt: time.Time{}, // resetting the time
				completed:   false,
			}

			fmt.Printf("\n======= Task %v Was Marked As Uncompleted =======\n", parsedSelectedId)

		} else if isDeleteCmd := strings.HasPrefix(command, DELETE_COMMAND+" "); isDeleteCmd { // delete task by id with handling edge-cases

			resolvedStr := strings.TrimSpace(strings.TrimPrefix(command, DELETE_COMMAND))

			selectedId := strings.Split(resolvedStr, " ")[0]

			if selectedId == "" {
				fmt.Println("======== The Task ID Should Be Valid Number Id ========")
				continue
			}

			parsedSelectedId, err := strconv.ParseInt(selectedId, 0, 32)

			if err != nil || parsedSelectedId <= 0 {
				fmt.Println("======== The Task ID Should Be Valid Number Id ========")
				continue
			}

			if int(parsedSelectedId) > len(*tasks) {
				fmt.Println("=========== Task With This ID Doesn't Exist ===========")
				continue
			}

			// removing the task from the tasks
			*tasks = append((*tasks)[:int(parsedSelectedId-1)], (*tasks)[int(parsedSelectedId):]...)

			fmt.Printf("\n======= Task %v Was Deleted  =======\n", parsedSelectedId)

		} else if command == EXIST_COMMAND { // existing
			clearUI()
			running = false
			break
		}

	}

}

// getting user input
func renderInputAndGetInput() string {

	reader := bufio.NewReader(os.Stdin)

	print("=> ")

	input, err := reader.ReadString('\n')

	// checking the error
	handleError(err)

	return strings.TrimSpace(strings.Trim(input, "\n"))
}
