package task

import (
	"fmt"
	"strconv"
)

type TaskMemoryRepository struct {
	data   []Task
	LastID uint32
}

func (repo *TaskMemoryRepository) NewTaskID() uint32 {
	repo.LastID++
	return repo.LastID
}

func (repo *TaskMemoryRepository) GetLastTaskID() uint32 {
	return repo.LastID
}

func (repo *TaskMemoryRepository) DeleteTask(taskId int) {
	repo.data = append(repo.data [:taskId], repo.data[taskId+1:]...)
}


func NewMemoryRepo() *TaskMemoryRepository {
	return &TaskMemoryRepository{
		data:   make([]Task, 0),
		LastID: 0,
	}
}


func (repo *TaskMemoryRepository) GetTask(id int) (Task, error) {
	if id >= int(repo.LastID) {
		err := fmt.Errorf("нет такой задачи")
		return Task{}, err
	}
	return repo.data[id - 1], nil
}


func (repo *TaskMemoryRepository) IsTaskContain(taskName string) bool {
	for _, task := range repo.data {
		if task.Name == taskName {
			return true
		}
	}
	return false
}

func (repo *TaskMemoryRepository) AddNewTask(newTask Task) {
	repo.data = append(repo.data, newTask)
	repo.LastID++
}

func (repo *TaskMemoryRepository)  PrintAllTasks(userName string) (res string, err error) {
	if len(repo.data) == 0 {
		err = fmt.Errorf("Нет задач")
		return "", err
	}

	for i, task := range repo.data {
		str := strconv.Itoa(task.Id) + ". " + task.Name + " by @" + task.Creator + "\n"
		if task.Assignee != "" { // задачу кто-то взял
			if task.Assignee == userName {
				str += "assignee: я\n"
				str += "/unassign_" + strconv.Itoa(task.Id) + " /resolve_" + strconv.Itoa(task.Id)
			} else {
				str += "assignee: @" + task.Assignee
			}

		} else { // задачу никто не взял
			str += "/assign_" + strconv.Itoa(task.Id)
		}
		res += str
		if i != len(repo.data)-1 {
			res += "\n" + "\n"
		}
	}
	return res, nil
}