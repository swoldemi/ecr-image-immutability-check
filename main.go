package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-xray-sdk-go/xray"
	log "github.com/sirupsen/logrus"
	"github.com/swoldemi/amazon-ecr-image-immutability-check/pkg/lib"
)

func main() {
	log.Info("Starting Lambda in live environment")
	sess, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
			SharedConfigState: session.SharedConfigEnable,
		},
	)
	if err != nil {
		log.Fatalf("Error creating session: %v\n", err)
		return
	}

	ecrSvc := ecr.New(sess)
	snsSvc := sns.New(sess)
	if err := xray.Configure(xray.Config{LogLevel: "trace"}); err != nil {
		log.Fatalf("Error configuring X-Ray: %v\n", err)
		return
	}

	xray.AWS(ecrSvc.Client)
	xray.AWS(snsSvc.Client)
	log.Info("Enabled request tracing on ECR and SNS API client")
	lambda.Start(lib.NewFunctionContainer(ecrSvc, snsSvc, lib.Production).GetHandler())
}
