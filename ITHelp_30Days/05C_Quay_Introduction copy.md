

前言
------------

-----


儘管我們可以在運行容器之前利用 Dockerfile 構建 Image，但是這樣會耗費很多時間，也無法確定建構好的 Image 是否跟測試環境一致。最普遍的做法將這些構建好的 Image 保留在一個「映像倉庫（Image Registry）」上，然後在運行容器之前再提取這些 Image 。映像倉庫的解決方案很多，其中最著名的是 Dockerhub。 您可以在 Dockerhub 上找到許多公開的 Image。 但是，如果要在私有開發環境或數據中心上構建本地的 Image Registry，那筆者建議您考慮使用 Red Hat 的解決方案 - Quay。

Quay 提供了相當完整的功能讓使用者在任何架構上都可以安全地存儲和部署容器。它可以作為獨立的服務使用，也可以部署在 OpenShift 上使用。


Compare Harbor and Quay
----------------------

-----


還有另一個類似且受歡迎的解決方案 - [Harbor]（https://github.com/goharbor/harbor）。 它是一個開放源代碼的映像倉庫（Image Registry）解決方案，由 Cloud Native Computing Foundation（CNCF）託管。 以下是 Harbor 和 Quay 之間的比較。


| Product | Quay | Harbor |
|---------|------|--------|
| Language                      | Python | Golang |
| Type                          | Public (quay.io), private | private |
| Authentication                | LDAP, OIDC, Google, Github  | LDAP, OIDC, DB |
| Robot Account                 | O | O |
| Permission Management         | O | O |
| Security Scan                 | O | O |
| Image Cleaning               | O | O |
| Helm Application Management   | O | O |
| Image Mirroring               | O | O |
| Notification                  | Webhook, Email, Slack | Webhook |


您可以看到 Quay 和 Harbor 之間的功能都非常完整且沒有太大區別。如果您未來打算在 OpenShift 上運行容器，我建議您使用Quay，因為它受 Red Hat 支持並與 OpenShift 高度相容。


Quay 安裝
-----------------------------

-----


要安裝開源版本，請參閱https://github.com/quay/quay/blob/master/docs/getting_started.md。 要在您的本地環境上安裝Red Hat Quay，請參考[https://access.redhat.com/documentation/en-us/red_hat_quay/3.3/](https://access.redhat.com/documentation/en-us /red_hat_quay/3.3/）。 本文將重點介紹 Red Hat Quay 3.2.1版的功能。


Quay 基本功能介紹
-------------------

-----


在Quay中，您可以為不同的業務部門或不同的用途創建不同的組織（Organization）。 每個組織包含獨立的團隊（Team），存儲庫（Registry），機器人賬號（Robot Account）和 Oauth Access Token。

1. 創建 Organization 和 Registry。

![](images/03_quay/02.png)


2. 利用 ``Skopeo`` 從 Red Hat Registry 複製 Images 到 Quay。

```
$ sudo docker login quay-eu-uat
$ sudo skopeo copy --src-tls-verify=false --dest-tls-verify=false docker://registry.redhat.io/rhscl/nginx-116-rhel7 docker://quay-eu-uat/application-images/test:1
```

3. Quay 會透過 Clair 針對上傳的 Image 做風險掃描，並在ＵＩ顯示結果。從下圖可以看到 Image 已通過安全掃描。

![](images/03_quay/03.png)



機器人賬號 （ Robot Account ） 介紹
--------------------

-----


在許多情況下，需要讓其他服務（例如CI / CD管道和Openshift）有權限存取特定 Registry 上的映像檔。 為了支持這種情況，Quay 允許使用用戶或組織擁有多個機器人賬號（ Robot Account ）來訪問 Registry。


1. 建立一個 robot account。

![](images/03_quay/04.png)

2. 給他一個合適的名稱。

![](images/03_quay/05.png)

3. 給予存取 Registry的權限。

![](images/03_quay/06.png)


4. 取得此 robot account 的密碼。

![](images/03_quay/07.png)

![](images/03_quay/08.png)

5. 從別的伺服器透過該機器人賬號來登入 Quay。

```
$ docker login -u="application-images+jenkins" -p="6BC26ZL0CUZQTJKL1SWKZIO9ZD58TDLS8O6VONE4VVNF9M1ZQGGMCVBXORNC0BNG" quay-eu-uat
```

6. 確定你可以透過該機器人賬號來存取 images。

```
$ docker pull quay-eu-uat/application-images/test:1
```

小結啦
------

-----



本文比較了兩個「映像倉庫（Image Registry）」的解決方案及介紹 Red Hat Quay 的最基本用法。 下一篇，我將繼續介紹 Quay 的進階用法，讓我們一起玩轉它的神秘面紗。


