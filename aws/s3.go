package aws

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
)

type bucket struct {
	bucketName   string
	bucketRegion string
	Context      *gin.Context
	Object       *bucketObject
}

type bucketObject struct {
	Body       io.Reader
	ObjectKey  string
	ObjectKeyS []string
	Delimiter  *string
}

func New(bucketName, bucketRegion string, context *gin.Context, Object *bucketObject) bucket {
	return bucket{bucketName, bucketRegion, context, Object}
}

func (b *bucket) CreateBucket() (res *s3.CreateBucketOutput, err error) {
	s3Client := b.Context.MustGet("s3Client").(*s3.Client)
	if res, err = s3Client.CreateBucket(b.Context, &s3.CreateBucketInput{
		Bucket: &b.bucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraintEuCentral1,
		},
	}); err != nil {
		panic(err)

	} else {
		fmt.Println(res)
	}
	return
}

func (b *bucket) GetObject() (res *s3.GetObjectOutput, err error) {
	s3Client := b.Context.MustGet("s3Client").(*s3.Client)
	if res, err = s3Client.GetObject(
		b.Context, &s3.GetObjectInput{
			Bucket: &b.bucketName,
			Key:    &b.Object.ObjectKey,
		},
	); err != nil {
		panic(err)
	} else {
		fmt.Println(res)
	}
	return
}

func (b *bucket) PutObject() (res *s3.PutObjectOutput, err error) {
	s3Client := b.Context.MustGet("s3Client").(*s3.Client)
	if res, err = s3Client.PutObject(b.Context, &s3.PutObjectInput{
		Bucket: &b.bucketName,
		Key:    &b.Object.ObjectKey,
		Body:   b.Object.Body,
	}); err != nil {
		panic(err)
	} else {
		fmt.Println(res)
	}
	return
}

func (b *bucket) DeleteObject() (res *s3.DeleteObjectOutput, err error) {
	s3Client := b.Context.MustGet("s3Client").(*s3.Client)
	if res, err = s3Client.DeleteObject(b.Context, &s3.DeleteObjectInput{
		Bucket: &b.bucketName,
		Key:    &b.Object.ObjectKey,
	}); err != nil {
		panic(err)
	} else {
		fmt.Println(res)
	}
	return
}

func (b *bucket) ListObjectsV2() ([]interface{}, error) {
	s3Client := b.Context.MustGet("s3Client").(*s3.Client)
	var output []interface{}
	var nextToken *string
	var res *s3.ListObjectsV2Output
	var err error

	for {
		if nextToken != nil {
			res, err = s3Client.ListObjectsV2(b.Context, &s3.ListObjectsV2Input{
				Bucket:            &b.bucketName,
				Delimiter:         b.Object.Delimiter,
				MaxKeys:           50,
				ContinuationToken: nextToken,
			})
		} else {
			res, err = s3Client.ListObjectsV2(b.Context, &s3.ListObjectsV2Input{
				Bucket:    &b.bucketName,
				Delimiter: b.Object.Delimiter,
				MaxKeys:   50,
			})
		}

		if err != nil {
			return nil, err
		} else {
			output = append(output, res.Contents)
			nextToken = res.NextContinuationToken
		}

		if nextToken == nil {
			break
		}
	}

	return output, nil
}
