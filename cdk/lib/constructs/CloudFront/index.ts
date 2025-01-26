import { Construct } from "constructs";
import * as cdk from "aws-cdk-lib";

type Props = {
  stage: string;
  lambdaFunction: cdk.aws_lambda.Function;
  lambdaFunctionUrl: cdk.aws_lambda.FunctionUrl;
};

export class CloudFront extends Construct {
  public readonly stage: string;
  public readonly distribution: cdk.aws_cloudfront.Distribution;

  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);

    this.stage = props.stage;

    const functionOriginAccessControl =
      new cdk.aws_cloudfront.FunctionUrlOriginAccessControl(
        this,
        "CloudFrontOAC",
        {
          signing: cdk.aws_cloudfront.Signing.SIGV4_ALWAYS,
          originAccessControlName: "AllowCloudFrontOAC",
        }
      );

    const functionUrlOrigin = new cdk.aws_cloudfront_origins.FunctionUrlOrigin(
      props.lambdaFunctionUrl,
      {
        originAccessControlId:
          functionOriginAccessControl.originAccessControlId,
      }
    );

    this.distribution = new cdk.aws_cloudfront.Distribution(
      this,
      "Distribution",
      {
        defaultBehavior: {
          origin: functionUrlOrigin,
          allowedMethods: cdk.aws_cloudfront.AllowedMethods.ALLOW_ALL,
        },
      }
    );

    props.lambdaFunction.addPermission("InvokeByCloudFront", {
      action: "lambda:InvokeFunctionUrl",
      principal: new cdk.aws_iam.ServicePrincipal("cloudfront.amazonaws.com"),
      sourceArn: `arn:aws:cloudfront::${cdk.Aws.ACCOUNT_ID}:distribution/${this.distribution.distributionId}`,
    });
  }
}
