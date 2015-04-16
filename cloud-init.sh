#!/bin/bash
set -x
sudo apt-get install -y -qq bzr build-essential devscripts python-setuptools
sudo easy_install pip
bzr branch lp:cloud-init
cd cloud-init
sudo pip install -r requirements.txt
sudo python setup.py install --init-system=sysvinit_deb
sudo rm -rf /var/lib/cloud
sudo shutdown -r now 

