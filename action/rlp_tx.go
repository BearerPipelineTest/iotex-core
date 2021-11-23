package action

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/go-pkgs/util"
	"github.com/iotexproject/iotex-address/address"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

type rlpTransaction interface {
	Nonce() uint64
	GasPrice() *big.Int
	GasLimit() uint64
	Recipient() string
	Amount() *big.Int
	Payload() []byte
}

func rlpRawHash(tx rlpTransaction, chainID uint32) (hash.Hash256, error) {
	rawTx, err := RlpToEthTx(tx)
	if err != nil {
		return hash.ZeroHash256, err
	}
	h := types.NewEIP155Signer(big.NewInt(int64(chainID))).Hash(rawTx)
	return hash.BytesToHash256(h[:]), nil
}

func rlpSignedHash(tx rlpTransaction, chainID uint32, sig []byte) (hash.Hash256, error) {
	signedTx, err := reconstructSignedRlpTxFromSig(tx, chainID, sig)
	if err != nil {
		return hash.ZeroHash256, err
	}
	h := sha3.NewLegacyKeccak256()
	rlp.Encode(h, signedTx)
	return hash.BytesToHash256(h.Sum(nil)), nil
}

func RlpToEthTx(act rlpTransaction) (*types.Transaction, error) {
	if act == nil {
		return nil, errors.New("nil action to generate RLP tx")
	}

	// generate raw tx
	if to := act.Recipient(); to != EmptyAddress {
		addr, err := address.FromString(to)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid recipient address %s", to)
		}
		ethAddr := common.BytesToAddress(addr.Bytes())
		return types.NewTransaction(act.Nonce(), ethAddr, act.Amount(), act.GasLimit(), act.GasPrice(), act.Payload()), nil
	}
	return types.NewContractCreation(act.Nonce(), act.Amount(), act.GasLimit(), act.GasPrice(), act.Payload()), nil
}

func reconstructSignedRlpTxFromSig(tx rlpTransaction, chainID uint32, sig []byte) (*types.Transaction, error) {
	if len(sig) != 65 {
		return nil, errors.Errorf("invalid signature length = %d, expecting 65", len(sig))
	}
	sc := make([]byte, 65)
	copy(sc, sig)
	if sc[64] >= 27 {
		sc[64] -= 27
	}

	rawTx, err := RlpToEthTx(tx)
	if err != nil {
		return nil, err
	}
	signedTx, err := rawTx.WithSignature(types.NewEIP155Signer(big.NewInt(int64(chainID))), sc)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

// DecodeRawTx decodes raw data string into eth tx
func DecodeRawTx(rawData string, chainID uint32) (tx *types.Transaction, sig []byte, pubkey crypto.PublicKey, err error) {
	//remove Hex prefix and decode string to byte
	dataInString, err := hex.DecodeString(util.Remove0xPrefix(rawData))
	if err != nil {
		return
	}

	// decode raw data into rlp tx
	tx = &types.Transaction{}
	if err = rlp.DecodeBytes(dataInString, tx); err != nil {
		return
	}

	// extract signature
	sig, err = getSignatureFromRLPTX(tx, chainID)
	if err != nil {
		return
	}

	// recover public key
	rawHash := types.NewEIP155Signer(big.NewInt(int64(chainID))).Hash(tx)
	pubkey, err = crypto.RecoverPubkey(rawHash[:], sig)
	return
}

// EncodeRawTx encodes action into the data string of eth tx
func EncodeRawTx(act Action, pvk crypto.PrivateKey, chainID uint32) (rawData string, err error) {
	rlpAct, err := actionToRLP(act)
	if err != nil {
		return
	}
	rawTx, err := RlpToEthTx(rlpAct)
	if err != nil {
		return
	}
	ecdsaPvk, ok := pvk.EcdsaPrivateKey().(*ecdsa.PrivateKey)
	if !ok {
		err = errors.New("private key is invalid")
		return
	}

	signer := types.NewEIP155Signer(big.NewInt(int64(chainID)))
	signedTx, err := types.SignTx(rawTx, signer, ecdsaPvk)
	if err != nil {
		return
	}
	encodedTx, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return
	}
	rawData = hex.EncodeToString(encodedTx[:])
	return
}

func getSignatureFromRLPTX(tx *types.Transaction, chainID uint32) ([]byte, error) {
	if tx == nil {
		return nil, errors.New("pointer is nil")
	}
	v, r, s := tx.RawSignatureValues()

	recID := uint32(v.Int64()) - 2*chainID - 8
	sig := make([]byte, 64, 65)
	rSize := len(r.Bytes())
	copy(sig[32-rSize:32], r.Bytes())
	sSize := len(s.Bytes())
	copy(sig[64-sSize:], s.Bytes())
	sig = append(sig, byte(recID))
	return sig, nil
}
