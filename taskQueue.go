package main

import (
	"fmt"
)

type Task struct {
	Taskid  int    `json:"task_id" binding:"required"`
	PicURL  string `json:"pic_url" binding:"required"`
	PicHash string `json:"pic_hash" binding:"required"`
	Style   string `json:"style" binding:"required"`
}

// TaskQueue two task queue
type TaskQueue struct {
	downloadQ chan Task
	transferQ chan Task
}

var taskQueue *TaskQueue

//InitTaskQueue before use
func InitTaskQueue() {
	taskQueue = &TaskQueue{downloadQ: make(chan Task, 100), transferQ: make(chan Task, 100)}
}

func (taskQ TaskQueue) addNewTask(taskid int, picURL string, picHash string, style string) {
	fmt.Printf("Received new task %v || %s || %s || %s\n", taskid, picURL, picHash, style)
	newTask := Task{taskid, picURL, picHash, style}
	taskQ.downloadQ <- newTask
}

func (taskQ TaskQueue) tDownload() {
	for {
		downloadTask := <-taskQ.downloadQ
		fmt.Printf("%v %v %v\n", downloadTask.Taskid, downloadTask.PicURL, downloadTask.Style)
		fmt.Printf("Add to transferQ\n")
		taskQ.transferQ <- downloadTask
	}
}

func (taskQ TaskQueue) tTransferArt() {
	for {
		transferTask := <-taskQ.transferQ
		fmt.Printf("%v %v %v\n", transferTask.Taskid, transferTask.PicURL, transferTask.Style)
		fmt.Printf("transferQ solved\n")
	}
}
