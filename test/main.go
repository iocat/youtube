package main

import (
	"github.com/iocat/youtube"
	"github.com/iocat/youtube/ytutil"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"

	"github.com/gopherjs/vecty/prop"
)

const playerID = "player"

// TestApp is the global application state container and a UI component
type TestApp struct {
	vecty.Core
	player *youtube.Player
}

func (a *TestApp) Render() *vecty.HTML {
	return elem.Body(
		elem.Div(
			prop.ID(playerID),
		),
	)
}

func (a *TestApp) rerender() {
	vecty.Rerender(a)
}

// -----------

func init() {
	ytutil.Load()
}

func main() {
	app := &TestApp{}
	vecty.RenderBody(app)
	ytutil.OnLoaded(func() {
		var props = youtube.NewProperties()
		props.Width = 640
		props.Height = 390
		props.VideoID = "b-tAiOVMYFY"
		props.PlayerVars.EnableJsAPI = 1
		props.PlayerEvents.OnReady = func(e *youtube.Event) {
			e.Target.PlayVideo()
		}
		app.player = youtube.NewPlayer(playerID, props)
	})
}
