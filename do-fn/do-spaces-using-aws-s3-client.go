package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	accessKeyID := flag.String("SPACES_ACCESS_KEY", "", "Spaces Access Key")
	//accessKeyID := "EXHT3UDL7KCJG3GASLS5"
	secretAccessKey := flag.String("SPACES_ACCESS_KEY_SECRET", "", "Spaces Access Key Secret")
	//secretAccessKey := "mDaCxUVcSUPPqy2Bh/6o2Fl6ImdBlTweJbE4QJeWQCs"
	endpoint := flag.String("SPACES_ENDPOINT", "", "Spaces Endpoint e.g. ${REGION}.digitaloceanspaces.com")
	createBucket := flag.String("CREATE_SPACE_NAME", "", "name of Space you want to create")

	flag.Parse()
	//endpoint := "sgp1.digitaloceanspaces.com"
	region := "us-east-1"

	useSSL := true
	config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(*accessKeyID, *secretAccessKey, ""),
		Endpoint:         aws.String(*endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(!useSSL),
		S3ForcePathStyle: aws.Bool(true),
	}
	client := s3.New(session.Must(session.NewSession(config)))

	_, err := client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(*createBucket)})
	if err != nil {
		fmt.Println("Failed to create Space " + err.Error())
	} else {
		fmt.Println("Space " + *createBucket + " created successfully! ")
	}

	fmt.Println("Listing all Spaces.....\n")
	buckets, _ := client.ListBuckets(&s3.ListBucketsInput{})
	//fmt.Println("buckets " + buckets.GoString())
	for _, bucket := range buckets.Buckets {
		fmt.Println("Space details\n", *bucket)
		loc, _ := client.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: aws.String(*bucket.Name)})
		fmt.Println("Space region - " + loc.String())
	}

	fmt.Println("Deleting all Spaces.....\n")
	for _, bucket := range buckets.Buckets {
		//fmt.Println("Deleting Space " + *bucket.Name)
		_, err := client.DeleteBucket(&s3.DeleteBucketInput{Bucket: bucket.Name})
		if err != nil {
			fmt.Println("Could not delete Space - " + *bucket.Name + " due to " + err.Error())
		} else {
			fmt.Println("Deleted Space " + *bucket.Name)
		}
	}

}
