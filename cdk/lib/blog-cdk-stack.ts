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
  }
}
