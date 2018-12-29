# Create persisting web application with SystemD
* Create configuration file
    * `cd /etc/systemd/system/`
    * Create file `myfile.service` with the following content
```
[Unit]
Description=Go Web Server

[Service]
ExecStart=/home/user/main.go
User=root
Group=root
Restart=always

[Install]
WantedBy=multi-user.target
```

* Add the service to systemd
```
sudo systemctl enable myfile.service
```
* Activate the service
```
sudo systemctl start myfile.service
```
* Ask status
```
sudo systemctl status myfile.service
```
* Stop service 
```
sudo systemctl stop myfile.service
```


