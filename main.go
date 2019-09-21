package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/rs/zerolog/log"
)

func main() {
	region := os.Getenv("PLUGIN_REGION")
	if region == "" {
		log.Fatal().
			Err(errors.New("env: aws region is required")).
			Send()
	}
	key, id := os.Getenv("PLUGIN_AWS_SECRET_ACCESS_KEY"),
		os.Getenv("PLUGIN_AWS_ACCESS_KEY_ID")
	if key == "" || id == "" {
		log.Fatal().
			Err(errors.New("auth: aws Credentials Key or ID is not provieded")).
			Send()
	}
	svc := lambda.New(session.New(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(key, id, ""),
	}))

	input := &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(os.Getenv("PLUGIN_FUNCTION_NAME")),
		Publish:      aws.Bool(true),
		S3Bucket:     aws.String(os.Getenv("PLUGIN_S3_BUCKET")),
		S3Key:        aws.String(os.Getenv("PLUGIN_FILE_NAME")),
	}

	result, err := svc.UpdateFunctionCode(input)
	if err == nil {
		log.Info().Msgf("[Deploy Success]: %s", result.GoString())
		return
	}
	e, ok := err.(awserr.Error)
	if !ok {
		log.Fatal().Err(fmt.Errorf("deploy failed: %q", err))
	}

	awsErrField := map[string]interface{}{
		e.Code():    e.Message(),
		"originErr": e.OrigErr(),
	}

	log.Fatal().Fields(awsErrField).Send()

}
