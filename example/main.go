package main

import (
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/darrenoakey/daz-golang-gio/persist"
)

func main() {
	// One line: window position and size are persisted automatically.
	w := persist.NewWindow("example",
		app.Title("Persist Example"),
		app.MinSize(unit.Dp(400), unit.Dp(300)),
	)
	defer w.Close()

	th := material.NewTheme()

	go func() {
		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.H4(th, "Resize or move this window — it remembers!")
					label.Color = color.NRGBA{R: 80, G: 80, B: 80, A: 255}
					label.Alignment = text.Middle
					return label.Layout(gtx)
				})
				e.Frame(gtx.Ops)
			case app.DestroyEvent:
				os.Exit(0)
			}
		}
	}()

	app.Main()
}
