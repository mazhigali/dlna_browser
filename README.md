# dlna_browser
This is simple dlna_browser for files in active directory

## Getting Started

My favorite mediaplayer MPV supports UPNP/DLNA network playback, but it hasn't support file browsering active directory. So i use this porgram.
You can browse your dlna files and copy file url to clipboard by pressing "Enter".
Quit : "Q" or "Esc"

Then you can give a file url of dlna file to MPV

```
mpv "paste url from clipboard"
```


### Gow to run

You can run compiled version

```
./dlna_browser
```
Or use Go

```
go run main.go
```
