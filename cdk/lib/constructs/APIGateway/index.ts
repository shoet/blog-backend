import { Construct } from "constructs";
import * as cdk from "aws-cdk-lib";

type Props = {
  lambdaFunction: cdk.aws_lambda.Function;
};

export class APIGateway extends Construct {
  public readonly httpAPI: cdk.aws_apigatewayv2.HttpApi;
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);

    const stack = cdk.Stack.of(this);

    this.httpAPI = new cdk.aws_apigatewayv2.HttpApi(this, "HTTPApi", {
      apiName: `${stack.stackName}-HttpApi`,
    });

    const lambdaIntegration =
      new cdk.aws_apigatewayv2_integrations.HttpLambdaIntegration(
        "HttpLambdaIntegration",
        props.lambdaFunction,
        { timeout: cdk.Duration.seconds(29) }
      );

    this.httpAPI.addRoutes({
      path: "/{proxy+}",
      methods: [cdk.aws_apigatewayv2.HttpMethod.ANY],
      integration: lambdaIntegration,
    });
  }
}
