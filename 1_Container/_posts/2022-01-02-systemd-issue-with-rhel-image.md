---
title: "#2 systemd issue with centos/redhat base image"
author: Brandon Tsai
---


systemd issue with centos/redhat base image
==============================================

Avoid to run app with systemd inside a container because systemd changes many host-level parameters.
You need to run it as --privileged and mount the /sys/fs/cgroup volume into the container. This breaks the Docker isolation, which is usually a bad idea.

>**Note:** A light-weight process manager like `supervisord` is a better match instead of systemd.


Workaround
----------

You can get the real start command from the service configuration after installing the package.

Use `jfrog-artifactory-oss` as example.

```bash
$ yum install -y jfrog-artifactory-oss

$ find / -name 'artifactory.service'
/opt/jfrog/artifactory/misc/service/artifactory.service

$ cat /opt/jfrog/artifactory/misc/service/artifactory.service
[Unit]
Description=Setup Systemd script for Artifactory in Tomcat Servlet Engine
After=network.target

[Service]
Type=forking
GuessMainPID=yes
Restart=always
RestartSec=5
PIDFile=/var/run/artifactory.pid
ExecStart=/opt/jfrog/artifactory/bin/artifactoryManage.sh start
ExecStop=/opt/jfrog/artifactory/bin/artifactoryManage.sh stop

[Install]
WantedBy=multi-user.target
```

If the start command does not keep artifactory run in the foreground,
implement an initial/entrypoint script which tail the log files to keep the container running.

```bash
#!/bin/bash

/opt/jfrog/artifactory/bin/artifactoryManage.sh start

tail -f /var/opt/jfrog/artifactory/logs/artifactory.log
```
