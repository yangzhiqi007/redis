package model

type Task struct {
	Name   string
	Group  string
	Server string
	Code   string
}

type TaskManager struct {
	TaskList []*Task
}

func (self *TaskManager) AddTask(task *Task) {
	self.TaskList = append(self.TaskList, task)
}

func (self *TaskManager) UpdateTask(name string, inTask *Task) {
	for index, task := range self.TaskList {
		if task.Name == name {
			self.TaskList[index] = inTask
			break
		}
	}
}

func (self *TaskManager) DeleteTask(name string) {
	for index, task := range self.TaskList {
		if task.Name == name {
			self.TaskList = append(self.TaskList[:index], self.TaskList[index+1:]...)
			break
		}
	}
}

func (self *TaskManager) TaskByName(name string) *Task {
	for _, task := range self.TaskList {
		if task.Name == name {
			return task
		}
	}

	return nil
}
