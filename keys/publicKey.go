package keys

import (
	"github.com/mghnmtt/neo3-gogogo/crypto"
	"github.com/mghnmtt/neo3-gogogo/helper"
	"github.com/mghnmtt/neo3-gogogo/sc"
)

// returns a NEO-specific hash of the public key.
func PublicKeyToScriptHash(p *crypto.ECPoint) *helper.UInt160 {
	b, _ := sc.CreateSignatureRedeemScript(p)
	hash := crypto.Hash160(b)
	return helper.UInt160FromBytes(hash)
}

// returns a base58-encoded NEO-specific address based on the key hash.
func PublicKeyToAddress(p *crypto.ECPoint, version byte) string {
	return crypto.ScriptHashToAddress(PublicKeyToScriptHash(p), version)
}
