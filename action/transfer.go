// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package action

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/iotexproject/iotex-core/pkg/util/byteutil"
	"github.com/iotexproject/iotex-core/pkg/version"
)

const (
	// TransferPayloadGas represents the transfer payload gas per uint
	TransferPayloadGas = uint64(100)
	// TransferBaseIntrinsicGas represents the base intrinsic gas for transfer
	TransferBaseIntrinsicGas = uint64(10000)
)

var _ hasDestination = (*Transfer)(nil)

// Transfer defines the struct of account-based transfer
type Transfer struct {
	AbstractAction

	amount    *big.Int
	recipient string
	payload   []byte
}

// NewTransfer returns a Transfer instance
func NewTransfer(
	nonce uint64,
	amount *big.Int,
	recipient string,
	payload []byte,
	gasLimit uint64,
	gasPrice *big.Int,
) (*Transfer, error) {
	return &Transfer{
		AbstractAction: AbstractAction{
			version:  version.ProtocolVersion,
			nonce:    nonce,
			gasLimit: gasLimit,
			gasPrice: gasPrice,
		},
		recipient: recipient,
		amount:    amount,
		payload:   payload,
		// SenderPublicKey and Signature will be populated in Sign()
	}, nil
}

// Amount returns the amount
func (tsf *Transfer) Amount() *big.Int { return tsf.amount }

// Payload returns the payload bytes
func (tsf *Transfer) Payload() []byte { return tsf.payload }

// Recipient returns the recipient address. It's the wrapper of Action.DstAddr
func (tsf *Transfer) Recipient() string { return tsf.recipient }

// Destination returns the recipient address as destination.
func (tsf *Transfer) Destination() string { return tsf.recipient }

// TotalSize returns the total size of this Transfer
func (tsf *Transfer) TotalSize() uint32 {
	size := tsf.BasicActionSize()
	if tsf.amount != nil && len(tsf.amount.Bytes()) > 0 {
		size += uint32(len(tsf.amount.Bytes()))
	}
	// 65 is the pubkey size
	return size + uint32(len(tsf.payload)) + 65
}

// Serialize returns a raw byte stream of this Transfer
func (tsf *Transfer) Serialize() []byte {
	return byteutil.Must(proto.Marshal(tsf.Proto()))
}

// Proto converts Transfer to protobuf's Action
func (tsf *Transfer) Proto() *iotextypes.Transfer {
	// used by account-based model
	act := &iotextypes.Transfer{
		Recipient: tsf.recipient,
		Payload:   tsf.payload,
	}

	if tsf.amount != nil {
		act.Amount = tsf.amount.String()
	}
	return act
}

// LoadProto converts a protobuf's Action to Transfer
func (tsf *Transfer) LoadProto(pbAct *iotextypes.Transfer) error {
	if pbAct == nil {
		return ErrNilProto
	}
	if tsf == nil {
		return ErrNilAction
	}
	*tsf = Transfer{}

	tsf.recipient = pbAct.GetRecipient()
	tsf.payload = pbAct.GetPayload()
	if pbAct.GetAmount() == "" {
		tsf.amount = big.NewInt(0)
	} else {
		amount, ok := new(big.Int).SetString(pbAct.GetAmount(), 10)
		// tsf amount gets zero when pbAct.GetAmount is empty string
		if !ok {
			return errors.Errorf("invalid amount %s", pbAct.GetAmount())
		}
		tsf.amount = amount
	}
	return nil
}

// IntrinsicGas returns the intrinsic gas of a transfer
func (tsf *Transfer) IntrinsicGas() (uint64, error) {
	payloadSize := uint64(len(tsf.Payload()))
	return CalculateIntrinsicGas(TransferBaseIntrinsicGas, TransferPayloadGas, payloadSize)
}

// Cost returns the total cost of a transfer
func (tsf *Transfer) Cost() (*big.Int, error) {
	intrinsicGas, err := tsf.IntrinsicGas()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get intrinsic gas for the transfer")
	}
	transferFee := big.NewInt(0).Mul(tsf.GasPrice(), big.NewInt(0).SetUint64(intrinsicGas))
	return big.NewInt(0).Add(tsf.Amount(), transferFee), nil
}

// SanityCheck validates the variables in the action
func (tsf *Transfer) SanityCheck() error {
	// Reject transfer of negative amount
	if tsf.Amount().Sign() < 0 {
		return ErrNegativeValue
	}
	return tsf.AbstractAction.SanityCheck()
}

// ToEthTx converts action to eth-compatible tx
func (tsf *Transfer) ToEthTx() (*types.Transaction, error) {
	addr, err := address.FromString(tsf.recipient)
	if err != nil {
		return nil, err
	}
	ethAddr := common.BytesToAddress(addr.Bytes())
	return types.NewTransaction(tsf.Nonce(), ethAddr, tsf.amount, tsf.GasLimit(), tsf.GasPrice(), tsf.payload), nil
}
