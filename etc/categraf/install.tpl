#!/bin/bash

cd /opt
curl -o categraf.tar.gz http://{{.ServerHost}}:{{.ServerPort}}/categraf/download/?token={{.Token}}
tar -xvf categraf.tar.gz
cd categraf
chmod 755 install
./install