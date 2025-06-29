#!/bin/bash

cd cmd/monitor
go build
cd ../../
rm -rf monitor_1.0
mkdir -p monitor_1.0/DEBIAN/ monitor_1.0/usr/sbin/ monitor_1.0/etc/systemd/system/

cd monitor_1.0
mv ../cmd/monitor/monitor usr/sbin/monitor

cat>DEBIAN/control<<EOF
Package: monitor
Version: 1.0
Architecture: amd64
Maintainer: Morten Jakobsen <morten@jakeobsen.com>
Description: A program for emitting system metrics to MQTT.
Priority: optional
Section: utils
EOF

cat>etc/systemd/system/monitor.service<<EOF
[Unit]
Description=A program for emitting system metrics to MQTT.
After=network.target

[Service]
Type=simple
User=root
Group=root
ExecStart=/usr/sbin/monitor
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

cd ..

rm monitor_1.0.deb
sudo apt-get -y --purge remove monitor
dpkg-deb --build monitor_1.0/
sudo dpkg -i monitor_1.0.deb
