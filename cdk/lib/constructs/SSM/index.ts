import { Construct } from "constructs";
import * as cdk from "aws-cdk-lib";

export function getAppParameter(scope: Construct, stage: string, key: string) {
  // デプロイ時に取得
  return cdk.aws_ssm.StringParameter.valueForStringParameter(
    scope,
    `/blog-api/${stage}/${key}`
  );
}

export function getInfraParameter(
  scope: Construct,
  stage: string,
  key: string
) {
  // デプロイ時に取得
  return cdk.aws_ssm.StringParameter.valueForStringParameter(
    scope,
    `/blog-api-infra/${stage}/${key}`
  );
}
