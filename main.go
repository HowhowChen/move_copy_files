package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jasonlvhit/gocron"
)

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func Task() {
	sourceLocation := ""
	copyLocation := ""
	moveLocation := ""
	now := time.Now()
	files, err := ioutil.ReadDir(sourceLocation)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		Copy(sourceLocation+file.Name(), copyLocation+file.Name())
		MoveFile(sourceLocation+file.Name(), moveLocation+file.Name())
		fmt.Printf("%s finish %s", now, file.Name())
	}
}

func main() {
	//task()
	gocron.Every(60).Second().Do(Task)
	// gocron.Every(1).Day().At("00:00").Do(task)
	<-gocron.Start()
}
