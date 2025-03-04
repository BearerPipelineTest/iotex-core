// Copyright (c) 2022 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package api

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iotexproject/go-pkgs/util"
	"github.com/iotexproject/iotex-address/address"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"google.golang.org/protobuf/proto"

	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/actpool"
	"github.com/iotexproject/iotex-core/blockchain"
	"github.com/iotexproject/iotex-core/blockchain/blockdao"
	"github.com/iotexproject/iotex-core/test/identityset"
	"github.com/iotexproject/iotex-core/testutil"
)

const (
	_evmNetworkID uint32 = 1
)

func TestGasPriceIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	ret, _ := svr.gasPrice()
	require.Equal(uint64ToHex(1000000000000), ret)
}

func TestGetChainIDIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	ret, _ := svr.getChainID()
	require.Equal(uint64ToHex(uint64(_evmNetworkID)), ret)
}

func TestGetBlockNumberIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	ret, _ := svr.getBlockNumber()
	require.Equal(uint64ToHex(4), ret)
}

func TestGetBlockByNumberIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	testData := []struct {
		data     string
		expected int
	}{
		{`{"params": ["1", true]}`, 1},
		{`{"params": ["1", false]}`, 2},
		{`{"params": ["10", false]}`, 0},
	}

	for i, v := range testData {
		t.Run(fmt.Sprintf("%d-%d", i, len(testData)-1), func(t *testing.T) {
			data := gjson.Parse(v.data)
			ret, err := svr.getBlockByNumber(&data)
			require.NoError(err)
			if v.expected == 0 {
				require.Nil(ret)
				return
			}
			blk, ok := ret.(*getBlockResult)
			require.True(ok)
			require.Equal(len(blk.transactions), v.expected)
		})
	}
}

func TestGetBalanceIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	testData := gjson.Parse(`{"params": ["0xDa7e12Ef57c236a06117c5e0d04a228e7181CF36", 1]}`)
	ret, _ := svr.getBalance(&testData)
	ans, _ := new(big.Int).SetString("9999999999999999999999999991", 10)
	require.Equal("0x"+fmt.Sprintf("%x", ans), ret)
}

func TestGetTransactionCountIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	testData := []struct {
		data     string
		expected int
	}{
		{`{"params": ["0xDa7e12Ef57c236a06117c5e0d04a228e7181CF36", "0x1"]}`, 2},
		{`{"params": ["0xDa7e12Ef57c236a06117c5e0d04a228e7181CF36", "pending"]}`, 2},
	}

	for i, v := range testData {
		t.Run(fmt.Sprintf("%d-%d", i, len(testData)-1), func(t *testing.T) {
			data := gjson.Parse(v.data)
			ret, _ := svr.getTransactionCount(&data)
			require.Equal(uint64ToHex(uint64(v.expected)), ret)
		})
	}
}

func TestCallIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	testData := []struct {
		data string
	}{
		{
			`{"params": [{
				"from":     "",
				"to":       "0x7c13866F9253DEf79e20034eDD011e1d69E67fe5",
				"gas":      "0x4e20",
				"gasPrice": "0xe8d4a51000",
				"value":    "0x1",
				"data":     "0x1"},
			1]}`,
		},
		{
			`{"params": [{
				"from":     "",
				"to":       "0xb1f8e55c7f64d203c1400b9d8555d050f94adf39",
				"gas":      "0x4e20",
				"gasPrice": "0xe8d4a51000",
				"value":    "0x1",
				"data":     "0x1"},
			1]}`,
		},
	}

	for i, v := range testData {
		t.Run(fmt.Sprintf("%d-%d", i, len(testData)-1), func(t *testing.T) {
			data := gjson.Parse(v.data)
			_, err := svr.call(&data)
			require.NoError(err)
		})
	}
}

func TestEstimateGasIntegrity(t *testing.T) {
	require := require.New(t)
	svr, bc, dao, actPool, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	// deploy a contract
	contractCode := "608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806360fe47b11461003b5780636d4ce63c14610057575b600080fd5b6100556004803603810190610050919061009d565b610075565b005b61005f61007f565b60405161006c91906100d9565b60405180910390f35b8060008190555050565b60008054905090565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea2646970667358221220c86a8c4dd175f55f5732b75b721d714ceb38a835b87c6cf37cf28c790813e19064736f6c63430008070033"
	contract, _ := deployContractV2(bc, dao, actPool, identityset.PrivateKey(13), 1, bc.TipHeight(), contractCode)

	fromAddr, _ := ioAddrToEthAddr(identityset.Address(0).String())
	toAddr, _ := ioAddrToEthAddr(identityset.Address(28).String())
	contractAddr, _ := ioAddrToEthAddr(contract)
	testData := []struct {
		input  string
		result uint64
	}{
		{
			input: fmt.Sprintf(`{"params": [{
					"from":     "%s",
					"to":       "%s",
					"gas":      "0x0",
					"gasPrice": "0x0",
					"value":    "0x0",
					"data":     "0x1123123c"},
				1]}`, fromAddr, toAddr),
			result: 21000,
		},
		{
			input: fmt.Sprintf(`{"params": [{
				"from":     "%s",
				"to":       "%s",
				"gas":      "0x0",
				"gasPrice": "0x0",
				"value":    "0x0",
				"data":      "344933be000000000000000000000000000000000000000000000000000be497a92e9f3300000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000f8be4046fd89199906ca348bcd3822c4b250e246000000000000000000000000000000000000000000000000000000006173a15400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000a00744882684c3e4747faefd68d283ea44099d030000000000000000000000000258866edaf84d6081df17660357ab20a07d0c80"},
				1]}`, fromAddr, toAddr),
			result: 36000,
		},
		{
			input: fmt.Sprintf(`{"params": [{
				"from":     "%s",
				"to":       "%s",
				"gas":      "0x0",
				"gasPrice": "0x0",
				"value":    "0x0",
				"data":     "0x6d4ce63c"},
			1]}`, fromAddr, contractAddr),
			result: 21000,
		},
	}

	for i, v := range testData {
		t.Run(fmt.Sprintf("%d-%d", i, len(testData)-1), func(t *testing.T) {
			input := gjson.Parse(v.input)
			ret, err := svr.estimateGas(&input)
			require.NoError(err)
			require.Equal(ret, uint64ToHex(v.result))
		})
	}
}

func TestSendRawTransactionIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	testData := gjson.Parse(`{"params": ["f8600180830186a09412745fec82b585f239c01090882eb40702c32b04808025a0b0e1aab5b64d744ae01fc9f1c3e9919844a799e90c23129d611f7efe6aec8a29a0195e28d22d9b280e00d501ff63525bb76f5c87b8646c89d5d9c5485edcb1b498"]}`)
	res, _ := svr.sendRawTransaction(&testData)
	require.Equal("0x778fd5a054e74e9055bf68ef5f9d559fa306e8ba7dee608d0a3624cca0b63b3e", res)
}

func TestGetCodeIntegrity(t *testing.T) {
	require := require.New(t)
	svr, bc, dao, actPool, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	// deploy a contract
	contractCode := "608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806360fe47b11461003b5780636d4ce63c14610057575b600080fd5b6100556004803603810190610050919061009d565b610075565b005b61005f61007f565b60405161006c91906100d9565b60405180910390f35b8060008190555050565b60008054905090565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea2646970667358221220c86a8c4dd175f55f5732b75b721d714ceb38a835b87c6cf37cf28c790813e19064736f6c63430008070033"
	contract, _ := deployContractV2(bc, dao, actPool, identityset.PrivateKey(13), 1, bc.TipHeight(), contractCode)

	contractAddr, _ := ioAddrToEthAddr(contract)
	testData := gjson.Parse(fmt.Sprintf(`{"params": ["%s", 1]}`, contractAddr))
	ret, _ := svr.getCode(&testData)
	require.Contains(contractCode, util.Remove0xPrefix(ret.(string)))
}

func TestGetNodeInfoIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	_, err := svr.getNodeInfo()
	require.NoError(err)
}

func TestGetBlockTransactionCountByHashIntegrity(t *testing.T) {
	require := require.New(t)
	svr, bc, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	header, err := bc.BlockHeaderByHeight(1)
	require.NoError(err)
	blkHash := header.HashBlock()
	testData := gjson.Parse(fmt.Sprintf(`{"params":["0x%s", 1]}`, hex.EncodeToString(blkHash[:])))
	ret, err := svr.getBlockTransactionCountByHash(&testData)
	require.NoError(err)
	require.Equal(uint64ToHex(2), ret)
}

func TestGetBlockByHashIntegrity(t *testing.T) {
	require := require.New(t)
	svr, bc, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	header, _ := bc.BlockHeaderByHeight(1)
	blkHash := header.HashBlock()

	testData := gjson.Parse(fmt.Sprintf(`{"params":["0x%s", false]}`, hex.EncodeToString(blkHash[:])))
	ret, err := svr.getBlockByHash(&testData)
	require.NoError(err)
	ans := ret.(*getBlockResult)
	require.Equal(blkHash, ans.blk.HashBlock())
	require.Equal(2, len(ans.transactions))

	testData2 := gjson.Parse(`{"params":["0xa2e8e0c9cafbe93f2b7f7c9d32534bc6fde95f2185e5f2aaa6bf7ebdf1a6610a", false]}`)
	ret, err = svr.getBlockByHash(&testData2)
	require.NoError(err)
	require.Nil(ret)
}

func TestGetTransactionByHashIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	testData := gjson.Parse(fmt.Sprintf(`{"params":["0x%s", false]}`, hex.EncodeToString(_transferHash1[:])))
	ret, err := svr.getTransactionByHash(&testData)
	require.NoError(err)
	require.Equal(_transferHash1, ret.(*getTransactionResult).receipt.ActionHash)

	testData2 := gjson.Parse(fmt.Sprintf(`{"params":["0x%s", false]}`, "0x58df1e9cb0572fea48e8ce9d9b787ae557c304657d01890f4fc5ea88a1f44c3e"))
	ret, err = svr.getTransactionByHash(&testData2)
	require.NoError(err)
	require.Nil(ret)
}

func TestGetLogsIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	testData := []struct {
		data   *filterObject
		logLen int
	}{
		{
			&filterObject{FromBlock: "0x1"},
			4,
		},
		{
			// empty log
			&filterObject{Address: []string{"0x8ce313ab12bf7aed8136ab36c623ff98c8eaad34"}},
			0,
		},
	}
	for i, v := range testData {
		t.Run(fmt.Sprintf("%d-%d", i, len(testData)-1), func(t *testing.T) {
			ret, err := svr.getLogs(v.data)
			require.NoError(err)
			require.Equal(len(ret.([]*getLogsResult)), v.logLen)
		})
	}
}

func TestGetTransactionReceiptIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	testData := gjson.Parse(fmt.Sprintf(`{"params":["0x%s", 1]}`, hex.EncodeToString(_transferHash1[:])))
	ret, err := svr.getTransactionReceipt(&testData)
	require.NoError(err)
	ans, ok := ret.(*getReceiptResult)
	require.True(ok)
	require.Equal(_transferHash1, ans.receipt.ActionHash)
	toAddr, _ := ioAddrToEthAddr(identityset.Address(30).String())
	require.Equal(identityset.Address(27), ans.from)
	require.Equal(toAddr, *ans.to)
	require.Nil(nil, ans.contractAddress)
	require.Equal(uint64(10000), ans.receipt.GasConsumed)
	require.Equal(uint64(1), ans.receipt.BlockHeight)

	testData2 := gjson.Parse(`{"params": ["0x58df1e9cb0572fea48e8ce9d9b787ae557c304657d01890f4fc5ea88a1f44c3e", 1]}`)
	ret, err = svr.getTransactionReceipt(&testData2)
	require.NoError(err)
	require.Nil(ret)
}

func TestGetBlockTransactionCountByNumberIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	testData := gjson.Parse(`{"params": ["0x1", 1]}`)
	ret, err := svr.getBlockTransactionCountByNumber(&testData)
	require.NoError(err)
	require.Equal(ret, uint64ToHex(2))
}

func TestGetTransactionByBlockHashAndIndexIntegrity(t *testing.T) {
	require := require.New(t)
	svr, bc, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	header, _ := bc.BlockHeaderByHeight(1)
	blkHash := header.HashBlock()

	testData := gjson.Parse(fmt.Sprintf(`{"params":["0x%s", "0x0"]}`, hex.EncodeToString(blkHash[:])))
	ret, err := svr.getTransactionByBlockHashAndIndex(&testData)
	ans := ret.(*getTransactionResult)
	require.NoError(err)
	require.Equal(_transferHash1, ans.receipt.ActionHash)
	toAddr, _ := ioAddrToEthAddr(identityset.Address(30).String())
	require.Equal(identityset.Address(27), ans.pubkey.Address())
	require.Equal(toAddr, *ans.to)
	require.Equal(uint64(20000), ans.ethTx.Gas())
	require.Equal(big.NewInt(0), ans.ethTx.GasPrice())

	testData2 := gjson.Parse(fmt.Sprintf(`{"params":["0x%s", "0x10"]}`, hex.EncodeToString(blkHash[:])))
	ret, err = svr.getTransactionByBlockHashAndIndex(&testData2)
	require.NoError(err)
	require.Nil(ret)

	testData3 := gjson.Parse(fmt.Sprintf(`{"params":["0x%s", "0x0"]}`, "0xa2e8e0c9cafbe93f2b7f7c9d32534bc6fde95f2185e5f2aaa6bf7ebdf1a6610a"))
	ret, err = svr.getTransactionByBlockHashAndIndex(&testData3)
	require.NoError(err)
	require.Nil(ret)
}

func TestGetTransactionByBlockNumberAndIndexIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	testData := gjson.Parse(`{"params": ["0x1", "0x0"]}`)
	ret, err := svr.getTransactionByBlockNumberAndIndex(&testData)
	ans := ret.(*getTransactionResult)
	require.NoError(err)
	require.Equal(_transferHash1, ans.receipt.ActionHash)
	toAddr, _ := ioAddrToEthAddr(identityset.Address(30).String())
	require.Equal(identityset.Address(27), ans.pubkey.Address())
	require.Equal(toAddr, *ans.to)
	require.Equal(uint64(20000), ans.ethTx.Gas())
	require.Equal(big.NewInt(0), ans.ethTx.GasPrice())

	testData2 := gjson.Parse(`{"params": ["0x1", "0x10"]}`)
	ret, err = svr.getTransactionByBlockNumberAndIndex(&testData2)
	require.NoError(err)
	require.Nil(ret)

	testData3 := gjson.Parse(`{"params": ["0x10", "0x0"]}`)
	ret, err = svr.getTransactionByBlockNumberAndIndex(&testData3)
	require.NoError(err)
	require.Nil(ret)
}

func TestNewfilterIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	testData := &filterObject{FromBlock: "0x1"}
	ret, err := svr.newFilter(testData)
	require.NoError(err)
	require.Equal(ret, "0xe10f7dd489b75a36de8e246eb974827fe86a02ed19d9b475a1600cf4f935feff")
}

func TestNewBlockFilterIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	ret, err := svr.newBlockFilter()
	require.NoError(err)
	require.Equal("0x71371f8dbaefc4c96d2534163a1b461951c88520cd32bc03b5bfdfe7340bc187", ret)
}

func TestGetFilterChangesIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	// filter
	filterReq := &filterObject{FromBlock: "0x1"}
	filterID1, _ := svr.newFilter(filterReq)
	filterID1Req := gjson.Parse(fmt.Sprintf(`{"params":["%s"]}`, filterID1.(string)))
	ret, err := svr.getFilterChanges(&filterID1Req)
	require.NoError(err)
	require.Equal(len(ret.([]*getLogsResult)), 4)
	// request again after last rolling
	ret, err = svr.getFilterChanges(&filterID1Req)
	require.NoError(err)
	require.Equal(len(ret.([]*getLogsResult)), 0)

	// blockfilter
	filterID2, _ := svr.newBlockFilter()
	filterID2Req := gjson.Parse(fmt.Sprintf(`{"params":["%s"]}`, filterID2.(string)))
	ret2, err := svr.getFilterChanges(&filterID2Req)
	require.NoError(err)
	require.Equal(1, len(ret2.([]string)))
	ret3, err := svr.getFilterChanges(&filterID2Req)
	require.NoError(err)
	require.Equal(0, len(ret3.([]string)))

}

func TestGetFilterLogsIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	filterReq := &filterObject{FromBlock: "0x1"}
	filterID, _ := svr.newFilter(filterReq)
	filterIDReq := gjson.Parse(fmt.Sprintf(`{"params":["%s"]}`, filterID.(string)))
	ret, err := svr.getFilterLogs(&filterIDReq)
	require.NoError(err)
	require.Equal(len(ret.([]*getLogsResult)), 4)
}

func TestLocalAPICacheIntegrity(t *testing.T) {
	require := require.New(t)

	testKey, testData := strconv.Itoa(rand.Int()), []byte(strconv.Itoa(rand.Int()))
	cacheLocal := newAPICache(1*time.Second, "")
	_, exist := cacheLocal.Get(testKey)
	require.False(exist)
	err := cacheLocal.Set(testKey, testData)
	require.NoError(err)
	data, _ := cacheLocal.Get(testKey)
	require.Equal(data, testData)
	cacheLocal.Del(testKey)
	_, exist = cacheLocal.Get(testKey)
	require.False(exist)
}

func TestGetStorageAtIntegrity(t *testing.T) {
	require := require.New(t)
	svr, bc, dao, actPool, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	// deploy a contract
	contractCode := "608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806360fe47b11461003b5780636d4ce63c14610057575b600080fd5b6100556004803603810190610050919061009d565b610075565b005b61005f61007f565b60405161006c91906100d9565b60405180910390f35b8060008190555050565b60008054905090565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea2646970667358221220c86a8c4dd175f55f5732b75b721d714ceb38a835b87c6cf37cf28c790813e19064736f6c63430008070033"
	contract, _ := deployContractV2(bc, dao, actPool, identityset.PrivateKey(13), 1, bc.TipHeight(), contractCode)

	contractAddr, _ := ioAddrToEthAddr(contract)
	testData := gjson.Parse(fmt.Sprintf(`{"params": ["%s", "0x0"]}`, contractAddr))
	ret, err := svr.getStorageAt(&testData)
	require.NoError(err)
	// the value of any contract at pos0 is be "0x0000000000000000000000000000000000000000000000000000000000000000"
	require.Equal("0x0000000000000000000000000000000000000000000000000000000000000000", ret)

	failData := []gjson.Result{
		gjson.Parse(`{"params": [1]}`),
		gjson.Parse(`{"params": ["TEST", "TEST"]}`),
	}
	for _, v := range failData {
		_, err := svr.getStorageAt(&v)
		require.Error(err)
	}
}

func TestGetNetworkIDIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	res, _ := svr.getNetworkID()
	require.Equal(fmt.Sprintf("%d", _evmNetworkID), res)
}

func setupTestWeb3Server() (*web3Handler, blockchain.Blockchain, blockdao.BlockDAO, actpool.ActPool, func()) {
	cfg := newConfig()

	// TODO (zhi): revise
	bc, dao, indexer, bfIndexer, sf, ap, registry, bfIndexFile, err := setupChain(cfg)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// Start blockchain
	if err := bc.Start(ctx); err != nil {
		panic(err)
	}
	// Add testing blocks
	if err := addTestingBlocks(bc, ap); err != nil {
		panic(err)
	}
	opts := []Option{WithBroadcastOutbound(func(ctx context.Context, chainID uint32, msg proto.Message) error {
		return nil
	})}
	core, err := newCoreService(cfg.API, bc, nil, sf, dao, indexer, bfIndexer, ap, registry, opts...)
	if err != nil {
		panic(err)
	}

	return &web3Handler{core, newAPICache(15*time.Minute, "")}, bc, dao, ap, func() {
		testutil.CleanupPath(bfIndexFile)
	}
}

func TestEthAccountsIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()
	res, _ := svr.ethAccounts()
	require.Equal(0, len(res.([]string)))
}

func TestWeb3StakingIntegrity(t *testing.T) {
	require := require.New(t)
	svr, _, _, _, cleanCallback := setupTestWeb3Server()
	defer cleanCallback()

	ecdsaPvk, ok := identityset.PrivateKey(28).EcdsaPrivateKey().(*ecdsa.PrivateKey)
	require.True(ok)

	type stakeData struct {
		testName         string
		stakeEncodedData []byte
	}
	testData := []stakeData{}
	toAddr, err := ioAddrToEthAddr(address.StakingProtocolAddr)
	require.NoError(err)

	// encode stake data
	act1, err := action.NewCreateStake(1, "test", "100", 7, false, []byte{}, 1000000, big.NewInt(0))
	require.NoError(err)
	data, err := act1.EncodeABIBinary()
	require.NoError(err)
	testData = append(testData, stakeData{"createStake", data})

	act2, err := action.NewDepositToStake(2, 7, "100", []byte{}, 1000000, big.NewInt(0))
	require.NoError(err)
	data2, err := act2.EncodeABIBinary()
	require.NoError(err)
	testData = append(testData, stakeData{"depositToStake", data2})

	act3, err := action.NewChangeCandidate(3, "test", 7, []byte{}, 1000000, big.NewInt(0))
	require.NoError(err)
	data3, err := act3.EncodeABIBinary()
	require.NoError(err)
	testData = append(testData, stakeData{"changeCandidate", data3})

	act4, err := action.NewUnstake(4, 7, []byte{}, 1000000, big.NewInt(0))
	require.NoError(err)
	data4, err := act4.EncodeABIBinary()
	require.NoError(err)
	testData = append(testData, stakeData{"unstake", data4})

	act5, err := action.NewWithdrawStake(5, 7, []byte{}, 1000000, big.NewInt(0))
	require.NoError(err)
	data5, err := act5.EncodeABIBinary()
	require.NoError(err)
	testData = append(testData, stakeData{"withdrawStake", data5})

	act6, err := action.NewRestake(6, 7, 7, false, []byte{}, 1000000, big.NewInt(0))
	require.NoError(err)
	data6, err := act6.EncodeABIBinary()
	require.NoError(err)
	testData = append(testData, stakeData{"restake", data6})

	act7, err := action.NewTransferStake(7, "io1xpq62aw85uqzrccg9y5hnryv8ld2nkpycc3gza", 7, []byte{}, 1000000, big.NewInt(0))
	require.NoError(err)
	data7, err := act7.EncodeABIBinary()
	require.NoError(err)
	testData = append(testData, stakeData{"transferStake", data7})

	act8, err := action.NewCandidateRegister(
		8,
		"test",
		"io1xpq62aw85uqzrccg9y5hnryv8ld2nkpycc3gza",
		"io1xpq62aw85uqzrccg9y5hnryv8ld2nkpycc3gza",
		"io1xpq62aw85uqzrccg9y5hnryv8ld2nkpycc3gza",
		"100",
		7,
		false,
		[]byte{},
		1000000,
		big.NewInt(0))
	require.NoError(err)
	data8, err := act8.EncodeABIBinary()
	require.NoError(err)
	testData = append(testData, stakeData{"candidateRegister", data8})

	act9, err := action.NewCandidateUpdate(
		9,
		"test",
		"io1xpq62aw85uqzrccg9y5hnryv8ld2nkpycc3gza",
		"io1xpq62aw85uqzrccg9y5hnryv8ld2nkpycc3gza",
		1000000,
		big.NewInt(0))
	require.NoError(err)
	data9, err := act9.EncodeABIBinary()
	require.NoError(err)
	testData = append(testData, stakeData{"candidateUpdate", data9})

	for i, test := range testData {
		t.Run(test.testName, func(t *testing.T) {
			// estimate gas
			gasLimit, err := estimateStakeGas(svr, identityset.Address(28).Hex(), toAddr, test.stakeEncodedData)
			require.NoError(err)

			// create tx
			rawTx := types.NewTransaction(
				uint64(9+i),
				common.HexToAddress(toAddr),
				big.NewInt(0),
				gasLimit,
				big.NewInt(0),
				test.stakeEncodedData,
			)
			tx, err := types.SignTx(rawTx, types.NewEIP155Signer(big.NewInt(int64(_evmNetworkID))), ecdsaPvk)
			require.NoError(err)
			BinaryData, err := tx.MarshalBinary()
			require.NoError(err)

			// send tx
			fmt.Println(hex.EncodeToString(BinaryData))
			rawData := gjson.Parse(fmt.Sprintf(`{"params": ["%s"]}`, hex.EncodeToString(BinaryData)))
			_, err = svr.sendRawTransaction(&rawData)
			require.NoError(err)
		})
	}
}

func estimateStakeGas(svr *web3Handler, fromAddr, toAddr string, data []byte) (uint64, error) {
	input := gjson.Parse(fmt.Sprintf(`{"params": [{
		"from":     "%s",
		"to":       "%s",
		"gas":      "0x0",
		"gasPrice": "0x0",
		"value":    "0x0",
		"data":     "%s"},
	1]}`, fromAddr, toAddr, hex.EncodeToString(data)))
	ret, err := svr.estimateGas(&input)
	if err != nil {
		panic(err)
	}
	return hexStringToNumber(ret.(string))
}
