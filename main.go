package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
)

type FileMap struct {
	fileName string
	Count    int
}

var fileItem = make(map[string]FileMap)

func main() {

	prevFileName := ""
	nextFileName := ""
	fileDir, ok := os.LookupEnv("PCAP_DIR")
	bucketName, buck_ok := os.LookupEnv("BUCKET_NAME")
	filePrefix := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	fmt.Println(filePrefix)
	fmt.Println("Target Bucket Name", bucketName)
	if buck_ok != true {
		fmt.Println("Assign BUCKET_NAME")
	}

	if ok != true {
		fmt.Println("Assign PCAP_DIR")

	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fmt.Println("Debug Events :: ", prevFileName)

				if prevFileName != event.Name && prevFileName != "" && fileExists(prevFileName) == true {
					fmt.Println("File has to get removed", prevFileName)
					uploadS3(bucketName, filePrefix, prevFileName)

					if err := os.Remove(prevFileName); err != nil {
						panic(err)
					}
				}
				if event.Op&fsnotify.Create == fsnotify.Create {

					fmt.Println("Created file:", event.Name)
					fmt.Println("prevFile", prevFileName)
					fmt.Println("nextFile", nextFileName)

				}
				prevFileName = event.Name

			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	fmt.Println(fileDir)
	err = watcher.Add(fileDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
