package main

import (
	"strings"
	"sync"

	"github.com/gopherjs/vecty/event"
	"github.com/iocat/youtube"
	"github.com/iocat/youtube/ytutil"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"

	"github.com/gopherjs/vecty/prop"

	"strconv"
)

const (
	playerID = "player"
	// statUpdateFrequency in miliseconds (use js default time unit)
	statUpdateFrequency = 500
)

// TestApp is the global application state container and a UI component
type TestApp struct {
	vecty.Core
	player      *youtube.Player
	seekToSec   float64
	volumeToSet int

	// decide to show stats or not ( only shown when the player is readily initialized)
	showStat bool
}

func (a *TestApp) play(e *vecty.Event) {
	a.player.PlayVideo()
}

func (a *TestApp) pause(e *vecty.Event) {
	a.player.PauseVideo()
}

func (a *TestApp) stop(e *vecty.Event) {
	a.player.StopVideo()
}

func (a *TestApp) seekTo(e *vecty.Event) {
	a.player.SeekTo(a.seekToSec, true)
}

func (a *TestApp) controllerPlay() *vecty.HTML {
	return elem.Div(
		prop.Class("ui labeled icon basic button "),
		elem.Italic(prop.Class("play icon")),
		vecty.Text("PlayVideo()"),
		event.Click(a.play),
	)
}

func (a *TestApp) controllerPause() *vecty.HTML {
	return elem.Div(
		prop.Class("ui labeled icon basic button "),
		elem.Italic(prop.Class("pause icon")),
		vecty.Text("PauseVideo()"),
		event.Click(a.pause),
	)
}

func (a *TestApp) controllerStop() *vecty.HTML {
	return elem.Div(
		prop.Class("ui labeled icon basic button "),
		elem.Italic(prop.Class("stop icon")),
		vecty.Text("StopVideo()"),
		event.Click(a.stop),
	)
}

func (a *TestApp) controllerSeekTo() *vecty.HTML {
	return elem.Div(
		prop.Class("ui action input"),
		elem.Input(
			prop.Type(prop.TypeText),
			prop.Value(strconv.FormatFloat(a.seekToSec, 'f', -1, 64)),
			prop.Placeholder("Seconds to Seek To"),
			event.Input(func(e *vecty.Event) {
				a.seekToSec, _ = strconv.ParseFloat(
					e.Target.Get("value").String(), 64)
			}),
		),
		elem.Button(
			prop.Class("ui right button"),
			vecty.Text("SeekTo()"),
			event.Click(a.seekTo),
		),
	)
}

func (a *TestApp) setVolume(e *vecty.Event) {
	a.player.SetVolume(a.volumeToSet)
}

func (a *TestApp) controllerVolume() *vecty.HTML {
	return elem.Div(
		prop.Class("ui action input"),
		elem.Input(
			prop.Type(prop.TypeText),
			prop.Value(strconv.FormatInt(int64(a.volumeToSet), 10)),
			prop.Placeholder("Volume (0-100)"),
			event.Input(func(e *vecty.Event) {
				v, _ := strconv.ParseInt(
					e.Target.Get("value").String(), 10, 32)
				a.volumeToSet = int(v)
			}),
		),
		elem.Button(
			prop.Class("ui right button"),
			vecty.Text("SetVolume()"),
			event.Click(a.setVolume),
		),
	)
}

func stat(desc string, value string) *vecty.HTML {
	return elem.Div(
		prop.Class(" statistic"),
		elem.Div(
			prop.Class("value"),
			vecty.Text(value),
		),
		elem.Div(
			prop.Class("label"),
			vecty.Text(desc),
		),
	)
}

func (a *TestApp) getVolume() string {
	return strconv.FormatInt(int64(a.player.Volume()), 10)
}

func (a *TestApp) getMuted() string {
	if a.player.IsMuted() {
		return "true"
	}
	return "false"

}

func (a *TestApp) playerState() string {
	switch st := a.player.PlayerState(); st {
	case youtube.Unstarted:
		return "-1 (Unstarted)"
	case youtube.Ended:
		return "0 (Ended)"
	case youtube.Playing:
		return "1 (Playing)"
	case youtube.Paused:
		return "2 (Paused)"
	case youtube.Buffering:
		return "3 (Buffering)"
	case youtube.VideoCued:
		return "5 (Video Cued)"
	default:
		return strings.Join([]string{
			strconv.FormatInt(int64(st), 10),
			"Unknown",
		}, " ")
	}
}

func (a *TestApp) playbackRate() string {
	return strconv.FormatFloat(a.player.PlaybackRate(), 'f', 2, 64)
}

func (a *TestApp) videoLoadedFraction() string {
	return strconv.FormatFloat(a.player.VideoLoadedFraction(), 'f', 2, 64)
}

func (a *TestApp) stats() *vecty.HTML {
	return elem.Div(
		prop.Class("ui mini statistics"),
		stat("Volume()", a.getVolume()),
		stat("IsMuted()", a.getMuted()),
		stat("PlayerState()", a.playerState()),
		stat("PlaybackRate()", a.playbackRate()),
		stat("VideoLoadedFraction()", a.videoLoadedFraction()),
	)
}

var once = &sync.Once{}

func (a *TestApp) Render() *vecty.HTML {
	once.Do(func() {
		// update stat every 1 second
		js.Global.Call("setInterval", func() { vecty.Rerender(a) }, statUpdateFrequency) // reset stat
	})

	var statsDiv *vecty.HTML
	if a.showStat {
		statsDiv = a.stats()
	}
	return elem.Body(
		elem.Div(
			prop.Class("ui center aligned container"),
			elem.Div(
				prop.ID(playerID),
			),
		),
		elem.Heading4(
			prop.Class("ui horizontal divider header"),
			elem.Italic(
				prop.Class(" setting icon"),
			),
			vecty.Text("Controls"),
		),
		elem.Div(
			prop.Class("ui fluid center aligned container"),
			elem.Div(
				prop.Class("ui container"),
				a.controllerPlay(),
				a.controllerPause(),
				a.controllerStop(),
			),
			elem.Div(
				prop.Class("ui container"),
				a.controllerSeekTo(),
			),
			elem.Div(
				prop.Class("ui container"),
				a.controllerVolume(),
			),
		),
		vecty.If(statsDiv != nil,
			elem.Heading4(
				prop.Class("ui horizontal divider header"),
				elem.Italic(
					prop.Class("bar chart icon"),
				),
				vecty.Text(
					strings.Join([]string{
						"Stats (every ",
						strconv.FormatInt(statUpdateFrequency, 10),
						" miliseconds)",
					}, ""),
				),
			),
			elem.Div(
				prop.Class("ui fluid center aligned container"),
				statsDiv,
			),
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

func setupSemanticUI() {
	addScript := func(url string) {
		script := js.Global.Get("document").Call("createElement", "script")
		script.Set("src", url)
		script.Set("type", "text/javascript")
		js.Global.Get("document").Get("head").Call("appendChild", script)
	}
	vecty.AddStylesheet("https://cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.2.13/semantic.min.css")
	addScript("https://code.jquery.com/jquery-3.2.1.min.js")
	addScript("https://cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.2.13/semantic.min.js")

}

func main() {
	setupSemanticUI()

	app := &TestApp{}
	app.showStat = false
	vecty.RenderBody(app)
	ytutil.OnLoaded(func() {
		var props = youtube.NewProperties()
		props.Width = 640
		props.Height = 390
		props.VideoID = "b-tAiOVMYFY"
		props.PlayerVars.EnableJsAPI = 1
		props.Events.OnReady = func(e *youtube.Event) {
			e.Target.PlayVideo()
		}

		app.player = youtube.NewPlayer(playerID, props)
		app.showStat = true
	})
}
