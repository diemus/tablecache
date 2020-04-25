package tablecache

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

func NewTableCache(endPoint, instanceName, namespace, accessKeyId, accessKeySecret string, options ...tablestore.ClientOption) *TableCache {
	client := tablestore.NewClient(endPoint, instanceName, accessKeyId, accessKeySecret, options...)
	return &TableCache{client: client, namespace: namespace}
}

type TableCache struct {
	client    *tablestore.TableStoreClient
	namespace string
}

func (t *TableCache) Get(key string) (string, error) {
	pk := new(tablestore.PrimaryKey)
	pk.AddPrimaryKeyColumn("key", key)

	criteria := new(tablestore.SingleRowQueryCriteria)
	criteria.TableName = t.namespace
	criteria.PrimaryKey = pk
	//MaxVersion必须有，否则汇报Invalid Input
	criteria.MaxVersion = 1

	getRowRequest := new(tablestore.GetRowRequest)
	getRowRequest.SingleRowQueryCriteria = criteria
	getResp, err := t.client.GetRow(getRowRequest)
	if err != nil {
		return "", err
	}

	//TableStore GetRow为空不会报错，但是可以通过主键为空检测出来
	if len(getResp.PrimaryKey.PrimaryKeys) == 0 {
		return "", KeyNotFoundError
	}

	//校验结果是否为空
	if len(getResp.Columns) == 0 {
		return "", NoFieldNamedValueError
	}

	//校验value字段是否为字符串
	value, ok := getResp.Columns[0].Value.(string)
	if !ok {
		return "", ValueNotStringError
	}

	return value, nil
}

func (t *TableCache) Set(key string, value string) error {
	pk := new(tablestore.PrimaryKey)
	pk.AddPrimaryKeyColumn("key", key)

	putRowChange := new(tablestore.PutRowChange)
	putRowChange.TableName = t.namespace
	putRowChange.PrimaryKey = pk
	putRowChange.AddColumn("value", value)
	putRowChange.SetCondition(tablestore.RowExistenceExpectation_IGNORE)

	putRowRequest := new(tablestore.PutRowRequest)
	putRowRequest.PutRowChange = putRowChange
	_, err := t.client.PutRow(putRowRequest)
	if err != nil {
		return err
	}

	return nil
}

func (t *TableCache) Del(key string) error {
	pk := new(tablestore.PrimaryKey)
	pk.AddPrimaryKeyColumn("key", key)

	deleteRowChange := new(tablestore.DeleteRowChange)
	deleteRowChange.TableName = t.namespace
	deleteRowChange.PrimaryKey = pk
	deleteRowChange.SetCondition(tablestore.RowExistenceExpectation_IGNORE)

	deleteRowReq := new(tablestore.DeleteRowRequest)
	deleteRowReq.DeleteRowChange = deleteRowChange
	_, err := t.client.DeleteRow(deleteRowReq)
	if err != nil {
		return err
	}

	return nil
}
