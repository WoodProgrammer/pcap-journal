package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func uploadS3(bucketName string, filePrefix int64, fileName string) {
	fmt.Println("Uploading file bucketNam %s %d %s", bucketName, filePrefix, fileName)

}

func main() {
	prevFileName := ""
	nextFileName := ""
	fileDir, ok := os.LookupEnv("PCAP_DIR")
	bucketName, buck_ok := os.LookupEnv("BUCKET_NAME")
	filePrefix := time.Now().Unix()
	fmt.Println(filePrefix)

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

				if event.Op&fsnotify.Create == fsnotify.Create {
					nextFileName = event.Name
					fmt.Println("Created file:", event.Name)
					fmt.Println("prevFile", prevFileName)
					fmt.Println("nextFile", nextFileName)

				}
				prevFileName = event.Name

				uploadS3(bucketName, filePrefix, prevFileName)

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
