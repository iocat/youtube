package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"

	"github.com/iocat/youtube"
)

const (
	youtubeIframeAPISrc = "https://www.youtube.com/iframe_api"
	playerID            = "playerID"
)

func main() {
	// 1. Add the Youtube script to your header
	yourAddScriptFunc(youtubeIframeAPISrc)

	// 2. Place the player container as a div in your html with a predefined id
	var app = &App{playerID: playerID}
	vecty.RenderBody(app)

	// 3. Initialize the player when the Youtube API's done loading
	js.Global.Set("onYouTubeIframeAPIReady", func() {
		// Create and set the initial properties of the player (check the document
		// for the specific fields)
		var props = youtube.NewProperties()
		props.Width = 640
		props.Height = 390
		props.VideoID = "b-tAiOVMYFY"
		props.PlayerVars.EnableJsAPI = 1
		props.PlayerEvents.OnReady = func(e *youtube.Event) {
			e.Target.PlayVideo()
		}
		// Create and cache the created player
		app.player = youtube.NewPlayer(playerID, props)
	})
}

type App struct {
	vecty.Core
	player   *youtube.Player
	playerID string
}

func (b *App) Render() *vecty.HTML {
	return elem.Body(
		elem.Div(
			prop.ID(b.playerID),
		),
	)
}

func yourAddScriptFunc(url string) {
	script := js.Global.Get("document").Call("createElement", "script")
	script.Set("src", url)
	script.Set("type", "text/javascript")
	js.Global.Get("document").Get("head").Call("appendChild", script)
}
