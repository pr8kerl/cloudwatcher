package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func main() {

	svc := cloudwatch.New(session.New())

	params := &cloudwatch.ListMetricsInput{
		Dimensions: []*cloudwatch.DimensionFilter{
			{ // Required
				Name:  aws.String("AutoScalingGroupName"), // Required
				Value: aws.String("vpc-mgmt-prod-SplunkUFAutoScalingGroup-1X0RAH90CM48E"),
			},
			// More values...
		},
		//		MetricName: aws.String("MetricName"),
		//		Namespace:  aws.String("Namespace"),
		//		NextToken:  aws.String("NextToken"),
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

}
