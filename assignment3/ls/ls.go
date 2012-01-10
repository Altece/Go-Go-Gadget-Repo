package ls

import (
	"os"
)

func Ls(dirname string, args []string) (string, os.Error) {
	// args can be any of the following:
	//		default
	//			list one line at a time display files in alphabetical order
	//		-n
	//			display with information
	//		-R
	//			go through directories recursively
	//		-t
	//			sort by timestamp
	lsdir := readdir
	sort := alphasort
	disp := namedisp
	for index, arg := range args {
		switch arg {
		case "-n":
			disp = infodisp
		case "-R":
			lsdir := recurdir
		case "-t":
			sort := timesort
		default:
			continue
		}
	}
	result := ""
	files := lsdir(dirname, sort, disp)
	for index, file := range files {
		result += file + "\n"
	}
	return result
}

// function to go through directories recursively
func recurdir(filename string, sort func([]string) []string, disp func(string) string) []string {
}
// function not go recursively through directories
func readdir(filename string, sort func([]string) []string, disp func(string) string) []string {
}

// display function to only display the name
func namedisp(filename string) string {
}
// display function to display all information
func infodisp(filename string) string {
}

// sort function to sort alphabetically
func alphasort(files []string) []string {
}
// sort function to sort by timestamp
func timesort(files []string) []string {
}
