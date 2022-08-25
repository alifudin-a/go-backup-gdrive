package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alifudin-a/go-backup-gdrive/pkg/gdrive"
	"github.com/joho/godotenv"
)

func init() {
	if len(os.Args) != 2 {
		log.Println("Production mode")
		_ = godotenv.Load(".env.production")
	} else if os.Args[1] == "dev" {
		log.Println("Development mode")
		_ = godotenv.Load(".env.development")
	}
}

func main() {
	gdrive.GetDriveService()

	var dir string

	if os.Getenv("ENV") == "prod" {
		dir = os.Getenv("DIR")
	} else {
		dir = os.Getenv("DIR")
	}

	var driveService = gdrive.DriveService

	file, err := findLastFileStartsWith(dir, "siakadonline")
	if err != nil {
		log.Println(err)
	}
	log.Println(dir + file.Name())
	fullpath := dir + file.Name()

	f, err := os.Open(fullpath)
	if err != nil {
		panic(fmt.Sprintf("Cannot open file: %v", err))
	}

	defer f.Close()

	fIDCreate := os.Getenv("FID_UPLOAD")
	_, err = gdrive.CreateFile(driveService, file.Name(), "application/octet-stream", f, fIDCreate)
	if err != nil {
		fmt.Printf("Could not create file: %v\n", err)
	}
}

func findLastFileStartsWith(dir, prefix string) (lastFile os.FileInfo, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.Mode().IsRegular() {
			continue
		}
		if strings.HasPrefix(file.Name(), prefix) {
			if lastFile == nil {
				lastFile = file
			} else {
				if lastFile.ModTime().Before(file.ModTime()) {
					lastFile = file
				}
			}
		}
	}

	if lastFile == nil {
		err = os.ErrNotExist
		return
	}
	return
}

func FindLastModifiedFileBefore(dir string, t time.Time) (path string, info os.FileInfo, err error) {
	isFirst := true
	min := 0 * time.Second
	err = filepath.Walk(dir, func(p string, i os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if !i.IsDir() && i.ModTime().Before(t) {
			if isFirst {
				isFirst = false
				path = p
				info = i
				min = t.Sub(i.ModTime())
			}
			if diff := t.Sub(i.ModTime()); diff < min {
				path = p
				min = diff
				info = i
			}
		}
		return nil
	})
	return
}

func newestFile() string {
	dir := "/home/bismillah/"
	files, _ := ioutil.ReadDir(dir)
	var newestFile string
	var newestTime int64 = 0
	for _, f := range files {
		fi, err := os.Stat(dir + f.Name())
		if err != nil {
			fmt.Println(err)
		}
		currTime := fi.ModTime().Unix()
		if currTime > newestTime {
			newestTime = currTime
			newestFile = f.Name()
		}
		log.Println(currTime)
	}

	fullpath := dir + newestFile
	log.Println(fullpath)
	log.Println(newestFile)
	log.Println(newestTime)
	return fullpath
}

func newerFile() []string {
	dir := `/home/bismillah` // Windows directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var modTime time.Time
	var names []string
	for _, fi := range files {
		if fi.Mode().IsRegular() {
			if !fi.ModTime().Before(modTime) {
				if fi.ModTime().After(modTime) {
					modTime = fi.ModTime()
					names = names[:0]
				}
				names = append(names, fi.Name())
			}
		}
	}
	if len(names) > 0 {
		fmt.Println(modTime, names)
		return names
	}

	return names
}
