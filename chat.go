package main

import (
	"math"
	"strings"

	"github.com/code-to-go/safepool.lib/api/chat"
	"github.com/code-to-go/safepool.lib/pool"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func printChatHelp() {
	color.White("commands: ")
	color.White("  \\r refresh chat content")
	color.White("  \\m ")
	color.White("  \\x back to App list")
}

func Chat(p *pool.Pool) {
	var lastId uint64
	c := chat.Get(p)

	identities, err := p.Identities()
	if err != nil {
		color.Red("cannot retrieve identities for pool '%s': %v", p.Name)
		return
	}

	id2nick := map[string]string{}
	for _, i := range identities {
		id2nick[i.Identity.Id()] = i.Nick
	}

	selfId := p.Self.Id()
	color.Green("Enter \\? for list of commands")
	for {
		messages, err := c.GetMessages(lastId, math.MaxInt64, 32)
		if err != nil {
			color.Red("cannot retrieve chat messages from pool '%s': %v", p.Name)
			return
		}
		for _, m := range messages {
			if m.Author == selfId {
				color.Blue(": %s", m.Content)
			} else {
				color.Green("%s: %s", id2nick[m.Author], m.Content)
			}
			lastId = m.Id
		}
		prompt := promptui.Prompt{
			Label:       "> ",
			HideEntered: true,
		}

		t, _ := prompt.Run()
		t = strings.Trim(t, " ")

		switch {
		case len(t) == 0:
		case strings.HasPrefix(t, "\\?"):
			printChatHelp()
		case strings.HasPrefix(t, "\\x"):
			return
		case strings.HasPrefix(t, "\\"):
			printChatHelp()
		default:
			_, err := c.SendMessage(t, "text/html", nil)
			if err != nil {
				color.Red("cannot send message: %s")
			}
		}
	}
}
