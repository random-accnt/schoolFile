package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"schoolFile/timetable"
	"strings"
	"time"
)

var startDate = time.Date(2025, time.February, 12, 0, 0, 0, 0, time.UTC)
var SCHOOL_DIR = "Documents/school/LS_2025_26/"

func responseToBool(response string) (bool, error) {
	response = strings.ToLower(response)
	switch response {
	case "":
		return true, nil
	case "a":
		return true, nil
	case "n":
		return false, nil
	default:
		return false, errors.New("Invalid response: " + response)
	}
}

func getLessonNumber() int {
	currentTime := time.Now()
	dayDiff := currentTime.Sub(startDate).Hours() / 24
	weekDiff := int(dayDiff / 7)

	return weekDiff + 1
}

func getLessonDir(lesson timetable.Lesson) string {
	lessonFilepath := SCHOOL_DIR
	lessonName := strings.ReplaceAll(lesson.Name, " ", "_")
	lessonFilepath += fmt.Sprintf("%s_%s/", lesson.Id, lessonName)

	return lessonFilepath
}

func main() {
	lessons := timetable.ParseTimetable()
	var currentLesson timetable.Lesson
	lessonFound := false
	for _, lesson := range lessons {
		if lesson.IsDuringLesson(true) {
			currentLesson = lesson
			lessonFound = true
			break
		}
	}

	if !lessonFound {
		fmt.Println("No lesson now, ", fmt.Sprint(len(lessons), " lesson/s in total"))
		return
	}

	lessonTypeStr := func(lt int) string {
		switch lt {
		case timetable.PREDNASKA:
			return "přednáška"
		case timetable.CVIKO:
			return "cvičení"
		default:
			return "undefined"
		}
	}

	fmt.Printf("Current lesson: %s (%s)\n", currentLesson.Name, lessonTypeStr(currentLesson.LessonType))

	scanner := bufio.NewScanner(os.Stdin)
	var (
		isCviko    = currentLesson.LessonType == timetable.CVIKO
		openEditor bool
	)

	// FILENAME
	filename := ""

	if isCviko {
		filename += "cv"
	} else {
		filename += "pr"
	}

	lessonNumber := getLessonNumber()
	filename += fmt.Sprint(lessonNumber)
	filename += ".md"

	lessonDir := getLessonDir(currentLesson)
	if lessonDir == "" {
		fmt.Println("No lesson dir")
		return
	}

	fmt.Println("Creating file: ", lessonDir+filename)

	// EDITOR
	ok := errors.New("")
	for ok != nil {
		fmt.Print("Otevřít editor (A/n): ")
		scanner.Scan()
		response := scanner.Text()
		openEditor, ok = responseToBool(response)
	}

	// create file
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	lessonDir = home + "/" + lessonDir

	_, err = os.Stat(lessonDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(lessonDir, 0755)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		} else {
			fmt.Println("Error: ", err)
		}
	}

	// add heading to new file
	_, err = os.Stat(lessonDir + filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(lessonDir + filename)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

			heading := "# "
			if isCviko {
				heading += "Cv"
			} else {
				heading += "Př"
			}
			heading += fmt.Sprintf(" %d", lessonNumber)

			_, err = file.WriteString(heading)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}
	}

	if openEditor {
		// TODO: fix editor open, for now just return
		return

		cmd := exec.Command("kitty", "@", "launch", "--type=tab", "nvim", lessonDir+filename, "--cmd", "cd %:h")
		cmd.Env = os.Environ()
		err := cmd.Start()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		// wait for finish
		err = cmd.Wait()
		if err != nil {
			fmt.Println("Error: ", err)
		}

	} else {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("File could not be created: ", err)
			return
		}
		defer file.Close()
	}
}
