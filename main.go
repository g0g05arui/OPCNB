package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var lastSource []byte
var isErrorLast bool = false
var token string
var dg *discordgo.Session

type config struct {
	Prefix  string `json:"Prefix"`
	URL     string `json:"URL"`
	Channel string `json:"Channel"`
	Online  bool   `json:"Online"`
}

var botCfg config

func init() {
	cfgStr, _ := ioutil.ReadFile("env.cfg")
	json.Unmarshal(cfgStr, &botCfg)
	fmt.Println(botCfg)
}

func main() {
	secs, err := strconv.ParseInt(os.Args[1], 10, 64)
	if len(os.Args) >= 3 {
		token = os.Args[2]
		dg, err = discordgo.New("Bot " + token)
		dg.AddHandler(onMessageCreate)
		defer dg.Close()
		if err != nil {
			panic(err)
		}
		err = dg.Open()
		if err != nil {
			panic(err)
		}
		fmt.Printf("openeed")
	}
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(time.Duration(secs) * time.Second)
	quit := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println(os.Args[1:])
	func() {
		for {
			select {
			case <-ticker.C:
				{
					resp, err := http.Get(botCfg.URL)

					if err != nil {
						if !isErrorLast {
							fmt.Printf("Err1 : %s %s \n", time.Now().String(), err.Error())
						}
						isErrorLast = true
					} else {
						newSource, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							if !isErrorLast {
								fmt.Printf("Err2 : %s %s \n", time.Now().String(), err.Error())
							}
							isErrorLast = true
						} else if string(newSource) != string(lastSource) {
							fmt.Printf("Page updated : %s\n", time.Now().String())
							if (len(os.Args) >= 3) && botCfg.Online {
								sendNotif(dg, fmt.Sprintf("@everyone Page updated : %s\n", time.Now().String()))
							}
							lastSource = newSource
						}
					}
				}
			case <-quit:
				{
					ticker.Stop()
					wg.Done()
					return
				}
			}
		}
	}()
	wg.Wait()
}

func sendNotif(s *discordgo.Session, text string) {
	for _, guild := range s.State.Guilds {
		channels, _ := s.GuildChannels(guild.ID)
		for _, c := range channels {
			if c.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			if c.Name == botCfg.Channel {
				s.ChannelMessageSend(c.ID, text)
			}
		}
	}
}

func checkURL(URL string) bool {
	_, err := http.Get(URL)
	if err != nil {
		return false
	}
	return true
}

func updateCfg() {
	f, _ := os.OpenFile("env.cfg", os.O_WRONLY, 0644)
	toJSON, _ := json.Marshal(botCfg)
	ioutil.WriteFile(f.Name(), toJSON, 0644)
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content[:len(botCfg.Prefix)] == botCfg.Prefix {
		if len(m.Content) > 4 && m.Content[:4] == "!set" {
			if checkURL(strings.TrimSpace(m.Content[4:])) {
				botCfg.URL = strings.TrimSpace(m.Content[4:])
				updateCfg()
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("URL changed to : %s", botCfg.URL))
			} else {
				s.ChannelMessageSend(m.ChannelID, "Invalid URL")
			}
		}
		if len(m.Content) >= 5 && m.Content[:5] == "!what" {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Now watching URL : %s", botCfg.URL))
		}
		if len(m.Content) >= 6 && m.Content[:6] == "!start" {
			botCfg.Online = true
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Bot starting..."))
			updateCfg()
		}
		if len(m.Content) >= 5 && m.Content[:5] == "!stop" {
			botCfg.Online = false
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Bot stoping..."))
			updateCfg()
		}
		if len(m.Content) >= 5 && m.Content[:5] == "!chch" {
			botCfg.Channel = strings.TrimSpace(m.Content[5:])
			s.ChannelMessageSend(m.ChannelID, "Update channel changed")
			updateCfg()
		}
	} else {
		return
	}
}
