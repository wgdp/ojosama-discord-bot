package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/jiro4989/ojosama"
)

func GetHelp() string {
	f, err := os.Open("help.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)

	text, err := ojosama.Convert(string(b), nil)
	if err != nil {
		panic(err)
	}

	return text
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	splitedContent := strings.Split(m.Content, " ")

	if splitedContent[0] != "/ojosama" {
		return
	}

	if splitedContent[0] == "/ojosama" {
		if splitedContent[1] == "--help" {
			s.ChannelMessageSend(m.ChannelID, GetHelp())
			return
		}
		text, err := ojosama.Convert(strings.Join(splitedContent[:len(splitedContent)], " "), nil)
		if err != nil {
			panic(err)
		}

		s.ChannelMessageSend(m.ChannelID, text)
	}
	return

}

func OjosamaErrorHandling(message string, err error) {
	errText, err := ojosama.Convert(message + "\n以下エラーメッセージです！", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(errText)
	fmt.Println(err)
}

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("CLIENT_TOKEN"))
	if err != nil {
		OjosamaErrorHandling("セッションの作成に失敗しました！", err)
		return
	}

	discord.AddHandler(onMessageCreate)

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()
	if err != nil {
		OjosamaErrorHandling("コネクションエラーです！", err)
		return
	}

	if runningText, err := ojosama.Convert("ボット実行中です！", nil); err == nil {
		fmt.Println(runningText)
	} else {
		panic(err)
	}


	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}
