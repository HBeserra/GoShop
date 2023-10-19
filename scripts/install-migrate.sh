#!/bin/su root
if [[ $(/usr/bin/id -u) -ne 0 ]]; then
    echo "Not running as root"
    exit 1
fi

# This program require super user to execute

curl -s -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
mv  migrate /usr/bin/migrate

# Verify than the program is installed in machine
if [[ -z $(which migrate) ]]; then
    echo "Error: Fail to install command"
    exit 2
fi

MIGRATE_VERSION=$(migrate -version 2>&1)

# Verify than the program runs by geting the version 
if [[ -z MIGRATE_VERSION  || $MIGRATE_VERSION =~ ^(v)[0-9]+.?([0-9]+)?.?([0-9]+)?$  ]]; then
    echo "Error: Fail to install command"
    exit 2
fi

if ! [[ $MIGRATE_VERSION =~ ^(v)?[0-9]+.?([0-9]+)?.?([0-9]+)? ]]; then
  echo "Error: Fail to install the command"
  echo $MIGRATE_VERSION
  exit 2
fi