package youtube

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	Red   = "red"
	White = "white"
)

const (
	ControlsNotDisplay = iota
	ControlsDisplayImmediately
	ControlsDisplayAfter
)

const (
	IvPolicyShown = iota + 1
	_
	IvPolicyNotShown
)

type ListType string

const (
	ListTypePlaylist    ListType = "playlist"
	ListTypeSearch      ListType = "search"
	ListTypeUserUploads ListType = "user_uploads"
)

type Quality string

const (
	Small   Quality = "small"
	Medium  Quality = "medium"
	Large   Quality = "large"
	HD720   Quality = "hd720"
	HD1080  Quality = "hd1080"
	HighRes Quality = "highres"
)

type PlayerStatus int

const (
	Unstarted PlayerStatus = iota - 1
	Ended
	Playing
	Paused
	Buffering
	VideoCued
)

type EventType string

const (
	OnReady                 EventType = "onReady"
	OnStateChange           EventType = "onStateChange"
	OnPlaybackQualityChange           = "onPlaybackQualityChange"
	OnPlaybackRateChange              = "onPlaybackRateChange"
	OnError                           = "onError"
	OnApiChange                       = "onApiChange"
)

type Event struct {
	*js.Object
	Target *Player    `js:"target"`
	Data   *js.Object `js:"data"`
}

type Player struct {
	*js.Object
}

type properties struct {
	*js.Object
	Width        int               `js:"width"`
	Height       int               `js:"height"`
	VideoID      string            `js:"videoId"`
	PlayerVars   *playerParameters `js:"playerVars"`
	PlayerEvents *playerEvents     `js:"events"`
}

func newObj() *js.Object {
	return js.Global.Get("Object").New()
}

func NewProperties() *properties {
	props := &properties{Object: newObj()}
	vars := &playerParameters{Object: newObj()}
	eves := &playerEvents{Object: newObj()}
	props.PlayerVars = vars
	props.PlayerEvents = eves
	return props
}

type playerEvents struct {
	*js.Object
	OnReady       func(*Event) `js:"onReady"`
	OnStateChange func(*Event) `js:"onStateChange"`
}

// The player parameter documented at
// https://developers.google.com/youtube/player_parameters
type playerParameters struct {
	*js.Object
	Autoplay       int      `js:"autoplay"`
	CcLoadPolicy   int      `js:"cc_load_policy"`
	Color          string   `js:"color"`
	Controls       string   `js:"controls"`
	DisableKB      int      `js:"disablekb"`
	EnableJsAPI    int      `js:"enablejsapi"`
	End            int      `js:"end"`
	Fs             int      `js:"fs"`
	Hl             string   `js:"hl"`
	IvLoadPolicy   int      `js:"iv_load_policy"`
	List           string   `js:"list"`
	ListType       ListType `js:"listType"`
	Loop           int      `js:"loop"`
	ModestBranding int      `js:"modestbranding"`
	Origin         string   `js:"origin"`
	Playlist       []string `js:"playlist"`
	PlaysInline    int      `js:"playsinline"`
	ShowInfo       int      `js:"showinfo"`
	Start          int      `js:"start"`
	WidgetReferrer string   `js:"widget_referrer"`
}

type loadByIDOptions struct {
	*js.Object
	VideoID          string  `js:"videoId"`
	StartSeconds     float64 `js:"startSeconds"`
	EndSeconds       float64 `js:"endSeconds"`
	SuggestedQuality Quality `js:"suggestedQuality"`
}

func NewLoadByIDOptions() *loadByIDOptions {
	return &loadByIDOptions{
		Object: newObj(),
	}
}

type loadByURLOptions struct {
	*js.Object
	MediaContentURL  string  `js:"mediaContentUrl"`
	StartSeconds     float64 `js:"startSeconds"`
	EndSeconds       float64 `js:"endSeconds"`
	SuggestedQuality Quality `js:"suggestedQuality"`
}

func NewLoadByURLOptions() *loadByURLOptions {
	return &loadByURLOptions{
		Object: newObj(),
	}
}

type cuePlaylistOptions struct {
	*js.Object
	ListType         ListType `js:"listType"`
	List             string   `js:"list"`
	Index            int      `js:"index"`
	startSeconds     float64  `js:"startSeconds"`
	SuggestedQuality Quality  `js:"suggestedQuality"`
}

func NewCuePlaylistOptions() *cuePlaylistOptions {
	return &cuePlaylistOptions{
		Object: newObj(),
	}
}

// Player creates a new youtube player by replacing the
// provided iframe with the id of iframeId
func NewPlayer(iframeID string, props *properties) *Player {
	np := js.Global.Get("YT").Get("Player").New(iframeID, props.Object)

	return &Player{
		Object: np,
	}
}

// UPDATE PLAYER CONTENT FUNCTIONS

func (p *Player) LoadVideoByID(vid string, startSec float64, q Quality) {
	p.Call("loadVideoById", vid, startSec, q)
}

func (p *Player) LoadVideoByID2(params *loadByIDOptions) {
	p.Call("loadVideoById", params)
}

func (p *Player) CueVideoByID(vid string, startSec float64, q Quality) {
	p.Call("cueVideoById", vid, startSec, q)
}

func (p *Player) CueVideoByID2(params *loadByIDOptions) {
	p.Call("cueVideoById", params)
}

func (p *Player) LoadVideoByURL(url string, startSec float64, q Quality) {
	p.Call("loadVideoByUrl", url, startSec, q)
}

func (p *Player) LoadVideoByURL2(params *loadByURLOptions) {
	p.Call("loadVideoByUrl", params)
}

func (p *Player) CuePlaylist(ids []string, index int, startSec float64, q Quality) {
	p.Call("cuePlaylist", ids, index, startSec, q)
}

func (p *Player) CuePlaylist2(params *cuePlaylistOptions) {
	p.Call("cuePlaylist", params)
}

func (p *Player) LoadPlaylist(ids []string, index int, startSec float64, q Quality) {
	p.Call("cuePlaylist", ids, index, startSec, q)
}

func (p *Player) LoadPlaylist2(params *cuePlaylistOptions) {
	p.Call("cuePlaylist", params)
}

// Playback controls and player settings

func (p *Player) PlayVideo() {
	p.Call("playVideo")
}

func (p *Player) PauseVideo() {
	p.Call("pauseVideo")
}

func (p *Player) StopVideo() {
	p.Call("stopVideo")
}

func (p *Player) SeekTo(seconds float64, allowSeekAhead bool) {
	p.Call("seekTo", seconds, allowSeekAhead)
}

func (p *Player) NextVideo() {
	p.Call("nextVideo")
}

func (p *Player) PreviousVideo() {
	p.Call("previousVideo")
}

func (p *Player) PlayVideoAt(index int) {
	p.Call("playVideoAt", index)
}

func (p *Player) Mute() {
	p.Call("mute")
}

func (p *Player) UnMute() {
	p.Call("unMute")
}

func (p *Player) IsMuted() bool {
	return p.Call("isMuted").Bool()
}

func (p *Player) SetVolume(vol int) {
	p.Call("setVolume", vol)
}

func (p *Player) Volume() int {
	return p.Call("getVolume").Int()
}

func (p *Player) SetSize(width int, height int) *js.Object {
	return p.Call("setSize", width, height)
}

func (p *Player) PlaybackRate() float64 {
	return p.Call("getPlaybackRate").Float()
}

func (p *Player) SetPlaybackRate(suggestedRate float64) {
	p.Call("setPlaybackRate", suggestedRate)
}

// AvailableRates returns the set of playback rates in which the current video
// is available
func (p *Player) AvailablePlaybackRates() []float64 {
	rates := p.Call("getAvailablePlaybackRates")
	length := rates.Length()
	res := make([]float64, 0, length)
	for i := 0; i < length; i++ {
		res = append(res, rates.Index(i).Float())
	}
	return res
}

func (p *Player) SetLoop(val bool) {
	p.Call("setLoop", val)
}

func (p *Player) SetShuffle(val bool) {
	p.Call("shufflePlaylist", val)
}

func (p *Player) VideoLoadedFraction() float64 {
	return p.Call("getVideoLoadedFraction").Float()
}

func (p *Player) PlayerState() PlayerStatus {
	return PlayerStatus(p.Call("getPlayerState").Int())
}

func (p *Player) CurrentTime() float64 {
	return p.Call("getCurrentTime").Float()
}

func (p *Player) PlaybackQuality() Quality {
	return Quality(p.Call("getPlaybackQuality").String())
}

func (p *Player) SetPlaybackQuality(suggested Quality) {
	p.Call("setPlaybackQuality", suggested)
}

func (p *Player) AvailableQualityLevels() []Quality {
	aql := p.Call("getAvailableQualityLevels")
	q := make([]Quality, 0, aql.Length())
	for i := 0; i < aql.Length(); i++ {
		q = append(q, Quality(aql.Index(i).String()))
	}
	return q
}

func (p *Player) Duration() float64 {
	return p.Call("getDuration").Float()
}

func (p *Player) VideoUrl() string {
	return p.Call("getVideoUrl").String()
}

func (p *Player) VideoEmbedCode() string {
	return p.Call("getVideoEmbedCode").String()
}

// Retrieve Playlist Info

func (p *Player) Playlist() []string {
	ids := p.Call("getPlaylist")
	res := make([]string, 0, ids.Length())
	for i := 0; i < ids.Length(); i++ {
		res = append(res, ids.Index(i).String())
	}
	return res
}

func (p *Player) PlaylistIndex() int {
	return p.Call("getPlaylistIndex").Int()
}

func (p *Player) AddEventListener(event EventType, listener func(event *Event)) {
	p.Call("addEventListener", listener)
}

func (p *Player) RemoveEventListener(event EventType, listener func(event *Event)) {
	p.Call("removeEventListener", listener)
}

func (p *Player) Iframe() *js.Object {
	return p.Call("getIframe")
}

func (p *Player) Destroy() {
	p.Call("destroy")
}
