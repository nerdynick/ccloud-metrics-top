package cmd

import (
	"time"

	"github.com/nerdynick/ccloud-metrics-top/widgets"
	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"

	ui "github.com/gizak/termui/v3"
	log "github.com/sirupsen/logrus"
)

const (
	UIDashboardCluster = "Cluster"
	UIDashboardTopic   = "Topic"
	UIDashboardRequest = "Request"
)

var (
	GloballHotKeys = []widgets.HotKey{
		{
			Key: "1", Name: string(ccloudmetrics.GranularityOneMin), Action: func() widgets.HotKeyAction {
				queryContext.Granulatory = ccloudmetrics.GranularityOneMin
				return widgets.HotKeyUpdate
			},
		},
		{
			Key: "2", Name: string(ccloudmetrics.GranularityFiveMin), Action: func() widgets.HotKeyAction {
				queryContext.Granulatory = ccloudmetrics.GranularityFiveMin
				return widgets.HotKeyUpdate
			},
		},
		{
			Key: "3", Name: string(ccloudmetrics.GranularityFifteenMin), Action: func() widgets.HotKeyAction {
				queryContext.Granulatory = ccloudmetrics.GranularityFifteenMin
				return widgets.HotKeyUpdate
			},
		},
		{
			Key: "4", Name: string(ccloudmetrics.GranularityThirtyMin), Action: func() widgets.HotKeyAction {
				queryContext.Granulatory = ccloudmetrics.GranularityThirtyMin
				return widgets.HotKeyUpdate
			},
		},
		{
			Key: "5", Name: string(ccloudmetrics.GranularityOneHour), Action: func() widgets.HotKeyAction {
				queryContext.Granulatory = ccloudmetrics.GranularityOneHour
				return widgets.HotKeyUpdate
			},
		},
	}
	queryContext = QueryContext{
		Granulatory: ccloudmetrics.GranularityOneMin,
	}
	uiContext = UIContext{
		activeDashboard: UIDashboardCluster,
	}
	lastUpdated time.Time
)

type UIDashboard string

type UIContext struct {
	activeDashboard UIDashboard
}

type QueryContext struct {
	ClusterId   string
	Topic       string
	Type        string
	Granulatory ccloudmetrics.Granularity
}

func UiMainLoop() error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	log.Info("Creating Cluster UI")
	cluster := NewClusterGraphSet(&queryContext, &apiClient)

	log.Info("Creating TOP UI")
	topUI := widgets.NewTopGrid(cluster, GloballHotKeys)
	termWidth, termHeight := ui.TerminalDimensions()
	topUI.SetRect(0, 0, termWidth, termHeight)

	log.WithField("TopGrid", topUI).Info("Drawing Top UI")
	topUI.ReDraw()

	log.Info("Doing initial update of data")
	lastUpdated = time.Now()
	topUI.Update()

	log.Info("Begining initial polling loop")
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			log.WithField("Event", e).Info("Got UX Event")
			action := topUI.HandleEvent(e)
			if action == widgets.HotKeyExit {
				return nil
			} else if action == widgets.HotKeyUpdate {
				topUI.Update()
			}
		// use Go's built-in tickers for updating and drawing data
		case t := <-ticker:
			topUI.ReDraw()

			if t.Sub(lastUpdated) >= time.Minute {
				lastUpdated = t
				topUI.Update()
			}
		}
	}
}
