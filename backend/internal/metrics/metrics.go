package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

type CustomGather struct {
	g       prometheus.Gatherer
	appName string
}

var (
	SyncImage = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sync_image_task",
		Help: "The total number of image sync task",
	})
	SyncImageSucc = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sync_image_succ",
		Help: "Successful number of image sync task",
	})
	SyncImageFail = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sync_image_fail",
		Help: "Failed number of image sync task",
	})

	JobsProcessed = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "sync_task_processed",
			Help: "Total number of jobs processed by the workers",
		},
	)
	RunningJobs = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "sync_task_running",
			Help: "Number of jobs inflight",
		},
	)
	RunningWorkers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "sync_worker_running",
			Help: "Number of jobs inflight",
		},
	)
	ProcessingTime = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name: "sync_process_time",
			Help: "Amount of time spent processing jobs",
		},
	)
)

func NewMetrics(appName string) http.Handler {
	// 定义一个采集数据的采集器集合，它可以合并多个不同的采集器数据导一个数据集合中
	cg := &CustomGather{
		g:       prometheus.DefaultGatherer,
		appName: appName,
	}

	return promhttp.HandlerFor(cg, promhttp.HandlerOpts{
		// ErrorLog:      logger,                   //采集过程中如果出现错误，记录日志
		ErrorHandling: promhttp.ContinueOnError, //采集过程中如果出现错误，继续采集其他指标，不会中断采集任务
	})
}

func (c *CustomGather) Gather() ([]*dto.MetricFamily, error) {
	if dtos, err := c.g.Gather(); err != nil {
		return nil, err
	} else {
		for _, d := range dtos {
			for _, m := range d.Metric {
				m.Label = append(m.Label, &dto.LabelPair{
					Name:  proto.String("app_name"),
					Value: proto.String(c.appName),
				})
			}
		}
		return dtos, nil
	}
}
