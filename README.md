# TableCache
use aliyun TableStore as a Cache or Session

## Usage
### import 

```
go get "github.com/diemus/tablecache"
```

### use as cache

```go
	client := NewTableCache(
		"https://xxxx",
		"instanceName",
		"namespace",
		"accessId",
		"accessSecret",
	)

	err := client.EnsureNamespaceExist()
	if err != nil {
		fmt.Println(err)
	}

	err = client.Set("key1", "test")
	if err != nil {
		fmt.Println(err)
	}

	v, err := client.Get("key1")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(v)

```

### use as gin session middlware
```go
	router := gin.Default()

	client := NewTableCache(
		"https://xxxx",
		"instanceName",
		"namespace",
		"accessId",
		"accessSecret",
	)
	store := sessionutils.NewTableCacheStoreForGin(client,
		[]byte("authkey"),  //auth key is requered
		[]byte("enckey1234567890"),)  // optional
	router.Use(sessions.Sessions("sessionId", store))

	router.Any("/*path", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("abc", "world")
		session.Save()
		c.String(200, msg+"\n")
	})
```
for more session usage example, you can go [https://github.com/gin-contrib/sessions](https://github.com/gin-contrib/sessions)