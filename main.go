package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

func compressFolder(path string, depth int, destDir string) {
	spath := strings.Split(path, "/")
	dirName := strings.Join(spath[len(spath)-depth-1:], "")

	cmd := exec.Command("sh", "-c", fmt.Sprintf("tar -cf - %s | pigz -p 10 > %s/%s.tar.gz", path, destDir, dirName))
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error compressing folder:", err)
	}
}

type DirDepth struct {
	Path  string
	Depth int
}

func main() {
	depth := flag.Int("depth", 2, "Depth to start packing from")
	thread := flag.Int("thread", 10, "Number of threads")
	path := flag.String("path", "", "Path to data")
	dest := flag.String("dest", "", "dest dir to save data")
	flag.Parse()

	var doneCount int32

	var wg sync.WaitGroup
	folderChan := make(chan string, *thread)

	doneC := make(chan struct{})
	go func() {
		defer close(doneC)

		queue := []DirDepth{{Path: *path, Depth: 0}}
		for len(queue) > 0 {
			dir := queue[0]
			queue = queue[1:]

			if dir.Depth == *depth {
				folderChan <- dir.Path
				continue
			}

			f, err := os.Open(dir.Path)
			if err != nil {
				log.Fatal(err)
			}

			files, err := f.Readdir(-1)
			f.Close()
			if err != nil {
				log.Fatal(err)
			}

			for _, file := range files {
				if file.IsDir() {
					newPath := filepath.Join(dir.Path, file.Name())
					queue = append(queue, DirDepth{Path: newPath, Depth: dir.Depth + 1})
				}
			}
		}
	}()

	if _, err := os.Stat(*dest); os.IsNotExist(err) {
		os.MkdirAll(*dest, os.ModePerm)
	}

	for i := 0; i < *thread; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for folder := range folderChan {
				fmt.Println(">>> compress folder:", i, folder)
				compressFolder(folder, *depth, *dest)
				dc := atomic.AddInt32(&doneCount, 1)
				fmt.Println("<<< thread", i, "is done", "total done:", dc)
			}
		}(i)
	}

	<-doneC
	close(folderChan)
	wg.Wait()
}
