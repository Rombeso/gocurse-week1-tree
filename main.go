package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var buffer bytes.Buffer
	//out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(&buffer, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
	content, _ := io.ReadAll(&buffer)
	fmt.Println(string(content))
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var tabSimbol string
	var depth int
	var emptyWall int
	err := treeTraversal(out, path, printFiles, tabSimbol, depth, emptyWall)
	if err != nil {
		return err
	}
	return nil
}

func treeTraversal(out io.Writer, path string, printFiles bool, tabSimbol string, depth int, emptyWall int) error {
	const simbols = "├───"
	const lastSimbols = "└───"
	files, err := filesSort(path)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}
	var dirCount int
	var dirAmount int

	for _, file := range files {
		if file.IsDir() {
			dirAmount++
		}
	}
	for i, file := range files {
		var symbol string
		var isLast bool
		if printFiles {
			isLast = i == len(files)-1
		} else {
			if file.IsDir() {
				dirCount++
				if dirCount == dirAmount {
					isLast = true
				}
			}
		}

		if isLast {
			symbol = lastSimbols
		} else {
			symbol = simbols
		}
		if !printFiles && !file.IsDir() {
			continue
		} else {
			var fileSize string
			if !file.IsDir() {
				fileSize = strconv.FormatInt(file.Size(), 10)
				if fileSize == "0" {
					fileSize = " (empty)"
				} else {
					fileSize = fmt.Sprintf(" (%vb)", file.Size())
				}
			}
			_, err = out.Write([]byte(tabSimbol + symbol + file.Name() + fileSize + "\n"))
			if err != nil {
				panic(err)
			}
		}
		if file.IsDir() {
			var newTabSimbol string

			if isLast && depth != 0 {
				if emptyWall == 0 {
					emptyWall = depth
				}
			}
			if depth == 0 && isLast {
				newTabSimbol = "\t"
			} else if !strings.Contains(tabSimbol, "│") && emptyWall != 0 && isLast {
				newTabSimbol = "\t"
			} else {
				newTabSimbol = "│\t"
			}
			for j := 1; j <= depth; j++ {
				if (emptyWall <= j && emptyWall != 0) || depth == 0 {
					newTabSimbol += "\t"
					continue
				}
				newTabSimbol += "│\t"
			}
			newPath := filepath.Join(path, file.Name())
			err := treeTraversal(out, newPath, printFiles, newTabSimbol, depth+1, emptyWall)
			if err != nil {
				panic(err.Error())
			}
		}
	}
	return nil
}

func filesSort(path string) ([]os.FileInfo, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	fileInfos := make([]os.FileInfo, len(files))
	for i, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			fmt.Println("Ошибка при получении информации о файле:", err)
			continue
		}
		fileInfos[i] = fileInfo
	}

	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Name() < fileInfos[j].Name()
	})

	return fileInfos, nil
}
