package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	if !printFiles {
		err := pathGetterWithoutFiles(out, path)
		return err
	}
	err := pathGetter(out, path)

	return err
}

func pathGetter(out io.Writer, path string) error {
	return printDirTree(out, path, "")
}

func printDirTree(out io.Writer, path string, prefix string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	for idx, file := range files {
		fmt.Fprintf(out, "%s", prefix)

		if idx == len(files)-1 {
			fmt.Fprintf(out, "└───")
		} else {
			fmt.Fprintf(out, "├───")
		}

		fmt.Fprintf(out, "%s", file.Name())
		if file.IsDir() {
			fmt.Fprintf(out, "\n")
			if idx == len(files)-1 {
				err := printDirTree(out, filepath.Join(path, file.Name()), prefix+"\t")
				if err != nil {
					return err
				}
			} else {
				err := printDirTree(out, filepath.Join(path, file.Name()), prefix+"│\t")
				if err != nil {
					return err
				}
			}
			continue
		}
		fileIn, _ := file.Info()
		if fileIn.Size() == 0 {
			fmt.Fprintf(out, " (empty)\n")
		} else {
			fmt.Fprintf(out, " (%db)\n", fileIn.Size())
		}
	}
	return nil
}

func pathGetterWithoutFiles(out io.Writer, path string) error {
	return printDirTreeWithoutFiles(out, path, "")
}

func printDirTreeWithoutFiles(out io.Writer, path string, prefix string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	var dirs []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
	}

	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	for idx, dir := range dirs {
		fmt.Fprintf(out, "%s", prefix)

		if idx == len(dirs)-1 {
			fmt.Fprintf(out, "└───")
		} else {
			fmt.Fprintf(out, "├───")
		}

		fmt.Fprintf(out, "%s\n", dir.Name())

		if idx == len(dirs)-1 {
			err := printDirTreeWithoutFiles(out, filepath.Join(path, dir.Name()), prefix+"\t")
			if err != nil {
				return err
			}
		} else {
			err := printDirTreeWithoutFiles(out, filepath.Join(path, dir.Name()), prefix+"│\t")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
