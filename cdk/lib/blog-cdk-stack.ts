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

    const s3BucketBlog = new cdk.aws_s3.CfnBucket(this, "s3BucketBlog", {
      publicAccessBlockConfiguration: {
        restrictPublicBuckets: false,
        ignorePublicAcls: false,
        blockPublicPolicy: false,
        blockPublicAcls: false,
      },
      bucketName: bucketName,
      ownershipControls: {
        rules: [
          {
            objectOwnership: "BucketOwnerPreferred",
          },
        ],
      },
      bucketEncryption: {
        serverSideEncryptionConfiguration: [
          {
            bucketKeyEnabled: false,
            serverSideEncryptionByDefault: {
              sseAlgorithm: "AES256",
            },
          },
        ],
      },
    });
    s3BucketBlog.cfnOptions.deletionPolicy = cdk.CfnDeletionPolicy.RETAIN;

    const s3BucketPolicyBlog = new cdk.aws_s3.CfnBucketPolicy(
      this,
      "S3BucketPolicyBlog",
      {
        bucket: s3BucketBlog.ref,
        policyDocument: {
          Version: "2012-10-17",
          Statement: [
            {
              Resource: [
                `${s3BucketBlog.attrArn}/thumbnail/*`,
                `${s3BucketBlog.attrArn}/content/*`,
              ],
              Action: "s3:GetObject",
              Effect: "Allow",
              Principal: "*",
            },
          ],
        },
      }
    );
    s3BucketPolicyBlog.cfnOptions.deletionPolicy = cdk.CfnDeletionPolicy.RETAIN;

    const cloudFrontCachePolicy = new cdk.aws_cloudfront.CfnCachePolicy(
      this,
      "CloudFrontCachePolicy",
      {
        cachePolicyConfig: {
          comment:
            "Policy with caching enabled. Supports Gzip and Brotli compression.",
          minTtl: 1,
          maxTtl: 31536000,
          parametersInCacheKeyAndForwardedToOrigin: {
            queryStringsConfig: {
              queryStringBehavior: "none",
            },
            enableAcceptEncodingBrotli: true,
            headersConfig: {
              headerBehavior: "none",
            },
            cookiesConfig: {
              cookieBehavior: "none",
            },
            enableAcceptEncodingGzip: true,
          },
          defaultTtl: 86400,
          name: "Managed-CachingOptimized",
        },
      }
    );
    cloudFrontCachePolicy.cfnOptions.deletionPolicy =
      cdk.CfnDeletionPolicy.RETAIN;

    const cloudFrontOriginAccessControl =
      new cdk.aws_cloudfront.CfnOriginAccessControl(
        this,
        "CloudFrontOriginAccessControl",
        {
          originAccessControlConfig: {
            signingBehavior: "always",
            description: "",
            signingProtocol: "sigv4",
            originAccessControlOriginType: "s3",
            name: s3BucketBlog.attrDomainName,
          },
        }
      );
    cloudFrontOriginAccessControl.cfnOptions.deletionPolicy =
      cdk.CfnDeletionPolicy.RETAIN;

    const cloudFrontOriginRequestPolicy =
      new cdk.aws_cloudfront.CfnOriginRequestPolicy(
        this,
        "CloudFrontOriginRequestPolicy",
        {
          originRequestPolicyConfig: {
            queryStringsConfig: {
              queryStringBehavior: "none",
            },
            comment: "Policy for S3 origin with CORS",
            headersConfig: {
              headerBehavior: "whitelist",
              headers: [
                "origin",
                "access-control-request-headers",
                "access-control-request-method",
              ],
            },
            cookiesConfig: {
              cookieBehavior: "none",
            },
            name: "Managed-CORS-S3Origin",
          },
        }
      );
    cloudFrontOriginRequestPolicy.cfnOptions.deletionPolicy =
      cdk.CfnDeletionPolicy.RETAIN;

    const cloudFrontCloudFrontOriginAccessIdentity =
      new cdk.aws_cloudfront.CfnCloudFrontOriginAccessIdentity(
        this,
        "CloudFrontCloudFrontOriginAccessIdentity",
        {
          cloudFrontOriginAccessIdentityConfig: {
            comment: `${s3BucketBlog.attrDomainName} CloudFront Origin Access Identity`,
          },
        }
      );
    cloudFrontCloudFrontOriginAccessIdentity.cfnOptions.deletionPolicy =
      cdk.CfnDeletionPolicy.RETAIN;

    const cloudFrontDistribution = new cdk.aws_cloudfront.CfnDistribution(
      this,
      "CloudFrontDistribution",
      {
        distributionConfig: {
          logging: {
            includeCookies: false,
            bucket: "",
            prefix: "",
          },
          comment: "",
          defaultRootObject: "",
          origins: [
            {
              connectionTimeout: 10,
              originAccessControlId: "",
              connectionAttempts: 3,
              originCustomHeaders: [],
              domainName: `${bucketName}.s3.${stack.region}.amazonaws.com`,
              originShield: {
                enabled: false,
              },
              s3OriginConfig: {
                originAccessIdentity: `origin-access-identity/cloudfront/${cloudFrontCloudFrontOriginAccessIdentity.attrId}`,
              },
              originPath: "",
              id: `${bucketName}.s3.${stack.region}.amazonaws.com`,
            },
          ],
          viewerCertificate: {
            minimumProtocolVersion: "TLSv1",
            sslSupportMethod: "vip",
            cloudFrontDefaultCertificate: true,
          },
          priceClass: "PriceClass_All",
          defaultCacheBehavior: {
            compress: true,
            functionAssociations: [],
            lambdaFunctionAssociations: [],
            targetOriginId: `${bucketName}.s3.${stack.region}.amazonaws.com`,
            viewerProtocolPolicy: "allow-all",
            responseHeadersPolicyId: "5cc3b908-e619-4b99-88e5-2cf7f45965bd", // Managed-CORS-With-Preflight
            trustedSigners: [],
            fieldLevelEncryptionId: "",
            trustedKeyGroups: [],
            allowedMethods: ["HEAD", "GET"],
            cachedMethods: ["HEAD", "GET"],
            smoothStreaming: false,
            originRequestPolicyId: cloudFrontOriginRequestPolicy.ref,
            cachePolicyId: cloudFrontCachePolicy.ref,
          },
          staging: false,
          customErrorResponses: [],
          continuousDeploymentPolicyId: "",
          originGroups: {
            quantity: 0,
            items: [],
          },
          enabled: true,
          aliases: [],
          ipv6Enabled: true,
          webAclId: "",
          httpVersion: "http2",
          restrictions: {
            geoRestriction: {
              locations: [],
              restrictionType: "none",
            },
          },
          cacheBehaviors: [],
        },
      }
    );
    cloudFrontDistribution.cfnOptions.deletionPolicy =
      cdk.CfnDeletionPolicy.RETAIN;
  }
}
