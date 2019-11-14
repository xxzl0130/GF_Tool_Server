#!/bin/sh
wget https://github.com/xxzl0130/GF_Tool_Server/releases/download/v1.1/GF_Tool_Server.zip
unzip -o -d ./GF_Tool_Server GF_Tool_Server_linux.zip
chmod +x ./GF_Tool_Server/GF_Tool_Server
cp -rf ./GF_Tool_Server /etc/
cp ./GF_Tool_Server/scrpits/GF_Server /etc/init.d
chmod 755 /etc/init.d/GF_Server
update-rc.d GF_Server defaults 95