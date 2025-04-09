package actions

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	appConf "github.com/adeo/lmru--devops--argocd-apps/go-configuration/config"
)

// проверяем наличие Bucket в S3
func checkBucket(client *s3.Client, bucketName string) error {
	result, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("Cannot connect to S3\n%w\n", bucketName, err)
	}

	for _, bucket := range result.Buckets {
		if aws.ToString(bucket.Name) == bucketName {
			return fmt.Errorf("Bucket already exist: %s", bucketName)
		}
	}
	return nil
}

func createBucket(client *s3.Client, bucketName string) error {
	input := &s3.CreateBucketInput{
		Bucket:                    aws.String(bucketName),
		ACL:                       types.BucketCannedACLPrivate,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{LocationConstraint: types.BucketLocationConstraintUsWest2},
	}
	_, err := client.CreateBucket(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("Bucket %s is not created\n%w\n", bucketName, err)
	}
	appConf.InfoLog.Printf("Bucket is created: %s", bucketName)
	appConf.InfoLog.Printf("Sleep 5 seconds")
	time.Sleep(5 * time.Second)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(appConf.BucketFolderName),
	})

	if err != nil {
		return fmt.Errorf("Folder %s in bucket %s is not created\n%w\n", appConf.BucketFolderName, bucketName, err)
	}
	appConf.InfoLog.Printf("Folder is created: %s", appConf.BucketFolderName)

	return nil
}

func MainS3(secretID, secretKey, bucketName string) error {
	// Создаем кастомный обработчик эндпоинтов, который для сервиса S3 и региона ru-central1 выдаст корректный URL
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   appConf.S3PartitionID,
			URL:           appConf.S3URL,
			SigningRegion: appConf.S3Region,
		}, nil
	})

	staticProvider := credentials.NewStaticCredentialsProvider(
		secretID,
		secretKey,
		"",
	)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(staticProvider),
	)
	if err != nil {
		return fmt.Errorf("Cannot create Config for S3\n%w\n", err)
	}
	// Создаем клиента для доступа к хранилищу S3
	client := s3.NewFromConfig(cfg)

	// Проверяем существует ли такой bucket
	err = checkBucket(client, bucketName)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// пробуем создать bucket в S3
	err = createBucket(client, bucketName)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
