package main

import (
	"fmt"
	"strings"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

func flag(c string) string {
	ct, _ := countries.FindCountryByName(c)
	f := strings.ToLower(ct.Codes.Alpha2)
	if f == "" {
		f = fmt.Sprintf("flag_%v", strings.ToLower(c[:3]))
		return fmt.Sprintf("<:%v:%v>", f, customFlags[f])
	}
	return fmt.Sprintf(":flag_%v:", f)
}

func matchEmbed(i int, m Match) (e dg.MessageEmbed) {

	htn := m.HomeTeam.Name
	atn := m.AwayTeam.Name

	htf := flag(htn)
	atf := flag(atn)

	var name string
	if strings.Contains(strings.ToLower(m.Group), "group") {
		name = fmt.Sprintf("**%v, %v | %v **", m.Stage, m.Competition.Name, m.Group)
	} else {
		name = fmt.Sprintf("**%v, %v**", m.Group, m.Competition.Name)
	}

	e.Title = fmt.Sprintf("__**Match %v**__", i)
	e.Color = 0xFF4136
	e.Description = fmt.Sprintf("**%v %v**  **-**  **%v %v**", htn, htf, atf, atn)

	if m.Status == "SCHEDULED" {

		sch, _ := time.Parse(time.RFC3339, m.UTCDate)
		schedule := sch.Format("Mon, Jan 2, 2006 at 3:04pm (MST)")
		e.Fields = []*dg.MessageEmbedField{
			&dg.MessageEmbedField{
				Name:   name,
				Value:  fmt.Sprintf("Match is scheduled for \n**%v**", schedule),
				Inline: false,
			},
		}
	} else {
		var gt string
		dur := m.Score.Duration

		if m.Status == "FINISHED" || dur == "EXTRA_TIME" || dur == "PENALTY_SHOOTOUT" {
			gt = "Full-Time"
		} else {
			gt = "Goals    "
		}
		hts := fmt.Sprintf("```%v : %v\n```", gt, m.Score.FullTime.HomeTeam)
		ats := fmt.Sprintf("```%v : %v\n```", gt, m.Score.FullTime.AwayTeam)

		if dur == "EXTRA_TIME" || dur == "PENALTY_SHOOTOUT" {
			hts = fmt.Sprintf("%v```Extra-Time: %v```", hts, m.Score.ExtraTime.HomeTeam)
			ats = fmt.Sprintf("%v```Extra-Time: %v```", ats, m.Score.ExtraTime.AwayTeam)
		}

		if dur == "PENALTY_SHOOTOUT" {
			hts = fmt.Sprintf("%v```Penalties : %v```", hts, m.Score.Penalties.HomeTeam)
			ats = fmt.Sprintf("%v```Penalties : %v```", ats, m.Score.Penalties.AwayTeam)
		}

		e.Fields = []*dg.MessageEmbedField{
			&dg.MessageEmbedField{
				Name:   name,
				Value:  fmt.Sprintf("__Status__: **%v**", m.Status),
				Inline: false,
			},
			&dg.MessageEmbedField{
				Name:   fmt.Sprintf("%v %v", htn, htf),
				Value:  hts,
				Inline: true,
			},
			&dg.MessageEmbedField{
				Name:   fmt.Sprintf("%v %v", atn, atf),
				Value:  ats,
				Inline: true,
			},
		}

		if m.Status == "FINISHED" {
			var result string
			w := m.Score.Winner
			if w == "HOME_TEAM" {
				w = htn
			} else if w == "AWAY_TEAM" {
				w = atn
			}
			if w != "DRAW" {
				result = fmt.Sprintf("**__%v__ won.**", w)
			} else if dur == "PENALTY_SHOOTOUT" {
				if m.Score.Penalties.HomeTeam > m.Score.Penalties.AwayTeam {
					w = m.HomeTeam.Name
				} else {
					w = m.AwayTeam.Name
				}
				result = fmt.Sprintf("**__%v__ won.**", w)
			} else {
				result = "__**DRAW**__"
			}
			e.Fields = append(e.Fields, &dg.MessageEmbedField{
				Name:   "Result:",
				Value:  result,
				Inline: false,
			})
		}
	}

	e.Footer = &dg.MessageEmbedFooter{
		Text: "Scores can be delayed. Real time updates are WIP",
	}

	return
}

func errorEmbed(er string) (e dg.MessageEmbed) {
	e.Title = fmt.Sprintf("Error")
	e.Color = 0xFF4136
	e.Description = er
	return
}

func msgEmbed(msg string) (e dg.MessageEmbed) {
	e.Color = 0xFF4136
	e.Description = msg
	return
}

func helpEmbed(gID string) (e dg.MessageEmbed) {
	e.Title = "Help"
	e.Color = 0xFF4136
	e.Description = fmt.Sprintf("Help for available commands.\nYou can either **mention the bot** or **use prefix** to invoke a command.\n__**%v**__ is the prefix.", getPrefix(gID))
	e.Fields = []*dg.MessageEmbedField{
		&dg.MessageEmbedField{
			Name:   "matches [day]",
			Value:  "Displays **all the matches for the day**.\nday = [today | tomorrow | yesterday]\n",
			Inline: false,
		},
		&dg.MessageEmbedField{
			Name:   "matches on [date]",
			Value:  "Displays **all match on that particular date**.\nDate format is **YYYY-MM-DD**\n",
			Inline: false,
		},
		&dg.MessageEmbedField{
			Name:   "score",
			Value:  "Displays **score for ongoing matches**\n",
			Inline: false,
		},
		&dg.MessageEmbedField{
			Name:   "prefix [ prefix ]",
			Value:  "Sets a **new prefix** for the guild. \nOnly members with administrator permissions can set new prefix.\n",
			Inline: false,
		},
		&dg.MessageEmbedField{
			Name:   "info",
			Value:  "Displays the info\n",
			Inline: false,
		},
	}
	return
}
