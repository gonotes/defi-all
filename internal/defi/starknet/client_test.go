package starknet

import (
	"context"
	"testing"

	v1 "github.com/hardstylez72/cry/internal/pb/gen/proto/go/v1"
	"github.com/hardstylez72/cry/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	c, err := NewClient(&ClientConfig{
		HttpCli:     tests.GetConfig().Cli,
		RPCEndpoint: MainnetRPC,
		Proxy:       "",
	})
	assert.NoError(t, err)
	assert.NotNil(t, c)

	ctx := context.Background()
	//r, err := c.Approve(ctx, &ApproveReq{
	//	Token:       v1.Token_ETH,
	//	Amount:      big.NewInt(1000000000),
	//	SpenderAddr: "0x7a6f98c03379b9513ca84cca1373ff452a7462a3b61598f0af5bb27ad7f76d1",
	//	PK:          tests.GetConfig().StarkNetPrivate,
	//})
	//assert.NotNil(t, r)

	//c.GetBalance(ctx, &defi.GetBalanceReq{
	//	WalletAddress: tests.GetConfig().StarkNetPuvlic,
	//	Token:         v1.Token_ETH,
	//})

	//assert.NoError(t, err)
	//assert.NotNil(t, r)
	res, err := c.Allowed(ctx, &AllowedReq{
		Token:       v1.Token_ETH,
		WalletAddr:  tests.GetConfig().StarkNetPuvlic,
		SpenderAddr: "0x7a6f98c03379b9513ca84cca1373ff452a7462a3b61598f0af5bb27ad7f76d1",
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestGen(t *testing.T) {
	pub, err := GetPublicKeyHash(tests.GetConfig().StarkNetPrivate)
	assert.NoError(t, err)
	assert.Equal(t, pub, tests.GetConfig().StarkNetPuvlic)
}

//println(err)
//println(res == nil)

//err = c.DeployAccount(ctx, &Account{
//	PublicKey:  tests.GetConfig().StarkNetPuvlic,
//	PrivateKey: tests.GetConfig().StarkNetPrivate,
//})
//assert.NoError(t, err)
//}
