package main

import (
	"fmt"
	"jojopoper/NBi/StressTest/api"
)

func main() {
	fmt.Println(" make order sign ....")
	mo := &api.MakeOrder{}
	src := &api.AccMakeOrder{}
	src.SrcAcc = &api.AccountInfo{}
	src.SrcAcc.Init("GBPFOGXEOF5KQZWB5PC3KTMUDVTDUXL2GB6UBUB7FZ6GDN5EFRUTPJ3S", "SDZ773LWRTF7LEGM4EAZOIKOOWT76HUVXFYFXKY7LOKXTO4DXJPBBL5M")
	// src.SrcAcc.GetInfo(nil)
	src.SetOrderInfo(
		&api.OrderInfo{},
		&api.OrderInfo{
			Code:   "NNN",
			Issuer: "GBPFOGXEOF5KQZWB5PC3KTMUDVTDUXL2GB6UBUB7FZ6GDN5EFRUTPJ3S",
		},
		"10",
		"100")
	b64s := mo.GetSignature("make order", src)

	fmt.Printf("%+v\r\n", b64s)
}
