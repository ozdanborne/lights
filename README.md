
```
git clone https://github.com/djosborne/dash-lights /opt/
virtualenv /opt/dash-lights/ 
source /opt/dash-lights/bin/activate
pip install -r /opt/dash-lights/requirements.txt
ln -s  /opt/dash-lights/dash-lights.service /lib/systemd/system/dash-lights.service
systemctl start dash-lights
```
