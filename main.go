package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func compressFolder(path string, destDir string) {
	fmt.Println("compress folder:", path)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tar -cf - %s | lz4 > %s/%s.tar.lz4", path, destDir, strings.ReplaceAll(path, "/", "")))
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error compressing folder:", err)
	}
}

func main() {
	depth := flag.Int("depth", 3, "Depth to start packing from")
	thread := flag.Int("thread", 4, "Number of threads")
	path := flag.String("path", "", "Path to data")
	dest := flag.String("dest", "", "dest dir to save data")
	flag.Parse()

	var wg sync.WaitGroup
	folderChan := make(chan string, *thread)

	doneC := make(chan struct{})
	go func() {
		defer close(doneC)
		filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && strings.Count(path, "/") == *depth {
				fmt.Println("push path:", path)
				folderChan <- path
			}
			return nil
		})
	}()

	for i := 0; i < *thread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for folder := range folderChan {
				fmt.Println("compress folder:", folder)
				compressFolder(folder, *dest)
			}
		}()
	}

	<-doneC
	close(folderChan)
	wg.Wait()
}
