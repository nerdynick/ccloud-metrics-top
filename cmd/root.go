package cmd

import (
	"fmt"
	"os"

	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	hotKeys = [][]string{
		[]string{"q", "Quit"},
	}
	rootCmd = &cobra.Command{
		Use:   "ccloud-metrics-top",
		Short: "Confluent Cloud Metrics Top",
		RunE: func(cmd *cobra.Command, args []string) error {
			uiContext.activeDashboard = UIDashboardCluster
			return UiMainLoop()
		},
	}
	topicCmd = &cobra.Command{
		Use:   "topic",
		Short: "Watch topic level metrics on a given cluster",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			uiContext.activeDashboard = UIDashboardTopic
			return UiMainLoop()
		},
	}
	typeCmd = &cobra.Command{
		Use:   "request",
		Short: "Watch a given request type on a given cluster",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			uiContext.activeDashboard = UIDashboardRequest
			return UiMainLoop()
		},
	}
)

//Global Vars
var (
	apiContext  ccloudmetrics.APIContext  = ccloudmetrics.NewAPIContext("", "")
	httpContext ccloudmetrics.HTTPContext = ccloudmetrics.NewHTTPContext()
	apiClient   ccloudmetrics.MetricsClient
)

type BlackholeWriter struct {
}

func (wr *BlackholeWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(&BlackholeWriter{})

	cobra.OnInitialize(func() {
		apiClient = ccloudmetrics.NewClientFromContext(apiContext, httpContext)
	})
	rootCmd.AddCommand(topicCmd, typeCmd)

	//Root Commands
	rootCmd.PersistentFlags().StringVarP(&apiContext.APIKey, "api-key", "k", "", "Service Account API Key")
	rootCmd.MarkPersistentFlagRequired("api-key")
	rootCmd.PersistentFlags().StringVarP(&apiContext.APISecret, "api-secret", "s", "", "Service Account API Secret")
	rootCmd.MarkPersistentFlagRequired("api-secret")
	rootCmd.PersistentFlags().StringVarP(&queryContext.ClusterId, "cluster", "c", "", "Confluent Cloud Cluster ID. Ex: lkc-ex123")
	rootCmd.MarkPersistentFlagRequired("cluster")
	rootCmd.PersistentFlags().StringVarP(&apiContext.BaseURL, "baseurl", "b", ccloudmetrics.DefaultBaseURL, "API Base Url")
	rootCmd.PersistentFlags().IntVarP(&httpContext.RequestTimeout, "timeout", "t", ccloudmetrics.DefaultRequestTimeout, "HTTP Request Timeout")
	rootCmd.PersistentFlags().StringVarP(&httpContext.UserAgent, "agent", "a", "ccloud-metrics-top/go-cli", "HTTP User Agent")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
