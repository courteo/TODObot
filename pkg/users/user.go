package users


import "TODobot/pkg/task"

type User struct {
	UserId 		uint32
	UserName     string
	CreatedTasks []int // Которые он создал
	UserTasks    []int // Которые ему задали
	ChatId       int64
}


type UserRepo interface {
	NewUserID() uint32
	AddNewTask(newTask task.Task, uid uint32)
	AddNewUser(newUser User)
	AddUserTask(taskId int, userId int)
	GetUserName(userId int) string 
	GetChatId(userid int) int64
	GetUserId(userName string) (int, error)
	DeleteTask(taskName int, uid uint32)
	DeleteUser(uid int)
	DeleteCreatedTask(taskName int, uid uint32)
	IsUserHasTask(taskName int, uid uint32) bool
	
}