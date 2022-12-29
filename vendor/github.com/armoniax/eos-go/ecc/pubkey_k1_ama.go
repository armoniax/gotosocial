package ecc

import (
	"fmt"
	"github.com/armoniax/eos-go/btcsuite/btcd/btcec"
)

type innerK1AMPublicKey struct {
}

func newInnerK1AMPublicKey() innerPublicKey {
	return &innerK1AMPublicKey{}
}

func (p *innerK1AMPublicKey) key(content []byte) (*btcec.PublicKey, error) {
	key, err := btcec.ParsePubKey(content, btcec.S256())
	if err != nil {
		return nil, fmt.Errorf("parsePubKey: %w", err)
	}

	return key, nil
}

func (p *innerK1AMPublicKey) prefix() string {
	return PublicKeyAMPrefix
}

func (p *innerK1AMPublicKey) keyMaterialSize() *int {
	return publicKeyDataSize
}
