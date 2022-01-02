# GoStressLib
A golang library for stress testing.

## Install
```shell
go get github.com/tenhan/gostresslib
```

## Usage
```golang
package main

import (
	"github.com/tenhan/gostresslib/measurer"
	"io"
	"net/http"
)

func main()  {
	m := measurer.NewJobMeasurer()
	total := 10000
	concurrency := 2
	m.Run(total,concurrency,[]string{"response_size(byte)"}, func(num int, metric *measurer.JobMetric) error {
		resp,err := http.Get("http://127.0.0.1:8000/ping")
		if err != nil{
			return err
		}
		defer resp.Body.Close()
		bytes,err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		metric.SetMetricsValue(float64(len(bytes)))
		return nil
	}).Print()
}
```

output:
```bash
GoStressLib version: v1.0.0
Running in 19s882ms(19.882s), count: 10000, concurrency: 4
TPS: 502.962/s

Metric: Latency(s)
Total: 79.483
Avg: 0.008
Min: 0.005
Max: 0.047
Stdev: 0.002
PerSec: 3.998

Metric: Error(count)
Total: 0.000
Avg: 0.000
Min: 0.000
Max: 0.000
Stdev: 0.000
PerSec: 0.000

Metric: response_size(byte)
Total: 40000.000
Avg: 4.000
Min: 4.000
Max: 4.000
Stdev: 3.992
PerSec: 2011.848
```