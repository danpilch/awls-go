package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var cli struct {
	Search string `required short:"s" help:"EC2 Instance Name search term"`
}

func main() {
	// Parse cli args
	kong.Parse(&cli)

	// Load session from shared config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create new EC2 client
	svc := ec2.New(sess)

	// Define search params
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(cli.Search + "*"),
				},
			},
		},
	}

	result, err := svc.DescribeInstances(params)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	if len(result.Reservations) < 1 {
		fmt.Println("no matching instances found")
		return
	}

	tableData := [][]string{}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			var nameTag string
			for _, t := range instance.Tags {
				if *t.Key == "Name" {
					nameTag = *t.Value
					break
				}
			}

			tableData = append(tableData, []string{
				string(nameTag),
				string(*instance.PrivateIpAddress),
				string(*instance.State.Name),
				string(*instance.Placement.AvailabilityZone),
				string(*instance.InstanceId),
				string(*instance.InstanceType),
				string(instance.LaunchTime.Format("2006-01-02 15:04:05")),
			})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "PrivateIp", "State", "AZ", "InstanceId", "InstanceType", "LaunchTime"})
	table.SetAutoFormatHeaders(true)
	table.AppendBulk(tableData)
	table.Render()
}
