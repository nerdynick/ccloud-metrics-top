package cmd

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/nerdynick/ccloud-metrics-top/widgets"
	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	log "github.com/sirupsen/logrus"
)

var (
	clusterHotKeys = []widgets.HotKey{}
)

type ClusterGraphSet struct {
	ui.Grid
	Granularity       ccloudmetrics.Granularity
	apiClient         *ccloudmetrics.MetricsClient
	Bytes             *widgets.CCloudMetricsPlot
	Records           *widgets.CCloudMetricsPlot
	RetainedBytes     *widgets.CCloudMetricsPlot
	Partitions        *widgets.CCloudMetricsPlot
	ActiveConnections *widgets.CCloudMetricsPlot
	queryContext      *QueryContext
	api               *ccloudmetrics.MetricsClient
	hotkeys           []widgets.HotKey
}

func (g *ClusterGraphSet) Graph() interface{} {
	return g
}

func (g *ClusterGraphSet) HotKeys() []widgets.HotKey {
	return g.hotkeys
}

func (g *ClusterGraphSet) GraphTitle() string {
	return fmt.Sprintf("Cluster: %s, Granularity: %s", g.queryContext.ClusterId, g.queryContext.Granulatory)
}

func (g *ClusterGraphSet) Update() {
	log.Info("Updating Cluster Metrics")
	g.Bytes.Update(g.api, g.queryContext.Granulatory, g.queryContext.ClusterId, g.queryContext.Topic, g.queryContext.Type)
	g.Records.Update(g.api, g.queryContext.Granulatory, g.queryContext.ClusterId, g.queryContext.Topic, g.queryContext.Type)
	g.RetainedBytes.Update(g.api, g.queryContext.Granulatory, g.queryContext.ClusterId, g.queryContext.Topic, g.queryContext.Type)
	g.Partitions.Update(g.api, g.queryContext.Granulatory, g.queryContext.ClusterId, g.queryContext.Topic, g.queryContext.Type)
	g.ActiveConnections.Update(g.api, g.queryContext.Granulatory, g.queryContext.ClusterId, g.queryContext.Topic, g.queryContext.Type)
}

func NewClusterGraphSet(queryContext *QueryContext, api *ccloudmetrics.MetricsClient) *ClusterGraphSet {
	graph := &ClusterGraphSet{
		Grid:              *ui.NewGrid(),
		Bytes:             widgets.NewCCloudMetricsPlot(ccloudmetrics.MetricReceivedBytes, ccloudmetrics.MetricSentBytes),
		Records:           widgets.NewCCloudMetricsPlot(ccloudmetrics.MetricReceivedRecords, ccloudmetrics.MetricSentRecords),
		RetainedBytes:     widgets.NewCCloudMetricsPlot(ccloudmetrics.MetricRetainedBytes),
		Partitions:        widgets.NewCCloudMetricsPlot(ccloudmetrics.MetricRetainedBytes),
		ActiveConnections: widgets.NewCCloudMetricsPlot(ccloudmetrics.MetricActiveConnections),
		queryContext:      queryContext,
		api:               api,
	}
	graph.Bytes.Title = "Bytes In/Out"
	graph.Bytes.LineColors = []ui.Color{ui.ColorGreen, ui.ColorRed}
	graph.Records.Title = "Records In/Out"
	graph.Records.LineColors = []ui.Color{ui.ColorGreen, ui.ColorRed}
	graph.RetainedBytes.Title = "Retained Bytes"
	graph.Partitions.Title = "Partitions"
	graph.ActiveConnections.Title = "ActiveConnections"

	graph.Set(ui.NewRow(1.0,
		ui.NewRow(1.0/3*2, ui.NewCol(0.5, graph.Bytes), ui.NewCol(0.5, graph.Records)),
		ui.NewRow(1.0/3, ui.NewCol(0.25, graph.RetainedBytes), ui.NewCol(0.25, graph.Partitions), ui.NewCol(0.25, graph.ActiveConnections)),
	))

	return graph
}
