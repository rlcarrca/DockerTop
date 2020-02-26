package main

import (
	"context"
	"fmt"
	"github.com/mum4k/termdash"
	tco "github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/text"
	"strings"
	"time"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	t, err := termbox.New()
	panicOnError(err)
	defer t.Close()

	ctx, cancel := context.WithCancel(context.Background())

	top, err := text.New()
	panicOnError(err)

	rolled, err := text.New(text.RollContent(), text.WrapAtWords())
	panicOnError(err)

	go updater(top, rolled)

	c, err := tco.New(
		t,
		tco.Border(linestyle.Light),
		tco.BorderTitle(" PRESS Q TO QUIT "),
		tco.SplitVertical(
			tco.Left(
				tco.PlaceWidget(top),
			),
			tco.Right(
				tco.Border(linestyle.Light),
				tco.BorderTitle(" HISTORY "),
				tco.PlaceWidget(rolled),
			),
			),
		)
	panicOnError(err)

	quit := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	err = termdash.Run(ctx, t, c,
		termdash.KeyboardSubscriber(quit))
	panicOnError(err)
}

func updater(top *text.Text, rolled *text.Text) {
	items_saved := []string{}
	for {
		err, items, _ := dockerList()
		panicOnError(err)

		add, remove := diff(items_saved, items)

		for _, item := range add {
			err := rolled.Write(
				fmt.Sprintf("New: %s\n", item))
			panicOnError(err)
		}
		for _, item := range remove {
			err := rolled.Write(
				fmt.Sprintf("Gone: %s\n", item))
			panicOnError(err)
		}

		content := strings.Join(items, "\n")
		if len(content) == 0 {
			content = " " // can't be empty
		}
		err = top.Write(content,
			text.WriteReplace())
		panicOnError(err)

		items_saved = items 
		time.Sleep(time.Second)
	}
}
