

Access public modules via Nexus in a restricted network
---------------------------------------------------------


If you just set up GOPROXY="https//nexus.example.local/repository/go-public-proxy/"
you will find out it always fail because go can not access the default GOSUMDB url `sum.golang.org`

There are two ways to tackle this issue

The first method is turn off GOSUMDB, which is not secure.

```go
# ENV GOSUMDB="off"
```

The better solution is using Both GOPRIVATE and GOPROXY

```
ENV GOPRIVATE=github.com,*.example.local
ENV GOPROXY="https//nexus.example.local/repository/go-public-proxy/"
ENV GONOPROXY=none
```



Publish private go modules with Nexus
-----------------------------------------

