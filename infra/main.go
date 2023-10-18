package main

import (
	"fmt"
	"log"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	baseResourceName := fmt.Sprintf("blog-%s", config.AWSAccountId)

	pulumi.Run(func(ctx *pulumi.Context) error {

		// S3 Bucket ------------------------------------
		resourceName := fmt.Sprintf("%s-s3_bucket", baseResourceName)
		bucketName := fmt.Sprintf("blog-%s", config.AWSAccountId)
		s3Bucket, err := s3.NewBucket(
			ctx,
			resourceName,
			&s3.BucketArgs{
				Bucket: pulumi.String(bucketName),
			},
		)
		if err != nil {
			return err
		}
		ctx.Export(resourceName, s3Bucket.ID())

		// S3 Bucket OwnerShip ---------------------------------
		resourceName = fmt.Sprintf("%s-s3_ownership", baseResourceName)
		bucketOwnership, err := s3.NewBucketOwnershipControls(
			ctx,
			resourceName,
			&s3.BucketOwnershipControlsArgs{
				Bucket: s3Bucket.ID(),
				Rule: &s3.BucketOwnershipControlsRuleArgs{
					ObjectOwnership: pulumi.String("BucketOwnerPreferred"),
				},
			})
		if err != nil {
			return err
		}
		ctx.Export(resourceName, bucketOwnership.ID())

		// S3 Bucket Policy CORS -------------------------
		// フロントでのPreFlightリクエストを許可する
		resourceName = fmt.Sprintf("%s-s3_cors", baseResourceName)
		s3CORS, err := s3.NewBucketCorsConfigurationV2(
			ctx,
			resourceName,
			&s3.BucketCorsConfigurationV2Args{
				Bucket: s3Bucket.ID(),
				CorsRules: s3.BucketCorsConfigurationV2CorsRuleArray{
					&s3.BucketCorsConfigurationV2CorsRuleArgs{
						AllowedHeaders: pulumi.StringArray{
							pulumi.String("*"),
						},
						AllowedMethods: pulumi.StringArray{
							pulumi.String("PUT"),
						},
						AllowedOrigins: pulumi.StringArray{
							pulumi.String("http://localhost:5173"),
						},
						MaxAgeSeconds: pulumi.Int(3000),
					},
				},
			})
		if err != nil {
			return err
		}
		ctx.Export(resourceName, s3CORS.ID())

		// S3 Bucket Public Access Block -----------------
		resourceName = fmt.Sprintf("%s-s3_public_access_block", baseResourceName)
		publicAccessBlock, err := s3.NewBucketPublicAccessBlock(
			ctx,
			resourceName,
			&s3.BucketPublicAccessBlockArgs{
				Bucket:                s3Bucket.ID(),
				BlockPublicAcls:       pulumi.Bool(false),
				BlockPublicPolicy:     pulumi.Bool(false),
				IgnorePublicAcls:      pulumi.Bool(false),
				RestrictPublicBuckets: pulumi.Bool(false),
			})
		if err != nil {
			return err
		}
		ctx.Export(resourceName, publicAccessBlock.ID())

		// S3 Bucket Policy ------------------------------
		// thumbnail配下のオブジェクトに対してGetObjectを許可する
		resourceName = fmt.Sprintf("%s-s3_bucket_policy", baseResourceName)
		bucketPolicy, err := s3.NewBucketPolicy(
			ctx,
			resourceName,
			&s3.BucketPolicyArgs{
				Bucket: s3Bucket.ID(), // refer to the bucket created earlier
				Policy: pulumi.Any(map[string]interface{}{
					"Version": "2012-10-17",
					"Statement": []map[string]interface{}{
						{
							"Effect":    "Allow",
							"Principal": "*",
							"Action": []interface{}{
								"s3:GetObject",
							},
							"Resource": []interface{}{
								pulumi.Sprintf("arn:aws:s3:::%s/thumbnail/*", s3Bucket.ID()),
							},
						},
					},
				}),
			})
		if err != nil {
			return err
		}
		ctx.Export(resourceName, bucketPolicy.ID())

		return nil

	})
}
