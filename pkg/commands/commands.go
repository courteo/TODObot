package commands

import (
	"TODobot/pkg/task"
	"TODobot/pkg/users"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

var numericKeyboardFirstLayer = tgbotapi.NewReplyKeyboard(
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("/tasks"),
        tgbotapi.NewKeyboardButton("/my"),
        tgbotapi.NewKeyboardButton("/owner"),
    ),
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("/help"),
    ),
)

var MyKeyboardSecondLayer = tgbotapi.NewReplyKeyboard(
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("/owner"),
        tgbotapi.NewKeyboardButton("/help"),
    ),
)

var OwnerKeyboardSecondLayer = tgbotapi.NewReplyKeyboard(
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("/my"),
        tgbotapi.NewKeyboardButton("/help"),
    ),
)
var TasksKeyboardSecondLayer = tgbotapi.NewReplyKeyboard(
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("/owner"),
		tgbotapi.NewKeyboardButton("/my"),
        tgbotapi.NewKeyboardButton("/help"),
    ),
)

func NewTask(taskName string, creator users.User, tasks task.TaskRepo, userRepo users.UserRepo) (res string) {
	if taskName == "" {
		return "Название задачи не может быть пустой"
	}

	if tasks.IsTaskContain(taskName) {
		return "the \"" + taskName + "\" task already exists"
	}

	Inc :=tasks.GetLastTaskID()
	newTask := task.Task{
		Name:    taskName,
		Creator: creator.UserName,
		Id:      int(Inc),
	}
	creator.CreatedTasks = append(creator.CreatedTasks, newTask.Id)
	tasks.AddNewTask(newTask)

	index, err := userRepo.GetUserId(creator.UserName)
	if err == nil {
		userRepo.DeleteUser(index)
	}
	userRepo.AddNewUser(creator)

	return "Задача \"" + taskName + "\" создана, id=" + strconv.Itoa(int(Inc))
}

func PrintTaskWithAssignee(currTask task.Task) string {
	return strconv.Itoa(currTask.Id) + ". " + currTask.Name + " by @" + currTask.Creator + "\n" +
		"/unassign_" + strconv.Itoa(currTask.Id) + " /resolve_" + strconv.Itoa(currTask.Id)
}

func PrintTaskWithoutAssignee(currTask task.Task) string {
	return strconv.Itoa(currTask.Id) + ". " + currTask.Name + " by @" + currTask.Creator + "\n" +
		"/assign_" + strconv.Itoa(currTask.Id)
}

func MyTask(user users.User, tasks task.TaskRepo) (res string) {
	for i, userTask := range user.UserTasks {
		currTask, err := tasks.GetTask(userTask)
		if err != nil {
			return "нет такой задачи"
		}
		res += PrintTaskWithAssignee(currTask)
		if i != len(user.UserTasks)-1 {
			res += "\n"
		}
	}
	if len(user.UserTasks) == 0 {
		return "на вас нет задач"
	}
	return res
}

func OwnerTask(user users.User, tasks task.TaskRepo) (res string) {
	if len(user.CreatedTasks) == 0 {
		return "вы не создали задачи"
	}

	for i, userTask := range user.CreatedTasks {
		currTask, err := tasks.GetTask(userTask)
		if err != nil {
			return "нет такой задачи"
		}

		if currTask.Assignee != "" {
			res += PrintTaskWithAssignee(currTask)
		} else {
			res += PrintTaskWithoutAssignee(currTask)
		}

		if i != (len(user.CreatedTasks) - 1) {
			res += "\n"
		}
	}
	return res
}

func Assign(user users.User, id int, tasks task.TaskRepo, userRepo users.UserRepo) (res []string, chatId []int64, err error) {
	currTask, errorID := tasks.GetTask(id)
	if errorID != nil {
		err = fmt.Errorf("нет такой задачи")
		return []string{}, []int64{}, err
	}

	if currTask.Assignee != "" || currTask.Creator != user.UserName {
		var userId int
		var errorUserID error

		if currTask.Assignee != "" {
			userId, errorUserID = userRepo.GetUserId(currTask.Assignee)
		} else {
			userId, errorUserID = userRepo.GetUserId(currTask.Creator)
		}

		if errorUserID != nil {
			return []string{}, []int64{}, errorUserID
		}

		userRepo.DeleteTask(currTask.Id, uint32(userId))
		str := "Задача \"" + currTask.Name + "\" назначена на @" + user.UserName // сообщение новому владельцу задачи
		res = append(res, str)
		chatId = append(chatId, userRepo.GetChatId(userId))
	}
	currTask.Assignee = user.UserName

	userId, errorUserID := userRepo.GetUserId(user.UserName)
	if errorUserID != nil {
		err = fmt.Errorf("")
		return []string{}, []int64{}, err
	}

	if !userRepo.IsUserHasTask(currTask.Id, user.UserId) {
		userRepo.AddUserTask(currTask.Id, userId)
	}

	str := "Задача \"" + currTask.Name + "\" назначена на вас" // сообщение новому владельцу задачи
	res = append(res, str)
	chatId = append(chatId, user.ChatId)

	return res, chatId, nil
}

func UnAssign(user users.User, id int, tasks task.TaskRepo, userRepo users.UserRepo) (res []string, chatId []int64, err error) {
	currTask, errorID := tasks.GetTask(id)
	if errorID != nil {
		err = fmt.Errorf("нет такой задачи")
		return []string{}, []int64{}, err
	}

	if !userRepo.IsUserHasTask(currTask.Id, user.UserId) {
		res = append(res, "Задача не на вас")
		chatId = append(chatId, user.ChatId)
		return res, chatId, nil
	}

	currTask.Assignee = ""
	userId, errorUserID := userRepo.GetUserId(user.UserName)
	if errorUserID != nil {
		return []string{}, []int64{}, errorUserID
	}

	userRepo.DeleteTask(currTask.Id, uint32(userId))
	str := "Принято" // сняли задачу с пользователя
	res = append(res, str)
	
	chatId = append(chatId, userRepo.GetChatId(userId))

	userId, errorUserID = userRepo.GetUserId(currTask.Creator)
	if errorUserID != nil {
		return []string{}, []int64{}, errorUserID
	}

	userRepo.DeleteTask(currTask.Id, uint32(userId))
	str = "Задача \"" + currTask.Name + "\" осталась без исполнителя" // сообщение автору задачи

	res = append(res, str)
	chatId = append(chatId, userRepo.GetChatId(userId))
	return res, chatId, nil
}

func Resolve(user users.User, id int, tasks task.TaskRepo, userRepo users.UserRepo) (res []string, chatId []int64, err error) {
	currTask, errorID := tasks.GetTask(id)
	if errorID != nil {
		err = fmt.Errorf("нет такой задачи")
		return []string{}, []int64{}, err
	}

	Assignee, errorUser := userRepo.GetUserId(currTask.Assignee)
	if errorUser != nil {
		errorUser = fmt.Errorf("Нет пользователя, которому задали эту задачу")
		return []string{}, []int64{}, errorUser
	}

	if userRepo.GetUserName(Assignee)!= user.UserName {
		err = fmt.Errorf("у вас нет доступка к этому")
		return []string{}, []int64{}, err
	}
	userRepo.DeleteTask(currTask.Id, uint32(Assignee)) // удаляем задачу у исполнителя
	str := "Задача \"" + currTask.Name + "\" выполнена"
	res = append(res, str)
	chatId = append(chatId, userRepo.GetChatId(Assignee))

	creator, errorUser := userRepo.GetUserId(currTask.Creator)
	if errorUser != nil {
		return []string{}, []int64{}, errorUser
	}

	if creator == Assignee {
		return res, chatId, nil
	}
	userRepo.DeleteCreatedTask(currTask.Id, uint32(creator)) // удаляем задачу у создателя
	str = "Задача \"" + currTask.Name + "\" выполнена @" + userRepo.GetUserName(Assignee)
	res = append(res, str)
	chatId = append(chatId, userRepo.GetChatId(creator))
	tasks.DeleteTask(id)
	// AllTasks = append(AllTasks[:taskId], AllTasks[taskId+1:]...)
	return res, chatId, nil
}

func ForCommand(bot tgbotapi.BotAPI, currUser users.User, update tgbotapi.Update, tasks task.TaskRepo, userRepo users.UserRepo) {
	var msg, command, body string
	var taskId int
	var errorConv error
	index := strings.Index(update.Message.Text, " ")
	if index != -1 {
		command = update.Message.Text[1:index]
		body = update.Message.Text[index+1:]
	} else {

		command = update.Message.Text[1:]
		taskIdTemp := strings.Index(command, "_")

		if taskIdTemp != -1 {
			taskId, errorConv = strconv.Atoi(command[taskIdTemp+1:])
			command = command[:taskIdTemp]

			if errorConv != nil {
				_, err := bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Следует вводить номер задачи",
				))
				if err != nil {
					fmt.Println("ошибка при отправке")
					return
				}
				return
			}
		}
	}
	var err1 error
	switch command {
	case "new":
		msg = NewTask(body, currUser, tasks, userRepo)
		
		_, err1 = bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			msg,
		))

	case "my":
		msg = MyTask(currUser, tasks)
		sendMsg :=tgbotapi.NewMessage(
			update.Message.Chat.ID,
			msg,
		)
		sendMsg.ReplyMarkup = MyKeyboardSecondLayer
		_, err1 = bot.Send(sendMsg)
	case "owner":
		msg = OwnerTask(currUser, tasks)
		sendMsg :=tgbotapi.NewMessage(
			update.Message.Chat.ID,
			msg,
		)
		sendMsg.ReplyMarkup = OwnerKeyboardSecondLayer
		_, err1 = bot.Send(sendMsg)
	case "assign":
		BotSend(bot, currUser, taskId, update, "assign", tasks, userRepo)
	case "unassign":
		BotSend(bot, currUser, taskId, update, "unassign", tasks, userRepo)
	case "resolve":
		BotSend(bot, currUser, taskId, update, "resolve", tasks, userRepo)
	case "tasks":
		msg1, err := tasks.PrintAllTasks(currUser.UserName)
		if err != nil {
			msg1 = "Нет задач"
		}
		sendMsg :=tgbotapi.NewMessage(
			update.Message.Chat.ID,
			msg1,
		)
		sendMsg.ReplyMarkup = TasksKeyboardSecondLayer
		_, err1 = bot.Send(sendMsg)
	case "start":
		_, err1 = bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Введите /help"))

	case "help":
		Help(bot, currUser, update)
	default:
		_, err1 = bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Команды не существует. Введите /help",
		))

	}
	if err1 != nil {
		fmt.Println("ошибка при отправке")
		return
	}
}

func Help(bot tgbotapi.BotAPI, currUser users.User, update tgbotapi.Update) {
	str := "Существующие команды:\n \t /tasks - выводит текущие задачи\n \t /new XXX - вы создаете новую задачу\n" +
		"\t /assign_$ID  - назначаете пользователя исполнителем задачи\n \t /unassign_$ID - снимаете задачу с текущего пользователя\n" +
		"\t /resolve_$ID - выполняется задача\n \t /my - выводит задачи, которые назначили на меня\n \t /owner - показывает задачи, созданные мной"


		msg :=tgbotapi.NewMessage(
			update.Message.Chat.ID,
			str,
		)

		
		msg.ReplyMarkup = numericKeyboardFirstLayer
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("ошибка при отправке")
		return
	}
}

func BotSend(bot tgbotapi.BotAPI, currUser users.User, taskId int, update tgbotapi.Update, name string, tasks task.TaskRepo, userRepo users.UserRepo) {
	var msgs []string
	var chatId []int64
	var err error
	switch name {
	case "assign":
		msgs, chatId, err = Assign(currUser, taskId, tasks, userRepo)
	case "unassign":
		msgs, chatId, err = UnAssign(currUser, taskId, tasks, userRepo)
	case "resolve":
		msgs, chatId, err = Resolve(currUser, taskId, tasks, userRepo)
	}

	if err != nil {
		_, err2 := bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"нет такой задачи",
		))
		if err2 != nil {
			fmt.Println("ошибка при отправке")
			return
		}
		return
	}

	for i := range msgs {
		_, err1 := bot.Send(tgbotapi.NewMessage(
			chatId[i],
			msgs[i],
		))
		if err1 != nil {
			fmt.Println("ошибка при отправке")
			return
		}
	}
}