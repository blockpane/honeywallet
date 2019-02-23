#!/bin/bash

if [[ "$USER" != "root" ]]; then
	echo "please run as root."
	exit 1
fi

mkdir -p /opt/honeywallet/data
chown -R 30303 /opt/honeywallet
chmod 01777 /opt/honeywallet/data
docker-compose build
docker-compose up -d

sleep 2
chown -R 30303 /opt/honeywallet

echo "done, watch logs using 'docker-compose logs -f --tail 50'"
echo "    the dashboard is running on this host on port 80"

