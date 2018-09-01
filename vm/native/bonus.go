package native

import (
	"github.com/iost-official/Go-IOS-Protocol/core/contract"
	"github.com/iost-official/Go-IOS-Protocol/vm/host"
)

var bonusABIs map[string]*abi

func init() {
	coinABIs = make(map[string]*abi)
	register(&coinABIs, createCoin)
}

var (
	claimBonus = &abi{
		name: "ClaimBonus",
		args: []string{"string", "number"},
		do: func(h *host.Host, args ...interface{}) (rtn []interface{}, cost *contract.Cost, err error) {
			cost = contract.Cost0()
			acc := args[0].(string)
			amount := args[1].(int64)

			totalServi, cost0 := h.TotalServi()
			cost.AddAssign(cost0)

			cost0, err = h.ConsumeServi(acc, amount)
			cost.AddAssign(cost0)
			if err != nil {
				return nil, cost, err
			}

			bl, cost0, err := h.GetBalance("iost.bonus")
			cost.AddAssign(cost0)
			if err != nil {
				return nil, cost, err
			}
			token := int64(amount * 1.0 / totalServi * bl)
			if token > bl {
				token = bl
			}

			cost0, err = h.Withdraw(acc, token)
			cost.AddAssign(cost0)

			return []interface{}{}, cost, err
		},
	}
)
