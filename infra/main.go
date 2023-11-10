package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var projectTag = "blog"

func CreateVPC(ctx *pulumi.Context, cidr string, resourceName string) (*ec2.Vpc, error) {
	return ec2.NewVpc(ctx, resourceName, &ec2.VpcArgs{
		CidrBlock:          pulumi.String(cidr),
		EnableDnsSupport:   pulumi.Bool(true),
		EnableDnsHostnames: pulumi.Bool(true),
		Tags:               createNameTag(resourceName),
	})
}

func CreateSubnet(
	ctx *pulumi.Context,
	vpc *ec2.Vpc,
	cidr string,
	// availabilityZone string,
	resourceName string,
) (*ec2.Subnet, error) {
	return ec2.NewSubnet(ctx, resourceName, &ec2.SubnetArgs{
		VpcId:     vpc.ID(),
		CidrBlock: pulumi.String(cidr),
		// AvailabilityZone: pulumi.String(availabilityZone),
		Tags: createNameTag(resourceName),
	})
}

func CreateIGW(
	ctx *pulumi.Context, vpc *ec2.Vpc, resourceName string,
) (*ec2.InternetGateway, error) {

	return ec2.NewInternetGateway(ctx, resourceName, &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags:  createNameTag(resourceName),
	})
}

func CreatePublicRouteTable(
	ctx *pulumi.Context, vpc *ec2.Vpc, igw *ec2.InternetGateway, resourceName string,
) (*ec2.RouteTable, error) {
	return ec2.NewRouteTable(
		ctx, resourceName, &ec2.RouteTableArgs{
			VpcId: vpc.ID(),
			Routes: ec2.RouteTableRouteArray{
				&ec2.RouteTableRouteArgs{
					CidrBlock: pulumi.String("0.0.0.0/0"),
					GatewayId: igw.ID(),
				},
			},
			Tags: createNameTag(resourceName),
		},
		pulumi.DependsOn([]pulumi.Resource{vpc, igw}),
	)
}

func CreateRouteTableAssociation(
	ctx *pulumi.Context, routeTable *ec2.RouteTable, subnet *ec2.Subnet, resourceName string,
) (*ec2.RouteTableAssociation, error) {
	return ec2.NewRouteTableAssociation(
		ctx,
		resourceName,
		&ec2.RouteTableAssociationArgs{
			RouteTableId: routeTable.ID(),
			SubnetId:     subnet.ID(),
		},
		pulumi.DependsOn([]pulumi.Resource{routeTable, subnet}),
	)
}

func CreateSecurityGroupForMaintenanceEC2(
	ctx *pulumi.Context, vpc *ec2.Vpc, resourceName string,
) (*ec2.SecurityGroup, error) {
	return ec2.NewSecurityGroup(
		ctx,
		resourceName,
		&ec2.SecurityGroupArgs{
			VpcId: vpc.ID(),
			Ingress: ec2.SecurityGroupIngressArray{
				&ec2.SecurityGroupIngressArgs{
					Description: pulumi.String("for ssh"),
					Protocol:    pulumi.String("tcp"),
					FromPort:    pulumi.Int(22),
					ToPort:      pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				&ec2.SecurityGroupEgressArgs{
					Description: pulumi.String("All outbound traffic"),
					Protocol:    pulumi.String("-1"),
					FromPort:    pulumi.Int(0),
					ToPort:      pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
				},
			},
			Tags: createNameTag(resourceName),
		})
}

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		// VPC //////////////////////////////////////////////////////////////////////////
		resourceName := fmt.Sprintf("%s-vpc", projectTag)
		vpc, err := CreateVPC(ctx, "10.1.0.0/16", resourceName)
		if err != nil {
			return fmt.Errorf("failed create vpc: %v", err)
		}
		ctx.Export(resourceName, vpc.ID())

		// Subnet /////////////////////////////////////////////////////////////////////////
		resourceName = fmt.Sprintf("%s-subnet-app-container", projectTag)
		subnetAppContainer, err := CreateSubnet(ctx, vpc, "10.1.0.0/24", resourceName)
		if err != nil {
			return fmt.Errorf("failed create subnet for App Container: %v", err)
		}
		ctx.Export(resourceName, subnetAppContainer.ID())

		resourceName = fmt.Sprintf("%s-subnet-maintenance", projectTag)
		subnetMaintenance, err := CreateSubnet(ctx, vpc, "10.1.1.0/24", resourceName)
		if err != nil {
			return fmt.Errorf("failed create subnet for Maintenance: %v", err)
		}
		ctx.Export(resourceName, subnetMaintenance.ID())

		// InternetGateway //////////////////////////////////////////////////////////
		resourceName = fmt.Sprintf("%s-igw", projectTag)
		igw, err := CreateIGW(ctx, vpc, resourceName)
		if err != nil {
			return fmt.Errorf("failed create igw: %v", err)
		}
		ctx.Export(resourceName, igw.ID())

		// ルートテーブル /////////////////////////////////////////////////////
		resourceName = fmt.Sprintf("%s-route-table-public", projectTag)
		publicRouteTable, err := CreatePublicRouteTable(ctx, vpc, igw, resourceName)
		if err != nil {
			return fmt.Errorf("failed create public route table: %v", err)
		}
		ctx.Export(resourceName, publicRouteTable.ID())

		// ルートテーブル 関連付け///////////////////////////////////////////////////////////
		resourceName = fmt.Sprintf("%s-route-table-association-app-container", projectTag)
		routeTableAssociationAppContainer, err := CreateRouteTableAssociation(
			ctx, publicRouteTable, subnetAppContainer, resourceName)
		if err != nil {
			return fmt.Errorf("failed create public route association for AppContainer: %v", err)
		}
		ctx.Export(resourceName, routeTableAssociationAppContainer.ID())

		resourceName = fmt.Sprintf("%s-route-table-association-maintenance", projectTag)
		routeTableAssociationMaintenance, err := CreateRouteTableAssociation(
			ctx, publicRouteTable, subnetMaintenance, resourceName)
		if err != nil {
			return fmt.Errorf("failed create public route association for maintenance: %v", err)
		}
		ctx.Export(resourceName, routeTableAssociationMaintenance.ID())

		// IAM //////////////////////////////////////////////////////////////
		// ロール メンテナンスEC2向け
		resourceName = fmt.Sprintf("%s-iam-role-for-maintenance-ec2", projectTag)
		iamMaintenanceEC2, err := iam.NewRole(
			ctx,
			resourceName,
			&iam.RoleArgs{
				AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Effect": "Allow",
					"Principal": {
						"Service": "ec2.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}]
			}`),
				Tags: createNameTag(resourceName),
			})
		if err != nil {
			return fmt.Errorf("failed create iam role for maintenance ec2: %v", err)
		}
		ctx.Export(resourceName, iamMaintenanceEC2.ID())

		// ポリシー メンテナンスEC2向け
		resourceName = fmt.Sprintf("%s-iam-policy-for-maintenance-ec2", projectTag)
		iamMaintenanceEC2Policy, err := iam.NewRolePolicy(
			ctx,
			resourceName,
			&iam.RolePolicyArgs{
				Role: iamMaintenanceEC2.Name,
				Policy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "ECRPermissions",
						"Effect": "Allow",
						"Action": [
							"ecr:GetAuthorizationToken",
							"ecr:BatchCheckLayerAvailability",
							"ecr:GetDownloadUrlForLayer",
							"ecr:GetRepositoryPolicy",
							"ecr:DescribeRepositories",
							"ecr:ListImages",
							"ecr:DescribeImages",
							"ecr:BatchGetImage",
							"ecr:InitiateLayerUpload",
							"ecr:UploadLayerPart",
							"ecr:CompleteLayerUpload",
							"ecr:PutImage"
						],
						"Resource": "*"
					}
				]
			}`),
			})
		if err != nil {
			return fmt.Errorf("failed create iam policy for maintenance ec2: %v", err)
		}
		ctx.Export(resourceName, iamMaintenanceEC2Policy.ID())

		// セキュリティグループ securitygroup ///////////////////////////////////////////////
		resourceName = fmt.Sprintf("%s-sg-public-maintenance", projectTag)
		securityGroupPublicMaintenance, err := CreateSecurityGroupForMaintenanceEC2(
			ctx, vpc, resourceName)
		if err != nil {
			return fmt.Errorf("failed create security group for public maintenance: %v", err)
		}
		ctx.Export(resourceName, securityGroupPublicMaintenance.ID())

		// EC2 ////////////////////////////////////////////////////////////////////
		// インスタンスプロファイル
		resourceName = fmt.Sprintf("%s-instance-profile-for-maintenance-ec2", projectTag)
		instanceProfileMaintenanceEC2, err := iam.NewInstanceProfile(
			ctx,
			resourceName,
			&iam.InstanceProfileArgs{
				Role: iamMaintenanceEC2.Name,
			})
		if err != nil {
			return fmt.Errorf("failed create iam instance profile for maintenance ec2: %v", err)
		}
		ctx.Export(resourceName, instanceProfileMaintenanceEC2.ID())

		// インスタンス
		userdataScript, err := loadFileToString("./maintenance_ec2_userdata.sh")
		if err != nil {
			return fmt.Errorf("failed load file: %v", err)
		}
		resourceName = fmt.Sprintf("%s-ec2-maintenance", projectTag)
		ec2MaintenanceInstance, err := ec2.NewInstance(
			ctx,
			resourceName,
			&ec2.InstanceArgs{
				InstanceType:             pulumi.String("t2.micro"),
				Ami:                      pulumi.String("ami-08a706ba5ea257141"),
				SubnetId:                 subnetMaintenance.ID(),
				KeyName:                  pulumi.String(config.BastionSSHKeyName),
				AssociatePublicIpAddress: pulumi.Bool(true),
				SecurityGroups: pulumi.StringArray{
					securityGroupPublicMaintenance.ID(),
				},
				IamInstanceProfile: instanceProfileMaintenanceEC2.Name,
				UserData:           pulumi.String(userdataScript),
				Tags:               createNameTag(resourceName),
			},
			pulumi.IgnoreChanges([]string{"securityGroups"}),
		)
		if err != nil {
			return fmt.Errorf("failed create new maintenance ec2 instance: %v", err)
		}
		ctx.Export(resourceName, ec2MaintenanceInstance.ID())

		// ElasticIPアドレス ///////////////////////////////////////////////////////////////////
		// メンテナンスEC2向け
		resourceName = fmt.Sprintf("%s-eip-for-maintenance", projectTag)
		eipForEc2MaintenanceInstance, err := ec2.NewEip(
			ctx,
			resourceName,
			&ec2.EipArgs{
				Domain:   pulumi.String("vpc"),
				Instance: ec2MaintenanceInstance.ID(),
			},
			pulumi.IgnoreChanges([]string{"instance"}),
		)
		if err != nil {
			return fmt.Errorf("failed create eip for maintenance ec2 instance: %v", err)
		}
		ctx.Export(resourceName, eipForEc2MaintenanceInstance.ID())

		// ElasticIP 紐づけ ////////////////////////////////////////////////////////////////////
		// メンテナンスEC2向け
		resourceName = fmt.Sprintf("%s-eip-associate-for-maintenance", projectTag)
		eipAssociate, err := ec2.NewEipAssociation(
			ctx,
			resourceName,
			&ec2.EipAssociationArgs{
				InstanceId:   ec2MaintenanceInstance.ID(),
				AllocationId: eipForEc2MaintenanceInstance.ID(),
			},
			pulumi.DependsOn([]pulumi.Resource{ec2MaintenanceInstance, eipForEc2MaintenanceInstance}),
		)
		ctx.Export(resourceName, eipAssociate.ID())

		// S3 ////////////////////////////////////////////////////////////////////////
		// bucket
		resourceName = fmt.Sprintf("%s-s3-bucket", projectTag)
		bucketName := fmt.Sprintf("blog-%s", config.AWSAccountId)
		s3Bucket, err := s3.NewBucket(
			ctx,
			resourceName,
			&s3.BucketArgs{
				Bucket: pulumi.String(bucketName),
			},
		)
		if err != nil {
			return err
		}
		ctx.Export(resourceName, s3Bucket.ID())

		// Bucket OwnerShip
		resourceName = fmt.Sprintf("%s-s3-ownership", projectTag)
		bucketOwnership, err := s3.NewBucketOwnershipControls(
			ctx,
			resourceName,
			&s3.BucketOwnershipControlsArgs{
				Bucket: s3Bucket.ID(),
				Rule: &s3.BucketOwnershipControlsRuleArgs{
					ObjectOwnership: pulumi.String("BucketOwnerPreferred"),
				},
			})
		if err != nil {
			return err
		}
		ctx.Export(resourceName, bucketOwnership.ID())

		// CORS Configuration
		// フロントでのPreFlightリクエストを許可する
		resourceName = fmt.Sprintf("%s-s3-cors", projectTag)
		s3CORS, err := s3.NewBucketCorsConfigurationV2(
			ctx,
			resourceName,
			&s3.BucketCorsConfigurationV2Args{
				Bucket: s3Bucket.ID(),
				CorsRules: s3.BucketCorsConfigurationV2CorsRuleArray{
					&s3.BucketCorsConfigurationV2CorsRuleArgs{
						AllowedHeaders: pulumi.StringArray{
							pulumi.String("*"),
						},
						AllowedMethods: pulumi.StringArray{
							pulumi.String("PUT"),
						},
						AllowedOrigins: pulumi.StringArray{
							pulumi.String("http://localhost:5173"), // TODO: whitelist
						},
						MaxAgeSeconds: pulumi.Int(3000),
					},
				},
			})
		if err != nil {
			return err
		}
		ctx.Export(resourceName, s3CORS.ID())

		// Public Access Block
		resourceName = fmt.Sprintf("%s-s3-public_access_block", projectTag)
		publicAccessBlock, err := s3.NewBucketPublicAccessBlock(
			ctx,
			resourceName,
			&s3.BucketPublicAccessBlockArgs{
				Bucket:                s3Bucket.ID(),
				BlockPublicAcls:       pulumi.Bool(false),
				BlockPublicPolicy:     pulumi.Bool(false),
				IgnorePublicAcls:      pulumi.Bool(false),
				RestrictPublicBuckets: pulumi.Bool(false),
			})
		if err != nil {
			return err
		}
		ctx.Export(resourceName, publicAccessBlock.ID())

		// Policy
		// thumbnail配下のオブジェクトに対してGetObjectを許可する
		resourceName = fmt.Sprintf("%s-s3-bucket_policy", projectTag)
		bucketPolicy, err := s3.NewBucketPolicy(
			ctx,
			resourceName,
			&s3.BucketPolicyArgs{
				Bucket: s3Bucket.ID(), // refer to the bucket created earlier
				Policy: pulumi.Any(map[string]interface{}{
					"Version": "2012-10-17",
					"Statement": []map[string]interface{}{
						{
							"Effect":    "Allow",
							"Principal": "*",
							"Action": []interface{}{
								"s3:GetObject",
							},
							"Resource": []interface{}{
								pulumi.Sprintf("arn:aws:s3:::%s/thumbnail/*", s3Bucket.ID()),
							},
						},
					},
				}),
			},
			pulumi.DependsOn([]pulumi.Resource{s3Bucket}),
		)
		if err != nil {
			return err
		}
		ctx.Export(resourceName, bucketPolicy.ID())

		return nil

	})
}

func createNameTag(tag string) pulumi.StringMap {
	return pulumi.StringMap{
		"Name": pulumi.String(tag),
	}
}

func loadFileToString(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed open file: %v", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed read file: %v", err)
	}
	return string(b), nil
}
