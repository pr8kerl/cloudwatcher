package main

import (
	//	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"os"
	"strings"
	"time"
)

var (
	prevTick     time.Time
	pollperiod   int64         = 5
	pollduration time.Duration = 5               // minutes
	granularity  int64         = pollperiod * 60 // seconds
	metrics      []*cloudwatch.Metric
	namespace    string
	namespacestr string
	svc          *cloudwatch.CloudWatch
	gprefix      string = "test.aws.cloudwatch."
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
	namespace = "AWS/EC2"
	namespacestr = "ec2"
}

func main() {

	metrics, err := getAvailableMetrics()
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Printf("error requesting available metrics: %s\n", err.Error())
		os.Exit(1)
	}

	for {
		//		timetunnel := make(chan string, 1)
		//		select {
		//		case res := <-timetunnel:
		//			fmt.Printf("result: %s\n", res)
		now := <-time.After(pollduration * time.Minute)
		fmt.Printf("\n\ntimeout time to poll for stats\n")
		svc = cloudwatch.New(session.New())
		for _, metric := range metrics {
			getMetric(metric, prevTick, now)
		}
		//		for _, metric := range metrics {
		//			fmt.Printf("metric: %v\n", metric)
		//			getStatistic(metric, prevTick, now)
		//		}
		prevTick = now
		//		}

	}

}

func getAvailableMetrics() ([]*cloudwatch.Metric, error) {

	svc := cloudwatch.New(session.New())

	var params *cloudwatch.ListMetricsInput
	if len(namespace) > 0 {
		params = &cloudwatch.ListMetricsInput{
			//		Dimensions: []*cloudwatch.DimensionFilter{
			//			{ // Required
			//				Name:  aws.String(""), // Required
			//				Value: aws.String(""),
			//			},
			//		},
			//		MetricName: aws.String("CPUUtilization"),
			Namespace: aws.String(namespace),
		}
	} else {
		// empty - gets everything
		params = &cloudwatch.ListMetricsInput{}
	}

	resp, err := svc.ListMetrics(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return nil, err
	}

	metrics = resp.Metrics
	for resp.NextToken != nil {

		fmt.Println("\n\nmore metrics available")
		// get more metrics
		// append resp.Metrics to metrics
		if len(namespace) > 0 {
			params = &cloudwatch.ListMetricsInput{
				//		Dimensions: []*cloudwatch.DimensionFilter{
				//			{ // Required
				//				Name:  aws.String(""), // Required
				//				Value: aws.String(""),
				//			},
				//		},
				//		MetricName: aws.String("CPUUtilization"),
				Namespace: aws.String(namespace),
				NextToken: resp.NextToken,
			}
		} else {
			// empty - gets everything
			params = &cloudwatch.ListMetricsInput{
				NextToken: resp.NextToken,
			}
		}
		resp, err = svc.ListMetrics(params)
		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return nil, err
		}

		//fmt.Println(resp.Metrics)
		metrics = append(metrics, resp.Metrics...)

	}

	//	fmt.Printf("num metrics: %d\n", len(metrics))

	// Pretty-print the response data.
	//fmt.Println(metrics)
	return metrics, nil

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
		EndTime:    aws.Time(to),      // Required
		MetricName: metric.MetricName, // Required
		Namespace:  metric.Namespace,
		Period:     aws.Int64(granularity), // Required
		StartTime:  aws.Time(from),         // Required
		Statistics: []*string{ // Required
			aws.String("Maximum"),
			aws.String("Average"),
			//aws.String("Sum"),
			// More values...
		},
		Dimensions: []*cloudwatch.Dimension{
			{ // Required
				Name:  metric.Dimensions[0].Name,
				Value: metric.Dimensions[0].Value,
			},
			// More values...
		},
	}
	resp, err := svc.GetMetricStatistics(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	//fmt.Println(resp)

	if resp.Datapoints != nil {
		for _, datapoint := range resp.Datapoints {
			gpoint := gprefix + namespacestr + "." + *metric.Dimensions[0].Name + "." + *metric.Dimensions[0].Value + "." + *metric.MetricName
			gpoint = gpoint + "." + *datapoint.Unit + "."
			tstamp := fmt.Sprintf("%d", datapoint.Timestamp.Unix())
			fmt.Printf("%sMaximum %f %s\n", gpoint, *datapoint.Maximum, tstamp)
			fmt.Printf("%sAverage %f %s\n", gpoint, *datapoint.Average, tstamp)
			//fmt.Printf("%sSum %f %s\n", gpoint, *datapoint.Sum, tstamp)
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
