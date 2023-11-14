package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type FileDef struct {
	Arguments []string `json:"arguments"`
	Directory string   `json:"directory"`
	File      string   `json:"file"`
}

var database []FileDef
var includeDirs map[string]bool = make(map[string]bool)
var directory string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func populateIncludeDirs(path string, di fs.DirEntry, err error) error {
	if filepath.Ext(path) == ".h" {
		abs, err := filepath.Abs(path)
		check(err)
		dir := filepath.Dir(abs)
		includeDirs[dir] = true
	}
	return nil
}

func generateDatabase(path string, di fs.DirEntry, err error) error {
	if filepath.Ext(path) == ".c" {
		absPath, err := filepath.Abs(path)
		check(err)

		// Get unique include dirs
		incDirs := make([]string, 0, len(includeDirs))
		for k := range includeDirs {
			incDirs = append(incDirs, "-I"+k)
		}

		args := []string{"arm-none-eabi-gcc", path} // NOTE: These are all optional. Go-to definition works without.
		// "-mcpu=cortex-m0plus",
		// "-mthumb",
		// "-DUSE_HAL_DRIVER",
		// "-DSTM32L011xx",

		args = append(args, os.Args[2:]...)
		args = append(args, incDirs...)

		def := FileDef{
			Directory: directory,
			Arguments: args,
			File:      absPath,
		}
		database = append(database, def)
	}
	return nil
}

func main() {
	dir := os.Args[1]
	var e error
	directory, e = filepath.Abs(dir)
	check(e)

	err := filepath.WalkDir(dir, populateIncludeDirs)
	check(err)

	err2 := filepath.WalkDir(dir, generateDatabase)
	check(err2)

	compDatabaseJson, err := json.Marshal(database)
	check(err)

	err = os.WriteFile("compile_commands.json", compDatabaseJson, 0644)
	check(err)

	fmt.Printf("SUCCESS! Generated compile_commands.json!")
}
