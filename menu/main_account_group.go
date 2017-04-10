package menu

import (
	"jojopoper/NBi/StressTest/api"
	"sync"
)

// MainAccountGroup 主账户分组定义
type MainAccountGroup struct {
	TestName             string
	GroupID              int
	MainAccounts         []*api.AccountInfo
	LengthOfMainAccounts int
	OperationCnt         int
	TestAccountCnt       int
	runType              api.SendType
	StartBalance         string
	itemList             []*GroupItem
}

// Do 执行创建测试账户操作
func (ths *MainAccountGroup) Do(wt *sync.WaitGroup) {
	if wt != nil {
		defer wt.Done()
	}
	ths.LengthOfMainAccounts = len(ths.MainAccounts)
	ths.itemList = make([]*GroupItem, 0)
	wait := &sync.WaitGroup{}
	// wait.Init()
	for i := 0; i < ths.LengthOfMainAccounts; i++ {
		grpItm := &GroupItem{
			Parent:         ths,
			TestName:       ths.TestName,
			GroupID:        ths.GroupID,
			ItemIndex:      i,
			MainAccounts:   ths.MainAccounts[i],
			OperationCnt:   ths.OperationCnt,
			TestAccountCnt: ths.TestAccountCnt,
			StartBalance:   ths.StartBalance,
		}
		if ths.runType == api.Linear {
			grpItm.DoneCreate(nil)
		} else {
			wait.Add(1)
			go grpItm.DoneCreate(wait)
		}
		ths.itemList = append(ths.itemList, grpItm)
	}
	wait.Wait()
}
