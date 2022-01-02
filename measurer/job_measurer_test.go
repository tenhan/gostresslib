package measurer

import (
	"fmt"
	"testing"
	"time"
)

func TestMeasurer_Run(t *testing.T) {
	m := NewJobMeasurer()
	type args struct {
		count       int
		concurrency int
		metricsName []string
		job         Job
	}
	tests := []struct {
		name   string
		args   args
	}{
		{
			name:   "0",
			args:   args{
				count:       0,
				concurrency: 0,
				metricsName: nil,
				job:         nil,
			},
		},
		{
			name:   "1",
			args:   args{
				count:       1,
				concurrency: 1,
				metricsName: nil,
				job: func(num int, metric *JobMetric) error {
					return nil
				},
			},
		},
		{
			name:   "2",
			args:   args{
				count:       100,
				concurrency: 1,
				metricsName: nil,
				job: func(num int, metric *JobMetric) error {
					return nil
				},
			},
		},
		{
			name:   "2",
			args:   args{
				count:       100,
				concurrency: 4,
				metricsName: nil,
				job: func(num int, metric *JobMetric) error {
					return nil
				},
			},
		},
		{
			name:   "2",
			args:   args{
				count:       100,
				concurrency: 4,
				metricsName: nil,
				job: func(num int, metric *JobMetric) error {
					return fmt.Errorf("test")
				},
			},
		},
		{
			name:   "2",
			args:   args{
				count:       100,
				concurrency: 4,
				metricsName: []string{"size"},
				job: func(num int, metric *JobMetric) error {
					metric.SetMetricsValue(float64(100))
					return fmt.Errorf("test")
				},
			},
		},
		{
			name:   "2",
			args:   args{
				count:       100,
				concurrency: 4,
				metricsName: []string{"step1","step2"},
				job: func(num int, metric *JobMetric) error {
					metric.SetMetricsValue(float64(time.Now().Unix()),float64(time.Now().Unix()))
					return fmt.Errorf("test")
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.Run(tt.args.count,tt.args.concurrency,tt.args.metricsName,tt.args.job)
		})
	}
}
