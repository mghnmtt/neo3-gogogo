package sc

import (
	"crypto/elliptic"
	"github.com/mghnmtt/neo3-gogogo/crypto"
	"github.com/mghnmtt/neo3-gogogo/helper"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var p256 = elliptic.P256()
var G, _ = crypto.CreateECPoint(p256.Params().Gx, p256.Params().Gy, &p256)

func TestCreateContract(t *testing.T) {
	script := make([]byte, 32)
	paramList := []ContractParameterType{Signature}
	c := CreateContract(paramList, script)
	assert.Equal(t, helper.BytesToHex(script), helper.BytesToHex(c.Script))
	assert.Equal(t, 1, len(c.ParameterList))
	assert.Equal(t, Signature, c.ParameterList[0])
}

func TestContract_GetScriptHash(t *testing.T) {
	c, err := CreateSignatureContract(G)
	assert.Nil(t, err)

	expectedArray := make([]byte, 40)
	expectedArray[0] = byte(PUSHDATA1)
	expectedArray[1] = 0x21
	tmp := G.EncodePoint(true)
	assert.Equal(t, 33, len(tmp))

	for i := 0; i < len(tmp); i++ {
		expectedArray[i+2] = tmp[i]
	}
	expectedArray[35] = byte(SYSCALL)
	tmp = helper.UInt32ToBytes(uint32(System_Crypto_CheckSig.ToInteropMethodHash()))
	assert.Equal(t, 4, len(tmp))

	for i := 0; i < len(tmp); i++ {
		expectedArray[i+36] = tmp[i]
	}
	assert.Equal(t, (crypto.BytesToScriptHash(expectedArray)).String(), c.GetScriptHash().String())
}

func TestByteSlice_GetVarSize(t *testing.T) {
	b := helper.HexToBytes("deadbeef")
	size := ByteSlice(b).GetVarSize()
	assert.Equal(t, 5, size)
}

func TestByteSlice_IsSignatureContract(t *testing.T) {
	script, err := CreateSignatureRedeemScript(G)
	assert.Nil(t, err)

	b := ByteSlice(script).IsSignatureContract()
	assert.Equal(t, true, b)
}

func TestByteSlice_IsMultiSigContract(t *testing.T) {
	pubKeys1 := make([]crypto.ECPoint, 20)
	for i := 0; i < 20; i++ {
		pubKeys1[i] = *G
	}
	script1, err := CreateMultiSigRedeemScript(20, pubKeys1)
	assert.Nil(t, err)
	b1, m1, n1 := ByteSlice(script1).IsMultiSigContractWithCounts()
	assert.Equal(t, true, b1)
	assert.Equal(t, 20, m1)
	assert.Equal(t, 20, n1)

	pubKeys2 := make([]crypto.ECPoint, 256)
	for i := 0; i < 256; i++ {
		pubKeys2[i] = *G
	}
	script2, err := CreateMultiSigRedeemScript(4, pubKeys2)
	assert.Nil(t, err)
	b2, m2, n2 := ByteSlice(script2).IsMultiSigContractWithCounts()
	assert.Equal(t, true, b2)
	assert.Equal(t, 4, m2)
	assert.Equal(t, 256, n2)

	pubKeys3 := make([]crypto.ECPoint, 256)
	for i := 0; i < 256; i++ {
		pubKeys3[i] = *G
	}
	script3, err := CreateMultiSigRedeemScript(4, pubKeys3)
	assert.Nil(t, err)
	script3[len(script3)-1] = 0x00
	b3, m3, p := ByteSlice(script3).IsMultiSigContractWithPoints()
	assert.Equal(t, false, b3)
	assert.Equal(t, 0, m3)
	assert.Nil(t, p)
}

func TestCreateMultiSigRedeemScript(t *testing.T) {
	pubKeyStrs := []string{
		"02cb6c347c568adc3cb749c84bbc7af26ea2d2347a1e8ab7ff33a505e36cfd3038",
		"029e1bcfaeb3a8463eb69bca9bc718b0001c4c5ff3074e4d53ff4e26e89da0b020",
		"02c9b2ae96f937dcea7d010286db14d8b052043e1b09614a4d4eadca74f57160d6",
		"032017f4b14afeb47e49f3bd32a9b84366cb3212b4d8de46934858c896ca6e299e",
	}
	pubKeys := make([]crypto.ECPoint, 4)
	for i, v := range pubKeyStrs {
		pk, err := crypto.NewECPointFromString(v)
		assert.Nil(t, err)
		pubKeys[i] = *pk
	}
	script, err := CreateMultiSigRedeemScript(3, pubKeys)
	assert.Nil(t, err)
	log.Println(helper.BytesToHex(script))
	sh := helper.UInt160FromBytes(crypto.Hash160(script))
	log.Println(sh.String())
}
