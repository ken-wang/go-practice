package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func execute(cmd *exec.Cmd, project string, channel chan<- string) {

	stdout, _ := cmd.CombinedOutput()
	result := project + "\n" + string(stdout)
	channel <- result
}

func trace() func() {
	start := time.Now()
	return func() {
		fmt.Println("Time spent:", time.Since(start))
	}
}

func main() {
	defer trace()()

	flag.Parse()
	commands := flag.Args()
	fmt.Println("Command:", strings.Join(flag.Args(), " "))

	root, _ := os.Getwd()
	entries, _ := ioutil.ReadDir(root)

	count := 0
	channel := make(chan string, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join(root, entry.Name())
			if _, err := os.Stat(path + "/.git"); err == nil {
				cmd := exec.Command(commands[0], strings.Join(commands[1:], " "))
				cmd.Dir = path
				go execute(cmd, entry.Name(), channel)
				count++
			}
		}
	}

	for i := count; i > 0; i-- {
		result := <-channel
		fmt.Println(result)
	}
	close(channel)

}
