# gifview

Animated GIFs for TView-based TUIs (powered by pixelview)

![demo](https://user-images.githubusercontent.com/622065/75621100-3d726680-5b45-11ea-9fa0-8a3153461789.gif)

## Usage

There are two ways to create a new GifView:

```go
// From an existing gif.GIF object
gifImg := &gif.GIF{}
img, err := gifview.FromImage(gifImg)

// From a file path
gifPath := "images/dancing-baby.gif"
img, err := gifview.FromImagePath(gifPath)
```

Once you have one or more GifViews, they will animate whenever the application re-draws.  You can force that to happen on a regular basis by using the `Animate` function.

```go
app := tview.NewApplication()
go gifview.Animate(app)
```

## Based on

* [tview](https://github.com/rivo/tview)
* [pixelview](https://github.com/Omnikron13/pixelview)
* [tcell](https://github.com/gdamore/tcell)
