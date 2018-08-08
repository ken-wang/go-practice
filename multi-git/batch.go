package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func execute(dir string, channel chan string) {
	cmd := exec.Command("git", "pull")
	cmd.Dir = dir
	stdout, _ := cmd.CombinedOutput()

	project := dir[strings.LastIndex(dir, "/")+1:]
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

	root, _ := os.Getwd()
	files, _ := ioutil.ReadDir(root)

	count := 0
	channel := make(chan string, len(files))

	for _, file := range files {
		if file.IsDir() {
			path := root + "/" + file.Name()
			if _, err := os.Stat(path + "/.git"); err == nil {
				count++
				go execute(path, channel)
			}
		}
	}
	fmt.Println("count:", count)
	for i := count; i > 0; i-- {
		result := <-channel
		fmt.Println(result)
	}
	close(channel)

}
