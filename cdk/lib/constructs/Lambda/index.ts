import { Construct } from "constructs";
import * as cdk from "aws-cdk-lib";
import { getAppParameter } from "../SSM";
import * as imagedeploy from "cdk-docker-image-deployment";

type Props = {
  stage: string;
  contentsBucketArn: string;
  ecrRepository: cdk.aws_ecr.IRepository;
  commitHash?: string;
};

export class Lambda extends Construct {
  public readonly stage: string;
  public readonly function: cdk.aws_lambda.Function;

  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);
    this.stage = props.stage;

    const stack = cdk.Stack.of(this);

    const cdkRoot = process.cwd();
    const { platform, architecture } = getPlatform();

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

    const imageTag = props.commitHash || "latest";

    /* 公式のCDKではcdk.aws_lambda.DockerImageCode.fromImageAsset()で作成したイメージは
     * CDK用のECRにまとめられてしまうため、cdk-docker-image-deploymentを利用して
     * 作成したイメージをECRにデプロイする
     */
    const deployedImage = new imagedeploy.DockerImageDeployment(
      this,
      "CDKDockerImageDeployment",
      {
        source: imagedeploy.Source.directory(`${cdkRoot}/../`, {
          platform: platform,
          buildArgs: {
            PORT: lambdaEnvironment["BLOG_APP_PORT"],
          },
          target: "production",
        }),
        destination: imagedeploy.Destination.ecr(props.ecrRepository, {
          tag: imageTag,
        }),
      }
    );

    const cloudWatchLogGroup = new cdk.aws_logs.LogGroup(
      this,
      "CloudWatchLogGroup",
      {
        logGroupName: `/aws/lambda/${stack.stackName}-function`,
        removalPolicy: cdk.RemovalPolicy.DESTROY,
        retention: cdk.aws_logs.RetentionDays.TWO_WEEKS,
      }
    );

    this.function = new cdk.aws_lambda.DockerImageFunction(
      this,
      "DockerImageFunction",
      {
        code: cdk.aws_lambda.DockerImageCode.fromEcr(props.ecrRepository, {
          tag: imageTag,
        }),
        architecture: architecture,
        role: functionRole,
        environment: lambdaEnvironment,
        timeout: cdk.Duration.seconds(30),
        logGroup: cloudWatchLogGroup,
      }
    );

    this.function.node.addDependency(deployedImage);
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
      env[key] = getAppParameter(this, this.stage, key);
    });
    // アプリケーション上のポートはLambdaWebAdapterと合わせるためにビルド時に固定
    env["BLOG_APP_PORT"] = "3000";
    return env;
  }
}

const getPlatform = () => {
  const architecture = process.arch;
  switch (architecture) {
    case "arm64":
      return {
        architecture: cdk.aws_lambda.Architecture.ARM_64,
        platform: cdk.aws_ecr_assets.Platform.LINUX_ARM64,
      };
    case "x64":
      return {
        architecture: cdk.aws_lambda.Architecture.X86_64,
        platform: cdk.aws_ecr_assets.Platform.LINUX_AMD64,
      };
    default:
      return {
        architecture: cdk.aws_lambda.Architecture.X86_64,
        platform: cdk.aws_ecr_assets.Platform.LINUX_AMD64,
      };
  }
};
