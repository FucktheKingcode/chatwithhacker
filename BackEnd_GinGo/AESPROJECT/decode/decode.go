package decode

import (
	"fmt"
	"math/big"
	"strings"
)

func BigIntToString(n *big.Int) string {
	return n.String()
}

func CreateBigIntArray(num1 *big.Int, num2 *big.Int) []string {
	arr := []string{BigIntToString(num1), BigIntToString(num2)}
	return arr
}

func getTwoStrings(input string) (string, string) {
	arr := strings.Split(input, " ")
	return arr[0], arr[1]
}

func processArray(arr [][]string) string {
	var s string

	for i := 0; i < len(arr)-1; i++ {
		if len(arr[i]) > 1 && len(arr[i+1]) > 1 {
			s = s + xorStrings(arr[i][0], arr[i+1][0])
		}
	}

	return s
}

func stringToArray(str string) [][]string {
	rows := strings.Split(str, " ")
	arr := make([][]string, len(rows))

	for i, row := range rows {
		arr[i] = strings.Split(row, " ")
	}

	return arr
}

func xorStrings(a, b string) string {
	result := ""

	// Đảm bảo độ dài của chuỗi a và b bằng nhau
	if len(a) != len(b) {
		panic("Độ dài của hai chuỗi không bằng nhau")
	}

	// Thực hiện phép XOR byte-to-byte giữa hai chuỗi
	for i := 0; i < len(a); i++ {
		result += string(a[i] ^ b[i])
	}

	return result
}

func decode(c *big.Int, d *big.Int, n *big.Int) *big.Int {
	var base, exponent, modulus big.Int

	// Tính toán (c^d) mod n
	base.Set(c)
	exponent.Set(d)
	modulus.Set(n)

	result := new(big.Int)
	result.Exp(&base, &exponent, &modulus)

	return result
}

func DecodeRSA(c string, privateKey string, aesM string) string {

	ciphertextKey := new(big.Int)
	ciphertextKey.SetString(c, 10)

	n, d := getTwoStrings(privateKey)

	N := new(big.Int)
	N.SetString(n, 10)

	modulus := new(big.Int)
	modulus.SetString(d, 10)

	m := decode(ciphertextKey, modulus, N)
	mStr := BigIntToString(m)
	resultStr := xorStrings(aesM, mStr)
	fmt.Printf(aesM)
	fmt.Printf(mStr)

	resultArr := stringToArray(resultStr)
	plaintext := processArray(resultArr)

	return plaintext
}
