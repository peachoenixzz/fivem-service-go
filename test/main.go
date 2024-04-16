package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const WatchMessageID = "1220328391165874176"

func main() {
	Token := "MTExMzkxNzY0NDQ2MzE2MTM2NA.GhvO0t.kCpXbMUNFslHX2EVKf2GEF4r458hXYW3ICzG4w" // Replace this with your actual bot token

	channelID := "1220328317769744404" // The ID of the channel where you want to send the message

	// Create a new Discord session using the provided bot token
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Open a WebSocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "modals_whitelist_"+i.User.ID {
			fmt.Println("Received a message component interaction")
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "modals_whitelist_" + i.User.ID,
					Title:    "Whitelist Form (ชื่อ - นามสกุล IC)",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "firstname",
									Label:       "ชื่อ (IC)",
									Style:       discordgo.TextInputShort,
									MaxLength:   30,
									MinLength:   1,
									Placeholder: "กรุณากรอก ชื่อ (IC)",
									Required:    true,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "lastname",
									Label:       "นามสกุล (IC)",
									Style:       discordgo.TextInputShort,
									MaxLength:   30,
									MinLength:   1,
									Placeholder: "กรุณากรอก นามสกุล (IC)",
									Required:    true,
								},
							},
						},
					},
				},
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
			}
		}
	})

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		// Check for the specific message and emoji reaction
		if r.MessageID == WatchMessageID && r.Emoji.Name == "✅" {
			user, err := s.User(r.UserID)
			if err != nil || user.Bot {
				return // Ignore if it's a bot or user fetch error
			}

			channel, err := s.UserChannelCreate(user.ID)
			if err != nil {
				log.Printf("Error creating DM channel: %v", err)
				return
			}

			button := discordgo.Button{
				Label:    "คลิกเพื่อกรอก ไวท์ลิสต์ (Whitelist)",
				Style:    discordgo.PrimaryButton,
				CustomID: "modals_whitelist_" + user.ID,
				Emoji: discordgo.ComponentEmoji{
					Name: "🚀", // Use a Unicode emoji or the name of a custom emoji
				},
			}

			messageSend := &discordgo.MessageSend{
				Content: "กรุณากรอก ชื่อ-นามสกุล ที่ต้องการเพื่อยืนยัน Whitelist :",
				Components: []discordgo.MessageComponent{
					&discordgo.ActionsRow{Components: []discordgo.MessageComponent{button}},
				},
			}

			_, err = s.ChannelMessageSendComplex(channel.ID, messageSend)
			if err != nil {
				log.Printf("Error sending button message: %v", err)
				return
			}
		}
	})

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionModalSubmit:
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "ขอบคุณนะคะที่กรอก ไวท์ลิส (Whitelist) เพื่อเล่นเมืองเรา",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
			}
			data := i.ModalSubmitData()

			if !strings.HasPrefix(data.CustomID, "modals_whitelist_"+i.User.ID) {
				return
			}

			userid := strings.Split(data.CustomID, "_")[2]
			_, err = s.ChannelMessageSend(channelID, fmt.Sprintf(
				"น้องหลินหลิน เพิ่มไวท์ลิส ของ <@%s>\n\n**โดยใช้ ชื่อ : %s นามสกุล : %s",
				userid,
				data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
				data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
			))
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}

			fullName := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value + " " + data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			roleID := "1213527063584313444"
			guildID := "1105607180264157204"
			err = s.GuildMemberNickname(guildID, i.User.ID, fullName)
			if err != nil {
				log.Printf("Failed to change nickname: %v", err)
				return
			}

			err = s.GuildMemberRoleAdd(guildID, i.User.ID, roleID)
			if err != nil {
				log.Printf("Failed to assign role: %v", err)
				return
			}

		}
	})

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session
	dg.Close()
}
