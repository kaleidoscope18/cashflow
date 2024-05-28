package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/rds"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
		env := conf.Get("cashflow:environmentName")

		ctx.Log.Info(env, nil)

		// TODO: use different files for different environments
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		err = Containerize(ctx)
		if err != nil {
			return err
		}

		defaultVPC, publicSubnet, err := Network(ctx)
		if err != nil {
			return err
		}

		ec2SecurityGroup, err := ec2.NewSecurityGroup(ctx, ServerSecurityGroup, &ec2.SecurityGroupArgs{
			VpcId:       pulumi.StringPtr(defaultVPC.Id),
			Description: pulumi.String("Bastion host security group"),
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Description: pulumi.String("SSH from my computer"),
					Protocol:    pulumi.String("tcp"),
					FromPort:    pulumi.Int(22),
					ToPort:      pulumi.Int(22),
					CidrBlocks:  pulumi.StringArray{pulumi.String("184.162.158.114/32")},
				},
				ec2.SecurityGroupIngressArgs{
					Description: pulumi.String("HTTP from my computer"),
					Protocol:    pulumi.String("tcp"),
					FromPort:    pulumi.Int(80),
					ToPort:      pulumi.Int(80),
					CidrBlocks:  pulumi.StringArray{pulumi.String("184.162.158.114/32")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Description: pulumi.String("Allows all outbound traffic"),
					Protocol:    pulumi.String("-1"),
					FromPort:    pulumi.Int(0),
					ToPort:      pulumi.Int(0),
					CidrBlocks:  pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Tags: pulumi.StringMap{
				"ApplicationName": pulumi.String("cashflow"),
			},
		})
		if err != nil {
			return err
		}

		// create a subnet group spanning the subnets for RDS instance
		subnetGroup, err := rds.NewSubnetGroup(ctx, VPCSubnetGroup, &rds.SubnetGroupArgs{
			SubnetIds: pulumi.StringArray{
				pulumi.String("subnet-0e33cbf0fa4156170").ToStringOutput(),
				pulumi.String("subnet-0bca036976b482b60").ToStringOutput(),
				pulumi.String("subnet-0d3bb1a6d4a115122").ToStringOutput(),
			},
			Tags: pulumi.StringMap{
				"ApplicationName": pulumi.String("cashflow"),
			},
		})
		if err != nil {
			return err
		}

		dbSecurityGroup, err := ec2.NewSecurityGroup(ctx, DbSecurityGroup, &ec2.SecurityGroupArgs{
			VpcId:       pulumi.String(defaultVPC.Id),
			Description: pulumi.String("Allow TCP access from bastion host"),
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Description:    pulumi.String("Allow tcp from ec2 bastion"),
					Protocol:       pulumi.String("tcp"),
					FromPort:       pulumi.Int(3306), // MySQL port, use 5432 for PostgreSQL
					ToPort:         pulumi.Int(3306),
					SecurityGroups: pulumi.StringArray{ec2SecurityGroup.ID().ToStringOutput()},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Description: pulumi.String("Allow all outbound traffic"),
					Protocol:    pulumi.String("-1"),
					FromPort:    pulumi.Int(0),
					ToPort:      pulumi.Int(0),
					CidrBlocks:  pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Tags: pulumi.StringMap{
				"ApplicationName": pulumi.String("cashflow"),
			},
		})
		if err != nil {
			return err
		}

		dbName := os.Getenv("DB_NAME")
		dbUsername := os.Getenv("DB_USERNAME")
		dbPassword := os.Getenv("DB_PASSWORD")

		rdsInstance, err := rds.NewInstance(ctx, DatabaseInstanceName, &rds.InstanceArgs{
			InstanceClass:           pulumi.String("db.t3.micro"), // Free tier eligible
			AllocatedStorage:        pulumi.Int(5),                // Free tier eligible
			Engine:                  pulumi.String("mysql"),
			EngineVersion:           pulumi.String("8.0.35"),
			AutoMinorVersionUpgrade: pulumi.Bool(true),
			DbName:                  pulumi.String(dbName),
			Username:                pulumi.String(dbUsername),
			Password:                pulumi.String(dbPassword),
			SkipFinalSnapshot:       pulumi.Bool(true),
			DbSubnetGroupName:       subnetGroup.Name,
			VpcSecurityGroupIds:     pulumi.StringArray{dbSecurityGroup.ID().ToStringOutput()},
			Tags: pulumi.StringMap{
				"ApplicationName": pulumi.String("cashflow"),
			},
		})
		if err != nil {
			return err
		}

		ctx.Export(DatabaseEndpoint, rdsInstance.Endpoint)

		// create a key pair resource to SSH access the bastion host
		keyPair, err := ec2.NewKeyPair(ctx, SSHKeyPair, &ec2.KeyPairArgs{
			KeyNamePrefix: pulumi.String(BastionHostKeynamePrefix),
			PublicKey:     pulumi.String(conf.Get("publicKey")),
			Tags: pulumi.StringMap{
				"ApplicationName": pulumi.String("cashflow"),
			},
		})
		if err != nil {
			return err
		}

		ec2Instance, err := ec2.NewInstance(ctx, ServerHostName, &ec2.InstanceArgs{
			InstanceType:             ec2.InstanceType("t2.micro"),
			KeyName:                  keyPair.KeyName,
			SubnetId:                 publicSubnet.ID(),
			Ami:                      pulumi.StringPtr("ami-0bb84b8ffd87024d8"), // Amazon Linux 2023 AMI 2023.4.20240513.0 x86_64 HVM kernel-6.1
			AssociatePublicIpAddress: pulumi.Bool(true),
			VpcSecurityGroupIds:      pulumi.StringArray{ec2SecurityGroup.ID().ToStringOutput()},
			Tags: pulumi.StringMap{
				"ApplicationName": pulumi.String("cashflow"),
			},
		}, pulumi.DependsOn([]pulumi.Resource{
			rdsInstance,
		}))
		if err != nil {
			return err
		}

		ctx.Export(BastionHostPublicIp, ec2Instance.PublicIp)

		return nil
	})
}
