# TableCache
use aliyun TableStore as a Cache or Session

## Usage
### import 

```
go get "github.com/diemus/tablecache"
```

### use as cache

```go
import (
	"fmt"
	"github.com/diemus/tablecache"
)

func main(){
	client := tablecache.NewTableCache(
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
}

```

### Use as gorilla session store
for more gorilla session usage example, you can go [https://github.com/gorilla/sessions](https://github.com/gorilla/sessions)

```go

	var store = sessionutils.NewTableCacheStore(client,
		[]byte("authkey"),
		[]byte("enckey1234567890"), //optional
	)

```

### Use as gin session middlware
for more gin session usage example, you can go [https://github.com/gin-contrib/sessions](https://github.com/gin-contrib/sessions)

```go

	import (
		"fmt"
		"github.com/diemus/tablecache"
		"github.com/diemus/tablecache/sessionutils"
		"github.com/gin-contrib/sessions"
		"github.com/gin-gonic/gin"
	)
	
	func main() {
		router := gin.Default()
	
		//create client
		client := tablecache.NewTableCache(
			"https://xxxx",
			"instanceName",
			"namespace",
			"accessId",
			"accessSecret",
		)
	
		//create gorilla session store
		store := sessionutils.NewTableCacheStoreForGin(client,
			[]byte("authkey"),
			[]byte("enckey1234567890")) //optional
	
		//add gin middleware
		router.Use(sessions.Sessions("sessionId", store))
	
		//user session in request
		router.Any("/*path", func(c *gin.Context) {
			session := sessions.Default(c)
	
			fmt.Println(session.Get("aa1"))
	
			session.Set("abc", "world")
			session.Save()
			fmt.Println(session.Get("aa1"))
			c.String(200, "ok")
		})
	
		router.Run()
	}
	
```