# TableCache
use aliyun TableStore as a Cache

## Usage

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