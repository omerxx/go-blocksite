package statefile

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

func createStateFileIfNeeded() {
	fileName := stateFilePath()
	if _, err := os.Stat(stateFilePath()); err == nil {
		fmt.Printf("Using state file: %s", fileName)
	} else if os.IsNotExist(err) {
		err := ioutil.WriteFile(fileName, []byte(""), 0755)
		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}
	}
}

func stateFilePath() string {
	fileName := viper.GetString("app.stateFile")
	return fmt.Sprintf("./%s", fileName)
}

func ListStateSites() []string {
	file, err := os.Open(stateFilePath())
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func AddToState(targets []string) {
	createStateFileIfNeeded()
	fileName := stateFilePath()
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for _, target := range targets {
		if _, err := f.WriteString(fmt.Sprintf("%s\n", target)); err != nil {
			panic(err)
		}
	}
}

func RemoveFromState() {
	createStateFileIfNeeded()
}
