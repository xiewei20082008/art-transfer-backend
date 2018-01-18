package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

type Task struct {
	Taskid  string `form:"task_id" binding:"required"`
	PicURL  string `form:"pic_url" binding:"required"`
	PicHash string `form:"pic_hash" binding:"required"`
	Style   string `form:"style" binding:"required"`
}

// TaskQueue two task queue
type TaskQueue struct {
	downloadQ chan Task
	transferQ chan Task
}

var taskQueue *TaskQueue
var workDir = "/home/weix/work-dir/"

//InitTaskQueue before use
func InitTaskQueue() {
	taskQueue = &TaskQueue{downloadQ: make(chan Task, 100), transferQ: make(chan Task, 100)}
}

func downloadFile(picURL string, picHash string) (bool, error) {
	downloadDir := path.Join(workDir, picHash)
	downloadPath := path.Join(downloadDir, "src.jpg")
	fmt.Println("downloading file to ", downloadPath)
	err := os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		fmt.Println("cannot create dir")
		return false, err
	}
	if _, err := os.Stat(downloadPath); err == nil {
		fmt.Println("file already existed.")
		return true, nil
	}
	out, err := os.Create(downloadPath)
	if err != nil {
		fmt.Println("cannot create target file")
		return false, err
	}
	defer out.Close()
	resp, err := http.Get(picURL)
	if err != nil {
		fmt.Println("Failed to get file")
		return false, err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Failed to copy file")
		return false, err
	}
	return true, nil
}

func (taskQ TaskQueue) addNewTask(taskid string, picURL string, picHash string, style string) {
	fmt.Printf("Received new task %v || %s || %s || %s\n", taskid, picURL, picHash, style)
	newTask := Task{taskid, picURL, picHash, style}
	taskQ.downloadQ <- newTask
}

func (taskQ TaskQueue) tDownload() {
	for {
		downloadTask := <-taskQ.downloadQ
		fmt.Printf("downloading %v\n", downloadTask.PicURL)
		_, err := downloadFile(downloadTask.PicURL, downloadTask.PicHash)
		if err != nil {
			fmt.Printf("Failed to download pic %v\n", downloadTask.PicURL)
			continue
			// taskQ.downloadQ <- downloadTask
		}
		taskQ.transferQ <- downloadTask
	}
}

func (taskQ TaskQueue) tTransferArt() {
	for {
		transferTask := <-taskQ.transferQ
		fmt.Printf("transferring %v with style %v\n", transferTask.PicURL, transferTask.Style)
		fmt.Printf("transferQ solved\n")
	}
}
