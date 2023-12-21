package main

import (
	"fmt"
	"log"
	"math/rand"
	"encoding/json"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type MemberInfo struct {
	ID   string `json:"memberID"`
	Nick string `json:"nick"`
}

type ListData struct {
	AuthorID    string       `json:"author_id"`
	GuildID    string       `json:"guild_id"`
	EventName    string       `json:"list_name"`
	NonBots     []MemberInfo `json:"non_bots"` // List of non-bot member IDs
	RandMapping []int        `json:"rand_mapping"`
}

func main() {
	// brokenCount := 0
	// for i := 0; i < 1000000; i++ {
	// 	randInts := randomMapping(20)

	// 	wrongOutput := false
	// 	index := 0
	// 	for j := 0; j < 20; j++ {
	// 		if randInts[j] == j {
	// 			wrongOutput = true
	// 			index = j
	// 		}
	// 	}
	// 	if wrongOutput {
	// 		brokenCount++
	// 		fmt.Printf("Repeat at: %d, arr: %v\n\n", index, randInts)
	// 	}
	// }
	// fmt.Printf("Broken Count: %d\n", brokenCount)
	// return

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token, isToken := os.LookupEnv("BOT_TOKEN")
	if !isToken {
		fmt.Println("Token not found")
		return
	}

	discordSesh, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:")
		panic(err)
	}
	discordSesh.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates | discordgo.IntentDirectMessages | discordgo.IntentGuildMembers | discordgo.IntentMessageContent | discordgo.IntentGuildMessageTyping | discordgo.IntentGuildMessages // don't mind these, I forget which ones I acc need

	discordSesh.AddHandler(ready)

	discordSesh.AddHandler(messageCreate)

	err = discordSesh.Open() 	// Open the websocket and begin listening.
	if err != nil {
		fmt.Println("Error opening Discord session: ")
		panic(err)
	}

	fmt.Println("SecretSanta is runnign now! Press CTRL-C to exit.") 	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discordSesh.Close() // Cleanly close down the Discord session.
}

func ready(sesh *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("The bot is now ready and running live")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { 	// Ignore all messages created by the bot itself
		return
	}
	mContent := strings.Fields(m.Content)

	if len(mContent) < 1 || strings.ToLower(mContent[0]) != "!secretsanta"{
		return
	}

	if len(mContent) != 2 {
		fmt.Println("Got a secret santa command")
		_, err := s.ChannelMessageSendReply(m.ChannelID, "usage: !secretSanta <Event Name>", m.Reference())
		if err != nil {
			panic(err)
		}
		return
	}

	members, err := s.GuildMembers(m.GuildID, "", 1000)
	if err != nil {
		panic(err)
	}

	eventName := mContent[1]
	var nonBots []*discordgo.Member
	var nonBotsInfo []MemberInfo

	for _, member := range members {
		if !member.User.Bot {
			nonBots = append(nonBots, member)
			nonBotsInfo = append(nonBotsInfo, MemberInfo{
				ID : member.User.ID,
				Nick: member.Nick,
			})
		}
	}
	randMapping := randomMapping(len(nonBots))

	for index, member := range nonBots {
		memberGiftTo := nonBots[randMapping[index]]

		memberChannel, err := s.UserChannelCreate(member.User.ID)
		if err != nil {
			panic(err)
		}

		msgString := fmt.Sprintf(
			"Merry Christmas!\nFor the %s, you are going to be the Secret Santa for %s, their username is <@%s>.",
			eventName, memberGiftTo.Nick, memberGiftTo.User.ID,
		)
		_, err = s.ChannelMessageSend(memberChannel.ID, msgString)
		if err != nil {
			panic(err)
		}
	}

	fileName := "secretSantaData.json"
	addListDataToJSON(fileName, ListData{
		AuthorID: m.Author.ID,
		GuildID: m.GuildID,
		EventName: eventName,
		NonBots: nonBotsInfo,
		RandMapping: randMapping,
	})

	_, err = s.ChannelMessageSendReply(m.ChannelID, "Everyone now has their Secret Santa! Happy Holidays!!", m.Reference())
	if err != nil {
		panic(err)
	}
}

func addListDataToJSON(fileName string, listData ListData){
	file, err := os.Open(fileName) // Open the file for reading
	if err != nil {
			if os.IsNotExist(err) { // Handle file not existing
					file, err = os.Create(fileName) // Create file 
					if err != nil {
							panic(err)
					}
					_, err = file.Write([]byte("[]")) // Write an empty JSON array
					if err != nil {
							panic(err)
					}
					file.Seek(0, 0) // Go back to the beginning of the file
			} else {
					panic(err) // Handle other errors
			}
	}
	defer file.Close()

	var listDatas []ListData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&listDatas)
	if err != nil {
			panic(err)
	}
	listDatas = append(listDatas, listData) 

	file, err = os.Create(fileName) // Open the file for writing (truncates the file)
	if err != nil {
			panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(listDatas)
	if err != nil {
			panic(err)
	}
}

func randomMapping(len int) []int {
	min := 0
	max := len - 1
	randInts := make([]int, len)

	for i := range randInts {
		randInts[i] = i
	}

	for i := 0; i < len; i++ {
		randomNum := min + rand.Intn(max-min+1)

		for randInts[randomNum] == i || randInts[i] == randomNum {
			randomNum = min + rand.Intn(max-min+1)
		}
		randInts[i], randInts[randomNum] = randInts[randomNum], randInts[i]
	}

	return randInts
}
