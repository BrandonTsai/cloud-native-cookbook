

從 DevOps 獸進化成 DevSecOps 獸：Kube-bench 跟 Falco 淺談
=========================================

最近瀏覽了一下其他先進的文章，發現大家都很有梗耶，今天標題來致敬一下這位大大:[Docker獸 究極進化 ～～ Kubernetes獸 系列](https://ithelp.ithome.com.tw/users/20129737/ironman/3501)。

回到正題，有些公司可能沒有獨立的 Security Team，因此這責任就變成 DevOps 要去幫忙注意，今天介紹兩個工具幫忙可憐的 DevOps 去檢查系統的安全性設定和監控系統中有風險的行為。


Kube-bench
------

-----


kube-bench 由 Go 語言寫成．它可以幫忙就查 Kubernetes 是否有根據 CIS Kubernetes Benchmark 做好安全性設定。

注意這工具只幫忙檢查 CIS 規範中標示為 Automation 的部分，標示為 Manual 的規範你仍須自己手動檢查，但沒魚蝦也好， 你可以在 [這裡](https://downloads.cisecurity.org/#/) 找到完整的 CIS Kubernetes Benchmark，然後手動檢查 kube-bench 沒有幫忙檢查到的部分。


### 安裝:

請參考 [官方文件](https://github.com/aquasecurity/kube-bench#download-and-install-binaries)，筆者本身是在運算節點安裝 Binary 檔。

```
$ curl -L https://github.com/aquasecurity/kube-bench/releases/download/v0.3.1/kube-bench_0.3.1_linux_amd64.rpm -o kube-bench_0.3.1_linux_amd64.rpm

$ sudo yum install kube-bench_0.3.1_linux_amd64.rpm -y
```


### 用 kube-bench 檢查 OpenShift 3.11 設定。

很不幸的 kube-bench 目前不支援 OpenShift 4.x
底下的測試範例是跑在我們之前的 OpenShift 3.11 系統。

在 Compute nodes:


```bash
$ /usr/local/bin/kube-bench node --version ocp-3.11 
[INFO] 2 Worker Node Security Configuration
[INFO] 7 Kubelet
[INFO] 7.1 Use Security Context Constraints to manage privileged containers as needed
[INFO] 7.2 Ensure anonymous-auth is not disabled
[PASS] 7.3 Verify that the --authorization-mode argument is set to WebHook
[PASS] 7.4 Verify the OpenShift default for the client-ca-file argument
[PASS] 7.5 Verify the OpenShift default setting for the read-only-port argument
[PASS] 7.6 Adjust the streaming-connection-idle-timeout argument
[INFO] 7.7 Verify the OpenShift defaults for the protect-kernel-defaults argument
[PASS] 7.8 Verify the OpenShift default value of true for the make-iptables-util-chains argument
[FAIL] 7.9 Verify that the --keep-terminated-pod-volumes argument is set to false
[INFO] 7.10 Verify the OpenShift defaults for the hostname-override argument
[PASS] 7.11 Set the --event-qps argument to 0
[PASS] 7.12 Verify the OpenShift cert-dir flag for HTTPS traffic
[PASS] 7.13 Verify the OpenShift default of 0 for the cadvisor-port argument
[PASS] 7.14 Verify that the RotateKubeletClientCertificate argument is set to true
[PASS] 7.15 Verify that the RotateKubeletServerCertificate argument is set to true
[INFO] 8 Configuration Files
[PASS] 8.1 Verify the OpenShift default permissions for the kubelet.conf file
[PASS] 8.2 Verify the kubeconfig file ownership of root:root
[FAIL] 8.3 Verify the kubelet service file permissions of 644
[WARN] 8.4 Verify the kubelet service file ownership of root:root
[PASS] 8.5 Verify the OpenShift default permissions for the proxy kubeconfig file
[PASS] 8.6 Verify the proxy kubeconfig file ownership of root:root
[PASS] 8.7 Verify the OpenShift default permissions for the certificate authorities file.
[PASS] 8.8 Verify the client certificate authorities file ownership of root:root
```


在 master nodes：

```bash
$ /usr/local/bin/kube-bench master --version ocp-3.11
[INFO] 1 Securing the OpenShift Master
[INFO] 1 Protecting the API Server
[INFO] 1.1 Maintain default behavior for anonymous access
[PASS] 1.2 Verify that the basic-auth-file method is not enabled
[INFO] 1.3 Insecure Tokens
[PASS] 1.4 Secure communications between the API server and master nodes
[PASS] 1.5 Prevent insecure bindings
[PASS] 1.6 Prevent insecure port access
[PASS] 1.7 Use Secure Ports for API Server Traffic
[INFO] 1.8 Do not expose API server profiling data
[PASS] 1.9 Verify repair-malformed-updates argument for API compatibility
[PASS] 1.10 Verify that the AlwaysAdmit admission controller is disabled
[FAIL] 1.11 Manage the AlwaysPullImages admission controller
...
```




Falco
------

-----


Falco 可以用來監控 Kubernetes runtime security。詳細說明請參閱 [Falco 官方說明](https://falco.org/docs/)。

我們可以透過 OperatorHub 來安裝 Falco Operator 到 OpenShift 上。

1. 安裝 Falco Operator。


![](falco1.PNG)

2. 安裝 Falco DaemonSet 


![](falco2.PNG)


3. 設定客製化的規則 (Rule)

例如，我們想要接收到通知當有人在 Pod 裡面運行  bash 指令時．

把下列 Rule 加到 Falco DaemonSet 的 YAML 檔。

```yaml
          - macro: shell_procs
            condition: proc.name in (shell_binaries)
            output: shell in a container (user=%user.name container_id=%container.id container_name=%container.name shell=%proc.name parent=%proc.pname cmdline=%proc.cmdline)
            priority: WARNING
```

![](flaco3.PNG)


3. 測試客製化的規則是否又用，

我們在別的 console 透過 "oc rsh" 連進某個 Pod, 然後你可以找到對應的 logs 出現在 Falco DaemonSet Pod 內，如果你有用 Splunk 或 EFK Stack 在監控 Pod logs，你可以很輕易的設定對應的 Alerts。


```bash
$ oc logs -f ds/example-falco

...

04:54:40.373142233: Notice A shell was spawned in a container with an attached terminal (user=<NA> k8s.ns=myproject k8s.pod=nginx-pod container=d165df35e317 shell=sh parent=runc cmdline=sh -c TERM="xterm-256color" /bin/sh terminal=34817 container_id=d165df35e317 image=<NA>) k8s.ns=myproject k8s.pod=nginx-pod container=d165df35e317
04:54:40.375631721: Notice A shell was spawned in a container with an attached terminal (user=<NA> k8s.ns=myproject k8s.pod=nginx-pod container=d165df35e317 shell=sh parent=runc cmdline=sh terminal=34817 container_id=d165df35e317 image=<NA>) k8s.ns=myproject k8s.pod=nginx-pod container=d165df35e317
```


結論啦
----------


-----


簡單整理一下, 要進化成 DevSecOps 獸來保護OpenShift 系統及應用程序，你應該

- 加密及備份 etcd 資料。
- 確定系統設定有遵守 CIS Benchmark (可藉由 kube-bench 輔助)。
- 在 CI/CD Pipeline 掃描映像檔的 vulnerabilities，沒有風險才可以 Release 到 OpenShift。
- 透過 container-security-operator 持續監控 Pod 的 vulnerabilities。
- 透過 Falco 持續監控系統的 runtime activity。
