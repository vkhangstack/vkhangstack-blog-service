package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	transfermanager "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type S3Config struct {
	Region          string
	Endpoint        string
	PublicURL       string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	UsePathStyle    bool
}

type PutInput struct {
	Key          string
	Body         io.Reader
	ContentType  string
	CacheControl string
	Metadata     map[string]string
	Bucket       string // Optional: override default bucket
}

type ObjectInfo struct {
	Key          string
	Size         int64
	ETag         string
	ContentType  string
	LastModified time.Time
	Metadata     map[string]string
}

type S3Adapter struct {
	client    *s3.Client
	presign   *s3.PresignClient
	uploader  *transfermanager.Client
	bucket    string
	endpoint  string
	publicURL string
	region    *string
	pathStyle *bool
}

func NewS3Adapter(ctx context.Context, cfg S3Config) (*S3Adapter, error) {
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket is required")
	}
	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("access key and secret key are required")
	}

	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			),
		),
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	// For S3-compatible endpoints (MinIO, LocalStack, etc.), BaseEndpoint pins the
	// URL and UsePathStyle forces /{bucket}/{key} path construction.
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = cfg.UsePathStyle
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	return &S3Adapter{
		client:    client,
		presign:   s3.NewPresignClient(client),
		uploader:  transfermanager.New(client),
		bucket:    cfg.Bucket,
		endpoint:  cfg.Endpoint,
		publicURL: cfg.PublicURL,
		region:    utils.StringPtr(cfg.Region),
		pathStyle: utils.BoolPtr(cfg.UsePathStyle),
	}, nil
}

func (a *S3Adapter) Put(ctx context.Context, in PutInput) error {
	if in.Key == "" {
		return fmt.Errorf("key is required")
	}
	if in.Body == nil {
		return fmt.Errorf("body is required")
	}
	if in.Bucket == "" {
		in.Bucket = a.bucket
	}

	input := &transfermanager.UploadObjectInput{
		Bucket:   aws.String(in.Bucket),
		Key:      aws.String(utils.CleanKey(in.Key)),
		Body:     in.Body,
		Metadata: in.Metadata,
	}
	if in.ContentType != "" {
		input.ContentType = aws.String(in.ContentType)
	}
	if in.CacheControl != "" {
		input.CacheControl = aws.String(in.CacheControl)
	}

	_, err := a.uploader.UploadObject(ctx, input)
	if err != nil {
		return fmt.Errorf("put object %q: %w", in.Key, err)
	}
	return nil
}

func (a *S3Adapter) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := a.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(utils.CleanKey(key)),
	})
	if err != nil {
		return nil, fmt.Errorf("get object %q: %w", key, err)
	}
	return out.Body, nil
}

func (a *S3Adapter) GetInfo(ctx context.Context, key string) (*ObjectInfo, error) {
	out, err := a.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(utils.CleanKey(key)),
	})
	if err != nil {
		return nil, fmt.Errorf("head object %q: %w", key, err)
	}

	info := &ObjectInfo{
		Key:          utils.CleanKey(key),
		Size:         aws.ToInt64(out.ContentLength),
		ETag:         aws.ToString(out.ETag),
		ContentType:  aws.ToString(out.ContentType),
		LastModified: aws.ToTime(out.LastModified),
		Metadata:     out.Metadata,
	}
	return info, nil
}

func (a *S3Adapter) Delete(ctx context.Context, key string, bucket string) error {
	if bucket == "" {
		bucket = a.bucket
	}
	logger.Log.WithFields(map[string]interface{}{
		"key":    key,
		"bucket": bucket,
	}).Debug("deleting s3 object")
	_, err := a.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket:                    aws.String(bucket),
		Key:                       aws.String(utils.CleanKey(key)),
		BypassGovernanceRetention: aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("delete object %q: %w", key, err)
	}
	return nil
}

func (a *S3Adapter) DeleteMany(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	objs := make([]types.ObjectIdentifier, 0, len(keys))
	for _, k := range keys {
		if strings.TrimSpace(k) == "" {
			continue
		}
		objs = append(objs, types.ObjectIdentifier{
			Key: aws.String(utils.CleanKey(k)),
		})
	}

	if len(objs) == 0 {
		return nil
	}

	_, err := a.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(a.bucket),
		Delete: &types.Delete{
			Objects: objs,
			Quiet:   aws.Bool(true),
		},
	})
	if err != nil {
		return fmt.Errorf("delete objects: %w", err)
	}
	return nil
}

func (a *S3Adapter) Exists(ctx context.Context, key string) (bool, error) {
	_, err := a.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(utils.CleanKey(key)),
	})
	if err == nil {
		return true, nil
	}

	var notFound *types.NotFound
	if ok := utils.ErrorAs(err, &notFound); ok {
		return false, nil
	}

	var apiErr interface{ ErrorCode() string }
	if ok := utils.ErrorAs(err, &apiErr); ok && (apiErr.ErrorCode() == "NotFound" || apiErr.ErrorCode() == "NoSuchKey") {
		return false, nil
	}

	return false, fmt.Errorf("head object %q: %w", key, err)
}

func (a *S3Adapter) Copy(ctx context.Context, srcKey, dstKey string) error {
	src := url.PathEscape(a.bucket + "/" + utils.CleanKey(srcKey))

	_, err := a.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(a.bucket),
		Key:        aws.String(utils.CleanKey(dstKey)),
		CopySource: aws.String(src),
	})
	if err != nil {
		return fmt.Errorf("copy object %q -> %q: %w", srcKey, dstKey, err)
	}
	return nil
}

func (a *S3Adapter) Rename(ctx context.Context, oldKey, newKey string) error {
	if utils.CleanKey(oldKey) == utils.CleanKey(newKey) {
		return nil
	}
	if err := a.Copy(ctx, oldKey, newKey); err != nil {
		return err
	}
	if err := a.Delete(ctx, oldKey, ""); err != nil {
		return fmt.Errorf("delete old key after rename %q: %w", oldKey, err)
	}
	return nil
}

func (a *S3Adapter) List(ctx context.Context, prefix string, limit int32) ([]ObjectInfo, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(a.bucket),
		Prefix: aws.String(utils.CleanKeyPrefix(prefix)),
	}
	if limit > 0 {
		input.MaxKeys = aws.Int32(limit)
	}

	out, err := a.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("list objects with prefix %q: %w", prefix, err)
	}

	items := make([]ObjectInfo, 0, len(out.Contents))
	for _, obj := range out.Contents {
		items = append(items, ObjectInfo{
			Key:          aws.ToString(obj.Key),
			Size:         aws.ToInt64(obj.Size),
			ETag:         aws.ToString(obj.ETag),
			LastModified: aws.ToTime(obj.LastModified),
		})
	}
	return items, nil
}

// PublicURL returns the direct public URL for an object.
// For path-style custom endpoints: {endpoint}/{bucket}/{key}
// For virtual-hosted custom endpoints: {scheme}://{bucket}.{host}/{key}
// For AWS S3 (no custom endpoint): https://{bucket}.s3.{region}.amazonaws.com/{key}
func (a *S3Adapter) PublicURL(key string, bucket string) string {
	cleanKey := utils.CleanKey(key)
	if a.publicURL != "" {
		ep := strings.TrimRight(a.publicURL, "/")
		if a.pathStyle != nil && *a.pathStyle {
			if bucket == "" {
				return ep + "/" + cleanKey
			}
			return ep + "/" + bucket + "/" + cleanKey
		}
		// virtual-hosted style: inject bucket as subdomain
		u, err := url.Parse(ep)
		if err == nil {
			return u.Scheme + "://" + bucket + "." + u.Host + "/" + cleanKey
		}
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", a.bucket, utils.StringVal(a.region), cleanKey)
}

func (a *S3Adapter) PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	req, err := a.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(utils.CleanKey(key)),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("presign get %q: %w", key, err)
	}
	return req.URL, nil
}

func (a *S3Adapter) PresignPut(ctx context.Context, key string, ttl time.Duration, contentType string) (string, error) {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(utils.CleanKey(key)),
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	req, err := a.presign.PresignPutObject(ctx, input, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("presign put %q: %w", key, err)
	}
	return req.URL, nil
}
