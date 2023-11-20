package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecs"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/rds"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/upstash/pulumi-upstash/sdk/go/upstash"
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
	availableZone string,
	resourceName string,
) (*ec2.Subnet, error) {
	return ec2.NewSubnet(ctx, resourceName, &ec2.SubnetArgs{
		VpcId:            vpc.ID(),
		CidrBlock:        pulumi.String(cidr),
		AvailabilityZone: pulumi.String(availableZone),
		Tags:             createNameTag(resourceName),
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

func CreateSecurityGroupForAppContainerBackend(
	ctx *pulumi.Context,
	vpc *ec2.Vpc,
	securityGroupForMaintenance *ec2.SecurityGroup,
	resourceName string,
) (*ec2.SecurityGroup, error) {
	return ec2.NewSecurityGroup(
		ctx,
		resourceName,
		&ec2.SecurityGroupArgs{
			VpcId: vpc.ID(),
			Ingress: ec2.SecurityGroupIngressArray{
				&ec2.SecurityGroupIngressArgs{
					SecurityGroups: pulumi.StringArray{
						securityGroupForMaintenance.ID(),
						// TODO: frontend container
					},
					Description: pulumi.String("for inbound"),
					Protocol:    pulumi.String("tcp"),
					FromPort:    pulumi.Int(3000),
					ToPort:      pulumi.Int(3000),
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

func CreateSecurityGroupForRDS(
	ctx *pulumi.Context,
	vpc *ec2.Vpc,
	securityGroupForMaintenance *ec2.SecurityGroup,
	securityGroupAppContainerBackend *ec2.SecurityGroup,
	resourceName string,
) (*ec2.SecurityGroup, error) {
	return ec2.NewSecurityGroup(
		ctx,
		resourceName,
		&ec2.SecurityGroupArgs{
			VpcId: vpc.ID(),
			Ingress: ec2.SecurityGroupIngressArray{
				&ec2.SecurityGroupIngressArgs{
					SecurityGroups: pulumi.StringArray{
						securityGroupForMaintenance.ID(),
					},
					Description: pulumi.String("for MaintenanceEC2Instance"),
					Protocol:    pulumi.String("tcp"),
					FromPort:    pulumi.Int(3306),
					ToPort:      pulumi.Int(3306),
				},
				&ec2.SecurityGroupIngressArgs{
					SecurityGroups: pulumi.StringArray{
						securityGroupAppContainerBackend.ID(),
					},
					Description: pulumi.String("for AppContainerBackend"),
					Protocol:    pulumi.String("tcp"),
					FromPort:    pulumi.Int(3306),
					ToPort:      pulumi.Int(3306),
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

func CreateRDSSubnetGroup(
	ctx *pulumi.Context,
	subnetZoneA *ec2.Subnet,
	subnetZoneC *ec2.Subnet,
	resourceName string,
) (*rds.SubnetGroup, error) {
	return rds.NewSubnetGroup(
		ctx,
		resourceName,
		&rds.SubnetGroupArgs{
			Name:        pulumi.String(resourceName),
			Description: pulumi.String(resourceName),
			SubnetIds: pulumi.StringArray{
				subnetZoneA.ID(),
				subnetZoneC.ID(),
			},
			Tags: createNameTag(resourceName),
		})
}

func CreateRDSInstance(
	ctx *pulumi.Context,
	dbName string,
	username string,
	password string,
	securityGroup *ec2.SecurityGroup,
	dbSubnetGroup *rds.SubnetGroup,
	resourceName string,
) (*rds.Instance, error) {
	return rds.NewInstance(
		ctx,
		resourceName,
		&rds.InstanceArgs{
			Engine:                pulumi.String("mysql"),
			EngineVersion:         pulumi.String("8.0.33"),
			InstanceClass:         pulumi.String("db.t3.micro"),
			AllocatedStorage:      pulumi.Int(20),
			DbSubnetGroupName:     dbSubnetGroup.Name,
			DbName:                pulumi.String(dbName),
			Username:              pulumi.String(username),
			Password:              pulumi.String(password),
			BackupRetentionPeriod: pulumi.Int(0),
			SkipFinalSnapshot:     pulumi.Bool(true),
			VpcSecurityGroupIds: pulumi.StringArray{
				securityGroup.ID(),
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
		// AccountID ///////////////////////////////////////////////////////////////////
		caller, err := aws.GetCallerIdentity(ctx, nil, nil)
		if err != nil {
			return err
		}
		accountId := caller.AccountId

		// Region ///////////////////////////////////////////////////////////////////////
		region, err := aws.GetRegion(ctx, nil, nil)
		if err != nil {
			return err
		}

		// VPC //////////////////////////////////////////////////////////////////////////
		resourceName := fmt.Sprintf("%s-vpc", projectTag)
		vpc, err := CreateVPC(ctx, "10.1.0.0/16", resourceName)
		if err != nil {
			return fmt.Errorf("failed create vpc: %v", err)
		}
		ctx.Export(resourceName, vpc.ID())

		// Subnet /////////////////////////////////////////////////////////////////////////
		// Public
		resourceName = fmt.Sprintf("%s-subnet-app-container", projectTag)
		subnetAppContainer, err := CreateSubnet(
			ctx, vpc, "10.1.0.0/24", "ap-northeast-1a", resourceName)
		if err != nil {
			return fmt.Errorf("failed create subnet for App Container: %v", err)
		}
		ctx.Export(resourceName, subnetAppContainer.ID())

		// Private
		resourceName = fmt.Sprintf("%s-subnet-private-1a", projectTag)
		subnetPrivate1a, err := CreateSubnet(
			ctx, vpc, "10.1.1.0/24", "ap-northeast-1a", resourceName)
		if err != nil {
			return fmt.Errorf("failed create subnet for private: %v", err)
		}
		ctx.Export(resourceName, subnetPrivate1a.ID())

		resourceName = fmt.Sprintf("%s-subnet-private-1c", projectTag)
		subnetPrivate1c, err := CreateSubnet(
			ctx, vpc, "10.1.2.0/24", "ap-northeast-1c", resourceName)
		if err != nil {
			return fmt.Errorf("failed create subnet for private: %v", err)
		}
		ctx.Export(resourceName, subnetPrivate1c.ID())

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

		// S3 ////////////////////////////////////////////////////////////////////////
		// bucket
		resourceName = fmt.Sprintf("%s-s3-bucket", projectTag)
		bucketName := fmt.Sprintf("blog-%s", config.AWSAccountId)
		bucketName := fmt.Sprintf("blog-%s", accountId)
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

		// ECSタスク実行用ロール
		secretsManagerArn := fmt.Sprintf("arn:aws:secretsmanager:%s:%s:secret:%s", region.Name, accountId, config.SecretsManagerSecretId)
		kmsArn := fmt.Sprintf("arn:aws:kms:%s:%s:key/%s", region.Name, accountId, config.KmsKeyId)
		resourceName = fmt.Sprintf("%s-iam-role-for-ecs-task-execute", projectTag)
		ecsTaskExecutionRole, err := iam.NewRole(
			ctx,
			resourceName,
			&iam.RoleArgs{
				AssumeRolePolicy: pulumi.String(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Effect": "Allow",
						"Principal": {
							"Service": "ecs-tasks.amazonaws.com"
						},
						"Action": "sts:AssumeRole"
					}]
				}`),
				ManagedPolicyArns: pulumi.StringArray{
					pulumi.String("arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"),
					// ECSからのECRをpullするために必要
					pulumi.String("arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"),
				},
				// ECSのcontainerDefinitionsのawslogs-create-groupに必要
				InlinePolicies: iam.RoleInlinePolicyArray{
					&iam.RoleInlinePolicyArgs{
						Name: pulumi.String("ecs-task-policy-logs"),
						Policy: pulumi.String(`{
							"Version": "2012-10-17",
							"Statement": [
								{
								   "Effect": "Allow",
								   "Action": [
										"logs:CreateLogGroup"
								   ],
								  "Resource": "*"
								}
							]
						}`),
					},
					&iam.RoleInlinePolicyArgs{
						Name: pulumi.String("ecs-task-policy-secretsmanager"),
						Policy: pulumi.String(`{
							"Version": "2012-10-17",
							"Statement": [
								{
									"Effect": "Allow",
									"Action": [
										"secretsmanager:GetSecretValue",
										"kms:Decrypt"
									],
									"Resource": [
										"` + secretsManagerArn + `",
										"` + kmsArn + `"
									]
								}
							]
						}`),
					},
				},
			})
		if err != nil {
			return fmt.Errorf("failed create iam role for ecs task execution: %v", err)
		}
		ctx.Export(resourceName, ecsTaskExecutionRole.ID())

		// ECSタスクロール AppContainer Backend
		resourceName = fmt.Sprintf("%s-iam-role-for-ecs-task-backend", projectTag)
		ecsTaskRole, err := iam.NewRole(
			ctx,
			resourceName,
			&iam.RoleArgs{
				AssumeRolePolicy: pulumi.String(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Effect": "Allow",
						"Principal": {
							"Service": "ecs-tasks.amazonaws.com"
						},
						"Action": "sts:AssumeRole"
					}]
				}`),
				InlinePolicies: iam.RoleInlinePolicyArray{
					&iam.RoleInlinePolicyArgs{
						Name: pulumi.String("ecs-task-policy"),
						Policy: pulumi.String(`{
							"Version": "2012-10-17",
							"Statement": [
								{
								   "Effect": "Allow",
								   "Action": [
										"ssmmessages:CreateControlChannel",
										"ssmmessages:CreateDataChannel",
										"ssmmessages:OpenControlChannel",
										"ssmmessages:OpenDataChannel"
								   ],
								  "Resource": "*"
								},
								{
								   "Effect": "Allow",
								   "Action": [
										"s3:GetObject"
								   ],
								   "Resource": "arn:aws:s3:::` + bucketName + `/*"
								}
							]
						}`),
					},
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed create iam role for ecs task: %v", err)
		}
		ctx.Export(resourceName, ecsTaskRole.ID())

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

		resourceName = fmt.Sprintf("%s-sg-public-app-container", projectTag)
		securityGroupAppContainerBackend, err := CreateSecurityGroupForAppContainerBackend(
			ctx,
			vpc,
			securityGroupPublicMaintenance,
			resourceName)
		if err != nil {
			return fmt.Errorf("failed create security group for AppContainer Backend: %v", err)
		}
		ctx.Export(resourceName, securityGroupAppContainerBackend.ID())

		resourceName = fmt.Sprintf("%s-sg-rds", projectTag)
		securityGroupRDS, err := CreateSecurityGroupForRDS(
			ctx,
			vpc,
			securityGroupPublicMaintenance,
			securityGroupAppContainerBackend,
			resourceName)
		if err != nil {
			return fmt.Errorf("failed create security group for RDS: %v", err)
		}
		ctx.Export(resourceName, securityGroupRDS.ID())

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
				SubnetId:                 subnetAppContainer.ID(),
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

		// KVS upstash
		resourceName = fmt.Sprintf("%s-kvs-redis-by-upstash", projectTag)
		redisKVS, err := upstash.NewRedisDatabase(ctx, resourceName, &upstash.RedisDatabaseArgs{
			DatabaseName: pulumi.String("blog-kvs"),
			Region:       pulumi.String("ap-northeast-1"),
			Tls:          pulumi.Bool(true),
			Eviction:     pulumi.Bool(true),
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		ctx.Export(resourceName, redisKVS.ID())

		// ECS ////////////////////////////////////////////////////////////////////////
		// TaskDefinition
		taskDefinition, err := loadEcsContainerDefinition(
			"./container_definition.json", accountId, region.Name, config.SecretsManagerSecretId)
		if err != nil {
			return fmt.Errorf("failed load ecs task definition: %v", err)
		}
		resourceName = fmt.Sprintf("%s-ecs-task-definition", projectTag)
		ecsTaskDefinition, err := ecs.NewTaskDefinition(
			ctx,
			resourceName,
			&ecs.TaskDefinitionArgs{
				Family:                  pulumi.String("blog-backend"),
				NetworkMode:             pulumi.String("awsvpc"),
				Cpu:                     pulumi.String("1024"),
				Memory:                  pulumi.String("3072"),
				TaskRoleArn:             ecsTaskRole.Arn,
				ExecutionRoleArn:        ecsTaskExecutionRole.Arn,
				RequiresCompatibilities: pulumi.StringArray{pulumi.String("FARGATE")},
				ContainerDefinitions:    pulumi.String(taskDefinition),
			})
		ctx.Export(resourceName, ecsTaskDefinition.ID())

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

func loadEcsContainerDefinition(
	path string, awsAccountId string, region string, secretsManagerId string) (string, error) {
	type Values struct {
		AwsAccountId     string
		Region           string
		SecretsManagerId string
	}
	definition, err := loadFileToString(path)
	if err != nil {
		return "", fmt.Errorf("failed load ecs container definition: %v", err)
	}
	tmpl, err := template.New("ecsTaskDefinition").Parse(definition)
	if err != nil {
		return "", fmt.Errorf("failed parse ecs container definition: %v", err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, Values{
		AwsAccountId:     awsAccountId,
		Region:           region,
		SecretsManagerId: secretsManagerId,
	})
	if err != nil {
		return "", fmt.Errorf("failed execute ecs container definition: %v", err)
	}
	return buffer.String(), nil
}
