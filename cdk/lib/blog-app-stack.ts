import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import { Lambda, APIGateway, ACM, Route53 } from "./constructs";
import { getAppParameter, getInfraParameter } from "./constructs/SSM";
import { Route53DomainNameWithDot } from "./constructs/Route53";

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

    const acm = new ACM(this, "ACM", {
      certificateArn: getInfraParameter(
        this,
        props.stage,
        "ACM_CERTIFICATE_ARN"
      ),
    });

    const domainName = getInfraParameter(this, props.stage, "DOMAIN_NAME");

    const apiGateway = new APIGateway(this, "APIGateway", {
      lambdaFunction: lambda.function,
      customDomainName: domainName,
      acmCertificate: acm.certificate,
    });

    const route53 = new Route53(this, "Route53", {
      hostedZoneId: getInfraParameter(
        this,
        props.stage,
        "ROUTE53_HOSTED_ZONE_ID"
      ),
      hostedZoneName: getInfraParameter(
        this,
        props.stage,
        "ROUTE53_HOSTED_ZONE_NAME"
      ),
    });

    route53.createAliasRecord(
      new Route53DomainNameWithDot(domainName),
      apiGateway.getRoute53AliasRecordTarget()
    );

    new cdk.CfnOutput(this, `APIGatewayUrl`, {
      value: apiGateway.httpAPI.url || "",
    });

    new cdk.CfnOutput(this, `APIUrl`, {
      exportName: `BlogAppStackApiUrl-${props.stage}`,
      value: `https://${domainName}`,
    });
  }
}
