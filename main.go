package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	prevFileName := ""
	nextFileName := ""
	fileDir, ok := os.LookupEnv("PCAP_DIR")
	bucketName, buck_ok := os.LookupEnv("BUCKET_NAME")
	filePrefix := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
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
