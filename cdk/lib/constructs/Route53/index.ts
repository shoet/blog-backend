import { Construct } from "constructs";
import * as cdk from "aws-cdk-lib";

type Props = {
  hostedZoneId: string;
  hostedZoneName: string;
};

export class Route53 extends Construct {
  public readonly hostedZone: cdk.aws_route53.IHostedZone;

  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);

    this.hostedZone = cdk.aws_route53.HostedZone.fromHostedZoneAttributes(
      this,
      "HostedZone",
      {
        hostedZoneId: props.hostedZoneId,
        zoneName: props.hostedZoneName,
      }
    );
  }

  public createAliasRecord(
    recordName: Route53DomainNameWithDot,
    alias: cdk.aws_route53.IAliasRecordTarget
  ): cdk.aws_route53.ARecord {
    return new cdk.aws_route53.ARecord(this, "AliasRecord", {
      zone: this.hostedZone,
      recordName: recordName.domainName,
      target: cdk.aws_route53.RecordTarget.fromAlias(alias),
    });
  }
}

// ドメイン名の末尾にドットのつけ忘れを防ぐためのクラス
export class Route53DomainNameWithDot {
  public readonly domainName: string;

  constructor(domainName: string) {
    if (domainName.endsWith(".")) {
      this.domainName = domainName;
    } else {
      this.domainName = `${domainName}.`;
    }
  }
}
