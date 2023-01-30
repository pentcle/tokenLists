package ethereum

import (
	"context"
	"math"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

var RPC = map[uint64]*ethclient.Client{}

// SUPPORTED_CHAIN_IDS lists the chainIDs supported by our program
var SUPPORTED_CHAIN_IDS = []uint64{1, 10, 56, 137, 250, 42161}

// RPC_ENDPOINTS contains the node endpoints to connect the blockchains
// Can be overwritten by env variables
var RPC_ENDPOINTS = map[uint64]string{
	1:     `https://eth.public-rpc.com`,
	10:    `https://mainnet.optimism.io`,
	56:    `https://bsc.rpc.blxrbdn.com`,
	137:   `https://polygon.llamarpc.com`,
	250:   `https://rpc.ftm.tools`,
	42161: `https://arbitrum.public-rpc.com`,
}

// MULTICALL_ADDRESSES contains the address of the multicall2 contract for a specific chainID
var MULTICALL_ADDRESSES = map[uint64]common.Address{
	1:     common.HexToAddress(`0x5ba1e12693dc8f9c48aad8770482f4739beed696`),
	10:    common.HexToAddress(`0xca11bde05977b3631167028862be2a173976ca11`),
	56:    common.HexToAddress(`0xca11bde05977b3631167028862be2a173976ca11`),
	137:   common.HexToAddress(`0xca11bde05977b3631167028862be2a173976ca11`),
	250:   common.HexToAddress(`0x470ADB45f5a9ac3550bcFFaD9D990Bf7e2e941c9`),
	42161: common.HexToAddress(`0x842eC2c7D803033Edf55E478F461FC547Bc54EB2`),
}

func init() {
	godotenv.Load(`.env`)
	for _, chainID := range SUPPORTED_CHAIN_IDS {
		RPC_ENDPOINTS[chainID] = useEnv(`RPC_URI_FOR_`+strconv.FormatUint(chainID, 10), RPC_ENDPOINTS[chainID])
	}
}

// GetRPC returns the current connection for a specific chain
func GetRPC(chainID uint64) *ethclient.Client {
	return RPC[chainID]
}

// GetRPCURI returns the URI to use to connect to the node for a specific chainID
func GetRPCURI(chainID uint64) string {
	return RPC_ENDPOINTS[chainID]
}

// MulticallClientForChainID holds the multicall client for a specific chainID
var MulticallClientForChainID = make(map[uint64]TEthMultiCaller)

func randomSigner() *bind.TransactOpts {
	privateKey, _ := crypto.GenerateKey()
	signer, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
	signer.NoSend = true
	signer.Context = context.Background()
	signer.GasLimit = math.MaxUint64
	signer.GasFeeCap = big.NewInt(0)
	signer.GasTipCap = big.NewInt(0)
	signer.GasPrice = big.NewInt(0)
	return signer
}

func useEnv(envName string, fallback string) string {
	envValue := os.Getenv(envName)
	if envValue == "" {
		os.Setenv(envName, fallback)
		return fallback
	}
	return os.Getenv(envName)
}
