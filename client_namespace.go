package tablecache

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

//查询namespace相关的表是否创建
func (t *TableCache) isNamespaceExist() (bool, error) {
	tables, err := t.client.ListTable()
	if err != nil {
		return false, err
	}

	for _, table := range tables.TableNames {
		if table == t.namespace {
			return true, nil
		}
	}

	return false, nil
}

//创建namespace相关的表，主键为key string
func (t *TableCache) createNamespace() error {
	tableMeta := new(tablestore.TableMeta)
	tableMeta.TableName = t.namespace
	tableMeta.AddPrimaryKeyColumn("key", tablestore.PrimaryKeyType_STRING)

	tableOption := new(tablestore.TableOption)
	tableOption.TimeToAlive = -1
	tableOption.MaxVersion = 1

	//预留吞吐量按小时收费，因此设为0，不预留
	reservedThroughput := new(tablestore.ReservedThroughput)
	reservedThroughput.Readcap = 0
	reservedThroughput.Writecap = 0

	createTableRequest := new(tablestore.CreateTableRequest)
	createTableRequest.TableMeta = tableMeta
	createTableRequest.TableOption = tableOption
	createTableRequest.ReservedThroughput = reservedThroughput

	_, err := t.client.CreateTable(createTableRequest)
	if err != nil {
		return err
	}

	return nil
}

//设置表数据的生命周期，不得低于86400秒（一天）, -1表示永不过期
func (t *TableCache) SetTTL(ttl int) error {
	updateTableReq := new(tablestore.UpdateTableRequest)
	updateTableReq.TableName = t.namespace
	updateTableReq.TableOption = new(tablestore.TableOption)
	updateTableReq.TableOption.TimeToAlive = ttl
	updateTableReq.TableOption.MaxVersion = 1

	_, err := t.client.UpdateTable(updateTableReq)
	if err != nil {
		return err
	}

	return nil
}

//检查namespace相关的表是否创建，如果没有创建，则创建一个默认的
func (t *TableCache) EnsureNamespaceExist() error {
	ok, err := t.isNamespaceExist()
	if err != nil {
		return err
	}

	if !ok {
		err = t.createNamespace()
		if err != nil {
			return err
		}
	}

	return nil
}
