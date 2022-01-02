package measurer

import (
	"fmt"
	"math"
	"sync"
	"time"
)


type Job func(num int,metric *JobMetric) error

type MetricStatistic struct {
	Name string
	Total float64
	TotalSeconds float64
	Avg float64
	Min float64
	Max float64
	Stdev float64
	PerSec float64
}
type JobStatistic struct {
	RunTime time.Duration
	TotalDuration time.Duration
	PerSeconds float64
	Count int
	Concurrency int
	MetricStatistic []*MetricStatistic
}
type JobMetric struct {
	values []float64
}

func (j *JobMetric)SetMetricsValue(value... float64)  {
	for i,v := range value{
		j.values[i] = v
	}
}

func newJobMetric(metricsCount int)*JobMetric{
	m := &JobMetric{}
	m.values = make([]float64,metricsCount)
	return m
}
type JobMeasurer struct {
	jobDurations []time.Duration
	interMetricsName []string
	interJobMetrics []*JobMetric
	outerJobMetrics []*JobMetric
}


func NewJobMeasurer() *JobMeasurer {
	m := &JobMeasurer{
		interMetricsName: []string{MetricNameLatency,MetricNameError},
	}
	return m
}

func (m *JobMeasurer) Run(count int,concurrency int,metricsName []string,job Job) (ret JobStatistic){
	if count <= 0 || concurrency <= 0 || job == nil{
		return ret
	}
	m.interJobMetrics = make([]*JobMetric,count)
	m.outerJobMetrics = make([]*JobMetric,count)
	m.jobDurations = make([]time.Duration,count)
	for i:=0;i<count;i++{
		m.interJobMetrics[i] = newJobMetric(len(m.interMetricsName))
		m.outerJobMetrics[i] = newJobMetric(len(metricsName))
	}
	wg := sync.WaitGroup{}
	start := time.Now()
	for i:=0;i<concurrency;i++{
		wg.Add(1)
		go func(n int) {
			for k:=n;k<count;k+=concurrency{
				begin := time.Now()
				err := job(k,m.outerJobMetrics[k])
				m.jobDurations[k] = time.Since(begin)
				m.interJobMetrics[k].values[MetricIndexLatency] = m.jobDurations[k].Seconds()
				if err != nil{
					m.interJobMetrics[k].values[MetricIndexError] = 1.0
				}else{
					m.interJobMetrics[k].values[MetricIndexError] = 0.0
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	ret.RunTime = time.Since(start)
	ret.Count = count
	ret.Concurrency = concurrency
	for _,v := range m.jobDurations{
		ret.TotalDuration += v
	}
	if ret.TotalDuration > 0{
		ret.PerSeconds = float64(ret.Count) / ret.RunTime.Seconds()
	}
	for metricIndex,metricName := range m.interMetricsName{
		ret.MetricStatistic = append(ret.MetricStatistic,m.calMetricStatistic(metricIndex,metricName,ret.RunTime.Seconds(),m.interJobMetrics))
	}
	for metricIndex,metricName := range metricsName{
		ret.MetricStatistic = append(ret.MetricStatistic,m.calMetricStatistic(metricIndex,metricName,ret.RunTime.Seconds(),m.outerJobMetrics))
	}
	return ret
}

func (m *JobMeasurer) calMetricStatistic(metricIndex int,metricName string,totalSeconds float64,metric []*JobMetric) *MetricStatistic  {
	ret :=&MetricStatistic{
		Name: metricName,
		TotalSeconds: totalSeconds,
	}
	itemCount := 0.0
	for k,v := range metric{

		itemCount += 1.0
		if k == 0{
			ret.Min = v.values[metricIndex]
			ret.Max = v.values[metricIndex]
		}else{
			if ret.Min > v.values[metricIndex]{
				ret.Min = v.values[metricIndex]
			}
			if ret.Max < v.values[metricIndex]{
				ret.Max = v.values[metricIndex]
			}
		}
		ret.Total += v.values[metricIndex]
	}
	if ret.Total == 0 || itemCount == 0{
		return ret
	}
	ret.Avg = ret.Total/itemCount
	for _,v := range m.interJobMetrics{
		ret.Stdev += (v.values[metricIndex] - ret.Avg) * (v.values[metricIndex] - ret.Avg)
	}
	ret.Stdev = math.Sqrt(ret.Stdev/itemCount)
	return ret
}

func (m MetricStatistic) Print()  {
	fmt.Printf("Metric: %s\n",m.Name)
	fmt.Printf("Total: %0.3f\n",m.Total)
	fmt.Printf("Avg: %0.3f\n",m.Avg)
	fmt.Printf("Min: %0.3f\n",m.Min)
	fmt.Printf("Max: %0.3f\n",m.Max)
	fmt.Printf("Stdev: %0.3f\n",m.Stdev)
	if m.TotalSeconds > 0{
		fmt.Printf("PerSec: %0.3f\n",m.Total/m.TotalSeconds)
	}
}
func (m JobStatistic) Print()  {
	fmt.Printf("GoStressLib version: %s\n",VERSION)
	fmt.Printf("Running in %s(%0.3fs), count: %d, concurrency: %d\n",NanosecondsToReadable(m.RunTime.Nanoseconds()),m.RunTime.Seconds(),m.Count,m.Concurrency)
	fmt.Printf("TPS: %0.3f/s\n",m.PerSeconds)
	for _,v := range m.MetricStatistic{
		fmt.Print("\n")
		v.Print()
	}
}

