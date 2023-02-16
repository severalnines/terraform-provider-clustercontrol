variable "awsprops" {
    default = {
    region = "eu-north-1"
#    vpc = "vpc-051a880b5109c71f3"
    vpc = "vpc-040d2949629b597f8"
    #ami = "ami-0333305f9719618c7"
    #ami = "ami-00c70b245f5354c0a"
    ami = "ami-0916e5af9db35b1a4"
    itype = "t3.micro"
#    subnet = "subnet-04fe882a8093e6cdd"
    subnet = "subnet-0d6ccc312ba7c4551"
    publicip = true
    secgroupname = "JOHANSECGROUPTEST"
  }
}

provider "aws" {
  region     = "eu-north-1"
  #region     = "eu-west-1"
  access_key = ""
  secret_key = ""
}


resource "aws_security_group" "project-iac-sg" {
  name = lookup(var.awsprops, "secgroupname")
  description = lookup(var.awsprops, "secgroupname")
  vpc_id = lookup(var.awsprops, "vpc")
  
  // To Allow SSH Transport
  ingress {
    from_port = 22
    protocol = "tcp"
    to_port = 22
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 22
    protocol = "tcp"
    to_port = 22
    cidr_blocks = ["90.230.41.169/32"]
  }

  ingress {
    from_port = 3306
    protocol = "tcp"
    to_port = 3306
    cidr_blocks = ["90.230.41.169/32"]
  }
  ingress {
    from_port = 3306
    protocol = "tcp"
    to_port = 3306
    cidr_blocks = ["0.0.0.0/0"]
#    security_group_id = aws_security_group.private.id
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }
}


resource "aws_security_group_rule" "mysql" {
  type              = "ingress"
  from_port         = 3306
  to_port           = 3306
  protocol          = "tcp"
#  cidr_blocks       = ["0.0.0.0"]
  source_security_group_id = aws_security_group.project-iac-sg.id
  security_group_id = aws_security_group.project-iac-sg.id
}

resource "aws_key_pair" "deployer" {
  key_name   = "deployer-key"
  public_key = ""
}

resource "aws_instance" "project-iac" {
  count = 3
  ami = lookup(var.awsprops, "ami")
  instance_type = lookup(var.awsprops, "itype")
  subnet_id = lookup(var.awsprops, "subnet") #FFXsubnet2
  associate_public_ip_address = lookup(var.awsprops, "publicip")
  key_name = aws_key_pair.deployer.key_name


  vpc_security_group_ids = [
    aws_security_group.project-iac-sg.id
  ]
  root_block_device {
    delete_on_termination = true
    volume_size = 50
    volume_type = "gp2"
  }
  tags = {
    Name ="SERVER"
    Environment = "DEV"
    OS = "UBUNTU"
    Managed = "S9s"
  }

#    user_data = <<-EOF
#cloud-config
#runcmd:
#  - sudo wget https://repo.percona.com/apt/percona-release_latest.$(lsb_release -sc)_all.deb
#  - sudo dpkg -i percona-release_latest.$(lsb_release -sc)_all.deb
#  - sudo percona-release enable-only tools release
#  - sudo apt-get update
#  - sudo percona-release setup ps80
#  - sudo apt-get install -y socat percona-server-server percona-server-client percona-xtrabackup-80
#EOF

  depends_on = [ aws_security_group.project-iac-sg ]
}

output "ec2instance" {
  value = aws_instance.project-iac.*.public_ip
}

