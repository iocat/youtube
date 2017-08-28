// Package ytutil contains utility functions for setting the Youtube Iframe API
// The application, however, can be set up without using this package.
package ytutil

import "github.com/gopherjs/gopherjs/js"

const youtubeIframeAPISrc = "https://www.youtube.com/iframe_api"

// Load loads the Youtube Iframe API script into the head of the HTML document
func Load() {
	addScript(youtubeIframeAPISrc)
}

// OnLoaded registers a callback executed when the Youtube Iframe API is ready
func OnLoaded(fn func()) {
	js.Global.Set("onYouTubeIframeAPIReady", fn)
}

func addScript(url string) {
	script := js.Global.Get("document").Call("createElement", "script")
	script.Set("src", url)
	script.Set("type", "text/javascript")
	js.Global.Get("document").Get("head").Call("appendChild", script)
}
