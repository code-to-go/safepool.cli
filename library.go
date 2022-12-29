package main

import (
	"fmt"
	"math"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/code-to-go/safepool.lib/api/library"
	"github.com/code-to-go/safepool.lib/core"
	"github.com/code-to-go/safepool.lib/pool"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func addDocument(l library.Library) {
	for {
		prompt := promptui.Prompt{
			Label: "Local Path",
		}

		localPath, _ := prompt.Run()
		if localPath == "" {
			return
		}

		stat, err := os.Stat(localPath)
		if err != nil {
			color.Red("invalid path '%s'", localPath)
			continue
		}
		if stat.IsDir() {
			color.Red("folders are not supported at the moment")
			continue
		}

		var items []string
		var item string
		parts := strings.Split(localPath, string(os.PathSeparator))
		sort.Slice(parts, func(i, j int) bool { return i > j })
		for _, p := range parts {
			item = path.Join(p, item)
			items = append(items, item)
		}

		sel := promptui.Select{
			Label: "Name in the pool",
			Items: items,
		}
		_, name, _ := sel.Run()

		prompt = promptui.Prompt{
			Label:   "Edit Name",
			Default: name,
		}
		name, _ = prompt.Run()
		h, err := l.Upload(localPath, name)
		if core.IsErr(err, "cannot upload %s: %v", localPath) {
			color.Red("cannot upload %s", localPath)
		} else {
			color.Green("'%s' uploaded to '%s:%s' with id %d", localPath, l.Pool.Name, name, h.Id)
			return
		}
	}

}

func actionsOnDocument(l library.Library, d library.Document) {
}

func Library(p *pool.Pool) {

	l := library.Get(p)

	for {
		documents, err := l.List(0, math.MaxInt)
		if core.IsErr(err, "cannot read document list: %v") {
			color.Red("something wrong")
			return
		}

		items := []string{"🔙 Back", "🗨 Add"}
		for _, d := range documents {
			item := fmt.Sprintf("⚀ %s @%s", d.Name, d.Author.Nick)
			items = append(items, item)
		}

		prompt := promptui.Select{
			Label: "Choose",
			Items: items,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return
		}
		switch idx {
		case 0:
			return
		case 1:
			addDocument(l)
		default:
			actionsOnDocument(l, documents[idx-2])
		}
	}
}
