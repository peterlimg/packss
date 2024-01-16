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

func compressFolder(path string) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tar -cvf - %s | lz4 > %s.tar.lz4", path, strings.ReplaceAll(path, "/", "-")))
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error compressing folder:", err)
	}
}

func main() {
	depth := flag.Int("depth", 3, "Depth to start packing from")
	thread := flag.Int("thread", 4, "Number of threads")
	path := flag.String("path", "", "Path to data")
	flag.Parse()

	var wg sync.WaitGroup
	folderChan := make(chan string, *thread)

	filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && strings.Count(path, "/") == *depth {
			folderChan <- path
		}
		return nil
	})

	for i := 0; i < *thread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for folder := range folderChan {
				compressFolder(folder)
			}
		}()
	}

	close(folderChan)
	wg.Wait()
}
