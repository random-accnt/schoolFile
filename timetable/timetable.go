package timetable

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	PREDNASKA     = 0
	CVIKO         = 1
	sep           = ";"
	timetablePath = "/home/jirka/projects/schoolFile/rozvrh"
)

var dayToNum map[string]int = map[string]int{
	"Ne": 0,
	"Po": 1,
	"Út": 2,
	"St": 3,
	"Čt": 4,
	"Pá": 5,
	"So": 6,
}

type HourMinute struct {
	Hour   int
	Minute int
}

type Lesson struct {
	Name       string
	Id         string
	Day        int
	Start      HourMinute
	End        HourMinute
	LessonType int
}

func (lesson *Lesson) IsDuringLesson(includeBreak bool) bool {
	timeNow := time.Now()
	if int(timeNow.Weekday()) != lesson.Day {
		return false
	}

	if timeNow.Hour() > lesson.End.Hour {
		return false
	} else if timeNow.Hour() == lesson.End.Hour {
		if timeNow.Minute() > lesson.End.Minute {
			return false
		}
	}

	// add 15 min to timeNow, because I'm lazy, it should simulate 15 mminute break before lesson
	if includeBreak {
		timeNow = timeNow.Add(15 * time.Minute)
	}

	if timeNow.Hour() < lesson.Start.Hour {
		return false
	} else if timeNow.Hour() == lesson.Start.Hour {
		if timeNow.Minute() < lesson.Start.Minute {
			return false
		}
	}

	return true
}

func parseHourMin(txt string) (HourMinute, error) {
	if len(txt) != 5 {
		return HourMinute{}, errors.New("Invalid time string length")
	}

	hour, err := strconv.Atoi(txt[:2])
	if err != nil {
		return HourMinute{}, errors.New("Can't convert to hours: " + txt[:2])
	}
	min, err := strconv.Atoi(txt[3:])
	if err != nil {
		return HourMinute{}, errors.New("Can't convert to minutes: " + txt[3:])
	}

	return HourMinute{Hour: hour, Minute: min}, nil
}

func parseLesson(splitted []string) (Lesson, error) {
	// data length
	if len(splitted) != 6 {
		return Lesson{}, errors.New("Invalid length of data, expected 6, got " + fmt.Sprint(len(splitted)))
	}

	// day
	day, ok := dayToNum[splitted[0]]
	if !ok {
		return Lesson{}, errors.New("Can't parse day " + splitted[0])
	}

	// start, end
	start, err := parseHourMin(splitted[1])
	if err != nil {
		return Lesson{}, err
	}
	end, err := parseHourMin(splitted[2])
	if err != nil {
		return Lesson{}, err
	}

	// ID
	if len(splitted[3]) != 6 {
		return Lesson{}, errors.New("Lesson ID has to be 6 chars long, not " + fmt.Sprint(len(splitted[3])))
	}

	// př / cv
	var lessonType int
	switch splitted[5] {
	case "Přednáška":
		lessonType = PREDNASKA
		break
	case "Cvičení":
		lessonType = CVIKO
		break
	default:
		return Lesson{}, errors.New("Don't know this lesson type: " + splitted[5])
	}

	return Lesson{Day: day, Start: start, End: end, Id: splitted[3], Name: splitted[4], LessonType: lessonType}, nil
}

func ParseTimetable() []Lesson {
	file, err := os.Open(timetablePath)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lessons := []Lesson{}

	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, sep)
		lesson, err := parseLesson(splitted)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		lessons = append(lessons, lesson)
	}

	return lessons
}
