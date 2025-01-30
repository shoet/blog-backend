import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import { Lambda, APIGateway } from "./constructs";
import { getAppParameter, getInfraParameter } from "./constructs/SSM";

export class BlogAppStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: cdk.StackProps) {
    super(scope, id, props);

    const s3Bucket = cdk.aws_s3.Bucket.fromBucketName(
      this,
      "ContentsS3Bucket",
      getAppParameter(this, props.stage, "BLOG_AWS_S3_BUCKET")
    );

    const lambda = new Lambda(this, "Lambda", {
      stage: props.stage,
      contentsBucketArn: s3Bucket.bucketArn,
    });

    const apiGateway = new APIGateway(this, "APIGateway", {
      lambdaFunction: lambda.function,
      domainName: getInfraParameter(this, props.stage, "DOMAIN_NAME"),
      acmCertificateArn: getInfraParameter(
        this,
        props.stage,
        "ACM_CERTIFICATE_ARN"
      ),
      route53HostedZoneId: getInfraParameter(
        this,
        props.stage,
        "ROUTE53_HOSTED_ZONE_ID"
      ),
      route53HostedZoneName: getInfraParameter(
        this,
        props.stage,
        "ROUTE53_HOSTED_ZONE_NAME"
      ),
    });

    new cdk.CfnOutput(this, "APIGatewayUrl", {
      value: apiGateway.httpAPI.url || "",
    });

    new cdk.CfnOutput(this, "APIUrl", {
      value: `https://${apiGateway.apiDomain}/`,
    });
  }
}
