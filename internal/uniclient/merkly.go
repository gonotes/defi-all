package uniclient

import (
	"github.com/hardstylez72/cry/internal/defi"
	"github.com/hardstylez72/cry/internal/defi/zksyncera"
	v1 "github.com/hardstylez72/cry/internal/pb/gen/proto/go/v1"
	"github.com/hardstylez72/cry/internal/socks5"
	"github.com/pkg/errors"
)

func NewMerklyMintAndBridgeNFT(network v1.Network, c *BaseClientConfig) (defi.MintAndBridgeNFT, error) {

	proxy, err := socks5.NewSock5ProxyString(c.ProxyString, c.UserAgentHeader)
	if err != nil {
		return nil, err
	}

	var cli defi.MintAndBridgeNFT
	switch network {
	//case v1.Network_ARBITRUM:
	//	cli, err = arbitrum.NewClient(&arbitrum.ClientConfig{HttpCli: proxy.CliL1, RPCEndpoint: c.RPCEndpoint})
	//case v1.Network_Etherium:
	//	cli, err = etherium.NewClient(&etherium.ClientConfig{HttpCli: proxy.CliL1, RPCEndpoint: c.RPCEndpoint})
	//case v1.Network_BinanaceBNB:
	//	cli, err = bnb.NewClient(&bnb.ClientConfig{HttpCli: proxy.CliL1, RPCEndpoint: c.RPCEndpoint})
	//case v1.Network_OPTIMISM:
	//	cli, err = optimism.NewClient(&optimism.ClientConfig{HttpCli: proxy.CliL1, RPCEndpoint: c.RPCEndpoint})
	//case v1.Network_POLIGON:
	//	cli, err = poligon.NewClient(&poligon.ClientConfig{HttpCli: proxy.CliL1, RPCEndpoint: c.RPCEndpoint})
	//case v1.Network_AVALANCHE:
	//	cli, err = avalanche.NewClient(&avalanche.ClientConfig{HttpCli: proxy.CliL1, RPCEndpoint: c.RPCEndpoint})
	case v1.Network_ZKSYNCERA:
		cli, err = zksyncera.NewMainNetClient(&zksyncera.ClientConfig{HttpCli: proxy.Cli, RPCEndpoint: c.RPCEndpoint})
	default:
		return nil, errors.New("network is not supported for Transfer")
	}

	if err != nil {
		return nil, err
	}

	return cli, nil
}
