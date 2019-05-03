// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// 4개의 int 를 encrypt 해서 key 를 생성하고, decrypt 해서 원 value를 얻을수 있는 serial key 생성기
package coupon

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/kasworld/uuidstr"
)

/* sample key or iv
eb6136fcf2e6ee6e3d5a03d85a82ad33
3600b48daa74fc83b791f08c149e9f63
57d470e8fb473901d04071859717f47e
82f1a60f780147129c695abd6a7772ce
3da24892255682460b5147a36926a982
21950280f6e9b33c408c3be2c6637a2e
eecccf51486975f4724e6b90e5d0cdb8
9c90b8a17384757c828619b920277803
7350bf25e27f85be0a30eeb5f8567493
bcddace2409e9c504279b5a5ffcc6aec
*/

type Coupon struct {
	key              []byte
	block            cipher.Block
	iv               []byte
	encryptBlockMode cipher.BlockMode
	decryptBlockMode cipher.BlockMode
}

func New(hexkey string, hexiv string) (*Coupon, error) {
	key, keyerr := hex.DecodeString(hexkey)
	if keyerr != nil {
		return nil, keyerr
	}
	iv, iverr := hex.DecodeString(hexiv)
	if iverr != nil {
		return nil, iverr
	}
	if len(key) != aes.BlockSize || len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("len key(%v) or len iv(%v) is not %v",
			len(key), len(iv), aes.BlockSize,
		)
	}
	rtn := &Coupon{
		key: []byte(key),
		iv:  []byte(iv),
	}

	block, err := aes.NewCipher(rtn.key)
	if err != nil {
		return nil, err
	}
	rtn.block = block

	rtn.encryptBlockMode = cipher.NewCBCEncrypter(block, rtn.iv)
	rtn.decryptBlockMode = cipher.NewCBCDecrypter(block, rtn.iv)

	return rtn, nil
}

// sender , receiver , sendtime, valid time
func (cpn *Coupon) Generate(uuid1, uuid2 string, i1, i2 int64) (string, error) {
	buf := new(bytes.Buffer)

	if uu1 := uuidstr.Parse(uuid1); uu1 != nil {
		buf.Write(uu1)
	} else {
		return "", fmt.Errorf("invalid uuidstr %v", uuid1)
	}
	if uu2 := uuidstr.Parse(uuid2); uu2 != nil {
		buf.Write(uu2)
	} else {
		return "", fmt.Errorf("invalid uuidstr %v", uuid2)
	}

	if err := binary.Write(buf, binary.LittleEndian, i1); err != nil {
		return "", fmt.Errorf("binary.Write failed %v %v", i1, err)
	}
	if err := binary.Write(buf, binary.LittleEndian, i2); err != nil {
		return "", fmt.Errorf("binary.Write failed %v %v", i2, err)
	}

	enced, err := cpn.encrypt(buf.Bytes())
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(enced), nil
}

func (cpn *Coupon) encrypt(plaintext []byte) ([]byte, error) {
	if len(plaintext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("plaintext is not a multiple of the block size")
	}
	ciphertext := make([]byte, len(plaintext))
	cpn.encryptBlockMode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func (cpn *Coupon) Parse(c string) (string, string, int64, int64, error) {
	cipherbytes, _ := hex.DecodeString(c)

	ori, err := cpn.decrypt(cipherbytes)
	if err != nil {
		return "", "", 0, 0, err
	}

	uuid1 := uuidstr.ToString(ori[0:16])
	uuid2 := uuidstr.ToString(ori[16:32])

	var i1, i2 int64
	buf := bytes.NewReader(ori[32:])

	if err := binary.Read(buf, binary.LittleEndian, &i1); err != nil {
		return uuid1, uuid2, i1, i2, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &i2); err != nil {
		return uuid1, uuid2, i1, i2, err
	}

	return uuid1, uuid2, i1, i2, nil
}

func (cpn *Coupon) decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}
	cpn.decryptBlockMode.CryptBlocks(ciphertext, ciphertext)
	return ciphertext, nil
}

func GenerateRandomString() (string, error) {
	randString := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, randString); err != nil {
		return "", err
	}
	return hex.EncodeToString(randString), nil
}
