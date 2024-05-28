package main

import (
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Network(ctx *pulumi.Context) (*ec2.LookupVpcResult, *ec2.Subnet, error) {
	// control tower creates a unique vpc and removes the default one
	defaultVPC, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{})
	if err != nil {
		return nil, nil, err
	}

	igw, err := ec2.NewInternetGateway(ctx, VPCInternetGateway, &ec2.InternetGatewayArgs{
		VpcId: pulumi.StringPtr(defaultVPC.Id),
		Tags: pulumi.StringMap{
			"ApplicationName": pulumi.String("cashflow"),
		},
	})
	if err != nil {
		return nil, nil, err
	}

	availabilityZone := os.Getenv("AZ")

	publicSubnet, err := ec2.NewSubnet(ctx, PublicSubnet, &ec2.SubnetArgs{
		VpcId:            pulumi.String(defaultVPC.Id),
		CidrBlock:        pulumi.String(PublicSubnetCIDR_1_1),
		AvailabilityZone: pulumi.String(availabilityZone),
		Tags: pulumi.StringMap{
			"ApplicationName": pulumi.String("cashflow"),
		},
	})
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	_, err = ec2.NewRouteTableAssociation(ctx, PublicRouteTableAssociation, &ec2.RouteTableAssociationArgs{
		SubnetId:     publicSubnet.ID(),
		RouteTableId: publicRouteTable.ID(),
	})
	if err != nil {
		return nil, nil, err
	}

	/*
		Summary: attached internet gateway to VPC
		created a public subnet + route table that directs all internet trafic to internet gateway
	*/

	return defaultVPC, publicSubnet, nil
}
