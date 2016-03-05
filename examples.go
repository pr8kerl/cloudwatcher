

svc := cloudwatch.New(session.New())

params := &cloudwatch.ListMetricsInput{
	Dimensions: []*cloudwatch.DimensionFilter{
		{ // Required
			Name:  aws.String("DimensionName"), // Required
			Value: aws.String("DimensionValue"),
		},
		// More values...
	},
	MetricName: aws.String("MetricName"),
	Namespace:  aws.String("Namespace"),
	NextToken:  aws.String("NextToken"),
}
resp, err := svc.ListMetrics(params)

if err != nil {
	// Print the error, cast err to awserr.Error to get the Code and
	// Message from an error.
	fmt.Println(err.Error())
	return
}

// Pretty-print the response data.
fmt.Println(resp)




/*
Gets statistics for the specified metric.

The maximum number of data points that can be queried is 50,850, whereas the maximum number of data points returned from a single GetMetricStatistics request is 1,440. If you make a request that generates more than 1,440 data points, Amazon CloudWatch returns an error. In such a case, you can alter the request by narrowing the specified time range or increasing the specified period. Alternatively, you can make multiple requests across adjacent time ranges. GetMetricStatistics does not return the data in chronological order.

Amazon CloudWatch aggregates data points based on the length of the period that you specify. For example, if you request statistics with a one-minute granularity, Amazon CloudWatch aggregates data points with time stamps that fall within the same one-minute period. In such a case, the data points queried can greatly outnumber the data points returned.

The following examples show various statistics allowed by the data point query maximum of 50,850 when you call GetMetricStatistics on Amazon EC2 instances with detailed (one-minute) monitoring enabled:

Statistics for up to 400 instances for a span of one hour Statistics for up to 35 instances over a span of 24 hours Statistics for up to 2 instances over a span of 2 weeks For information about the namespace, metric names, and dimensions that other Amazon Web Services products use to send metrics to CloudWatch, go to Amazon CloudWatch Metrics, Namespaces, and Dimensions Reference (http://docs.aws.amazon.com/AmazonCloudWatch/latest/DeveloperGuide/CW_Support_For_AWS.html) in the Amazon CloudWatch Developer Guide.

Examples:

Calling the GetMetricStatistics operation
*/

svc := cloudwatch.New(session.New())

params := &cloudwatch.GetMetricStatisticsInput{
	EndTime:    aws.Time(time.Now()),     // Required
	MetricName: aws.String("MetricName"), // Required
	Namespace:  aws.String("Namespace"),  // Required
	Period:     aws.Int64(1),             // Required
	StartTime:  aws.Time(time.Now()),     // Required
	Statistics: []*string{ // Required
		aws.String("Statistic"), // Required
		// More values...
	},
	Dimensions: []*cloudwatch.Dimension{
		{ // Required
			Name:  aws.String("DimensionName"),  // Required
			Value: aws.String("DimensionValue"), // Required
		},
		// More values...
	},
	Unit: aws.String("StandardUnit"),
}
resp, err := svc.GetMetricStatistics(params)

if err != nil {
	// Print the error, cast err to awserr.Error to get the Code and
	// Message from an error.
	fmt.Println(err.Error())
	return
}

// Pretty-print the response data.
fmt.Println(resp)
