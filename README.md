# Background
This repository contains a collection of recipes for building containers with SAS Viya software and other tools.

* [viya-programming](viya-programming/README.md) is a set of tools for creating SAS Viya 3.4 containers.
* [addons](addons/README.md) is a set of tools for adding to the SAS Viya 3.4 containers.
* [utilities](utilities/README.md) is a set of files that can be used when creating SAS Viya 3.4 containers.
* [all-in-one](all-in-one/README.md) is a set of tools for creating a SAS Viya 3.3 single container, for convenience and simplicity.

# Prerequisites
## SAS Software
* SAS Viya 3.4 software: a SAS Viya on Linux order and the SAS_Viya_deployment_data.zip from the Software Order Email (SOE).
* It is strongly recommended that you create a mirror repository of the SAS Viya 3.4 software. For more information, see [Create a Mirror Repository](https://go.documentation.sas.com/?docsetId=dplyml0phy0lax&docsetTarget=p1ilrw734naazfn119i2rqik91r0.htm&docsetVersion=3.4) in the SAS Viya for Linux: Deployment Guide.

**Note:** To include any future updates of the SAS Viya 3.4 software, you must rebuild  recipes with the updated SAS Viya 3.4 software.

## Other Software
* A [supported version](https://success.docker.com/article/maintenance-lifecycle) of Docker is required.
* Git is required.

# How to Clone the Repository

Here is an example of the `git clone` command for a Linux host that has Git and Docker installed. 
The command assumes that you will clone the repository in the $HOME directory.

```
cd $HOME
git clone https://github.com/sassoftware/sas-container-recipes.git
```

**Note:** On Windows, you can clone the repository by using PowerShell. Also, you can clone the repository on a Mac.

# How to Build the Images and Start the Container

A script named `build.sh` is at the repository root level. After the sassoftware/sas-container-recipes project is cloned, run `build.sh` to build a set of Docker images for SAS Viya 3.4.

The following example assumes that you are in the 
/$HOME/sas-container-recipes directory and that an [addon](addons/README.md) layer, auth-demo, is included in the build.

```
cp /path/to/SAS_Viya_deployment_data.zip viya-programming/viya-single-container
build.sh addons/auth-demo
```
**Notes:**

* The auth-demo addon sets up a demo user, which allows you to log on to SAS Studio after the container is running.
* Running `build.sh` prints to the console and to a file
named build_sas_container.log. 

After the build completes, execute `docker images` to see the images. Here is an example of the output:

```
REPOSITORY                                 TAG                 IMAGE ID            CREATED             SIZE
svc-auth-demo                              latest              965b522a213d        21 hours ago        8.52GB
viya-single-container                      latest              2351f556f15a        2 days ago          8.52GB
centos                                     latest              5182e96772bf        6 weeks ago         200MB
```

**Notes:** 

* In the preceding example, a single addon was included. You can include multiple [addons](addons/README.md) for the build.sh command. Here is an example build.sh command that includes three addons:

```
build.sh addons/auth-demo addons/ide-jupyter-python3 addons/access-pcfiles
```

* The sizes of the   viya-single-container image and the svc-auth-demo image vary depending on your order.
* In the example output, the identical size for two images can be misleading. There is an image that is 8.52 GB, which includes the three images. The svc-auth-demo image is a small image layer stacked on the viya-single-container image, which is a large image layer stacked on the centos image.
* If an [addon](addons/README.md) does not include a Dockerfile, an image is not created.  

After the build is complete, use the `docker run` command to start the container.

```
docker run \
    --detach \
    --rm \
    --env CASENV_CAS_VIRTUAL_HOST=<host-name-where-docker-is-running> \
    --env CASENV_CAS_VIRTUAL_PORT=8081 \
    --publish-all \
    --publish 8081:80 \
    --name svc-auth-demo \
    --hostname svc-auth-demo \
    svc-auth-demo

# Check the status
docker ps --filter name=svc-auth-demo --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}} \t{{.Ports}}"
CONTAINER ID        NAMES               IMAGE               STATUS              PORTS
4b426ce49b6b        svc-auth-demo       svc-auth-demo       Up 2 minutes        0.0.0.0:8081->80/tcp, 0.0.0.0:33221->443/tcp, 0.0.0.0:33220->5570/tcp    
```

After the container has started, log on to SAS Studio with the user name `sasdemo` and the password `sasdemo` at:
 
 http://_host-name-where-docker-is-running_:8081

 **Note:** The user name `sasdemo` and the password `sasdemo` are the credentials for the demo user that is set up by the auth-demo addon. 

# Troubleshooting
## Errors When Building the Docker Image
### ERRO[0000] failed to dial gRPC: cannot connect to the Docker daemon

For Linux hosts, make sure that the Docker daemon is running: 

```
sudo systemctl status docker
● docker.service - Docker Application Container Engine
   Loaded: loaded (/usr/lib/systemd/system/docker.service; disabled; vendor preset: disabled)
   Active: active (running) since Wed 2018-08-08 06:57:53 EDT; 1 months 12 days ago
     Docs: https://docs.docker.com
 Main PID: 24833 (dockerd)
    Tasks: 104
```

If the process is running, then you might need to run the Docker commands using sudo. For information on running the Docker commands without using sudo, see the [Docker documentation](https://docs.docker.com/v17.12/install/linux/linux-postinstall/).

### COPY failed: stat /var/lib/docker/tmp/docker-builderXXXXXXXXXX/\<file name\>: no such file or directory

The `docker build` command expects a Dockerfile and a build "context," which is a set of files in a specified path or URL. If the files are not present, the Docker build will display the error
message. To resolve it, make sure that the files are in the directory where the Docker build takes place.

**Notes:**

* For this project, the build context is where the files are copied from. In the examples,  `.` represents the build context.  
* In some recipes, the user is expected to copy the files into the current directory before running the `docker build` command. For example, copying files is required for building the [viya-single-container](viya-programming/viya-single-container/README.md) image and some of the 
[addon](addons/README.md) images.

### Ansible Playbook Fails 

This error might indicate that Docker is running out of space on the host where the Docker
daemon is running. To find out if more space is needed, look in the Ansible output for a message similar to the following example:

```
        "Error Summary",
        "-------------",
        "Disk Requirements:",
        "  At least 6344MB more space needed on the / filesystem."
```

If more space is needed, try pruning the Docker system:

```
docker system prune --force --volumes
```

If the error persists after the pruning, check to see if the Device Mapper storage driver is used:

```
docker system info 2>/dev/null | grep "Storage Driver"
```

If the output is _Storage Driver: devicemapper_, then the Device Mapper storage driver is used. The Device Mapper storage driver has a default layer size of 10 GB, and the SAS Viya 
image is typically larger. Possible workarounds to free up space are to change the layer size or to switch to
the [overlay2 storage driver](https://docs.docker.com/storage/storagedriver/overlayfs-driver/).


## Warnings When Building the Docker Image
### warning: /var/cache/yum/x86_64/7/**/*.rpm: Header V3 RSA/SHA256 Signature, key ID \<key\>: NOKEY

It is safe to ignore this warning. This warning indicates that the Gnu Privacy Guard (gpg) key is not available on the host, and it is followed by a call to retrieve the missing key.

# Copyright

Copyright 2018 SAS Institute Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

&nbsp;&nbsp;&nbsp;&nbsp;https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
