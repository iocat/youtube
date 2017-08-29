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
)

// TestApp is the global application state container and a UI component
type TestApp struct {
	vecty.Core

	// The frequency in miliseconds with which the stats is updated
	StatUpdateFreq int
	player         *youtube.Player
	seekToSec      float64
	volumeToSet    int

	idToLoad        string
	startsSecond    float64
	selectedQuality youtube.Quality

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
		prop.Class("ui labeled icon tiny basic button"),
		elem.Italic(prop.Class("play icon")),
		vecty.Text("PlayVideo()"),
		event.Click(a.play),
	)
}

func (a *TestApp) controllerPause() *vecty.HTML {
	return elem.Div(
		prop.Class("ui labeled icon tiny basic button "),
		elem.Italic(prop.Class("pause icon")),
		vecty.Text("PauseVideo()"),
		event.Click(a.pause),
	)
}

func (a *TestApp) controllerStop() *vecty.HTML {
	return elem.Div(
		prop.Class("ui labeled icon tiny basic button "),
		elem.Italic(prop.Class("stop icon")),
		vecty.Text("StopVideo()"),
		event.Click(a.stop),
	)
}

func (a *TestApp) controllerSeekTo() *vecty.HTML {
	return elem.Div(
		prop.Class("ui action tiny input"),
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
		prop.Class("ui action tiny input"),
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

func (a *TestApp) loadVideoByID(e *vecty.Event) {
	a.player.LoadVideoByID(a.idToLoad, a.startsSecond, a.selectedQuality)
}

func (a *TestApp) controllerLoadVideoByID() *vecty.HTML {
	// returns a list of options in markup
	getQualityOptions := func() []vecty.MarkupOrComponentOrHTML {
		qualities := []youtube.Quality{youtube.Small, youtube.Medium, youtube.Large, youtube.HD720, youtube.HD1080, youtube.HighRes}
		qhtml := make([]vecty.MarkupOrComponentOrHTML, 0, len(qualities))
		for _, quality := range qualities {
			chosen := false
			if quality == a.selectedQuality {
				chosen = true
			}
			qhtml = append(qhtml,
				elem.Option(
					vecty.If(chosen, vecty.Property("selected", true)),
					prop.Value(string(quality)),
					vecty.Text(string(quality)),
				),
			)
		}
		return qhtml
	}

	return elem.Div(
		prop.Class("ui action input"),
		elem.Input(
			prop.Type(prop.TypeText),
			prop.Value(a.idToLoad),
			prop.Placeholder("o_Ay_iDRAbc"),
			event.Input(func(e *vecty.Event) {
				a.idToLoad = e.Target.Get("value").String()
			}),
		),
		elem.Input(
			prop.Type(prop.TypeText),
			prop.Value(strconv.FormatFloat(a.startsSecond, 'f', 0, 64)),
			prop.Placeholder("0"),
			event.Input(func(e *vecty.Event) {
				a.startsSecond, _ = strconv.ParseFloat(
					e.Target.Get("value").String(), 64)
			}),
		),
		elem.Select(
			append(
				getQualityOptions(),
				prop.Class("ui compact selection dropdown"),
				event.Change(func(e *vecty.Event) {
					a.selectedQuality = youtube.Quality(e.Target.Get("value").String())
				}),
			)...,
		),
		elem.Button(
			prop.Class("ui right button"),
			vecty.Text("LoadVideoByID()"),
			event.Click(a.loadVideoByID),
		),
	)
}

func stat(desc string, value string) *vecty.HTML {
	return elem.TableRow(
		elem.TableData(
			prop.Class(""),
			vecty.Text(desc),
		),
		elem.TableData(vecty.Text(value)),
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

func (a *TestApp) currentTime() string {
	return strconv.FormatFloat(a.player.CurrentTime(), 'f', 2, 64)
}

func (a *TestApp) playbackQuality() string {
	return string(a.player.PlaybackQuality())
}

func (a *TestApp) duration() string {
	return strconv.FormatFloat(a.player.Duration(), 'f', 2, 64)
}

func (a *TestApp) availableQualityLevels() string {
	convert := func(qs []youtube.Quality) []string {
		res := make([]string, 0, len(qs))
		for _, q := range qs {
			res = append(res, string(q))
		}
		return res
	}
	return strings.Join(
		append(
			[]string{"["},
			strings.Join(convert(a.player.AvailableQualityLevels()), ", "),
			"]",
		), " ",
	)

}

func (a *TestApp) playlist() string {
	return strings.Join(
		append([]string{"["}, strings.Join(a.player.Playlist(), ", "), "]"),
		" ",
	)
}

func (a *TestApp) playlistIndex() string {
	return strconv.FormatInt(int64(a.player.PlaylistIndex()), 10)
}

func (a *TestApp) stats() *vecty.HTML {
	return elem.Table(
		prop.Class("ui padded small red table"),
		elem.TableHead(
			elem.TableRow(
				elem.TableHeader(
					vecty.Text("Caller (player.*)"),
				),
				elem.TableHeader(
					vecty.Text("Results"),
				),
			),
		),
		elem.TableBody(
			stat("Volume()", a.getVolume()),
			stat("IsMuted()", a.getMuted()),
			stat("PlayerState()", a.playerState()),
			stat("PlaybackRate()", a.playbackRate()),
			stat("VideoLoadedFraction()", a.videoLoadedFraction()),
			stat("CurrentTime()", a.currentTime()),
			stat("PlaybackQuality()", a.playbackQuality()),
			stat("Duration()", a.duration()),
			stat("VideoURL()", a.player.VideoURL()),
			stat("VideoEmbedCode()", a.player.VideoEmbedCode()),
			stat("AvailableQualityLevels()", a.availableQualityLevels()),
			stat("Playlist()", a.playlist()),
			stat("PlaylistIndex()", a.playlistIndex()),
			stat("Iframe()", a.player.Iframe().String()),
		),
	)
}

func (a *TestApp) getControllersColumn() *vecty.HTML {
	return elem.Div(
		prop.Class("column"),
		elem.Div(
			prop.Class("ui basic segment"),
			elem.Heading4(
				prop.Class("ui horizontal divider header"),
				elem.Italic(
					prop.Class(" setting icon"),
				),
				vecty.Text("Controls"),
			),
		),
		elem.Div(
			prop.Class("ui segments"),
			elem.Div(
				prop.Class("ui segment"),
				a.controllerPlay(),
				a.controllerPause(),
				a.controllerStop(),
			),
			elem.Div(
				prop.Class("ui segment"),
				a.controllerSeekTo(),
			),
			elem.Div(
				prop.Class("ui segment"),
				a.controllerVolume(),
			),
			elem.Div(
				prop.Class("ui segment"),
				a.controllerLoadVideoByID(),
			),
		),
	)
}

func (a *TestApp) getStatsColumn() *vecty.HTML {
	var statsDiv *vecty.HTML
	if a.showStat {
		statsDiv = a.stats()
	}
	if statsDiv == nil {
		return nil
	}
	return elem.Div(
		prop.Class("column"),
		elem.Div(
			prop.Class("ui basic segment"),
			elem.Heading4(
				prop.Class("ui horizontal divider header"),
				elem.Italic(
					prop.Class("bar chart icon"),
				),
				vecty.Text(
					strings.Join([]string{
						"Stats (every ",
						strconv.FormatInt(int64(a.StatUpdateFreq), 10),
						" milisecond)",
					}, ""),
				),
			),
			elem.Div(
				prop.Class("ui container"),
				statsDiv,
			),
		),
	)
}

var once = &sync.Once{}

func (a *TestApp) Render() *vecty.HTML {
	once.Do(func() {
		a.selectedQuality = youtube.Large
		var updFrq int
		if a.StatUpdateFreq == 0 {
			updFrq = 500
		} else {
			updFrq = a.StatUpdateFreq
		}
		// update stat every 1 second
		js.Global.Call("setInterval", func() { vecty.Rerender(a) }, updFrq) // reset stat
	})

	return elem.Body(
		elem.Div(
			prop.Class("two column stackable ui grid"),
			elem.Div(
				prop.Class("column"),
				elem.Div(
					prop.ID(playerID),
				),
			),
			a.getControllersColumn(),
			a.getStatsColumn(),
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

	app := &TestApp{
		StatUpdateFreq: 400,
	}
	app.showStat = false
	vecty.RenderBody(app)
	ytutil.OnLoaded(func() {
		var props = youtube.NewProperties()
		props.Width = 640
		props.Height = 360
		props.VideoID = "b-tAiOVMYFY"
		props.PlayerVars.EnableJsAPI = 1
		props.Events.OnReady = func(e *youtube.Event) {
			e.Target.PlayVideo()
		}

		app.player = youtube.NewPlayer(playerID, props)
		app.showStat = true
	})
}
