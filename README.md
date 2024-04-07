# School file
A simple app that solves a following problem. I have to take notes during lectures. I have a directory for each lecture and use week numbers to name my notes. This requires too much thinking, especially the week number.

This app reads the timetable from file *rozvrh* and if a lecture is in progress it opens a file at SCHOOL_DIR/\<lecture\>/\<cv|př\>\<week_num\>.md

The week_num is calculated from start of semester which is set in schoolFile.go

## Timetable format
Timetable is a CSV file in format:

day;lecture_start;lecture_end;lecture_id;lecture_name;lecture_type

| field         | Description
| :---:         | :---------
| day           | Short name of the day in Czech (Po, Út...) 
| lecture_start | Time in format \<hour\>:\<minutes\>
| lecture_end   | Time in format \<hour\>:\<minutes\>
| lecture_id    | 6 character long ID for example 1IT234
| lecture_name  | Name of the lecture used to create a directory in SCHOOL_DIR
| lecture_type  | Lecture type (*přednáška* or *cvičení*)

## Install
Clone the repository

```
git clone https://github.com/random-accnt/schoolFile.git
```

Compile the program and optionally save it in some better location

```
cd schoolFile
go build schoolFile
ln -s path/to/schoolFile ~/bin/sf
```

Add timetable (default name is *rozvrh*, path can be changed inside timetable/timetable.go)

## Config
This program is intended for my personal use, so there are not many options for configuration. There are some consts in schoolFile.go and timetable/timetable.go that can be edited. Better configuration can be added in the future.

The program currently only works in **Kitty** terminal with **NeoVim**
