package main

import (
	"fmt"
	"time"

	"github.com/iost-official/go-iost/account"
	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-iost/core/tx"
	"github.com/iost-official/go-iost/crypto"
	"github.com/iost-official/go-iost/ilog"
	"github.com/iost-official/go-iost/vm"
	v8 "github.com/iost-official/go-iost/vm/v8vm"
)

var testID = []string{
	"IOST4wQ6HPkSrtDRYi2TGkyMJZAB3em26fx79qR3UJC7fcxpL87wTn", "EhNiaU4DzUmjCrvynV3gaUeuj2VjB1v2DCmbGD5U2nSE",
	"IOST558jUpQvBD7F3WTKpnDAWg6HwKrfFiZ7AqhPFf4QSrmjdmBGeY", "8dJ9YKovJ5E7hkebAQaScaG1BA8snRUHPUbVcArcTVq6",
	"IOST7ZNDWeh8pHytAZdpgvp7vMpjZSSe5mUUKxDm6AXPsbdgDMAYhs", "7CnwT7BXkEFAVx6QZqC7gkDhQwbvC3d2CkMZvXHZdDMN",
	"IOST54ETA3q5eC8jAoEpfRAToiuc6Fjs5oqEahzghWkmEYs9S9CMKd", "Htarc5Sp4trjqY4WrTLtZ85CF6qx87v7CRwtV4RRGnbF",
	"IOST7GmPn8xC1RESMRS6a62RmBcCdwKbKvk2ZpxZpcXdUPoJdapnnh", "Bk8bAyG4VLBcrsoRErPuQGhwCy4C1VxfKE4jjX9oLhv",
	"IOST7ZGQL4k85v4wAxWngmow7JcX4QFQ4mtLNjgvRrEnEuCkGSBEHN", "546aCDG9igGgZqVZeybajaorP5ZeF9ghLu2oLncXk3d6",
	"IOST59uMX3Y4ab5dcq8p1wMXodANccJcj2efbcDThtkw6egvcni5L9", "DXNYRwG7dRFkbWzMNEbKfBhuS8Yn51x9J6XuTdNwB11M",
	"IOST8mFxe4kq9XciDtURFZJ8E76B8UssBgRVFA5gZN9HF5kLUVZ1BB", "AG8uECmAwFis8uxTdWqcgGD9tGDwoP6CxqhkhpuCdSeC",
	"IOST7uqa5UQPVT9ongTv6KmqDYKdVYSx4DV2reui4nuC5mm5vBt3D9", "GJt5WSSv5WZi1axd3qkb1vLEfxCEgKGupcXf45b5tERU",
	"IOST6wYBsLZmzJv22FmHAYBBsTzmV1p1mtHQwkTK9AjCH9Tg5Le4i4", "7U3uwEeGc2TF3Xde2oT66eTx1Uw15qRqYuTnMd3NNjai",
}

func MakeTxWithAuth(act tx.Action, ac *account.KeyPair) (*tx.Tx, error) {
	trx := tx.NewTx([]*tx.Action{&act}, nil, int64(100000), int64(1), int64(10000000))
	trx, err := tx.SignTx(trx, ac)
	if err != nil {
		return nil, err
	}
	return trx, nil
}

func main() { //629123ns(local) vs 1048236(server)
	js := vm.NewJSTester(nil)
	defer js.Clear()
	f, err := vm.ReadFile("go-iost/test/performance/transfer.js")
	if err != nil {
		ilog.Info(err)
	}
	js.SetJS(string(f))
	js.SetAPI("transfer", "string", "string", "number")
	js.DoSet()

	js.Visitor().SetBalance(testID[0], 100000000)

	act2 := tx.NewAction(js.CName(), "transfer", fmt.Sprintf(`["%v","%v",%v]`, testID[0], testID[2], 100))

	ac, err := account.NewKeyPair(common.Base58Decode(testID[1]), crypto.Secp256k1)
	if err != nil {
		panic(err)
	}

	trx2, err := MakeTxWithAuth(act2, ac)
	if err != nil {
		ilog.Info(err)
	}

	//b.ResetTimer()
	for i := 0; i < 400; i++ {
		t1 := time.Now()
		r, err := js.Engine().Exec(trx2, time.Second)
		T := time.Since(t1)
		ilog.Info(i, T, v8.T1.Nanoseconds(), v8.T2.Nanoseconds(), float64(v8.T1.Nanoseconds())/float64(T.Nanoseconds()), float64(v8.T2.Nanoseconds())/float64(T.Nanoseconds()))

		if r.Status.Code != 0 || err != nil {
			ilog.Fatal(r.Status.Message, err)
		}
	}
	//b.StopTimer()
}
