package users

import (
	"TODobot/pkg/task"
	"errors"
	"fmt"
)

var (
	ErrNoUser  = errors.New(" No user found")
)

type UsersMemoryRepository struct {
	data []User
	LastID uint32
}

func (repo *UsersMemoryRepository) NewUserID() uint32 {
	repo.LastID++
	return repo.LastID
}

func NewMemoryRepo() *UsersMemoryRepository {
	return &UsersMemoryRepository{
		data:   make([]User, 0),
		LastID: 0,
	}
}


func (repo *UsersMemoryRepository) AddNewTask(newTask task.Task, uid uint32) {
	repo.data[uid].CreatedTasks = append(repo.data[uid].CreatedTasks, newTask.Id)
}

func (repo *UsersMemoryRepository) AddNewUser(newUser User) {
	repo.data = append(repo.data, newUser)
	repo.LastID++
}

func (repo *UsersMemoryRepository) AddUserTask(taskId int, userId int) {
	repo.data[userId].CreatedTasks = append(repo.data[userId].CreatedTasks, taskId)
}

func (repo *UsersMemoryRepository) GetChatId(userid int) int64 {
	return repo.data[userid].ChatId
}

func (repo *UsersMemoryRepository) GetUserName(userid int) string {
	return repo.data[userid].UserName
}


func (repo *UsersMemoryRepository) DeleteTask(taskName int, uid uint32) {
	index := -1
	for i, task := range repo.data[uid].UserTasks {
		if task == taskName {
			index = i
			break
		}
	}
	if index != -1 {
		repo.data[uid].UserTasks = append(repo.data[uid].UserTasks[:index], repo.data[uid].UserTasks[index+1:]...)
	}
}

func (repo *UsersMemoryRepository) DeleteUser(uid int) {
	repo.data= append(repo.data[:uid], repo.data[uid+1:]...)
}

func (repo *UsersMemoryRepository) DeleteCreatedTask(taskName int, uid uint32) {
	index := -1
	for i, task := range repo.data[uid].CreatedTasks {
		if task == taskName {
			index = i
			break
		}
	}
	if index != -1 {
		repo.data[uid].CreatedTasks = append(repo.data[uid].CreatedTasks[:index], repo.data[uid].CreatedTasks[index+1:]...)
	}
}

func (repo *UsersMemoryRepository) IsUserHasTask(taskName int, uid uint32) bool {
	for _, userTask := range repo.data[uid].UserTasks {
		if userTask == taskName {
			return true
		}
	}
	return false
}

func (repo *UsersMemoryRepository) GetUserId(userName string) (int, error) {
	for i, user := range repo.data {
		if user.UserName == userName {
			return i, nil
		}
	}
	err := fmt.Errorf("нет пользователя")
	return -1, err
}

func (repo *UsersMemoryRepository) GetUser(userName string) (User, error) {
	id, err := repo.GetUserId(userName)
	if err != nil {
		return User{}, fmt.Errorf("нет пользователя")
	}
	return repo.data[id], nil
}