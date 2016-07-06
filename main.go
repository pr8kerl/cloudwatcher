package main

import (
	//	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	prevTick time.Time
	interval time.Duration
	metrics  map[string][]*cloudwatch.Metric
	svc      *cloudwatch.CloudWatch
	config   Config
	done     chan bool
	sigs     chan os.Signal
	mu       sync.Mutex
)

/*
		{ // Required
			Name:  aws.String("AutoScalingGroupName"), // Required
			Value: aws.String("vpc-mgmt-prod-SplunkUFAutoScalingGroup-1X0RAH90CM48E"),
		},
		// More values...
	},
	//		MetricName: aws.String("CPUUtilization"),
	Namespace: aws.String("AWS-EC2"),
	//		NextToken:  aws.String("NextToken"),
*/

func init() {
	prevTick = time.Now()
	metrics = make(map[string][]*cloudwatch.Metric)
	done = make(chan bool, 1)
	sigs = make(chan os.Signal, 1)
}

func main() {

	err := InitialiseConfig("config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing config: %s\n", err)
		os.Exit(1)
	}
	interval, err = time.ParseDuration(config.PollInterval)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing interval: %s\n", err)
		os.Exit(1)
	}

	awscfg := aws.NewConfig().WithRegion(config.Region).WithCredentials(credentials.NewSharedCredentials("", config.Profile))
	svc = cloudwatch.New(session.New(awscfg))

	for namespace, shortnamespace := range config.Namespaces {
		metrics[namespace], err = getAvailableMetrics(namespace)
		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Fprintf(os.Stderr, "error requesting available %s metrics: %s\n", shortnamespace, err.Error())
			os.Exit(1)
		}
	}
	// Pretty-print the response data.
	if config.Debug {
		fmt.Fprintf(os.Stderr, "%s\n", metrics)
		fmt.Fprintf(os.Stderr, "time to gather some stats\n")
	}

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go updateAvailableMetrics()

	go func() {
		for {
			//		timetunnel := make(chan string, 1)
			//		select {
			//		case res := <-timetunnel:
			//			fmt.Printf("result: %s\n", res)
			select {
			case <-done:
				if config.Debug {
					fmt.Fprintf(os.Stderr, "goroutine: time to exit\n")
				}
				return
			case <-time.After(interval):
				fmt.Fprintf(os.Stderr, "timeout: time to poll for stats\n")
				now := time.Now()
				mu.Lock()
				for namespace, _ := range config.Namespaces {
					for _, metric := range metrics[namespace] {
						getMetric(metric, prevTick, now)
					}
				}
				mu.Unlock()
				//		for _, metric := range metrics {
				//			fmt.Printf("metric: %v\n", metric)
				//			getStatistic(metric, prevTick, now)
				//		}
				prevTick = now
				//		}
			}
		}
	}()

	sig := <-sigs
	close(done)
	time.Sleep(1 * time.Second)
	fmt.Fprintf(os.Stderr, "signal %v - bye\n", sig)

}

func getAvailableMetrics(namespace string) ([]*cloudwatch.Metric, error) {

	var params *cloudwatch.ListMetricsInput
	params = &cloudwatch.ListMetricsInput{
		Namespace: aws.String(namespace),
	}

	resp, err := svc.ListMetrics(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return nil, err
	}

	tmetrics := resp.Metrics
	for resp.NextToken != nil {

		// get more metrics
		// append resp.Metrics to metrics
		params = &cloudwatch.ListMetricsInput{
			Namespace: aws.String(namespace),
			NextToken: resp.NextToken,
		}

		resp, err = svc.ListMetrics(params)
		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return nil, err
		}

		//fmt.Println(resp.Metrics)
		tmetrics = append(tmetrics, resp.Metrics...)

	}

	return tmetrics, nil

}

func getMetric(metric *cloudwatch.Metric, from time.Time, to time.Time) {

	if metric.Dimensions == nil {
		return
	}

	// ignore statuschecks
	if strings.HasPrefix(*metric.MetricName, "StatusCheck") {
		return
	}

	params := &cloudwatch.GetMetricStatisticsInput{
		EndTime:    aws.Time(to.Add(-interval)), // Required
		MetricName: metric.MetricName,           // Required
		Namespace:  metric.Namespace,
		Period:     aws.Int64(int64(interval.Seconds())), // Required
		StartTime:  aws.Time(from.Add(-interval)),        // Required
		Statistics: []*string{ // Required
			aws.String("Maximum"),
			aws.String("Average"),
			aws.String("Sum"),
			// More values...
		},
		Dimensions: metric.Dimensions,
	}
	resp, err := svc.GetMetricStatistics(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	// Pretty-print the response data.
	//fmt.Printf(os.Stderr, "%s\n", resp)

	if resp.Datapoints != nil {
		for _, datapoint := range resp.Datapoints {
			for _, dim := range metric.Dimensions {
				gpoint := config.Prefix + "." + config.Namespaces[*metric.Namespace] + "." + *dim.Name + "." + *dim.Value + "." + *metric.MetricName
				gpoint = gpoint + "." + *datapoint.Unit + "."
				tstamp := fmt.Sprintf("%d", datapoint.Timestamp.Unix())
				fmt.Printf("%sMaximum %f %s\n", gpoint, *datapoint.Maximum, tstamp)
				fmt.Printf("%sAverage %f %s\n", gpoint, *datapoint.Average, tstamp)
				if datapoint.Sum != nil {
					fmt.Printf("%sSum %f %s\n", gpoint, *datapoint.Sum, tstamp)
				}
			}
		}

		/*
		   {
		     Datapoints: [{
		         Average: 2.701533e+06,
		         Maximum: 3.642292e+06,
		         Timestamp: 2016-03-06 07:05:00 +0000 UTC,
		         Unit: "Bytes"
		       }],
		     Label: "NetworkOut"
		   }
		*/

	}

}

func updateAvailableMetrics() int {

	var err error
	tick := time.Tick(time.Duration(config.AvailableMetricsInterval) * time.Minute)
	for {
		select {
		case <-done:
			if config.Debug {
				fmt.Fprintf(os.Stderr, "goroutine: time to exit\n")
			}
			return 0
		case <-tick:
			if config.Debug {
				fmt.Fprintf(os.Stderr, "updating available metrics\n")
			}
			for namespace, shortnamespace := range config.Namespaces {
				mu.Lock()
				metrics[namespace], err = getAvailableMetrics(namespace)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error requesting available %s metrics: %s\n", shortnamespace, err.Error())
				}
				mu.Unlock()
			}
		}
	}

}
