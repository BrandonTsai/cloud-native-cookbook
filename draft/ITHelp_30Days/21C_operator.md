



蝦米係 Operator？
------

對於 Kubernetes 早期的傳統觀點是，它非常擅長管理無狀態（Stateless）應用程序 。但是對於有狀態（Stateful）的應用程序而言，例如數據庫，情況並非如此：這些應用程序需要更多的手動操作。在 Helm hub 裡面，很多資料庫的 charts 其實是單節點的，因為分布式且高可靠性的資料庫設定起來會較為麻煩，而且此類服務在添加或刪除節點可能需要準備和/或後置步驟-例如，更改其內部配置，與其他節點進行通信，與DNS等外部系統進行溝通等，從歷史上看，這通常需要手動設定，從而增加了營運工程師的負擔並增加了出錯的可能性。

Operator 就是為了解決這一個問題而存在的一個工具。 Operators 可以：

- 提供安裝和升級的可重複性。
- 持續檢查每個系統組件的運行狀況。
- OpenShift 組件和ISV內容的無線更新（OTA）。
- 每個服務提供商可應用該領域知識，用軟體的方式管理複雜的應用並將其傳播給所有用戶。



蝦米係 Operator Framework ?
------------------
Operator Framework 是一系列開源工具，旨在以更有效，自動化和可擴展的方式管理 Operators。 這不僅僅是寫代碼； 測試，交付和更新 Operator 也同樣重要。 Operator Framework 組件包含用於解決以下問題的開源工具：

### Operator SDK
Operator SDK 可根據其專業知識協助 Operator 開發者構建，測試和打包自己的 Operator，而無需了解 Kubernetes API 的複雜性。

### Operator Lifecycle Manager（OLM）
Operator Lifecycle Manager 控制 Operators的安裝，升級和訪問控制權限（RBAC）。 預設內建在 OpenShift 4.5。

### Operator Registry

Operator Registry 用於儲存在群集中創建的 ClusterServiceVersions（CSV）和自定義資源定義（CRD），並存儲有關 Package 和 Channel的metadata。 它在Kubernetes或OpenShift 平台中運行，以將該Operator 的數據提供給 Operator Lifecycle Manager （OLM）。

### OperatorHub
OperatorHub 是一個Web控制台，系統管理員可以通過該控制台發現並選擇要在其平台上安裝的 Operators。同樣預設內建在 OpenShift。

### Operator Metering
Operator Metering 收集系統上有關 Operators 運行上的相關數據（Metrics），以進行管理和彙總。

這些工具被設計為可組合的，因此您可以使用任何對您有用的工具。


透過 OperatorHub 來安裝第一個 Operators
---------------------


以 kubeadmin 登入後, 你就口以找到 OperatorHub ，如下圖

![](images/Operator/1.png)

搜尋 『 Web Terminal Operator 』

![](images/Operator/2.png)

使用預設值來安裝 Web Terminal Operator

![](images/Operator/3.png)

確定 Operator 狀態

![](images/Operator/4.png)

該 Operator 可以讓使用者直接在網頁透過 Terminal 連進 Pod 裡面，讓偶們確認是不是有用。

![](images/Operator/5.png)



結論
--------

用 Operator 來部署和管理一些需要 Production level 及 High Availaity 的服務是非常有用的，
接下來幾天我將會繼續探索及介紹其他有用的 Operator 。
