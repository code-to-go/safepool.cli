package main

import (
	"github.com/code-to-go/safepool.lib/api"
	"github.com/code-to-go/safepool.lib/core"
	"github.com/code-to-go/safepool.lib/pool"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func AddExisting() {

	for {
		prompt := promptui.Prompt{
			Label:       "Token from your host. Empty to cancel",
			HideEntered: true,
		}

		t, _ := prompt.Run()
		if len(t) == 0 {
			return
		}

		token, err := pool.DecodeToken(&api.Self, t)
		if core.IsErr(err, "invalid token: %v") {
			continue
		}

		if core.IsErr(pool.Define(token.Config), "cannot save pool in db: %s") {
			continue
		}

		color.Green("Pool %s added. Host %s is trusted", token.Config.Name, token.Host.Nick)
	}

}

func AddPool() {

	items := []string{"Add Existing", "Create New", "Cancel"}
	prompt := promptui.Select{
		Label: "Choose",
		Items: items,
	}

	idx, _, _ := prompt.Run()
	switch idx {
	case 0:
	}

}
