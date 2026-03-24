package main

import (
	"fmt"
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
	w := persist.NewWindow("example",
		app.Title("Persist Example"),
		app.MinSize(unit.Dp(400), unit.Dp(300)),
	)

	th := material.NewTheme()

	go func() {
		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				frame := w.Frame()
				layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					info := fmt.Sprintf("Persist Example\n\nPosition: (%.0f, %.0f)\nSize: %.0f x %.0f\n\nMove or resize me — I remember!",
						frame.X, frame.Y, frame.Width, frame.Height)
					label := material.H6(th, info)
					label.Color = color.NRGBA{R: 30, G: 30, B: 30, A: 255}
					label.Alignment = text.Middle
					return label.Layout(gtx)
				})
				e.Frame(gtx.Ops)
			case app.DestroyEvent:
				w.Close()
				os.Exit(0)
			}
		}
	}()

	app.Main()
}
