
高富帥版的K8S - OpenShift
====================


Openshift 介紹
---------------------

OpenShift（也稱為OCP）是基於 Kubernetes 的容器應用程序平台，適用於企業應用程序的開發和部署。 但是，為什麼我們需要另一個企業級的 Kubernetes 來運行我們的容器化應用程序呢？ 我列出了以下3個 Openshift 的優點：

1. OpenShift 是針對需要更可靠和穩定解決方案的企業組織的產品 (Product)，而 Kubernetes 比較偏向任何人都適用的 Project。 OpenShift 允許用戶通過訂閱方式來提供付費支援以即時解決客戶的問題，而 Kubernetes 則需要仰賴社群的回答。

2. OpenShift 內建了一些有用的功能，因為它們是維護企業級平台所必需的。而 kKubernetes 需要你自己手動安裝及設置這些附加功能。當然，如果你是使用 [kubematic](https://www.kubermatic.com/products/kubermatic/) 提供的解決方案，那麼差異就沒那麼大。

3. OpenShift 的安全策略比 Kubernetes 嚴格。 人們在 Kubernetes 上運行簡單的應用程序很容易，但是 OpenShift的安全策略限制了他們這樣做。 需要一定級別的權限才能維護最低安全級別，因此用戶別無選擇，必須學習這些規定以部署應用程序在 OpenShift。


安裝 單機版 OpenShift 4.x 在 Mac 上
----------------------------------

我們可以透過 Red Hat CodeReady Containers 來運行 OpenShift 4.x 在Linux，Mac， 甚至是 Windows上，方便開發者開發及測試。

CodeReady Containers 需要以下最低系統資源才能運行 Red Hat OpenShift：

- 4 virtual CPUs (vCPUs)
- 8 GB 記憶體
- 35 GB 儲存空間

您還將需要主機作業系統支援的虛擬機器監視器（Hypervisor）。
CodeReady Containers 以下 Hypervisors:

- Linux 的 libvirt
- macOS 的 HyperKit
- Windows 的 Hyper-V


1) **安裝 hyperkit**

```
$ brew install hyperkit
$ brew link --overwrite hyperkit
```

2) **下載 CodeReady Containers**

從 [Red Hat CodeReady Containers product page](https://developers.redhat.com/products/codeready-containers/overview) 下載 CodeReady Containers 壓縮檔。

3) **解壓縮後然後將 `crc` binary 檔放到 `$PATH` 下**

```
$ cp crc /usr/local/bin/
```

4) **設置主機環境**

一旦安裝了 CodeReady Containers，在啟動 OpenShift 之前，必須運行 `crc setup` 來設置主機環境。


5) **啟動 OpenShift 4.x cluster**


使用`crc setup`命令設置主機環境後，可以使用`crc start`命令啟動OpenShift。

出現提示時，請輸入user pull secret。 此密碼可以從 Red Hat CodeReady Containers 產品頁面的 Pull Secret部分下複製。 注意，您需要紅帽帳戶才能訪問此 pull secret。

啟動過程完成後，它將顯示可用於登錄OpenShift的用戶憑據。

![](https://lh3.googleusercontent.com/Cj8dsFgSD20czZbS0mQkgEWw5kYM4Uhx0th_Ox_NyUwlo_YMTl040P-4U4ZEkGPVxxbjEMj4VP9cPK8dNGgq2AC357PW-pgm2-Up2IUHjTkNu4GdiMJUhXszB7n2RKQZl124w3Lu8nAHR5D3QaxP0DIVyzbAenQFMwbHLWUENdwArjOq4elkSOvmtJfVuCmwHMvdO0-lxPtxdgIbt-kBNSp3cPKGcUu6h8rO6u22_zxQj5DvOQ6-ZfwS6N4uTSWxzO7Xq585zTrljyMuCcOAWrbjQBp3GvAF5z-PrXxhVbCdtfY_LsqeLXdC3ebqdof5cI3wXaEQpdyaeqw2Ln2AurL0mMQ3EakagWRO3kDKbD_pun1Pii5TBVKl6qPvmCA8z3Vjk3NHM3jB0j8nH86VaYpNxccKxw0BcSUQc84z_Yjgd8lH1YppqYYM4jKAaJlpSuJRaxDvk5Agr5_6bsW9A-5bcx0ZvabuHxIqgKRwVBwMuLjwHUhEqJTtKKXalasgTN3-Kz14TBAou7LRt3987J575mq1L6CqJV0YHkXalaz1bOIX7-vRFSNHGwzgGXC0xmwqrj8qrJXTk9Q1XK4In5n2YrXySEXojAXC0WLDkEphEEAe8wXKb6Pu2XfL4gpAML0V8Bhd7wUc9Bwk8P4rJNri4yfyfnYSo1YTzhgVPNF-5dBkvEkpcJ_bU1YjkQ=w2400-h238-no?authuser=0)


訪問 OpenShift 
--------------

要訪問 OpenShift 前，請首先按照 `crc oc-env` 說明進行設置。

```
$ Brandons-MacBook-Pro:crc-macos-1.16.0-amd64 brandon$ crc oc-env
export PATH="/Users/brandon/.crc/bin/oc:$PATH"
# Run this command to configure your shell:
# eval $(crc oc-env)
```

已開發者身份登入 OpenShift

```
$ oc login -u developer -p developer https://api.crc.testing:6443
Login successful.

You don't have any projects. You can try to create a new project, by running

    oc new-project <projectname>
```


建立第一個 Ｐroject

```
$ oc new-project myproject
```

> `Project` 是用來對 `Namespace` 附加一些註釋(Annotation)，例如 對 Project 的描述等， 一個 Project 會對應到一個 Namespace，他們就像連體嬰，你刪除一個 Project，對應的 Namespace 也會刪除。 而管理員也必須向用戶授予對 Project 的訪問權限，或者允許使用者自己創建 Project，則使用者會自動有權訪問其自己的 Project。與 Kubernetes 不同，普通用戶無法一開始就在 Default Project/Namespace 中創建 Pod。 必須先創建第一個新的 Project，然後在新的 Project 內創建 Pod。

列出所有 Ｐroject

```
$ oc projects
You have one project on this server: "myproject".

Using project "myproject" on server "https://api.crc.testing:6443".
```

在 Kebernetes 中，您需要安裝插件 `kubens` 才能在 Namespace 之間切換。 但OpenShift已經內置此功能，您只需要運行 `oc project <my-project>`。


切換到剛剛新增的 Project。

```
$ oc project myproject
Already on project "myproject" on server "https://api.crc.testing:6443".
```


訪問使用者介面
-----------

運行`crc console`並使用這些啟動OpenShift時給的密碼來登入使用者介面。


1) **已開發者身份登入**

選擇 "htpasswd_provider"

![](https://lh3.googleusercontent.com/QT2Ubxy_cJVaisklaAcnqavWNQ7e7BWIw4uizE5QMjwQkIV_3Ny82U5gE27hDNDy3UKJixkbTgNgBcmxgtLHbhJ2g1rtYgN1RysmJamORA8T5ZRhGtNwF_JW0dfz3bzZjjblpKrMkYjdcmvXt90qu_nRDhd5ItqWiqr8jWWoTO02sOWf2Wrm7t-rKHT2WdWJE448-CJ7hax8GkmwMcTB8iKPDkfRxVWzycXRmEHTI1veBudfMTWpnULYtmO-1h3Ofk41xX7Ojuu0J63IKue3bIOUNaTqRliZK3zHqBnS6u28misoNfGYk49ozlLTXJx4yk-7MDwBHnBLSBrON8txVW4QstW8WGK4AF-O-mh7sNdVeiWDXxHJNPCQGXQ4LIbiU2RfMWaIZB0LNYx9DqQ8Djfpdcn7dQQbJmLe0aQuFfcPV_7HEMAJdj6IEQpTJGb6v_nMXgExnei6nd6lxZm7_TWeVWv-BGimABa9nwWCf-6aBe372YCxOpXlrNNsExerJnb5MRInT-o3dm3Iy3_pF6UY0x8PPiLMcNLwhMV-i00z2OkDPKzqv1ngR9I-zua6D8LyjsswOgnjqPWF4_0DuEn3zBLc8hjQzTTZen-p3eNLvX5dvZas2UfNYVkHsN_VeGJS1RgpzwigrPbDN17ERgRg7kJiXsux9Y9Chb5fKQHMTFsToH1KIewMyhEbhQ=w2208-h998-no?authuser=0)

輸入密碼

![](https://lh3.googleusercontent.com/vQkCt_eXPfKrQi2qKCUou8jAodjfkCGyhPGQCDZpM1kweuwwbiGpYQewcF_fO_YT6SbH1ZAFBHqGvnc61zSSOj2k1LbcIp9Q6FKAV3JIfBsg5aGsSonJSrFoKirjYusSYMC0CTrezpZqXKXjM9VJcad0H_HyWosS-0neuxYH0UtHsJ-j8RGW3YGW_XGQO4OUJoXYaXM8hjOLShUsCnfhL6KSdMFiXpfqHQ2yUzbTg6xHOC4hK6Upj8mKpltK6dAETay8jkwamF3cWUgNg4d1qdeIMT6WHkRS-XxmNWLc2_CMvPB2mEvrCugNrcK0LWPSPiTEGea-GvlBR0CgJZdUN23KewBIcME-5xFxC6CT_VGgIb2dH4SKf09bYNQ9qkuNO69Nd3HBMXes2ep5cyIy_Pk-uhArGb_TQHnob2DmdoC-RQv6C45KXDaxcf7Q8dczlkJ3a-Vb90sGU0kqveXf0EEt-t5mnaTNj37a-v63qWl95obTmLlATmSVxpg6bS8CmjUca-Y-dRXaPWy4bhDIbBa3pMLjzCQhSHhrAJkEcCYu_DWcdWjoWA6hxM3nOyqu94ek25VoWP88w1_WSd7n5c_1Mq2eR9QH_L3Pi5IzBFvSyBBo2ZqPMtMa35-ZjvJiCRZtgtHA-EyIbkX26DcOd-CXjCJJtdbzI2pDgYQRufh3Z7S5ji_FYpSvmFGOew=w2146-h1104-no?authuser=0)


成功登入開發者的使用者介面啦

![](https://lh3.googleusercontent.com/somhA2CDnW9xAk_U_f5EtGRBoeAWZu8mimH6pMzQfCFzpBRAJLYErWEhd0BfFfYNgQMGOauik1Hd7F6QJPuSjbGtET1nzUd9OCCBQfKJnYcN22wS1n4fS_qWN20tFcoDugNZFuG6JaNqNcPN_gGyGAwlxb7x6sjbyqGBlTf4Ra5XmypJEsILodt4aV4kLjGffH1GdmFQkkc2ewGzyASxcHbqxXdqyw6S7pmOo5_2lOw0fEdfLoMKerW5l1JTRvWeOZFto7JCWjug_mnCRS-jYyO9R8w0X_mPmuhqx7JqcpA5OWzomK2bUod0vTYFXnmVTngRCcWo3Rsf1voRmwd6wixRtZN70JnQ5h8fFd87utQXaqqc89W3eiZciORHp_wGvH2DXKsTj8h3e2MWceDTqQP6PClDzH5m2oGoFw3TwwHT3-_aoNirFnE0srJdo94FG9OJBxSfjx7APWl5_oHiJw5c96ZqBJiBZkmMcC_27iN8V4YAT9rYi2ev2uhpvYWVvcBC5KzkHEzPvtju2MvpqF8VvEwqu7UXQ8wE-2hmv6g9VAsdsPbstpl1NQd6joQLivElFurC_tvEFSAjY-LzV-mNDqFynuZjoV7Xq6b2-qGo1JqKSilOVhtsCXXolZfOPaWEbsYtbE0Mz0J0w1kF5edSwKdcc6zl0A8OnZYc2pPVLoOrbei9cLr28u1ROw=w2858-h1406-no?authuser=0)



2) **以系統管理員身份登入**

選擇 "kube:admin"

![](https://lh3.googleusercontent.com/lI-VTKyem3h253ZFPvqh7kdQCTjO3lxPwT4FSoj1FgsakVWpfvmI8taBv-zA1GB76hJLWWH3x9mRw__ZrZT9j0uLGlIH0Rs8rzfmZR3XHfk87SSFZspERhYYIwXceo3W2rqcFmvmuOdW45tl5cyzH7ta2A61ELpVv5sUraGcY2G46g5Bq_dfw1ZNgsDBcpSExjqS_SPx08xKojlSzM-uMHUr_-XxA9BNUtX24InW8ugCEbYJKdumYpMoBC6arzZ9rpvtu_CLK9yLVPO-hBd74C5WfsT9z2ljYPVbZ3XMYkfHNs4Kd6YAt4BfH15WzGRRDBHNdnma-pbw2Aycov49CE01hbC-6vyaLaadjt9gcVs_vQ6Kc9W44bODi3bLz2yO_YM1E-EPMFAywD-pJhCSo1R8W7CZdkKy5S0TBpDYVIWEmdvT0Diufyvb66bH3Orz-B2hhC8N6oGLIgb7gGAoJH9h-llUspGTi7V0yfOEONNDxoVoW_szzTbaJ3JPkDPAzJ4bF-k-rR7esUbCx-FnL8UDqXUKr2sM_bYv1O30o_fAwVwB55WCVRDV-Nb8jS_zg48awaOGPKvwpPw4NTbsUlSvGBC0ZdJGbJbj3uXyXmzkhJ_hNyqMBnNVgLIe1JB5x1Qc-b9XKQsYgRBcZAiDzOqBxPaaimeN2Idx5j4VdvpgD732OTNmXOYQKjYXMQ=w2208-h998-no?authuser=0)

輸入管理者密碼

![](https://lh3.googleusercontent.com/REMQYQhvgkfH_TmDNmlueaOz0TKFz1VTVpwC_eAKJYZuod92R5GJbHFBcekII1gLClzZ9lgA7SAZkQZ5U6wFBhyIHLpsAMIs0wzdP98tb0iGTCz_iPEgYo58mhsSLtzixLI57hq8prCE3kk5rznWztBzHVL2K-q_5qP7Z6BFQae2rqCDMEgrQ_awd-S81EmmmHXqVYX5-QeizHKP6_4fOqJJOT7Kkc9Ce3sh8SWxItwJE2TE7KZEybjlZ86c9ilvnPs07hLvCEg5REi9jSfx2Sf8mN6t35xYseNzqfi6h_ShF3IJbGC4M8uUvDHNTkLIDjQL00fYscW32u4WroKl5bM9q_gH-kffqIi0TknRvMtZ2Je13nNuc5VfgH7BzEiQJ8jt22N7EbJHy3Cim9cyjQMI2vwgonr64BMWPy0OPfE-ilHZ00Dna1iWiY_Tt4Q9ZCS1di7rOhARvU3DG0WspldJHSpuJU_vpTuJyKehAkhAqec9u4FqHAc32GejA5IgECvvVf0XfWUVKLZPrIyCTYGXZbwGYdjQDBT1ly3vZRIHDrtaR71lizQuuEz4e_mll2EJWPpDoSGHZ0aoGcB7z_zYlnal7e4uOe2E_D9HogWOXGR9-WfJnWZeFj-lHzbIClECdOeNFPDIwwpmEnhHwrIQat8JBIysiFJmzJDFatyEwO1xyQ-X_FBU_yYqmA=w711-h367-no?authuser=0)


確定有登入管理者頁面

![](https://lh3.googleusercontent.com/RArp4lUsRydUzHyHFID3H-UuOd057k0CdWbHxPK8DvWpQXIgFiZ69FY5m13vRCxI0xcIS6MDQ_JfGDNxbb59Y--KWfkLUC_SDNqunAdNQE95jg89XIR6hplsmCRAG6YZyewugF34Ysb9l3cmntbHxjuXcX4GJO2Y9916DTN8-LXLJ81CpuFJZSOqRLYMLiB9F0RZZHYbdkX4j3KQX4BNKoqJFJ5E8_9ZC85ssV1gaWPapMRbKpVcG3HGoVIeBc2Rd6IMqG1pX6vE_2x_Tykw_4y1NCQoFW3ucD_89GrcxC_YMm7DqonAULLPFE-BFL4RKh-zAf3FGMEd_FQKRJvSkisBjHgGdUbA-m0lfwNiJpu3Ig3TL19PX1oIryEn9zuAkF_PWX7OKQXUrtbNHmR1WjveG8y9B4GpmPfcd_ktUnPIEOkY71lrIULjWqA7pi9lht4SARdiuGgNAvPvUunqkZLI7WaYWLCNnX3vGJFJRjWctV1U9usi-2xhejfIEz_PvMWs53K9Hb-KtmZ6dUFC_pSNP3TRB7M_okPMpSairyvII5gIGEm-pDdDsWH60Fsbzj12b42NfVg2TBVr9Hs3L0iaTwbNBsMWRG_3u3G-tnVRc2HefPRpDfJAwe98WDdQxZtDaf-XEXMu-lXOiM1mCVU6ekNR6sSJVE2MQ_LnjHTKgMOjg_FhQGCeEUOnnQ=w2872-h1558-no?authuser=0)


停止 OpenShift
-------------


```
$ crc stop
```

確定 OpenShift 狀態

```
$ crc status
CRC VM:          Stopped
OpenShift:       Stopped
Disk Usage:      0B of 0B (Inside the CRC VM)
Cache Usage:     13.05GB
Cache Directory: /Users/brandon/.crc/cache
```


結論
-----1

由於 Kubernetes 是構成 OpenShift 的基礎，因此可以在兩者之間找到許多共同點。 儘管 Kubernetes 是開放且快速發展的平台，但 OpenShift 關注不同企業的需求，並根據這些企業級的需求應運而生，讓更多企業願意接納 Kubernetes 對應的平台跟服務。 Kubernetes 跟 OpenShift 可以說是相輔相成。