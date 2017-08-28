# youtube
Gopherjs Bindings for Youtube's Iframe API documented @ https://developers.google.com/youtube/iframe_api_reference

#### This binding library is currently under development. There will be breaking changes. 
The test is being written.


Usage:
Sample code is in the [example package](https://github.com/iocat/youtube/blob/master/example/main.go)

```go

const playerID            = "playerID"

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
		props.VideoID = "dQw4w9WgXcQ"
		props.PlayerVars.EnableJsAPI = 1
		props.PlayerEvents.OnReady = func(e *youtube.Event) {
			e.Target.PlayVideo()
		}
		// Create and cache the created player
		app.player = youtube.NewPlayer(playerID, props)
	})
}

```


Please feel free to contribute!
