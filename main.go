package main

import (
	"flag"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"tg-keyword-reply-bot/common"
	"tg-keyword-reply-bot/db"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron"
)

var bot *api.BotAPI
var gcron *cron.Cron

var (
	debug       bool
	superUserId int
)

func main() {
	botToken := flag.String("t", "", "your bot Token")
	flag.IntVar(&superUserId, "s", 0, "super manager Id")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.Parse()
	token := db.Init(*botToken)
	gcron = cron.New()
	gcron.Start()
	//开始工作
	start(token)
}

package main

import (
	"log"
	"strconv"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

func start(botToken string) {
	var err error
	bot, err = api.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = debug
	log.Printf("Authorized on account: %s  ID: %d", bot.Self.UserName, bot.Self.ID)

	u := api.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic("Can't get Updates")
	}
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		// 检查消息的群组ID是否为特定群组的ID
		if update.Message.Chat.ID == -SpecialGroupID {
			go processUpdate(&update)
		} else {
			go processUpdateAll(&update)
		}
	}
}

// 发送消息
func sendMessage(msg *api.MessageConfig) {
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Failed to send message:", err)
	}
}

func main() {
	// 在这里添加你的机器人初始化代码

	for update := range updates {
		if update.Message != nil {
			go processUpdate(&update)
		}
	}

	// ...
}


func processUpdate(update *api.Update) {
    // ...
    if update.Message != nil { 
        if update.Message.IsCommand() {
            processCommand(update)
            continue
        }
        // ...
    }
    // ...
}

func processCommand(update *api.Update) {
    if update.Message != nil && update.Message.IsCommand() {
        command := update.Message.Command()
        switch command {
        case "start":
            // 处理 /start 命令
            // ...
        case "help":
            // 处理 /help 命令
            // ...
        case "setcommands":
            if checkAdmin(update.Message.From) {
                setBotCommands()
            }
        default:
            // 处理其他命令
            // ...
        }
    }
}

func processCommond(update *api.Update) {
	var msg api.MessageConfig
	upmsg := update.Message
	gid := upmsg.Chat.ID
	uid := upmsg.From.ID
	msg = api.NewMessage(update.Message.Chat.ID, "")
	_, _ = bot.DeleteMessage(api.NewDeleteMessage(update.Message.Chat.ID, upmsg.MessageID))
	switch upmsg.Command() {
	case "start", "help":
		msg.Text = "本机器人能够自动回复特定关键词"
		sendMessage(msg)
	case "add":
		if checkAdmin(gid, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" {
				addRule(gid, order)
				msg.Text = "规则添加成功: " + order
			} else {
				msg.Text = addText
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
			}
			sendMessage(msg)
		}
	case "del":
		if checkAdmin(gid, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" {
				delRule(gid, order)
				msg.Text = "规则删除成功: " + order
			} else {
				msg.Text = delText
				msg.ParseMode = "Markdown"
			}
			sendMessage(msg)
		}
	case "list":
		if checkAdmin(gid, *upmsg.From) {
			rulelists := getRuleList(gid)
			msg.Text = "ID: " + strconv.FormatInt(gid, 10)
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true
			sendMessage(msg)
			for _, rlist := range rulelists {
				msg.Text = rlist
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				sendMessage(msg)
			}
		}
	case "admin":
		msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(uid) + ") 请求管理员出来打屁股\r\n\r\n" + getAdmins(gid)
		msg.ParseMode = "Markdown"
		sendMessage(msg)
		banMember(gid, uid, 30)
	case "me":
		myuser := upmsg.From
		msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ") 的账号信息" +
			"\r\nID: " + strconv.Itoa(uid) +
			"\r\nUseName: [" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ")" +
			"\r\nLastName: " + myuser.LastName +
			"\r\nFirstName: " + myuser.FirstName +
			"\r\nIsBot: " + strconv.FormatBool(myuser.IsBot)
		msg.ParseMode = "Markdown"
		sendMessage(msg)
	default:
	}
}
