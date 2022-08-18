package wallet

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/startfellows/tongo"
	"github.com/startfellows/tongo/boc"
	"github.com/startfellows/tongo/tlb"
	"testing"
)

func TestGetCodeByVer(t *testing.T) {
	for ver := V1R1; ver < HighLoadV2; ver++ {
		_ = GetCodeByVer(ver)
	}
}

func TestVersionToString(t *testing.T) {
	testData := map[Version]string{
		V1R1:       "v1R1",
		V3R1:       "v3R1",
		V3R2:       "v3R2",
		V4R1:       "v4R1",
		V4R2:       "v4R2",
		HighLoadV2: "highload_v2",
	}
	for ver, name := range testData {
		if ver.ToString() != name {
			t.Fatal("invalid mapping version to string")
		}
	}
}

func TestGenerateWalletAddress(t *testing.T) {
	type walletData struct {
		Address   string
		PublicKey string
	}
	testData := map[Version]walletData{
		// TODO: add other versions
		V3R2: {"0:f3a069b7fc4631da4401de03eddd7cd30caca618c6ad0e3ac3fa454370b73a96",
			"f96db56e72de2e84e0aef780428e439a6c84e0b27bc2b2591075785479f2e9c3"},
		V4R1: {"0:17afeaaa61cb575e3e340a296da6bf55bc6b996cfab1d9f87840b2b6dc4cf613",
			"6f58b9fecb87e847825a7ecf3ae1f32b5578eee156ac10b398e2f1d67c12ca05"},
		V4R2: {"0:8f2983152d1480ba6af25e087d672232080b294dc8992525e35e4ff6d601f405",
			"7843fd9de6cd858154d9a914b8c3cd0bf1dc5af3a0c1dd273586568fc4d1c002"},
	}
	for ver, data := range testData {
		key, _ := hex.DecodeString(data.PublicKey)
		publicKey := ed25519.PublicKey(key)
		address, err := GenerateWalletAddress(publicKey, ver, 0, nil)
		if err != nil {
			t.Fatalf("address generation failed: %v", err)
		}
		if address.ToRaw() != data.Address {
			t.Fatal("addresses mismatch")
		}
	}
}

func TestLongCommentSerialization(t *testing.T) {
	// TODO: add real serialized data
	longComment := `
		The Quick Brown Fox Jumps Over The Lazy Dog
		The Quick Brown Fox Jumps Over The Lazy Dog
		The Quick Brown Fox Jumps Over The Lazy Dog
		The Quick Brown Fox Jumps Over The Lazy Dog
		The Quick Brown Fox Jumps Over The Lazy Dog
		The Quick Brown Fox Jumps Over The Lazy Dog`

	cell := boc.NewCell()
	err := tlb.Marshal(cell, TextComment(longComment))
	if err != nil {
		t.Fatalf("long comment serialization error: %v", err)
	}
	var text TextComment
	err = tlb.Unmarshal(cell, &text)
	if err != nil {
		t.Fatalf("long comment deserialization error: %v", err)
	}
	if string(text) != longComment {
		t.Fatal("TextComment invalid serialization/deserialization")
	}
}

func TestGenerateTonTransferMessage(t *testing.T) {
	// TODO: implement for other versions
	testMessage := "b5ee9c720101040100b70001458801bd55cca31423fa4983b538a75a715dba6be94cda37de9e3c61028e5b24ed52400c01019c0f95a1f93ce745d54c2c837927f4b2f5eeff4b4119b63f465adba5c1eaf861b29d037adca6dd325d441b02f06d54f3292f3e9ade15d67effcf945c74a28da10d29a9a317ffffffff0000000d00010201644200283ef53eb037916cf42b3c69f76f1cddf099d434695071f03faa8165b85c5289113880000000000000000000000000010300120000000068656c6c6f"
	recipientAddr, _ := tongo.AccountIDFromRaw("0:507dea7d606f22d9e85678d3eede39bbe133a868d2a0e3e07f5502cb70b8a512")
	pk, _ := base64.StdEncoding.DecodeString("OyAWIb4FeP1bY1VhALWrU2JN9/8O1Kv8kWZ0WfXXpOM=")
	privateKey := ed25519.NewKeyFromSeed(pk)

	w, err := NewWallet(privateKey, V4R2, 0, nil)
	if err != nil {
		t.Fatalf("Unable to create wallet: %v", err)
	}

	tonTransfer := TonTransfer{
		Recipient: *recipientAddr,
		Amount:    10000,
		Comment:   "hello",
		Bounce:    false,
		Mode:      1,
	}

	msg, err := w.GenerateTonTransferMessage(13, 0xFFFFFFFF, []TonTransfer{tonTransfer})
	if err != nil {
		t.Fatalf("Unable to generate transfer message: %v", err)
	}
	if fmt.Sprintf("%x", msg) != testMessage {
		t.Fatalf("messages mismatch")
	}
}