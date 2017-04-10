package operator

import (
	_DB "jojopoper/NBi/StressTest/db"
	"sync"
)

// SaveStaticsToDB 保存统计结果到数据库
func SaveStaticsToDB(wg *sync.WaitGroup, src []*_DB.TStatic) (err error) {
	if wg != nil {
		defer wg.Done()
	}
	opera := _DB.DataBaseInstance.GetOperation(_DB.KeyStatics)
	_, err = opera.Query(_DB.QtAddRecords, src)
	return
}
