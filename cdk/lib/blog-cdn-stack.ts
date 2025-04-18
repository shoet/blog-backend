import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";

export class BlogCDNStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: cdk.StackProps) {
    super(scope, id, props);

    const stack = cdk.Stack.of(this);

    const bucketName = cdk.aws_ssm.StringParameter.valueForStringParameter(
      this,
      `/blog-backend/infra/${props.stage}/CONTENTS_BUCKET_NAME`
    );

    const s3Bucket = new cdk.aws_s3.Bucket(this, "s3BucketBlog", {
      bucketName: bucketName,
      blockPublicAccess: cdk.aws_s3.BlockPublicAccess.BLOCK_ALL,
      objectOwnership: cdk.aws_s3.ObjectOwnership.BUCKET_OWNER_PREFERRED,
      encryption: cdk.aws_s3.BucketEncryption.S3_MANAGED,
      bucketKeyEnabled: false,
    });
    (s3Bucket.node.defaultChild as cdk.aws_s3.CfnBucket).overrideLogicalId(
      "s3BucketBlog"
    );

    const originAccessControl = new cdk.aws_cloudfront.S3OriginAccessControl(
      this,
      "CloudFrontOriginAccessControl",
      {
        originAccessControlName: s3Bucket.bucketDomainName,
        signing: cdk.aws_cloudfront.Signing.SIGV4_ALWAYS,
      }
    );
    (
      originAccessControl.node
        .defaultChild as cdk.aws_cloudfront.CfnOriginAccessControl
    ).overrideLogicalId("CloudFrontOriginAccessControl");

    const distribution = new cdk.aws_cloudfront.Distribution(
      this,
      "CloudFrontDistribution",
      {
        defaultBehavior: {
          origin:
            cdk.aws_cloudfront_origins.S3BucketOrigin.withOriginAccessControl(
              s3Bucket,
              {
                originAccessControl: originAccessControl,
              }
            ),
          viewerProtocolPolicy:
            cdk.aws_cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
          cachePolicy: cdk.aws_cloudfront.CachePolicy.CACHING_OPTIMIZED,
          compress: true,
          cachedMethods: cdk.aws_cloudfront.CachedMethods.CACHE_GET_HEAD,
          allowedMethods: cdk.aws_cloudfront.AllowedMethods.ALLOW_GET_HEAD,
        },
      }
    );
    (
      distribution.node.defaultChild as cdk.aws_cloudfront.CfnDistribution
    ).overrideLogicalId("CloudFrontDistribution");
  }
}
