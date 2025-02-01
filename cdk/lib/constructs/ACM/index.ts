import { Construct } from "constructs";
import * as cdk from "aws-cdk-lib";

type Props = {
  certificateArn: string;
};

export class ACM extends Construct {
  public readonly certificate: cdk.aws_certificatemanager.ICertificate;

  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);

    this.certificate =
      cdk.aws_certificatemanager.Certificate.fromCertificateArn(
        this,
        "ACMCertificate",
        props.certificateArn
      );
  }
}
