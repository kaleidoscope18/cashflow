# Pulumi Infra

## SSH

To connect to the public ec2 instance (bastion host)

```sh
ssh -i ~/.ssh/id_rsa ec2-user@<ec2 instance public ip>
```

Install mysql client with:

```sh
sudo su -
dnf -y localinstall https://dev.mysql.com/get/mysql80-community-release-el9-4.noarch.rpm
dnf -y install mysql mysql-community-client
exit
```

Optional: verify rds endpoint

```sh
sudo yum install -y nc
nc -zv <rds endpoint> 3306
```

Then you can connect to the mysql database with the pre-established credentials, i.e.:

```sh
mysql -h cashflow-db9b3ad36.c9gia6eo0ryf.us-east-1.rds.amazonaws.com -u admin -p
```

type `exit` to exit all instances.
