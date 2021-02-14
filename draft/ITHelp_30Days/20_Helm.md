Helm
============================


Helm is a package manager for Kubernetes (think apt or yum). It works by combining several manifests into a single package that is called a chart. Helm also supports chart storage in remote or local Helm repositories that function like package registries such as Maven Central, Ruby Gems, npm registry, etc.

Helm is currently the only solution that supports

The grouping of related Kubernetes manifests in a single entity (the chart)
Basic templating and value support for Kubernetes manifests
Dependency declaration between applications (chart of charts)
A registry of available applications to be deployed (Helm repository)
A view of a Kubernetes cluster in the application/chart level
Management of installation/upgrades of charts as a whole
Built-in rollback of a chart to a previous version without running a CI/CD pipeline again
You can find a list of public curated charts in the default Helm repository. 

find helm on https://hub.helm.sh/

Several third-party tools support Helm chart creation such as Draft. Local Helm development is also supported by garden.io and/or skaffold. Check your favorite tool for native Helm support.



Helm Basic Usage
---------------


### Add a repository of Helm charts to your local Helm client:

```
$ helm repo add stable https://kubernetes-charts.storage.googleapis.com/
```


### Update the repository:

```
$ helm repo update
```

### helm search repo:

```
$ helm search repo mysql
NAME                            	CHART VERSION	APP VERSION	DESCRIPTION                                       
stable/mysql                    	1.6.7        	5.7.30     	Fast, reliable, scalable, and easy to use open-...
stable/mysqldump                	2.6.1        	2.4.1      	A Helm chart to help backup MySQL databases usi...
stable/prometheus-mysql-exporter	0.7.1        	v0.11.0    	DEPRECATED A Helm chart for prometheus mysql ex...
stable/percona                  	1.2.1        	5.7.26     	free, fully compatible, enhanced, open source d...
stable/percona-xtradb-cluster   	1.0.5        	5.7.19     	free, fully compatible, enhanced, open source d...
stable/phpmyadmin               	4.3.5        	5.0.1      	DEPRECATED phpMyAdmin is an mysql administratio...
stable/gcloud-sqlproxy          	0.6.1        	1.11       	DEPRECATED Google Cloud SQL Proxy                 
stable/mariadb                  	7.3.14       	10.3.22    	DEPRECATED Fast, reliable, scalable, and easy t...

```

### Search with available versions

```
$ helm search repo -l stable/mysql
```

### Get value information

```
$ helm show values stable/mysql
```


### Install an example MySQL chart:

```

# Install with default values
$ helm install mydb stable/mysql

# Install with custom yaml file
$ helm install mydb stable/mysql -f myvalue.yaml

# Install with value overrides
$ helm install mydb stable/mysql --set mysqlRootPassword=helmtest123

# INstall with particular version
$ helm install mydb stable/mysql --version "1.6.6"
```

### Verify that the chart has installed successfully:

```
$ helm list
NAME      	NAMESPACE	REVISION	UPDATED                              	STATUS  	CHART            	APP VERSION  
mydb      	myproject	1       	2020-10-05 11:58:23.154884 +1100 AEDT	deployed	mysql-1.6.6      	5.7.30       
```


### Upgrade to new version

```
$ helm upgrade mydb stable/mysql --version "1.6.7"
```

### Helm rollbacks

```
# Current chart version is 1.6.7
$ helm list
NAME      	NAMESPACE	REVISION	UPDATED                              	STATUS  	CHART            	APP VERSION  
mydb      	myproject	2       	2020-10-05 12:00:49.252883 +1100 AEDT	deployed	mysql-1.6.7      	5.7.30       

# Rollback to revision 1
$ helm rollback mydb 1
Rollback was a success! Happy Helming!

# verify chart version is rollback to 1.6.6
$ helm list
NAME      	NAMESPACE	REVISION	UPDATED                              	STATUS  	CHART            	APP VERSION  
mydb      	myproject	3       	2020-10-05 12:02:47.831528 +1100 AEDT	deployed	mysql-1.6.6      	5.7.30    
```


### Remove package

```
$ helm uninstall mydb
release "mydb" uninstalled

```


Conclusion: Helm vs OpenShift templates
-----------------


Helm is a package manager that also happens to include templating capabilities. Unfortunately, a lot of people focus only on the usage of Helm as a template manager and nothing else.

Technically Helm can be used as only a templating engine by stopping the deployment process in the manifest level. It is perfectly possible to use Helm only to create plain Kubernetes manifests and then install them on the cluster using the standard methods (such as kubectl). But then you miss all the advantages of Helm (especially the application registry aspect).

At the time of writing Helm is the only package manager for Kubernetes, so if you want a way to group your manifests and a registry of your running applications, there are no off-the-shelf alternative apart from Helm.

Here is a table that highlights the comparison:

| Feature | Helm  | Template |
|---------|-------|----------|
| Templating | V | V |
| Direct rollbacks and Upgrades | V |  |
| Uninstall Command | V |  |
| Registry of applications | V |  |
| Application/package dependencies| V |  |

Although Helm seems like that it is more useful, however, most helm charts in Helmhub does not suitable for OpenShift.

Although red hat have a tool called [template2helm](https://github.com/redhat-cop/template2helm) to translate current template to helm charts, but it is a very new and underdeveloped. In most cases, you still need to modify the helm chart to make it works on OpenShift platform.

