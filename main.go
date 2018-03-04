package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/atotto/clipboard"
	"github.com/huin/goupnp/dcps/av1"
	"github.com/marcusolsson/tui-go"
)

type DIDLLite struct {
	XMLName xml.Name
	DC      string   `xml:"xmlns:dc,attr"`
	UPNP    string   `xml:"xmlns:upnp,attr"`
	XSI     string   `xml:"xmlns:xsi,attr"`
	XLOC    string   `xml:"xsi:schemaLocation,attr"`
	Objects []Object `xml:"item"`
	Folders []Object `xml:"container"`
}

type Object struct {
	ID          string `xml:"id,attr"`
	Parent      string `xml:"parentID,attr"`
	Restricted  string `xml:"restricted,attr"`
	Title       string `xml:"title"`
	Creator     string `xml:"creator"`
	Class       string `xml:"class"`
	Date        string `xml:"date"`
	Description string `xml:"description"`
	Results     []Res  `xml:"res"`
}

type Res struct {
	Resolution      string `xml:"resolution,attr"`
	Size            uint64 `xml:"size,attr"`
	ProtocolInfo    string `xml:"protocolInfo,attr"`
	Duration        string `xml:"duration,attr"`
	Bitrate         string `xml:"bitrate,attr"`
	SampleFrequency uint64 `xml:"sampleFrequency"`
	NrAudioChannels uint64 `xml:"nrAudioChannels"`
	Value           string `xml:",chardata"`
}

func main() {

	progress := tui.NewProgress(100)

	Library := tui.NewTable(0, 0)
	Library.SetColumnStretch(0, 1)
	Library.SetColumnStretch(1, 1)
	Library.SetColumnStretch(2, 3)
	Library.SetColumnStretch(3, 2)

	Library.AppendRow(
		tui.NewLabel("ItemName"),
		tui.NewLabel("Type"),
		tui.NewLabel("Description"),
		tui.NewLabel("URL"),
	)

	Library.SetSelected(0)

	status := tui.NewStatusBar("")
	status.SetPermanentText(`Permanent`)

	root := tui.NewVBox(
		Library,
		tui.NewSpacer(),
		progress,
		status,
	)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}
	X := browsing("V_F^FOL*R1")
	showItems(X, Library)

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("q", func() { ui.Quit() })
	ui.SetKeybinding("Down", func() {
		numselectitem := Library.Selected()
		n := numselectitem + 1
		Library.SetSelected(n)
	})
	ui.SetKeybinding("Up", func() {
		numselectitem := Library.Selected()
		n := numselectitem - 1
		Library.SetSelected(n)
	})
	ui.SetKeybinding("Enter", func() {
		numselectitem := Library.Selected()
		//status.SetPermanentText(u.Title)
		if numselectitem == 0 {
			u := X[0]
			X = browsing(u.Parent)
			Library.RemoveRows()
			//showItems(X, Library)
			fmt.Println(u.Parent)
		} else {

			u := X[numselectitem-1]
			if u.Class == "object.container.storageFolder" {
				X = browsing(u.ID)
				Library.RemoveRows()
				Library.AppendRow(
					tui.NewLabel("UP"),
					tui.NewLabel("object.container.storageFolder"),
				)
				showItems(X, Library)
			}
			v := u.Results[0].Value
			status.SetText("COPYED TO CLIPBOARD " + v)
			clipboard.WriteAll(v)
		}
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func browsing(dirId string) []Object {

	clients, _, _ := av1.NewContentDirectory1Clients()
	client := clients[0]
	res, _, _, _, _ := client.Browse(dirId, "BrowseDirectChildren", "*", 0, 30, "")

	r := &DIDLLite{}
	err2 := xml.Unmarshal([]byte(res), r)
	X := []Object{}
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "Error while parsing result xml %v\n", err2)
	} else {
		X = append(r.Folders, r.Objects...)
	}
	return X
}

func showItems(x []Object, library *tui.Table) {

	for _, item := range x {
		value := "nothing to display"
		if len(item.Results) > 0 {
			value = item.Results[0].Value
		}
		library.AppendRow(
			tui.NewLabel(item.Title),
			tui.NewLabel(item.Class),
			tui.NewLabel(item.Description),
			tui.NewLabel(value),
		)
	}
}
