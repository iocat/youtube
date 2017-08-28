# youtube
Gopherjs Bindings for Youtube's Iframe API documented @ https://developers.google.com/youtube/iframe_api_reference

#### This binding library is currently under development. There will be breaking changes. 
The test is being written.


Usage: check the example package

```go

const (
	youtubeIframeAPISrc = "https://www.youtube.com/iframe_api"
	playerID            = "playerID"
)

func main() {
	// 1. Add the Youtube script to your header
	YourAddScriptFunc(youtubeIframeAPISrc)

	// 2. Place the player container as a div in your html with a predefined id
	var app = &Body{}
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
		// Cache the created player
		app.player = youtube.NewPlayer(playerID, props)
	})
}

type Body struct {
	vecty.Core
	player *youtube.Player
}

func (b *Body) Render() *vecty.HTML {
	return elem.Body(
		elem.Div(
			prop.ID(playerID),
		),
	)
}

```


Please feel free to contribute!
