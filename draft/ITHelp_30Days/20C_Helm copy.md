Helm 淺談
==========


Helm 我想就不用多做介紹了吧，它是 Kubernetes 的套件管理工具（類似 linux 系統中的 apt 或 yum ）。 它可以：

- 把多個 Kubernetes manifests 包成單一套件。
- 提供基本的模板功能（Templating)，你可以透過 “value.yaml” 來插入適合的值到 Kubernetes manifests 中。
- 有 Helm repository 可以提供可用的 Helm Charts 讓你直接使用來部署。你可以在 https://hub.helm.sh/ 找到許多公開且有用的 Helm Charts。
- 方便管理跟應用程序相關的所有 Kubernetes manifests，例如安裝，升新版本及刪除．
- 內建的 rollback 機制，讓你有問題時可以快速地回到前一個版本，而不用再透過 CI/CD Pipeline 來重新部署。


基於安全性問題，我建議使用 Helm 3 而不是 Helm 2。


Helm 3 基本用法複習
---------------


新增 Helm Repository

```
$ helm repo add stable https://kubernetes-charts.storage.googleapis.com/

$ helm repo update
```

搜尋可以用的 Repository

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

搜尋可以用的 Repository 並列出所有可用的版本

```
$ helm search repo -l stable/mysql
```

取得 value 預設值資訊

```
$ helm show values stable/mysql
```


不同的安裝方式

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

確認套件已經成功安裝:

```
$ helm list
NAME      	NAMESPACE	REVISION	UPDATED                              	STATUS  	CHART            	APP VERSION  
mydb      	myproject	1       	2020-10-05 11:58:23.154884 +1100 AEDT	deployed	mysql-1.6.6      	5.7.30       
```


升級到新版本

```
$ helm upgrade mydb stable/mysql --version "1.6.7"
```

Rollback 到之前的版本．

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


移除套件

```
$ helm uninstall mydb
release "mydb" uninstalled

```


Helm vs  Templates
-----------------

Helm 是一個套件管理工具，也恰好包含模板功能（templating），結果不幸的，很多人其實只使用 Helm 的模板功能而已。
如果你只是需要模板功能，那麼 OpenShift Templates 就很夠用了，（甚至其實你也可以用 kdeploy，sed 或 Ansible 等其他工具都可以達到相同目的），尤其在 OpenShift 上，你可以新增自己的 Templates 供其他人從 Web UI 或從指令直接使用，而 Helm 則必須自己另外架設 Helm Registry， OpenShift 目前不提供保存 Helm Charts 並分享的功能。

但如果你想用好好的管理你的應用程序，包含升級及 Rollback，那麼 Helm 絕對是首選．筆者認為 Helm 很適合用在測試環境，用來測試各種版本，因為它可以很輕易地移除應用程序所有相關的資源，測試完後很容易就可以清空專案空間。

以下是 Helm 跟 Template 的比較表

| Feature | Helm  | Template |
|---------|-------|----------|
| Templating | V | V |
| Registry | V ( external Registry) | Ｖ （Built-in） |
| Direct rollbacks and Upgrades | V |  |
| Uninstall Command | V |  |
| Application/package dependencies| V |  |



儘管 Red Hat 大大發表了 [template2helm](https://github.com/redhat-cop/template2helm) 這個工具可以將 Template 轉成 Helm Charts， 但實際測試後，這個工具還是沒有很成熟，轉完後的 Helm Charts 並不能直接用來部署到 OpenShift 平台，你還是得做一定程度的修改。

現階段我認為如果你可以在 OpenShift 內建的 98 個 Template 中找到你想要部署測試的應用程序，那就用 Template。不然的話就用 Helm 吧。

