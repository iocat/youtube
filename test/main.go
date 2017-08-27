package main

import (
	"time"

	"github.com/iocat/youtube"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"

	"github.com/gopherjs/vecty/prop"
)

const playerID = "player"

// The global application state container and UI component
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

func (a *TestApp) Play(e *youtube.Event) {
	go e.Target.PlayVideo()
}

// -----------

func newObj() *js.Object {
	return js.Global.Get("Object").New()
}

func main() {
	app := &TestApp{}
	vecty.RenderBody(app)
	Delay(1*time.Second, func() {
		var props = youtube.NewProperties()
		props.Width = 640
		props.Height = 390
		props.VideoID = "M7lc1UVf-VE"
		props.PlayerVars.EnableJsAPI = 1
		props.PlayerEvents.OnReady = func(e *youtube.Event) {
			e.Target.PlayVideo()
		}
		app.player = youtube.NewPlayer(playerID, props)
	})
}

// setTimeout(fn, 0)
func Schedule(fn func()) {
	Delay(0, fn)
}

func Delay(d time.Duration, fn func()) {
	time.AfterFunc(d, fn)
}

func init() {
	AddScript("https://www.youtube.com/iframe_api")
}

func AddScript(url string) {
	script := js.Global.Get("document").Call("createElement", "script")
	script.Set("src", url)
	script.Set("type", "text/javascript")
	js.Global.Get("document").Get("head").Call("appendChild", script)
}
