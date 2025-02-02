import { Construct } from "constructs";
import * as cdk from "aws-cdk-lib";

type Props = {
  lambdaFunction: cdk.aws_lambda.Function;
  acmCertificate: cdk.aws_certificatemanager.ICertificate;
  customDomainName: string;
};

export class APIGateway extends Construct {
  public readonly httpAPI: cdk.aws_apigatewayv2.HttpApi;
  public readonly customDomainName: cdk.aws_apigatewayv2.DomainName;

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

    this.customDomainName = new cdk.aws_apigatewayv2.DomainName(
      this,
      "CustomDomain",
      {
        certificate: props.acmCertificate,
        domainName: props.customDomainName,
      }
    );

    new cdk.aws_apigatewayv2.ApiMapping(this, "BasePathMapping", {
      domainName: this.customDomainName,
      api: this.httpAPI,
      stage: this.httpAPI.defaultStage,
    });
  }

  public getRoute53AliasRecordTarget(): cdk.aws_route53.IAliasRecordTarget {
    return new cdk.aws_route53_targets.ApiGatewayv2DomainProperties(
      this.customDomainName.regionalDomainName,
      this.customDomainName.regionalHostedZoneId
    );
  }
}
