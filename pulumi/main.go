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
		availabilityZone := os.Getenv("AZ")

		// control tower creates a unique vpc and removes the default one
		defaultVPC, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{})
		if err != nil {
			return err
		}

		igw, err := ec2.NewInternetGateway(ctx, VPCInternetGateway, &ec2.InternetGatewayArgs{
			VpcId: pulumi.StringPtr(defaultVPC.Id),
			Tags: pulumi.StringMap{
				"ApplicationName": pulumi.String("cashflow"),
			},
		})
		if err != nil {
			return err
		}

		publicSubnet, err := ec2.NewSubnet(ctx, PublicSubnet, &ec2.SubnetArgs{
			VpcId:            pulumi.String(defaultVPC.Id),
			CidrBlock:        pulumi.String(PublicSubnetCIDR_1_1),
			AvailabilityZone: pulumi.String(availabilityZone),
			// MapPublicIpOnLaunch: pulumi.Bool(true),
			Tags: pulumi.StringMap{
				"ApplicationName": pulumi.String("cashflow"),
			},
		})
		if err != nil {
			return err
		}

		// To make a subnet public (internet accessible):
		// its route table must have a route that sends all destined internet traffic (0.0.0.0/0) to an internet gateway
		publicRouteTable, err := ec2.NewRouteTable(ctx, PublicRouteTable, &ec2.RouteTableArgs{
			VpcId: pulumi.String(defaultVPC.Id),
			Routes: ec2.RouteTableRouteArray{
				ec2.RouteTableRouteArgs{
					CidrBlock: pulumi.String("0.0.0.0/0"),
					GatewayId: igw.ID(),
				},
			},
			Tags: pulumi.StringMap{
				"ApplicationName": pulumi.String("cashflow"),
			},
		})
		if err != nil {
			return err
		}

		_, err = ec2.NewRouteTableAssociation(ctx, PublicRouteTableAssociation, &ec2.RouteTableAssociationArgs{
			SubnetId:     publicSubnet.ID(),
			RouteTableId: publicRouteTable.ID(),
		})
		if err != nil {
			return err
		}

		/*
			Summary: created an internet gateway for the VPC then created a public subnet
			with a route table that directs all internet trafic to the internet gateway
		*/

		// create a key pair resource to SSH access the bastion host
		keyPair, err := ec2.NewKeyPair(ctx, BastionHostKeyPair, &ec2.KeyPairArgs{
			KeyNamePrefix: pulumi.String(BastionHostKeynamePrefix),
			PublicKey:     pulumi.String(conf.Get("public-key")),
			Tags: pulumi.StringMap{
				"ApplicationName": pulumi.String("cashflow"),
			},
		})
		if err != nil {
			return err
		}

		// Allow SSH access from internet -> bastion host & bastion host -> everywhere
		bastionSecurityGroup, err := ec2.NewSecurityGroup(ctx, BastionHostSecurityGroup, &ec2.SecurityGroupArgs{
			VpcId: pulumi.StringPtr(defaultVPC.Id),
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Description: pulumi.String("Allows SSH inbound traffic from anywhere"),
					Protocol:    pulumi.String("tcp"),
					FromPort:    pulumi.Int(22),
					ToPort:      pulumi.Int(22),
					CidrBlocks:  pulumi.StringArray{pulumi.String("0.0.0.0/0")},
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

		// create the EC2 instance that will serve as the bastion host
		bastionHost, err := ec2.NewInstance(ctx, BasionHostName, &ec2.InstanceArgs{
			InstanceType:             ec2.InstanceType("t2.micro"),
			VpcSecurityGroupIds:      pulumi.StringArray{bastionSecurityGroup.ID().ToStringOutput()},
			KeyName:                  keyPair.KeyName,
			SubnetId:                 publicSubnet.ID(),
			Ami:                      pulumi.StringPtr("ami-0bb84b8ffd87024d8"), // Amazon Linux 2023 AMI 2023.4.20240513.0 x86_64 HVM kernel-6.1
			AssociatePublicIpAddress: pulumi.Bool(true),
			Tags: pulumi.StringMap{
				"ApplicationName": pulumi.String("cashflow"),
			},
		})
		if err != nil {
			return err
		}

		ctx.Export(BastionHostPublicIp, bastionHost.PublicIp)

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
					Description:    pulumi.String("Allow tcp from bastion"),
					Protocol:       pulumi.String("tcp"),
					FromPort:       pulumi.Int(3306), // MySQL port, use 5432 for PostgreSQL
					ToPort:         pulumi.Int(3306),
					SecurityGroups: pulumi.StringArray{bastionSecurityGroup.ID().ToStringOutput()},
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

		// create the RDS MySQL database in subnet group
		dbName := os.Getenv("DB_NAME")
		dbUsername := os.Getenv("DB_USERNAME")
		dbPassword := os.Getenv("DB_PASSWORD")

		rdsInstance, err := rds.NewInstance(ctx, DatabaseInstanceName, &rds.InstanceArgs{
			InstanceClass:       pulumi.String("db.t3.micro"), // Free tier eligible
			AllocatedStorage:    pulumi.Int(5),                // Free tier eligible
			Engine:              pulumi.String("mysql"),
			EngineVersion:       pulumi.String("8.0.35"),
			DbName:              pulumi.String(dbName),
			Username:            pulumi.String(dbUsername),
			Password:            pulumi.String(dbPassword),
			SkipFinalSnapshot:   pulumi.Bool(true),
			DbSubnetGroupName:   subnetGroup.Name,
			VpcSecurityGroupIds: pulumi.StringArray{dbSecurityGroup.ID().ToStringOutput()},
			Tags: pulumi.StringMap{
				"ApplicationName": pulumi.String("cashflow"),
			},
		})
		if err != nil {
			return err
		}

		ctx.Export(DatabaseEndpoint, rdsInstance.Endpoint)

		return nil
	})
}
