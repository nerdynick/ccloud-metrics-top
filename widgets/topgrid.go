package widgets

import (
	"fmt"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	log "github.com/sirupsen/logrus"
)

const (
	HotKeyNothing HotKeyAction = "Nil"
	HotKeyExit    HotKeyAction = "Exit"
	HotKeyUpdate  HotKeyAction = "Update"
)

var (
	commonHotKeys = []HotKey{
		{
			Key: "q", Name: "Quit/Exit", Action: func() HotKeyAction {
				return HotKeyExit
			},
		},
	}
)

type TopGrid struct {
	ui.Grid
	HotKeys        []HotKey
	Title          *widgets.Paragraph
	titleInfo      map[string]string
	Logs           *LogrusList
	logsShow       bool
	ActiveGraphSet GraphSet
}

func NewTopGrid(graph interface{}, hotkeys []HotKey) *TopGrid {
	g := graph.(GraphSet)

	top := &TopGrid{
		Grid:           *ui.NewGrid(),
		HotKeys:        append(commonHotKeys, hotkeys...),
		Title:          widgets.NewParagraph(),
		Logs:           NewLogrusList(log.ErrorLevel, log.WarnLevel, log.InfoLevel),
		logsShow:       false,
		ActiveGraphSet: g,
	}
	log.AddHook(top.Logs)
	top.Logs.Title = "Debug Logs"

	top.HotKeys = append(top.HotKeys, HotKey{
		Key: "d", Name: "Toggle Debug Logs", Action: func() HotKeyAction {
			top.ToggleLogs()
			return HotKeyNothing
		},
	})

	top.updateGrid()

	return top
}

func (g *TopGrid) updateGrid() {
	g.Items = []*ui.GridItem{}

	if g.logsShow {
		g.Set(
			ui.NewRow(1.0/10, g.Title),
			ui.NewRow(1.0/10*8, g.ActiveGraphSet.Graph()),
			ui.NewRow(1.0/10, g.Logs),
		)
	} else {
		g.Set(
			ui.NewRow(1.0/10, g.Title),
			ui.NewRow(1.0/10*9, g.ActiveGraphSet.Graph()),
		)
	}
}

//ToggleLogs Hides and Shows the Logs widget at the bottom of the screen
func (g *TopGrid) ToggleLogs() bool {
	if g.logsShow {
		g.logsShow = false
	} else {
		g.logsShow = true
	}
	ui.Clear()
	g.updateGrid()
	g.ReDraw()
	return g.logsShow
}

//ReDraw Triggers a ReDrawing/ReRendering of the UI
func (g *TopGrid) ReDraw() {
	//Update Title String
	g.Title.Title = g.ActiveGraphSet.GraphTitle()

	//Update HotKey Help
	helpText := []string{}
	for _, k := range g.HotKeys {
		helpText = append(helpText, fmt.Sprintf("(%s) %s", k.Key, k.Name))
	}
	for _, k := range g.ActiveGraphSet.HotKeys() {
		helpText = append(helpText, fmt.Sprintf("(%s) %s", k.Key, k.Name))
	}
	g.Title.Text = strings.Join(helpText, " ")

	ui.Render(g)
}

func (g *TopGrid) HandleEvent(e ui.Event) HotKeyAction {
	switch e.ID { // event string/identifier
	// case "<MouseLeft>":
	// 	payload := e.Payload.(ui.Mouse)
	// 	x, y := payload.X, payload.Y
	case "<Resize>":
		payload := e.Payload.(ui.Resize)
		g.SetRect(0, 0, payload.Width, payload.Height)
		ui.Clear()
		g.ReDraw()
	}
	switch e.Type {
	case ui.KeyboardEvent: // handle all key presses
		for _, k := range g.HotKeys {
			if k.Key == e.ID {
				return k.Action()
			}
		}
		for _, k := range g.ActiveGraphSet.HotKeys() {
			if k.Key == e.ID {
				return k.Action()
			}
		}
	}

	return HotKeyNothing
}

//Update Triggers any updates to fetch metrics/stats from the Active GraphSet
func (g *TopGrid) Update() {
	g.ActiveGraphSet.Update()
	g.ReDraw()
}

type GraphSet interface {
	HotKeys() []HotKey
	GraphTitle() string
	Update()
	Graph() interface{}
}

type HotKeyAction string

type HotKey struct {
	Key    string
	Name   string
	Action func() HotKeyAction
}
