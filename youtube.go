package youtube

import (
	"github.com/gopherjs/gopherjs/js"
)

type ProgessBarColor string

const (
	Red   ProgessBarColor = "red"
	White ProgessBarColor = "white"
)

type ControlsMode int

const (
	// ControlsNotDisplay has the player controls not displayed in the player
	ControlsNotDisplay ControlsMode = iota
	// ControlsDisplayImmediately has the controls displayed immediately and
	// the Flash player also loads immediately.
	ControlsDisplayImmediately
	// ControlsDisplayAfter has the controls displayed, and the Flash player
	// loaded after the user initiates the video playback.
	ControlsDisplayAfter
)

const (
	// IvPolicyShown shows video's annotation
	IvPolicyShown = iota + 1
	_
	// IvPolicyNotShown hides video's annotation
	IvPolicyNotShown
)

type ListType string

const (
	// ListTypePlaylist represents playlist
	ListTypePlaylist ListType = "playlist"
	// ListTypeSearch represents search list
	ListTypeSearch ListType = "search"
	// ListTypeUserUploads represents user uploads
	ListTypeUserUploads ListType = "user_uploads"
)

// Quality represents the player video's quality
type Quality string

const (
	Small   Quality = "small"
	Medium  Quality = "medium"
	Large   Quality = "large"
	HD720   Quality = "hd720"
	HD1080  Quality = "hd1080"
	HighRes Quality = "highres"
)

//
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
	OnPlaybackQualityChange EventType = "onPlaybackQualityChange"
	OnPlaybackRateChange    EventType = "onPlaybackRateChange"
	OnError                 EventType = "onError"
	OnApiChange             EventType = "onApiChange"
)

type Event struct {
	*js.Object
	Target *Player    `js:"target"`
	Data   *js.Object `js:"data"`
}

// Player represents the Youtube Iframe player
type Player struct {
	*js.Object
}

// Properties represents a set of video properties feeded to NewPlayer(id, properties) call
// to create the player. NewProperties() is recommended to create the properties.
type Properties struct {
	*js.Object
	Width        int           `js:"width"`
	Height       int           `js:"height"`
	VideoID      string        `js:"videoId"`
	PlayerVars   *PlayerParams `js:"playerVars"`
	PlayerEvents *PlayerEvents `js:"events"`
}

func newObj() *js.Object {
	return js.Global.Get("Object").New()
}

// NewProperties creates a new Property JS object
// with all inner objects properly initialized
func NewProperties() *Properties {
	props := &Properties{Object: newObj()}
	vars := &PlayerParams{Object: newObj()}
	eves := &PlayerEvents{Object: newObj()}
	props.PlayerVars = vars
	props.PlayerEvents = eves
	return props
}

// PlayerEvents contains a set of callbacks assigned at the creation of the
// player. This struct's fields correspond to each youtube.EventType
type PlayerEvents struct {
	*js.Object
	OnReady                 func(*Event) `js:"onReady"`
	OnStateChange           func(*Event) `js:"onStateChange"`
	OnPlaybackQualityChange func(*Event) `js:"onPlaybackQualityChange"`
	OnPlaybackRateChange    func(*Event) `js:"onPlaybackRateChange"`
	OnError                 func(*Event) `js:"onError"`
	OnAPIChange             func(*Event) `js:"onApiChange"`
}

// PlayerParams represents the player parameter documented at
// https://developers.google.com/youtube/player_parameters
type PlayerParams struct {
	*js.Object
	Autoplay       int             `js:"autoplay"`
	CcLoadPolicy   int             `js:"cc_load_policy"`
	Color          ProgessBarColor `js:"color"`
	Controls       ControlsMode    `js:"controls"`
	DisableKB      int             `js:"disablekb"`
	EnableJsAPI    int             `js:"enablejsapi"`
	End            int             `js:"end"`
	Fs             int             `js:"fs"`
	Hl             string          `js:"hl"`
	IvLoadPolicy   int             `js:"iv_load_policy"`
	List           string          `js:"list"`
	ListType       ListType        `js:"listType"`
	Loop           int             `js:"loop"`
	ModestBranding int             `js:"modestbranding"`
	Origin         string          `js:"origin"`
	Playlist       []string        `js:"playlist"`
	PlaysInline    int             `js:"playsinline"`
	ShowInfo       int             `js:"showinfo"`
	Start          int             `js:"start"`
	WidgetReferrer string          `js:"widget_referrer"`
}

// LoadByIDOptions represents an argument for Player.LoadVideoByID2(arg)
// and Player.CueVideoByID2(arg)
type LoadByIDOptions struct {
	*js.Object
	VideoID          string  `js:"videoId"`
	StartSeconds     float64 `js:"startSeconds"`
	EndSeconds       float64 `js:"endSeconds"`
	SuggestedQuality Quality `js:"suggestedQuality"`
}

// NewLoadByIDOptions prepares an argument for
// Player.LoadVideoByID2(arg) and Player.CueVideoByID2(arg)
// with all inner objects properly initialized
func NewLoadByIDOptions() *LoadByIDOptions {
	return &LoadByIDOptions{
		Object: newObj(),
	}
}

// LoadByURLOptions represents an argument for Player.LoadVideoByUrl2(arg)
type LoadByURLOptions struct {
	*js.Object
	MediaContentURL  string  `js:"mediaContentUrl"`
	StartSeconds     float64 `js:"startSeconds"`
	EndSeconds       float64 `js:"endSeconds"`
	SuggestedQuality Quality `js:"suggestedQuality"`
}

// NewLoadByURLOptions returns a prepared argument for Player.LoadeVideoByUrl2(arg)
// with all inner objects properly initialized
func NewLoadByURLOptions() *LoadByURLOptions {
	return &LoadByURLOptions{
		Object: newObj(),
	}
}

// CuePlaylistOptions represents an argument for Player.CuePlaylist2(arg)
type CuePlaylistOptions struct {
	*js.Object
	ListType ListType `js:"listType"`
	// the ID of the list
	List             string  `js:"list"`
	Index            int     `js:"index"`
	startSeconds     float64 `js:"startSeconds"`
	SuggestedQuality Quality `js:"suggestedQuality"`
}

// NewCuePlaylistOptions prepares an argument for Player.CuePlaylist2(arg)
// with all inner objects properly initialized
func NewCuePlaylistOptions() *CuePlaylistOptions {
	return &CuePlaylistOptions{
		Object: newObj(),
	}
}

// NewPlayer creates a new youtube player by replacing the
// provided iframe with the id of iframeId
// This call is equivalent to new YT.Player(id, props)
func NewPlayer(iframeID string, props *Properties) *Player {
	np := js.Global.Get("YT").Get("Player").New(iframeID, props.Object)

	return &Player{
		Object: np,
	}
}

// UPDATE PLAYER CONTENT FUNCTIONS

func (p *Player) LoadVideoByID(vid string, startSec float64, q Quality) {
	p.Call("loadVideoById", vid, startSec, q)
}

func (p *Player) LoadVideoByID2(params *LoadByIDOptions) {
	p.Call("loadVideoById", params)
}

func (p *Player) CueVideoByID(vid string, startSec float64, q Quality) {
	p.Call("cueVideoById", vid, startSec, q)
}

func (p *Player) CueVideoByID2(params *LoadByIDOptions) {
	p.Call("cueVideoById", params)
}

func (p *Player) LoadVideoByURL(url string, startSec float64, q Quality) {
	p.Call("loadVideoByUrl", url, startSec, q)
}

func (p *Player) LoadVideoByURL2(params *LoadByURLOptions) {
	p.Call("loadVideoByUrl", params)
}

func (p *Player) CuePlaylist(ids []string, index int, startSec float64, q Quality) {
	p.Call("cuePlaylist", ids, index, startSec, q)
}

func (p *Player) CuePlaylist2(params *CuePlaylistOptions) {
	p.Call("cuePlaylist", params)
}

func (p *Player) LoadPlaylist(ids []string, index int, startSec float64, q Quality) {
	p.Call("cuePlaylist", ids, index, startSec, q)
}

func (p *Player) LoadPlaylist2(params *CuePlaylistOptions) {
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
