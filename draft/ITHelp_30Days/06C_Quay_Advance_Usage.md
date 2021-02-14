Quay 進階功能介紹
===================


本篇將介紹4個非常實用的進階功能。

進階功能 1: Tag Expiration
-------------

-----


可以使用 Tag Expiration 的功能將映像檔設置為在選定的日期和時間從 Red Hat Quay 中自動刪除。以下是有關這個功能的一些注意事項：

- 標籤過期後，Quay 將自動刪除該標籤。但不會將映像檔立即將其從 Quay 中刪除。 Time Machine 的值 定義了刪除標籤的實際刪除時間和垃圾回收時間。默認情況下，該值為14天。在此之前，可以將標籤重新指向已過期或已刪除的圖像。
- 這功能是針對 Tag 設置的，而不是針對整個 Repository。

當您只想構建並推送一些映像檔進行短期測試用時，標籤過期可以幫忙減少這些老舊的映像檔佔據太多空間。 標籤過期可以通過兩種方式設置：

- 從 Dockerfile 設置標籤過期。通過 Dockerfile LABEL 命令添加類似 ``quay.expires-after = 20h`` 的標籤將導致標籤在指定時間後自動過期。時間值可以分別為小時，天和周，例如1h，2d，3w。
- 從 UI 來 設置標籤過期。在 Repository 的 “標籤” 頁面上，有一個名為 EXPIRES 的 UI 列，指示標籤何時過期。 用戶可以通過單擊過期時間或單擊右側的“設置”按鈕並選擇“更改過期時間”來進行設置。


進階功能 2: Mirror Repository
-----------------

-----


如果您在不同區域上有多個不同的 Quay 服務器，或者想要將 DockerHub 或 Red Hat Registry 的最新官方映像檔同步到本地 Repository 的話，那麼這個 Mirror Repository 功能很適合您。要注意的是，一旦設定為 Mirror Repository，則無法透過一般的方式來將映像檔上傳到該 Repository。而且通常同步是單向的，你不應該讓兩個 Repository 雙向同步。

要鏡像外部Repository，操作如下：


1. 創建一個機械人帳號以存取 Mirror Repository 的圖像：

2. 選擇新建 Repository 並為其命名。

3. 選擇設置按鈕，然後將 Repository 狀態更改為 MIRROR。

![](images/03_quay/m01.png)


4. 打開新的 Repository，然後在左列中選擇 “鏡像” 按鈕。填寫要鏡像的 Repository 資訊：

- Registry URL：您要鏡像的 Repository 的位置。
- Tags：此為必填。您可以輸入以逗號分隔的單個標籤（1-1,1-2，latest）或用``*``來選擇多個標籤（1-*）。您至少必須明確輸入一個標籤。
- 同步間隔（Sync Interval）：默認為每24小時同步一次。您可以根據小時或天數進行更改。
- 機械人帳號：選擇您之前創建的機械人帳號進行鏡像。
- 用戶名：用於登錄要鏡像的外部Repository用戶名。
- 密碼：與用戶名關聯的密碼。請注意，密碼不能包含需要轉義字符（\）的字符。
- 開始日期：鏡像開始的日期。默認情況下使用的當前日期和時間。
- 驗證TLS：如果要驗證外部 Repository 的TLS，請選中此框。
- HTTP Proxy：如果需要，則設定 HTTP 代理服務器。

![](images/03_quay/m02.png)

5. 從``使用日誌``中檢查同步狀態

![](images/03_quay/m03.png)


進階功能 3: 事件通知（ Notification ）
------------

-----


Quay 可以為 Repository 生命週期中發生事件時，利用不同方式通知使用者或透過 Webhook 呼叫CI/CD服務進行自動化處理。Quay 常見的事件通知有：

- 新映像檔被成功推送到 Repository。
- 檢測到漏洞：在現有映像檔中檢測到新的漏洞。


事件發生時，Quay 可以使用以下方法通知用戶：

- 電子郵件：電子郵件將發送到指定的地址，以描述發生的事件。
- Flowdock 通知：將消息發佈到 Flowdock。
- Hipchat 通知：將消息發佈到 HipChat。
- Slack 通知：將消息發佈到 Slack。
- Webhook POST：將使用事件數據對指定的 URL 進行 HTTP POST 調用。 此方法可用於觸發CI/CD服務（例如 Jenkins ）以實現某些程度的自動化處理。


進階功能 4: Oauth Application
----------------

-----


Red Hat Quay offers programmatic access via an OAuth2 compatible API. It is very useful when you want to do some automation to set up and manage Quay, such as set up ``Mirror-Repository`` or ``Notification`` automatically. Generation of an OAuth access token must normally be done via either an OAuth web approval flow, or via the Generate Token tab in the Application settings with Quay's UI.

Quay 提供 OAuth2 兼容的 API。 當您要進行一些自動化的設置及管理 Quay 時非常有用，例如自動設置 `` Mirror-Repository``或 ``Notification``。 OAuth access token 的生成通常必須通過 UI 在 “應用程序” 設置中的 “Generate Token” 選項完成。步驟如下：

1. 新建一個新的 Oauth 應用程序。

![](images/03_quay/09.png)

2. 新增 Access Token。

![](images/03_quay/10.png)


3. 給予適當的權限。

![](images/03_quay/12.png)

4. 取得 Access Token 內容, 它只能展示一次。

![](images/03_quay/13.png)
![](images/03_quay/14.png)

5. 利用 Access Token 來呼叫 Quay API。 你可以參考 https://docs.quay.io/api/swagger/ 看 API 如何使用。 例如:

```
curl -s -X GET -k https://quay-eu-uat/api/v1/repository/${QUAY_REPO}/image/${IMAGE_ID}/security?vulnerabilities=true -H "Authorization: Bearer ${QUAY_ACCESS_TOKEN}" -H "Content-Type: application/json"
```

結論
-----------

-----


作為私有映像倉庫解決方案，Quay 企業版雖然超級宇宙無敵昂貴，但也提供了相當完整且有用的功能。之後我們將繼續介紹如何在 OpenShift 存取 Quay 上面的映像檔。