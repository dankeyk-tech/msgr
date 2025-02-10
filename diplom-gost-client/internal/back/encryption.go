package back

import "github.com/Theo730/gogost/gost3412128"

func EncryptK(plainText []byte, password string) []byte {
	cryptoKuz := gost3412128.NewCipher([]byte(password))
	blockSize := cryptoKuz.BlockSize()

	dst := []byte{}
	bufdst := make([]byte, blockSize)
	for n := 0; n < len(plainText); n += blockSize {
		if n+blockSize > len(plainText) {
			cryptoKuz.Encrypt(bufdst, plainText[n:])
			dst = append(dst, bufdst...)
			continue
		}
		cryptoKuz.Encrypt(bufdst, plainText[n:n+blockSize])
		dst = append(dst, bufdst...)
	}

	return dst
}

func DecryptK(cipherText []byte, password string) []byte {
	cryptoKuz := gost3412128.NewCipher([]byte(password))
	blockSize := cryptoKuz.BlockSize()

	dst := []byte{}
	bufdst := make([]byte, blockSize)
	for n := 0; n < len(cipherText); n += blockSize {
		if n+blockSize > len(cipherText) {
			cryptoKuz.Decrypt(bufdst, cipherText[n:])
			dst = append(dst, bufdst...)
			continue
		}
		cryptoKuz.Decrypt(bufdst, cipherText[n:n+blockSize])
		dst = append(dst, bufdst...)
	}

	return dst
}
