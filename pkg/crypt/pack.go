package crypt

import (
	. "github.com/haxqer/gofunc"
	"math/rand"
)

type EData struct {
	EK []byte
	ED []byte
	IV []byte
}

type PData struct {
	Data []byte
	Key  []byte
}

func Packv2(p *PData) (*EData, error) {
	en, err := AesCBCEncrypt(p.Data, p.Key, PKCS5Padding)
	if err != nil {
		return nil, err
	}
	iv := en[:16]
	ed := en[16:]
	biv := Base64EncodeByte(iv)
	bed := Base64EncodeByte(ed)
	ek, err := GenRandomBytes(rand.Intn(10) + 120)
	if err != nil {
		ek = nil
	}
	bek := Base64EncodeByte(ek)
	return &EData{
		EK: bek,
		ED: bed,
		IV: biv,
	}, nil
}

func UnPackOnlyAES(r *EData, key []byte) (*PData, error) {
	iv, err := Base64DecodeByte(r.IV)
	if err != nil {
		return nil, err
	}
	ed, err := Base64DecodeByte(r.ED)
	if err != nil {
		return nil, err
	}
	realEd := append(iv, ed...)

	b, err := AesCBCDecrypt(realEd, key, PKCS5UnPadding)
	if err != nil {
		return nil, err
	}
	return &PData{
		Data: b,
		Key:  key,
	}, nil
}

func Unpackv2(r *EData, privateKey []byte) (*PData, error) {
	ek, err := Base64DecodeByte(r.EK)
	if err != nil {
		return nil, err
	}
	key, err := RsaDecode(ek, privateKey)
	if err != nil {
		return nil, err
	}
	iv, err := Base64DecodeByte(r.IV)
	if err != nil {
		return nil, err
	}
	ed, err := Base64DecodeByte(r.ED)
	if err != nil {
		return nil, err
	}
	realEd := append(iv, ed...)

	b, err := AesCBCDecrypt(realEd, key, PKCS5UnPadding)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%s", b)
	return &PData{
		Data: b,
		Key:  key,
	}, nil
}

func UnClientPack(ebk, bed, biv, privateKey []byte) ([]byte, error) {
	ek, err := Base64DecodeByte(ebk)
	if err != nil {
		return nil, err
	}
	key, err := RsaDecode(ek, privateKey)
	if err != nil {
		return nil, err
	}
	iv, err := Base64DecodeByte(biv)
	if err != nil {
		return nil, err
	}
	ed, err := Base64DecodeByte(bed)
	if err != nil {
		return nil, err
	}
	realEd := append(iv, ed...)

	//fmt.Printf("%d\n", len(realEd))

	b, err := AesCBCDecrypt(realEd, key, PKCS5UnPadding)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%s \n", b)
	return b, nil
}
