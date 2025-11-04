// VPC
resource "aws_vpc" "vpn_onprem_vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "kfupm-vpn-onprem-vpc"
    project = "master"
  }
}

resource "aws_vpc" "vpn_remote_vpc" {
  cidr_block = "10.10.0.0/16"

  tags = {
    Name = "kfupm-vpn-remote-vpc"
  }
}

// Subnets
resource "aws_subnet" "vpn_onprem_subnet_public" {
  vpc_id     = aws_vpc.vpn_onprem_vpc.id
  cidr_block = "10.0.0.0/20"

  tags = {
    Name = "kfupm-vpn-onprem-subnet-public1-me-south-1a"
    project = "master"
  }
}

resource "aws_subnet" "vpn_onprem_subnet_private" {
  vpc_id     = aws_vpc.vpn_onprem_vpc.id
  cidr_block = "10.0.128.0/20"

  tags = {
    Name = "kfupm-vpn-onprem-subnet-private1-me-south-1a"
    project = "master"
  }
}

resource "aws_subnet" "vpn_remote_subnet_public" {
  vpc_id     = aws_vpc.vpn_remote_vpc.id
  cidr_block = "10.10.0.0/20"

  tags = {
    Name = "kfupm-vpn-remote-subnet-public1-me-south-1a"
  }
}

resource "aws_subnet" "vpn_remote_subnet_private" {
  vpc_id     = aws_vpc.vpn_remote_vpc.id
  cidr_block = "10.10.128.0/20"

  tags = {
    Name = "kfupm-vpn-remote-subnet-private1-me-south-1a"
  }
}