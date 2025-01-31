import { Construct } from "constructs";
import * as cdk from "aws-cdk-lib";

type Props = {
  lambdaFunction: cdk.aws_lambda.Function;
  acmCertificateArn: string;
  domainName: string;
  route53HostedZoneId: string;
  route53HostedZoneName: string;
};

export class APIGateway extends Construct {
  public readonly httpAPI: cdk.aws_apigatewayv2.HttpApi;
  public readonly domainName: cdk.aws_apigatewayv2.DomainName;
  public readonly apiDomain: string;

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

    const acmCertificate =
      cdk.aws_certificatemanager.Certificate.fromCertificateArn(
        this,
        "ACMCertificate",
        props.acmCertificateArn
      );

    this.domainName = new cdk.aws_apigatewayv2.DomainName(
      this,
      "CustomDomain",
      {
        certificate: acmCertificate,
        domainName: props.domainName,
      }
    );

    new cdk.aws_apigatewayv2.ApiMapping(this, "BasePathMapping", {
      domainName: this.domainName,
      api: this.httpAPI,
      stage: this.httpAPI.defaultStage,
    });

    const hostedZone = cdk.aws_route53.HostedZone.fromHostedZoneAttributes(
      this,
      "HostedZone",
      {
        hostedZoneId: props.route53HostedZoneId,
        zoneName: props.route53HostedZoneName,
      }
    );

    new cdk.aws_route53.ARecord(this, "AliasRecord", {
      zone: hostedZone,
      recordName: `${props.domainName}.`,
      target: cdk.aws_route53.RecordTarget.fromAlias(
        new cdk.aws_route53_targets.ApiGatewayv2DomainProperties(
          this.domainName.regionalDomainName,
          this.domainName.regionalHostedZoneId
        )
      ),
    });

    this.apiDomain = props.domainName;
  }
}
