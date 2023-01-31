package liteclient

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo"
)

func (c *Client) GetSeqno(ctx context.Context, account tongo.AccountID) (uint32, error) {
	errCode, stack, err := c.RunSmcMethod(ctx, 4, account, "seqno", tongo.VmStack{})
	if err != nil {
		return 0, err
	}
	if errCode == 0xFFFFFF00 {
		return 0, nil
	} else if errCode != 0 && errCode != 1 {
		return 0, fmt.Errorf("method execution failed with code: %v", errCode)
	}
	if len(stack) != 1 || stack[0].SumType != "VmStkTinyInt" {
		return 0, fmt.Errorf("invalid stack")
	}
	return uint32(stack[0].VmStkTinyInt), nil
}
