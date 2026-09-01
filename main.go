package main

import (
	"bufio"
	"encoding/json"
	"errors"
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
	CREATE_COMMAND         = "create"       // needs the content
	UPDATE_COMMAND         = "update"       // needs the task number
	DELETE_COMMAND         = "delete"       // needs the task number
	COMPLETE_COMMAND       = "complete"     // needs the task number
	UNCOMPLETE_COMMAND     = "uncomplete"   // needs the task number
	CLEAR_COMMAND          = "clear"        // clears the screen
	CLEAR_COMMAND2         = "cls"          // clears the screen (for windows users)
	EXIT_COMMAND           = "exit"         // exits the program
	HELP_COMMAND           = "help"         // prints all commands
	DATA_FILE_PATH         = "./tasks.json" // the dat file path
)

var (
	errored = false
)

// todo item Struct

type Task struct {
	Content     string    `json:"content"`
	Completed   bool      `json:"isCompleted"`
	CreatedAt   time.Time `json:"createAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CompletedAt time.Time `json:"completedAt"`
}

func main() {

	// tasks slice
	tasksList := handleInitializeTasksFromFile()

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
	fmt.Println("== create: 'Your Task Info': to create a new task ")
	fmt.Println("== update: 'task-number-id' 'Your Task Info': to update an existing task ")
	fmt.Println("== complete: 'task-number-id' : to mark a task as completed")
	fmt.Println("== uncomplete: 'task-number-id' : to mark a task as uncompleted")
	fmt.Println("== delete: 'task-number-id' : to delete a task")
	fmt.Println("=======================================================")
	fmt.Println("== clear/cls: to clear the CLI")
	fmt.Println("== exit: to exit the program")
	fmt.Println("== help: all commands")
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

	fmt.Print("\n")
	// printing the app ui header
	printUIHeader()

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

		command := renderInputAndGetInput()

		// clearing the ui
		clearUI()

		if command == CLEAR_COMMAND || command == CLEAR_COMMAND2 {
			clearUI()
			fmt.Println("======================= Cleared =======================")

		} else if command == LIST_COMMAND { // Listing Tasks

			if len(*tasks) == 0 {
				fmt.Println("========================================================")
				fmt.Println("===              No tasks yet, Add new               ===")
				fmt.Println("========================================================")
				continue
			}

			for idx, task := range *tasks {

				completedStr := ""

				if task.Completed {
					completedStr = "Completed"
				} else {
					completedStr = "Pending"
				}

				fmt.Print(idx+1, " - ", task.Content, " (", completedStr, ")", " (", task.UpdatedAt.Format(time.RFC822), ")", "\n")
			}

		} else if command == LIST_COMPLETED_COMMAND { // Listing Completed Tasks

			for idx, task := range *tasks {

				if !task.Completed {
					continue
				}

				completeDate := task.CompletedAt.Format(time.RFC822)

				fmt.Print(idx+1, " - ", task.Content, " (", completeDate, ")", "\n")
			}

		} else if isCreateCmd := strings.HasPrefix(command, CREATE_COMMAND+" "); isCreateCmd { // create task

			*tasks = append(*tasks, Task{
				Content:   strings.TrimSpace(strings.TrimPrefix(command, CREATE_COMMAND)),
				CreatedAt: time.Now(),
				Completed: false,
				UpdatedAt: time.Now(),
			})

			fmt.Println("================== New Task Created ===================")

			// backing up data
			handleBackupData(DATA_FILE_PATH, tasks)
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
				Content:     newContent,
				UpdatedAt:   time.Now(),
				CreatedAt:   selectedTask.CreatedAt,
				CompletedAt: selectedTask.CompletedAt,
				Completed:   selectedTask.Completed,
			}

			fmt.Printf("\n================ Task %v Was Updated =================\n", parsedSelectedId)

			// backing up data
			handleBackupData(DATA_FILE_PATH, tasks)
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
				Content:     selectedTask.Content,
				UpdatedAt:   time.Now(),
				CreatedAt:   selectedTask.CreatedAt,
				CompletedAt: time.Now(),
				Completed:   true,
			}

			fmt.Printf("\n========== Task %v Was Marked As Completed ===========\n", parsedSelectedId)

			// backing up data
			handleBackupData(DATA_FILE_PATH, tasks)

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

				Content:     selectedTask.Content,
				UpdatedAt:   time.Now(),
				CreatedAt:   selectedTask.CreatedAt,
				CompletedAt: time.Time{}, // resetting the time
				Completed:   false,
			}

			fmt.Printf("\n======= Task %v Was Marked As Uncompleted =======\n", parsedSelectedId)

			// backing up data
			handleBackupData(DATA_FILE_PATH, tasks)
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
			// backing up data
			handleBackupData(DATA_FILE_PATH, tasks)
		} else if command == EXIT_COMMAND { // existing
			clearUI()
			running = false
			break
		} else if command == HELP_COMMAND { // help
			clearUI()
			printUIHeader()
		} else {
			fmt.Println("========================================================")
			fmt.Println("===          Type 'help' for all commands            ===")
			fmt.Println("========================================================")
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

// handle handleInitializeTasksFromFile

func handleInitializeTasksFromFile() *[]Task {

	// reading the data file content and returns the list
	list := handleReadFileAndInitializationData(DATA_FILE_PATH)

	return list
}

func handleReadFileAndInitializationData(path string) *[]Task {

	uncodedFileBytes, err := os.ReadFile(path)

	if err != nil {
		// if the error isn't os.ErrNotExist
		if !errors.Is(err, os.ErrNotExist) {
			// printing error
			fmt.Println("=======================================================")
			fmt.Println("=== Something Went Wrong Went Loading The File Data ===")
			fmt.Println("=======================================================")

			panic(err)
		}

		// default tasks list
		defaultList := &[]Task{
			{
				Content:     "Your First Task :)",
				Completed:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				CompletedAt: time.Time{},
			},
		}

		// else creating and initializing the file with default tasks
		handleBackupData(path, defaultList)

		return defaultList
	}

	return validatingAndConvertingTasksDataFile(uncodedFileBytes)
}

func validatingAndConvertingTasksDataFile(content []byte) *[]Task {

	// validating the json format of the tasks
	isValidJson := json.Valid(content)

	if !isValidJson {
		// printing error
		fmt.Println("=======================================================")
		fmt.Println("=== Something Went Wrong Went Loading The File Data ===")
		fmt.Println("=======================================================")

		panic("Content Of The File Isn't Valid")
	}

	list := &[]Task{}

	json.Unmarshal(content, list)

	return list
}

// handles backup data (a func that backups the data into the data file)
func handleBackupData(path string, tasks *[]Task) {

	// encoding the tasks as json
	encodedJson, err := json.Marshal(*tasks)

	if err != nil {
		// printing error
		fmt.Println("=====================================================")
		fmt.Println("===== Something Went Wrong When Backing Up Data =====")
		fmt.Println("=====================================================")

		panic("Error When Backing Up Data Info The File")
	}

	// writing to the file (will create the file it doesn't exist)
	writerErr := os.WriteFile(path, encodedJson, 0666)

	if writerErr != nil {
		// printing error
		fmt.Println("=====================================================")
		fmt.Println("===== Something Went Wrong When Backing Up Data =====")
		fmt.Println("=====================================================")

		panic("Error When Backing Up Data Info The File")
	}
}
