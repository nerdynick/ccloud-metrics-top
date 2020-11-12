package widgets

import (
	"fmt"
	"time"

	"github.com/gizak/termui/v3/widgets"
	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	log "github.com/sirupsen/logrus"
)

const (
	MaxDatapointsBar  = 20
	MaxDatapointsPlot = 120
)

var (
	GranularityIncr = map[ccloudmetrics.Granularity]time.Duration{
		ccloudmetrics.GranularityOneMin:     1 * time.Minute,
		ccloudmetrics.GranularityFiveMin:    5 * time.Minute,
		ccloudmetrics.GranularityFifteenMin: 15 * time.Minute,
		ccloudmetrics.GranularityThirtyMin:  30 * time.Minute,
		ccloudmetrics.GranularityOneHour:    1 * time.Hour,
	}
)

type CCloudMetricWidget interface {
	Update(*ccloudmetrics.MetricsClient, ccloudmetrics.Granularity, string, string, string)
}

type CCloudMetricsBarChart struct {
	widgets.BarChart
	metric ccloudmetrics.Metric
}

func NewCCloudMetricsBarChart(metric ccloudmetrics.Metric) *CCloudMetricsBarChart {
	m := &CCloudMetricsBarChart{
		BarChart: *widgets.NewBarChart(),
		metric:   metric,
	}
	m.Data = ProcessResult(nil, MaxDatapointsBar)

	return m
}

func (w *CCloudMetricsBarChart) Update(api *ccloudmetrics.MetricsClient, gran ccloudmetrics.Granularity, cluster, topic, requestType string) {
	go func() {
		data := ProcessResult(getMetric(w.metric, api, gran, cluster, topic, requestType, MaxDatapointsBar), MaxDatapointsBar)
		w.Data = data
	}()
}

type CCloudMetricsPlot struct {
	widgets.Plot
	metrics []ccloudmetrics.Metric
}

func NewCCloudMetricsPlot(metrics ...ccloudmetrics.Metric) *CCloudMetricsPlot {
	m := &CCloudMetricsPlot{
		Plot:    *widgets.NewPlot(),
		metrics: metrics,
	}
	m.DataLabels = make([]string, len(metrics))
	m.Data = make([][]float64, len(metrics))

	for i, me := range metrics {
		m.DataLabels[i] = me.ShortName()
		m.Data[i] = ProcessResult(nil, MaxDatapointsPlot)
	}

	return m
}

func (w *CCloudMetricsPlot) Update(api *ccloudmetrics.MetricsClient, gran ccloudmetrics.Granularity, cluster, topic, requestType string) {
	go func() {
		labels, data := getMetrics(w.metrics, api, gran, cluster, topic, requestType, MaxDatapointsPlot)
		w.Data = data
		w.DataLabels = labels
	}()
}

// type CCloudMetricsSparklineGroup struct {
// 	widgets.SparklineGroup
// 	metrics []ccloudmetrics.Metric
// }

func GetTimeRange(gran ccloudmetrics.Granularity, maxDatapoints int) (time.Time, time.Time) {
	now := time.Now().Round(time.Minute)
	return now.Add(-(GranularityIncr[gran] * time.Duration(maxDatapoints))), now
}

func ProcessResult(result []ccloudmetrics.QueryData, maxDatapoints int) []float64 {
	data := make([]float64, maxDatapoints)
	for i, d := range result {
		data[i] = d.Value
	}
	return data
}

func getMetric(metric ccloudmetrics.Metric, api *ccloudmetrics.MetricsClient, gran ccloudmetrics.Granularity, cluster, topic, requestType string, maxDatapoints int) []ccloudmetrics.QueryData {
	log.Info("Updating Metric: " + metric.Name)

	sTime, eTime := GetTimeRange(gran, maxDatapoints)
	exSTime := time.Now()
	var (
		results []ccloudmetrics.QueryData
		err     error
	)

	if topic != "" {
		results, err = api.QueryMetricAndTopic(cluster, metric, topic, gran, sTime, eTime, false)
	} else if requestType != "" {
		results, err = api.QueryMetricAndType(cluster, metric, requestType, gran, sTime, eTime)
	} else {
		results, err = api.QueryMetric(cluster, metric, gran, sTime, eTime)
	}

	exTTime := time.Now().Sub(exSTime)
	if err != nil {
		log.WithError(err).Error("Failed to get metrics")
		return nil
	}
	log.Info(fmt.Sprintf("Results Returned in %f secs for %s. Parsing and Rendering", exTTime.Seconds(), metric.Name))
	return results
}

func getMetrics(metrics []ccloudmetrics.Metric, api *ccloudmetrics.MetricsClient, gran ccloudmetrics.Granularity, cluster, topic, requestType string, maxDatapoints int) ([]string, [][]float64) {
	log.Info("Updating Metrics: " + metrics[0].Name)

	sTime, eTime := GetTimeRange(gran, maxDatapoints)
	exSTime := time.Now()

	results, err := api.QueryMetrics(cluster, metrics, gran, sTime, eTime)

	exTTime := time.Now().Sub(exSTime)
	if err != nil {
		log.WithError(err).Error("Failed to get metrics")
		return nil, nil
	}
	log.Info(fmt.Sprintf("Results Returned in %f secs for %s. Parsing and Rendering", exTTime.Seconds(), metrics[0].Name))

	stepResults := map[string][]ccloudmetrics.QueryData{}
	finalResults := [][]float64{}
	labels := []string{}

	for _, r := range results {
		stepResults[r.Metric] = append(stepResults[r.Metric], r)
	}

	for i, r := range stepResults {
		labels = append(labels, i)
		finalResults = append(finalResults, ProcessResult(r, maxDatapoints))
	}

	return labels, finalResults
}
