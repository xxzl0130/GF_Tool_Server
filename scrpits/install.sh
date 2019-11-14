#!/bin/sh
wget https://github.com/xxzl0130/GF_Tool_Server/releases/download/v1.0/GF_Tool_Server_linux.zip
unzip -o -d ./GF_Tool_Server GF_Tool_Server_linux.zip
mkdir /etc/GF_Tool_Server
chmod +x ./GF_Tool_Server/GF_Tool_Server
cp -f ./GF_Tool_Server /etc/GF_Tool_Server
cp ./GF_Tool_Server/GF_Server /etc/init.d
chmod 755 /etc/init.d/GF_Server
update-rc.d GF_Server defaults 95