package Connection

import (
	"context"
	"encoding/base64"
	"log"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/digitalocean/godo"
	Log "github.com/distribRuted/framework/library/log"
	Parameters "github.com/distribRuted/framework/library/parameters"
	"github.com/fatih/color"
)

func AWS_Servers(server_count int) {
	var awsAccessKeyID string = "YOUR_AWS_CREDENTIAL_HERE"
	var awsSecretAccessKey string = "YOUR_AWS_CREDENTIAL_HERE"
	var awsRegion string = "us-west-1"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := ec2.NewFromConfig(cfg)

	startupScript := "#!/bin/bash\n\nsudo apt update\nsudo apt-get install nmap -y\nwget https://yourwebsite.com/source_codes.zip\npython3 -m zipfile -e source_codes.zip /\ncd distribRuted\nchmod +x distribRuted\ntmux new-session -d -s distribRuted '" + Parameters.ShowCommand() + "'"

	runInstancesInput := &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-023e8dfe2208927a7"),
		InstanceType: types.InstanceTypeT3aMicro,
		MinCount:     aws.Int32(int32(server_count)),
		MaxCount:     aws.Int32(int32(server_count)),
		UserData:     aws.String(base64.StdEncoding.EncodeToString([]byte(startupScript))),
	}

	result, err := svc.RunInstances(context.TODO(), runInstancesInput)
	if err != nil {
		Log.PrintMsg(color.RedString("Error creating instance: %s", err)) // %v
	}

	for _, instance := range result.Instances {
		Log.PrintMsg(color.GreenString("Created instance: ") + *instance.InstanceId)
	}
}

func DO_Servers(server_count int) {
	token := "YOUR_DO_CREDENTIAL_HERE"
	client := godo.NewFromToken(token)
	ctx := context.TODO()

	var wg sync.WaitGroup
	wg.Add(server_count)

	for i := 0; i < server_count; i++ {
		createRequest := &godo.DropletCreateRequest{
			Name:   "client-" + strconv.Itoa(i+1),
			Region: "nyc3",
			Size:   "s-1vcpu-1gb",
			Image: godo.DropletCreateImage{
				Slug: "ubuntu-20-04-x64",
			},
			UserData: "#!/bin/bash\n\nsudo apt-get install nmap -y\nwget https://yourwebsite.com/source_codes.zip\npython3 -m zipfile -e source_codes.zip /\ncd distribRuted\nchmod +x distribRuted\ntmux new-session -d -s distribRuted '" + Parameters.ShowCommand() + "'",
		}
		go func() {
			defer wg.Done()
			droplet, _, err := client.Droplets.Create(ctx, createRequest)
			if err != nil {
				Log.PrintMsg(color.RedString("Error creating droplet: %s", err))
				return
			}
			Log.PrintMsg(color.GreenString("Created droplet: ") + droplet.Name)
		}()
	}
	wg.Wait()
}
