import { Construct } from "constructs";
import * as cdk from "aws-cdk-lib";

type Props = {
  repositoryName: string;
};

export class ECR extends Construct {
  public readonly repository: cdk.aws_ecr.Repository;
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);

    this.repository = new cdk.aws_ecr.Repository(this, "Repository", {
      repositoryName: props.repositoryName,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      lifecycleRules: [
        {
          maxImageCount: 3,
        },
      ],
    });
  }
}
