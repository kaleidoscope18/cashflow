package main

var VPCSubnetGroup = "cashflow-db-subnet-group"

var VPCInternetGateway = "cashflow-default-vpc-internet-gateway"
var PublicSubnet = "cashflow-public-subnet"
var PublicRouteTable = "cashflow-public-route-table"
var PublicRouteTableAssociation = "cashflow-public-route-table"

var BasionHostName = "cashflow-bastion-host"
var SSHKeyPair = "cashflow-bastion-host-keypair"
var BastionHostKeynamePrefix = "cashflow-bastion-host-key-"
var BastionHostSecurityGroup = "cashflow-bastion-host-security-group"

var DbSecurityGroup = "cashflow-db-security-group"
var DatabaseInstanceName = "cashflow-db"

// Exports
var DatabaseEndpoint = "cashflow-db-endpoint"
var BastionHostPublicIp = "cashflow-bastion-host-public-ip"

// CIDR blocks
// current VPC CIDR block is 172.31.0.0/16
// The allowed IPv4 CIDR block size for a subnet is between a /28 netmask and /16 netmask.
var PublicSubnetCIDR_1_1 = "172.31.0.0/28"

// Public Subnet 1 in AZ1: 172.31.0.0/28  // 172.31.0.0 - 172.31.0.15
// Public Subnet 2 in AZ2: 172.31.0.16/28 // 172.31.0.16 - 172.31.0.31
// Private Subnet 1 in AZ1: 172.31.1.0/28
// Private Subnet 2 in AZ2: 172.31.1.16/28
