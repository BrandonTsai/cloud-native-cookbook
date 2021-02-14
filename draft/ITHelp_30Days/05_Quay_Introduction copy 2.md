

Introduction
------------

Despite the fact that we can build Image from Dockerfile before running a container, it is better to keep these Images on "Image Registry" and pull these images before running the container.

There are many Image Registry solutions, the most famous is Dockerhub. You can find numerous public registries and images on Dockerhub. However, if you want to build a local Image Registry on your private development environment or in your datacenter, I would suggest you consider the Red Hat solution - Quay.

Red Hat Quay container and application registry provides secure storage, distribution, and deployment of containers on any infrastructure. It is available as a standalone component or in conjunction with OpenShift.

Compare Harbor and Quay
----------------------

There is another similar and popular project - [Harbor](https://github.com/goharbor/harbor). It is an open source trusted cloud native registry project and hosted by the Cloud Native Computing Foundation (CNCF). Following is the comparison between Harbor and Quay.


| Product | Quay | Harbor |
|---------|------|--------|
| Language                      | Python | Golang |
| Type                          | Public (quay.io), private | private |
| Authentication                | LDAP, OIDC, Google, Github  | LDAP, OIDC, DB |
| Robot Account                 | O | O |
| Permission Management         | O | O |
| Security Scan                 | O | O |
| Image Signing                 | 0 | 0 |
| Image Clearning               | O | O |
| Helm Application Management   | O | O |
| Image Mirroring               | O | O |
| Notification                  | Webhook, Email, Slack | Webhook |

You can see there are not much different between Quay and Harbor.
If you are planing to run containers on OpenShift, I would recommend you using Quay, because it is supported by Red Hat and highly integrated with OpenShift.


Quay Installation
-----------------------------

To install open source version, please refer https://github.com/quay/quay/blob/master/docs/getting_started.md. To install Red Hat Quay on your local environment, please refer [https://access.redhat.com/documentation/en-us/red_hat_quay/3.3/](https://access.redhat.com/documentation/en-us/red_hat_quay/3.3/). This article will focus on introducing the functionality of Red Hat Quay version 3.2.1.


Quay Basic Usage
-------------------

In Quay, you can create a different organization for different business units or different business usage. Each organization contains isolated teams, repositories, robot accounts and API Tokens.

1. Create organizations and repositories

![](images/03_quay/02.png)


2. Copy Images from Red Hat Registry to Quay via Skopeo

```
$ sudo docker login quay-eu-uat
$ sudo skopeo copy --src-tls-verify=false --dest-tls-verify=false docker://registry.redhat.io/rhscl/nginx-116-rhel7 docker://quay-eu-uat/application-images/test:1
```

3. Check Scan Result

![](images/03_quay/03.png)

From above image you can see the images has pass the security scan.


Robot Account Introduce
-------------------------

There are many circumstances where permissions for repositories need to be shared with other services, such as CI/CD pipeline and Openshift. To support this case, Quay allows the use of robot accounts that are owned by a user or organization to access multiple repositories.


1. Create a robot account

![](images/03_quay/04.png)

2. Give the robot account a name

![](images/03_quay/05.png)

3. Give robot account permission to access your repository

![](images/03_quay/06.png)


4. Get the credential of the robot account.

![](images/03_quay/07.png)

![](images/03_quay/08.png)

5. Logging in with a robot account on other Server

```
$ docker login -u="application-images+jenkins" -p="6BC26ZL0CUZQTJKL1SWKZIO9ZD58TDLS8O6VONE4VVNF9M1ZQGGMCVBXORNC0BNG" quay-eu-uat
```

6. Make sure you can pull images with this robot account

```
$ docker pull quay-eu-uat/application-images/test:1
```

Conclusion
-----------

This article compare tow Image registry solutions: Harbor and Quay. I also introduce the very basic usage of Red Hat Quay. Next Day, I will talk more advanced usage of Quay.


