package utils

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

var AccessKeyID string
var SecretAccessKey string
var Region string

var sess *session.Session

const RenamePartSize int64 = 1024 * 1024 * 1024 * 1

type RenamerInput struct {
	Bucket    string
	SourceKey string
	DestKey   string
}

func deleteAfterCopy(input *RenamerInput) {
	if sess == nil {
		sess = ConnectAws()
	}
	svc := s3.New(sess)
	dparams := &s3.DeleteObjectInput{
		Bucket: aws.String(input.Bucket),
		Key:    aws.String(input.SourceKey),
	}
	_, err := svc.DeleteObject(dparams)
	if err != nil {
		msg := "* S3 DID NOT DEL A FILE!"
		fmt.Println(err, msg, input.Bucket, input.SourceKey)
	}
}

func RenameS3(input *RenamerInput) error {
	if sess == nil {
		sess = ConnectAws()
	}
	svc := s3.New(sess)

	copyInput := &s3.CopyObjectInput{
		Bucket:     aws.String(input.Bucket),
		CopySource: aws.String(fmt.Sprintf("/%s/%s", input.Bucket, input.SourceKey)),
		Key:        aws.String(input.DestKey),
	}

	_, err := svc.CopyObject(copyInput)
	if err != nil {
		return err
	}
	deleteAfterCopy(input)
	return nil

	// params := &s3.CreateMultipartUploadInput{
	// 	Bucket: aws.String(input.DestBucket),
	// 	Key:    aws.String(input.DestKey),
	// }
	// resp, err := svc.CreateMultipartUpload(params)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	// uid := string(*resp.UploadId)
	// s := input.Size
	// part := int64(1)
	// parts := int(math.Ceil(float64(input.Size) / float64(RenamePartSize)))
	// var theParts []*s3.CompletedPart = make([]*s3.CompletedPart, parts)
	// for {
	// 	offset := RenamePartSize * (part - 1)
	// 	endbyte := offset + RenamePartSize - 1
	// 	if endbyte >= input.Size {
	// 		endbyte = offset + s - 1
	// 	}
	// 	source, err := svc.UploadPartCopy(&s3.UploadPartCopyInput{
	// 		Bucket:          aws.String(input.DestBucket),
	// 		Key:             aws.String(input.DestKey),
	// 		CopySource:      aws.String(input.SourceBucket + "/" + input.SourceKey),
	// 		CopySourceRange: aws.String(fmt.Sprintf("bytes=%d-%d", offset, endbyte)),
	// 		PartNumber:      aws.Int64(part),
	// 		UploadId:        &uid,
	// 	})

	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return err
	// 	}

	// 	etag := string(*source.CopyPartResult.ETag)
	// 	etagl := len(etag)
	// 	etag = etag[1 : etagl-1]

	// 	fooi := int64(part)
	// 	partn := &fooi
	// 	theParts[part-1] = &s3.CompletedPart{ETag: &etag, PartNumber: partn}

	// 	part++
	// 	s -= RenamePartSize
	// 	if s <= 0 {
	// 		break
	// 	}
	// }

	// cparams := &s3.CompleteMultipartUploadInput{
	// 	Bucket:   aws.String(input.DestBucket),
	// 	Key:      aws.String(input.DestKey),
	// 	UploadId: &uid,
	// 	MultipartUpload: &s3.CompletedMultipartUpload{
	// 		Parts: theParts,
	// 	},
	// }
	// _, err = svc.CompleteMultipartUpload(cparams)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	// deleteAfterCopy(input)
	// return nil
}

func ConnectAws() *session.Session {
	AccessKeyID = viper.GetString("BUCKET_ACCESS_KEY")
	SecretAccessKey = viper.GetString("BUCKET_ACCESS_SECRET")
	Region = viper.GetString("AWS_REGION")
	var err error = nil // no
	sess, err = session.NewSession(
		&aws.Config{
			Region: aws.String(Region),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}

func GetObjectInS3(key string, expiryTime time.Duration) (string, error) {
	if sess == nil {
		sess = ConnectAws()
	}
	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("AWS_BUCKET")),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(expiryTime)
	if err != nil {
		return "", err
	}
	return urlStr, nil
}

func ListObjectInS3(key string) []string {
	if sess == nil {
		sess = ConnectAws()
	}
	svc := s3.New(sess)
	params := &s3.ListObjectsInput{
		Bucket: aws.String(viper.GetString("BUCKET_NAME")),
		Prefix: aws.String(key),
	}
	keys := []string{}
	resp, _ := svc.ListObjects(params)
	for _, key := range resp.Contents {
		keys = append(keys, *key.Key)
	}

	return keys
}
