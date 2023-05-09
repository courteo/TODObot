package task


type Task struct {
	Name     string
	Assignee string
	Creator  string
	Id       int
}

type TaskRepo interface{
	NewTaskID() uint32
	GetLastTaskID() uint32
	DeleteTask(taskId int)
	GetTask(id int) (Task, error)
	IsTaskContain(taskName string) bool
	AddNewTask(newTask Task)
	PrintAllTasks(userName string) (res string, err error)
}