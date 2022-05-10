package cmd

import (
	"bufio"
	"io"
	"os"
)

// getLines 打开文件filename, 把每一行放到c channel中去
func getLines(filename string) (c chan string, err error) {

	f, oErr := os.Open(filename)
	if oErr != nil {
		err = oErr
		return
	}
	c = getLinesFromReader(f)
	return
}

func getLinesFromReader(r io.Reader) (c chan string) {

	c = make(chan string)
	scanner := bufio.NewScanner(r)

	go func() {
		for scanner.Scan() {
			c <- scanner.Text()
		}

		close(c)
	}()

	return c
}
