package main

import (
	"fmt"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

func hi(s *dg.Session, g *dg.GuildCreate) {
	// fmt.Println(g.Name, ":", g.JoinedAt)

	if flag, err := isNewGuild(g.ID); err == nil && !flag {
		return
	} else if err != nil {
		fmt.Println("err in isNewGuild ", err)
		return
	}

	for _, v := range g.Channels {

		p, err := s.State.UserChannelPermissions(botID, v.ID)
		psm := err == nil && p&dg.PermissionSendMessages == dg.PermissionSendMessages
		prm := p&dg.PermissionReadMessages == dg.PermissionReadMessages

		if v.Type == dg.ChannelTypeGuildText && psm && prm {

			if strings.ToLower(v.Name) == "general" || v.ID == g.ID || !(strings.Contains(v.Name, "rule") || strings.Contains(v.Name, "info")) {
				if _, err := s.ChannelMessageSend(v.ID, "Thank you for inviting Iniesto in this guild. This bot was written as a simple way to keep tabs on soccer score. The default prefix is **>**. Type `> help` for list of commands."); err == nil {
					break
				}
			}
		}
	}

	_ = addNewGuild(g.ID, g.Name)

}

func bye(s *dg.Session, g *dg.GuildDelete) {

	if g.Unavailable {
		return
	}
	removeGuild(g.ID)
	// owner, err := s.User(g.OwnerID)
	// fmt.Println("Owner: ", owner.Username)
	dm, err := s.UserChannelCreate(g.OwnerID)

	if err != nil {
		fmt.Printf("err in dming owner: %v\n", err)
		return
	}

	_, err = s.ChannelMessageSend(dm.ID, "Thanks for using Iniesto. If you have any feedback or suggestions on how to improve Iniesto, please consider leaving a message on the support server. You can use `> info` in this dm channel to get the support server invite. Hope you have a good day. Bye!")
	if err != nil {
		fmt.Printf("err in dming owner: %v\n", err)
	}
}

func pingPong(s *dg.Session, m *dg.MessageCreate) {

	if m.Content == "ping" && m.Author.ID == "399951813237014528" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

func commands(s *dg.Session, m *dg.MessageCreate) {
	if u, _ := s.User(m.Author.ID); m.Author.ID == botID || u.Bot {
		return
	}

	ch, _ := s.Channel(m.ChannelID)
	prefix := getPrefix(ch.GuildID)

	ms := m.Mentions
	mentionsBot := len(ms) == 1 && ms[0].ID == botID
	hasPrefix := strings.HasPrefix(m.Content, prefix)

	if !mentionsBot && !hasPrefix {
		return
	}

	trim := strings.TrimSpace
	cnt := trim(m.Content)
	if hasPrefix {
		cnt = trim(m.Content[len(prefix):])
	}

	cmd := strings.Split(cnt, " ")
	if mentionsBot {
		temp := []string{}
		for _, w := range cmd {
			if !strings.HasPrefix(w, "<@") {
				temp = append(temp, trim(w))
			}
		}

		cmd = temp
	}

	f := func(*dg.Session, string) {}

	switch trim(cmd[0]) {
	case "info":
		f = info
	case "matches":
		args := []string{""}
		if len(cmd) > 1 {
			args = cmd[1:]
		}
		f = matches(args)
	case "score":
		f = score
	case "prefix":

		f = func(s *dg.Session, cID string) {
			c, err := s.Channel(cID)
			if err != nil {
				e := errorEmbed("Internal Error. Try again later or report it in support server.")
				fmt.Println("Err in getting Guild Channel", err)
				s.ChannelMessageSendEmbed(cID, &e)
				return
			}

			p, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
			// fmt.Println(err)
			pmc := err == nil && p&dg.PermissionManageChannels == dg.PermissionManageChannels

			if !pmc {
				e := errorEmbed("Only higher level members can set new prefix.")
				s.ChannelMessageSendEmbed(cID, &e)
				return
			}
			if len(cmd) != 2 {
				e := errorEmbed("`prefix` needs a second argument to set as the new prefix.")
				s.ChannelMessageSendEmbed(cID, &e)
				return
			}

			err = setPrefix(c.GuildID, cmd[1])
			if err != nil {
				e := errorEmbed("Error setting new prefix for Guild. Try again later.")
				s.ChannelMessageSendEmbed(cID, &e)
				return
			}
			e := msgEmbed(fmt.Sprintf("**%v** is the new prefix.", cmd[1]))
			s.ChannelMessageSendEmbed(cID, &e)
			return
		}

	case "help":
		f = func(s *dg.Session, cID string) {
			c, err := s.Channel(cID)
			if err != nil {
				e := errorEmbed("Internal Error. Try again later or report it in support server.")
				fmt.Println("Err in getting Guild Channel", err)
				s.ChannelMessageSendEmbed(cID, &e)
				return
			}

			e := helpEmbed(c.GuildID)
			s.ChannelMessageSendEmbed(cID, &e)
		}
	}

	go f(s, m.ChannelID)
}
