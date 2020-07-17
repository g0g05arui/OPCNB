package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var lastSource []byte
var isErrorLast bool = false
var token string
var dg *discordgo.Session

func main() {
	secs, err := strconv.ParseInt(os.Args[2], 10, 64)
	if len(os.Args) >= 4 {
		token = os.Args[3]
		dg, err = discordgo.New("Bot " + token)

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
					resp, err := http.Get(os.Args[1])

					if err != nil {
						if !isErrorLast {
							fmt.Printf("Err : %s %s \n", time.Now().String(), err.Error())
						}
						isErrorLast = true
					}
					newSource, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						if !isErrorLast {
							fmt.Printf("Err : %s %s \n", time.Now().String(), err.Error())
						}
						isErrorLast = true
					} else if string(newSource) != string(lastSource) {
						fmt.Printf("Page updated : %s\n", time.Now().String())
						if len(os.Args) >= 4 {
							sendNotif(dg, fmt.Sprintf("@everyone Page updated : %s\n", time.Now().String()))
						}
						lastSource = newSource
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
			if c.Name == "general" {
				s.ChannelMessageSend(c.ID, text)
			}
		}
	}
}
