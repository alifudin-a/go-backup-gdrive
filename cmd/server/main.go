package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/alifudin-a/go-backup-gdrive/pkg/gdrive"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load(".env")
}

func main() {
	gdrive.GetDriveService()

	filename := newestFile()
	var driveService = gdrive.DriveService

	// Step 1: Open  file
	f, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot open file: %v", err))
	}
	defer f.Close()

	// fIDList := os.Getenv("FID_LIST")
	fIDCreate := os.Getenv("FID_UPLOAD")

	_, err = gdrive.CreateFile(driveService, f.Name(), "application/octet-stream", f, fIDCreate)
	if err != nil {
		fmt.Printf("Could not create file: %v\n", err)
	}
}

func newestFile() string {
	dir := "/home/fariz/Documents/Puskom/Backup/DB/"
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
	}

	fullpath := dir + newestFile
	log.Println(fullpath)
	return fullpath
}
