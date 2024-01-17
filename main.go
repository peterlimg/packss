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
)

func compressFolder(path string, destDir string) {
	fmt.Println("compress folder:", path)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tar -cf - %s | lz4 > %s/%s.tar.lz4", path, destDir, strings.ReplaceAll(path, "/", "")))
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error compressing folder:", err)
	}
}

func walkDir(path string, depth, currentDepth int, folderChan chan<- string) error {
	if currentDepth > depth {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			newPath := filepath.Join(path, file.Name())
			if strings.Count(newPath, "/") == depth {
				folderChan <- newPath
			}
			if err := walkDir(newPath, depth, currentDepth+1, folderChan); err != nil {
				return err
			}
		}
	}

	return nil
}

type DirDepth struct {
	Path  string
	Depth int
}

func main() {
	depth := flag.Int("depth", 3, "Depth to start packing from")
	thread := flag.Int("thread", 4, "Number of threads")
	path := flag.String("path", "", "Path to data")
	dest := flag.String("dest", "", "dest dir to save data")
	flag.Parse()

	var wg sync.WaitGroup
	folderChan := make(chan string, *thread)

	origin := strings.Count(*path, "/")
	doneC := make(chan struct{})
	// go func() {
	// 	defer close(doneC)
	// 	if err := walkDir(*path, *depth+origin, 0, folderChan); err != nil {
	// 		log.Fatal(err)
	// 	}

	// filepath.WalkDir(*path, func(path string, d fs.DirEntry, err error) error {
	// 	fmt.Println("walk in path:", path)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if d.IsDir() && strings.Count(path, "/") == *depth+origin {
	// 		fmt.Println("push path:", path)
	// 		folderChan <- path
	// 	}
	// 	return nil
	// })
	// }()

	// doneC := make(chan struct{})
	go func() {
		defer close(doneC)

		queue := []DirDepth{{Path: *path, Depth: 0}}
		for len(queue) > 0 {
			dir := queue[0]
			queue = queue[1:]

			if dir.Depth == *depth+origin {
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
