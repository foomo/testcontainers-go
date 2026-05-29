package rustfs_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/foomo/testcontainers-go/modules/rustfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	t.Cleanup(cancel)

	container, err := rustfs.Run(ctx, "rustfs/rustfs:latest")
	require.NoError(t, err)

	endpoint, err := container.Endpoint(ctx)
	require.NoError(t, err)

	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(endpoint),
		Region:       "us-east-1",
		UsePathStyle: true,
		Credentials: credentials.NewStaticCredentialsProvider(
			container.AccessKey,
			container.SecretKey,
			"",
		),
	})

	callCtx, callCancel := context.WithTimeout(ctx, 30*time.Second)
	t.Cleanup(callCancel)

	resp, err := client.ListBuckets(callCtx, &s3.ListBucketsInput{})
	require.NoError(t, err)

	assert.NotNil(t, resp)
}
