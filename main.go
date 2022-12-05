package main

import (
	"github.com/code-to-go/safepool.lib/api"
	"github.com/code-to-go/safepool.lib/pool"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
)

func selectPool() {
	for {
		pools := pool.List()

		items := []string{"Add", "Settings", "Exit"}
		items = append(items, pools...)
		prompt := promptui.Select{
			Label: "Choose",
			Items: items,
		}

		idx, _, _ := prompt.Run()
		switch idx {
		case 0:
			AddExisting()
		case 1:
			Settings()
		case 2:
			return
		default:
			ChooseFunction(items[idx])
		}
	}
}

func ChooseFunction(poolName string) {
	p, err := pool.Open(api.Self, poolName)
	if err != nil {
		color.Red("cannot open pool '%s': %v", err)
	}
	defer p.Close()

	for {
		items := []string{"Chat", "Library", "Cancel"}
		prompt := promptui.Select{
			Label: "Choose App",
			Items: items,
		}

		idx, _, _ := prompt.Run()
		switch idx {
		case 0:
			Chat(p)
		case 1:
			Library(p)
		default:
			return
		}
	}
}

func Settings() {
	key, _ := api.Self.Public().Base64()
	color.Green("Public Key: %s", key)
}

func main() {
	logrus.SetLevel(logrus.FatalLevel)
	api.Start()
	selectPool()
}
