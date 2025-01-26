import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import { Lambda } from "./constructs";

export class BlogAppStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: cdk.StackProps) {
    super(scope, id, props);

    const BucketName = cdk.aws_ssm.StringParameter.valueForStringParameter(
      this,
      `/blog-api/${props.stage}/BLOG_AWS_S3_BUCKET`
    );

    const s3Bucket = cdk.aws_s3.Bucket.fromBucketName(
      this,
      "ContentsS3Bucket",
      BucketName
    );

    s3Bucket.bucketArn;

    const lambda = new Lambda(this, "Lambda", {
      stage: props.stage,
      contentsBucketArn: s3Bucket.bucketArn,
    });

    new cdk.CfnOutput(this, "FunctionUrl", {
      value: lambda.functionUrl.url,
    });
  }
}
