/usr/local/bind-9.11.37/bin/nsupdate -y chinanet-key:dMZeJ1vhhjjoig+sFYjZvg==
server 127.0.0.1
zone guahao-local.com

update add jason77.guahao-local.com 86400 A 192.168.99.100
send


/usr/local/bind-9.11.37/bin/nsupdate -y defaultisp-key:dkMHny2EEQSmvfiGKeL9YA==
server 127.0.0.1
zone guahao-local.com

update add jason77.guahao-local.com 86400 A 115.115.115.115
send

10.10.14.34 是master