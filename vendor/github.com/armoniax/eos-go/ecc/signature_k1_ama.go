package ecc

import (
	"github.com/armoniax/eos-go/btcsuite/btcd/btcec"
	"github.com/armoniax/eos-go/btcsuite/btcutil/base58"
)

type innerK1AMASignature struct {
}

func newInnerK1AMASignature() innerSignature {
	return &innerK1AMASignature{}
}

// verify checks the signature against the pubKey. `hash` is a sha256
// hash of the payload to verify.
func (s *innerK1AMASignature) verify(content []byte, hash []byte, pubKey PublicKey) bool {
	recoveredKey, _, err := btcec.RecoverCompact(btcec.S256(), content, hash)
	if err != nil {
		return false
	}
	key, err := pubKey.Key()
	if err != nil {
		return false
	}
	if recoveredKey.IsEqual(key) {
		return true
	}
	return false
}

func (s *innerK1AMASignature) publicKey(content []byte, hash []byte) (out PublicKey, err error) {

	recoveredKey, _, err := btcec.RecoverCompact(btcec.S256(), content, hash)

	if err != nil {
		return out, err
	}

	return PublicKey{
		Curve:   CurveK1,
		Content: recoveredKey.SerializeCompressed(),
		inner:   &innerK1AMPublicKey{},
	}, nil
}

func (s innerK1AMASignature) string(content []byte) string {
	checksum := ripemd160checksumHashCurve(content, CurveK1)
	buf := append(content[:], checksum...)
	return "SIG_K1_" + base58.Encode(buf)
}

func (s innerK1AMASignature) signatureMaterialSize() *int {
	return signatureDataSize
}
