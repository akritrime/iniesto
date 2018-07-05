package main

import (
	"fmt"
	"net/url"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

func info(s *dg.Session, cID string) {
	e := dg.MessageEmbed{}
	e.Title = "About"
	e.Color = 0xFF4136
	e.Description = "A Discord bot to follow football(soccer) scores, named after the great [Iniesto](https://youtu.be/pCVF0CSRTYA?t=95)"
	e.Fields = []*dg.MessageEmbedField{
		&dg.MessageEmbedField{
			Name:   "Creator",
			Value:  "akritrime#7920",
			Inline: true,
		},
		&dg.MessageEmbedField{
			Name:   "Links",
			Value:  "[Invite](https://discordapp.com/api/oauth2/authorize?client_id=463635074265513995&permissions=383040&scope=bot) | [Support](https://discord.gg/HaPHVY2) | [Github](https://github.com/akritrime/iniesto)",
			Inline: true,
		},
	}
	e.Footer = &dg.MessageEmbedFooter{
		Text: "v0.0.1",
	}
	e.Thumbnail = &dg.MessageEmbedThumbnail{
		URL: "https://cdn.discordapp.com/avatars/463635074265513995/139f21151cb63ed4c1ab7be2ef26f432.png?size=128",
	}
	s.ChannelMessageSendEmbed(cID, &e)
}

func matches(args []string) (cmd func(s *dg.Session, cID string)) {
	day := time.Now().UTC()
	var t1, t2 string
	switch args[0] {
	case "today", "":
		t1 = day.Format(DATE_LAYOUT)
		t2 = t1
	case "tomorrow", "tom", "t":
		t1 = day.AddDate(0, 0, 1).Format(DATE_LAYOUT)
		t2 = t1
	case "yesterday", "yes", "y":
		t1 = day.AddDate(0, 0, -1).Format(DATE_LAYOUT)
		t2 = t1
	case "on":

		fn := func(s *dg.Session, cID string) {
			e := errorEmbed("`matches on` needs a date argument in the format YYYY-MM-DD")
			s.ChannelMessageSendEmbed(cID, &e)
		}
		if len(args) != 2 {
			return fn
		}

		t1 = args[1]
		t2 = t1
		_, err := time.Parse(DATE_LAYOUT, t1)
		if err != nil {
			return fn
		}

	}

	v := url.Values{}
	v.Add("dateFrom", t1)
	v.Add("dateTo", t2)
	j, err := f.Get("matches", v.Encode())
	if err != nil {
		fmt.Println("err in my getting matches: ", err)
		return func(s *dg.Session, cID string) {
			e := errorEmbed("Error in getting matches. Try again later. If problem persists, report it on support server")
			s.ChannelMessageSendEmbed(cID, &e)
		}
	}

	var mr ResponseMatches
	err = j.Decode(&mr)
	if err != nil {
		fmt.Println("err in decoding", err)
		return func(s *dg.Session, cID string) {
			e := errorEmbed("Error in getting matches. Try again later. If problem persists, report it on support server")
			s.ChannelMessageSendEmbed(cID, &e)
		}
	}

	return func(s *dg.Session, cID string) {
		if len(mr.Matches) == 0 {
			e := errorEmbed("Sorry, there are no matches on that day.")
			_, err = s.ChannelMessageSendEmbed(cID, &e)
			if err != nil {
				fmt.Println(err)
			}

		}
		for i, m := range mr.Matches {
			e := matchEmbed(i+1, m)
			_, err = s.ChannelMessageSendEmbed(cID, &e)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func score(s *dg.Session, cID string) {
	sendErr := func(err error) {
		fmt.Println("err in my getting matches: ", err)
		e := errorEmbed("Error in getting matches. Try again, later. If problem persists, report it on support server.")
		s.ChannelMessageSendEmbed(cID, &e)
	}

	v := url.Values{}
	v.Add("status", "LIVE")
	j, err := f.Get("matches", v.Encode())
	if err != nil {
		sendErr(err)
		return
	}

	var mr ResponseMatches

	err = j.Decode(&mr)
	if err != nil {
		fmt.Println("err in decoding", err)
		sendErr(err)
		return

	}

	matches := mr.Matches

	if len(matches) == 0 {
		e := errorEmbed("Sorry, there are no matches being played now.")
		s.ChannelMessageSendEmbed(cID, &e)
		return
	}

	for i, m := range matches {
		e := matchEmbed(i+1, m)
		_, err = s.ChannelMessageSendEmbed(cID, &e)
		if err != nil {
			fmt.Println(err)
		}
	}
	// s.ChannelMessageSend(cID, buf.String())
}
