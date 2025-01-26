import { Construct } from "constructs";
import * as cdk from "aws-cdk-lib";

type Props = {
  stage: string;
  contentsBucketArn: string;
};

export class Lambda extends Construct {
  public readonly stage: string;
  public readonly function: cdk.aws_lambda.Function;
  public readonly functionUrl: cdk.aws_lambda.FunctionUrl;

  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);
    this.stage = props.stage;

    const cdkRoot = process.cwd();

    const functionRole = new cdk.aws_iam.Role(this, "FunctionRole", {
      assumedBy: new cdk.aws_iam.ServicePrincipal("lambda.amazonaws.com"),
      inlinePolicies: {
        s3: new cdk.aws_iam.PolicyDocument({
          statements: [
            new cdk.aws_iam.PolicyStatement({
              actions: ["s3:GetObject", "s3:PutObject"],
              resources: [props.contentsBucketArn],
            }),
          ],
        }),
        cloudwatch_logs: new cdk.aws_iam.PolicyDocument({
          statements: [
            new cdk.aws_iam.PolicyStatement({
              actions: [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents",
              ],
              resources: ["*"],
            }),
          ],
        }),
      },
    });

    const lambdaEnvironment = this.getLambdaEnvironment();

    this.function = new cdk.aws_lambda.DockerImageFunction(
      this,
      "DockerImageFunction",
      {
        code: cdk.aws_lambda.DockerImageCode.fromImageAsset(`${cdkRoot}/../`, {
          platform: cdk.aws_ecr_assets.Platform.LINUX_ARM64,
          buildArgs: {
            PORT: lambdaEnvironment["BLOG_APP_PORT"],
          },
        }),
        architecture: cdk.aws_lambda.Architecture.ARM_64,
        role: functionRole,
        environment: lambdaEnvironment,
        timeout: cdk.Duration.seconds(30),
      }
    );

    this.functionUrl = new cdk.aws_lambda.FunctionUrl(this, "FunctionUrl", {
      function: this.function,
      authType: cdk.aws_lambda.FunctionUrlAuthType.NONE,
    });
  }

  getLambdaEnvironment(): { [key: string]: string } {
    const lambdaEnvironmentKeys = [
      "BLOG_ENV",
      "BLOG_LOG_LEVEL",
      "BLOG_DB_HOST",
      "BLOG_DB_PORT",
      "BLOG_DB_USER",
      "BLOG_DB_PASS",
      "BLOG_DB_NAME",
      "BLOG_DB_TLS_ENABLED",
      "BLOG_KVS_HOST",
      "BLOG_KVS_PORT",
      "BLOG_KVS_USER",
      "BLOG_KVS_PASS",
      "BLOG_KVS_TLS_ENABLED",
      "BLOG_AWS_S3_BUCKET",
      "BLOG_AWS_S3_THUMBNAIL_DIRECTORY",
      "BLOG_AWS_S3_CONTENT_IMAGE_DIRECTORY",
      "ADMIN_NAME",
      "ADMIN_EMAIL",
      "ADMIN_PASSWORD",
      "JWT_SECRET",
      "SITE_DOMAIN",
      "CORS_WHITE_LIST",
      "CDN_DOMAIN",
      "GITHUB_PERSONAL_ACCESS_TOKEN",
    ];

    let env: { [key: string]: string } = {};
    lambdaEnvironmentKeys.forEach((key) => {
      const value = cdk.aws_ssm.StringParameter.valueForStringParameter(
        this,
        `/blog-api/${this.stage}/${key}`
      );
      env[key] = value;
    });
    // アプリケーション上のポートはLambdaWebAdapterと合わせるためにビルド時に固定
    env["BLOG_APP_PORT"] = "3000";
    return env;
  }
}
