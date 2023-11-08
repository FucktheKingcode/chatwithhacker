package rsaAES

import (
	"io"
	"path/filepath"
	"archive/zip"
	"os"
	"io/ioutil"
	"aesproject/aesctr"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// func main() {
// 	result := aesctr.Aesctr("1234567890qwertyuiopasdfghjklzxc", "1234567890qwertyuiopasdfghjklzxc")
// 	// In kết quả
// 	for _, row := range result {
// 		fmt.Println(row)
// 	}
// }

// Kiểm tra xem một số có phải là số nguyên tố hay không
func isPrime(n *big.Int) bool {
	return n.ProbablyPrime(20)
}

// Tạo một số nguyên ngẫu nhiên có độ dài là digits chữ số
func generateRandomNumber(digits int) *big.Int {
	// Tính toán giới hạn trên và dưới của số nguyên có digits chữ số
	min := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits-1)), nil)
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)

	// Tạo một số nguyên ngẫu nhiên trong khoảng từ min đến max-1
	num, _ := rand.Int(rand.Reader, new(big.Int).Sub(max, min))

	// Thêm giá trị min để đảm bảo số nguyên có digits chữ số
	num.Add(num, min)

	return num
}

// Tìm hai số nguyên tố sao cho tích của chúng có digits chữ số
func findPrimeFactors(digits int) (*big.Int, *big.Int, *big.Int) {
	n := generateRandomNumber(digits)
	sqrtN := new(big.Int).Sqrt(n)

	// Chọn hai số nguyên tố bắt đầu từ sqrtN và kiểm tra tính nguyên tố của chúng
	p := new(big.Int).Set(sqrtN)
	for !isPrime(p) {
		p = new(big.Int).Sub(p, big.NewInt(1))
	}

	q := new(big.Int).Set(sqrtN)
	for !isPrime(q) {
		q = new(big.Int).Add(q, big.NewInt(1))
	}

	return p, q, new(big.Int).Mul(p, q)
}

func Qofn(p *big.Int, q *big.Int) *big.Int {
	one := big.NewInt(1)
	pMinusOne := new(big.Int).Sub(p, one)
	qMinusOne := new(big.Int).Sub(q, one)
	Qn := new(big.Int).Mul(pMinusOne, qMinusOne)
	return Qn
}

// Tìm số nguyên ngẫu nhiên e thỏa mãn yêu cầu
func findRandomE(max *big.Int) *big.Int {
	one := big.NewInt(1)
	e := big.NewInt(0)

	for e.Cmp(one) <= 0 || e.Cmp(max) >= 0 {
		randomBytes := make([]byte, max.BitLen()/8+1)
		_, _ = rand.Read(randomBytes)
		e = new(big.Int).SetBytes(randomBytes)
	}

	// Kiểm tra e và max có ước chung lớn nhất là 1
	for gcd(e, max).Cmp(one) != 0 {
		e.Add(e, one)
		if e.Cmp(max) >= 0 {
			e.Set(one)
		}
	}

	return e
}

// Tính ước chung lớn nhất của hai số
func gcd(a, b *big.Int) *big.Int {
	for b.Sign() != 0 {
		a, b = b, new(big.Int).Mod(a, b)
	}
	return a
}

func findD(e, n *big.Int) *big.Int {
	d := new(big.Int)
	d.ModInverse(e, n)
	return d
}

func encrypt(m, e, n *big.Int) *big.Int {
	c := new(big.Int)
	c.Exp(m, e, n)
	return c
}

func arrayToString(arr [][]string) string {
	var result []string

	for _, row := range arr {
		result = append(result, strings.Join(row, " "))
	}

	return strings.Join(result, " ")
}

func stringToArray(str string) [][]string {
	rows := strings.Split(str, " ")
	arr := make([][]string, len(rows))

	for i, row := range rows {
		arr[i] = strings.Split(row, " ")
	}

	return arr
}

func xorStrings(array []string, str string) []string {
	result := make([]string, len(array))
	for i, s := range array {
		result[i] = xorString(s, str)
	}
	return result
}

func xorString(s string, str string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		result.WriteByte(s[i] ^ str[i%len(str)])
	}
	return result.String()
}


func BigIntToString(n *big.Int) string {
	return n.String()
}

func CreateBigIntArray(num1 *big.Int, num2 *big.Int) []string {
	arr := []string{BigIntToString(num1), BigIntToString(num2)}
	return arr
}

func RsaAES(plaintextFile *os.File, key string) error {
	result, err := aesctr.Aesctr(plaintextFile, key)
	if err != nil {
		return fmt.Errorf("Failed to perform AES-CTR encryption: %v", err)
	}

	p, q, n := findPrimeFactors(618)
	if err != nil {
		return fmt.Errorf("Failed to find prime factors: %v", err)
	}

	Qn := Qofn(p, q)
	e := findRandomE(Qn)
	d := findD(e, Qn)
	m := generateRandomNumber(len(result))
	mStr := m.String()

	cipherFile := xorStrings(result, mStr)
	cipherKey := BigIntToString(encrypt(m, e, n))
	publicKey := CreateBigIntArray(n, e)
	privateKey := CreateBigIntArray(n, d)

	// Create directory to store files
	err = os.MkdirAll("output", 0755)
	if err != nil {
		return fmt.Errorf("Failed to create output directory: %v", err)
	}

	// Write files
	err = ioutil.WriteFile("output/cipherFile.txt", []byte(strings.Join(cipherFile, "")), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write cipherFile.txt: %v", err)
	}

	err = ioutil.WriteFile("output/cipherKey.txt", []byte(cipherKey), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write cipherKey.txt: %v", err)
	}

	err = ioutil.WriteFile("output/publicKey.txt", []byte(strings.Join(publicKey, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write publicKey.txt: %v", err)
	}

	err = ioutil.WriteFile("output/privateKey.txt", []byte(strings.Join(privateKey, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write privateKey.txt: %v", err)
	}

	// Create zip file
	zipFile, err := os.Create("output.zip")
	if err != nil {
		return fmt.Errorf("Failed to create output.zip: %v", err)
	}
	defer zipFile.Close()

	// Create a new zip archive
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add files to the zip archive
	files := []string{"output/cipherFile.txt", "output/cipherKey.txt", "output/publicKey.txt", "output/privateKey.txt"}
	for _, file := range files {
		err = addFileToZip(zipWriter, file)
		if err != nil {
			return fmt.Errorf("Failed to add file to zip: %v", err)
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	// Get file info to set attributes in the zip entry
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Failed to get file info for %s: %v", filename, err)
	}

	// Create a new zip file entry
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("Failed to create zip header for %s: %v", filename, err)
	}

	// Set the name of the file inside the zip archive
	header.Name = filepath.Base(filename)

	// Create the zip file entry
	entry, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("Failed to create zip entry for %s: %v", filename, err)
	}

	// Copy the file content to the zip entry
	_, err = io.Copy(entry, file)
	if err != nil {
		return fmt.Errorf("Failed to copy file content to zip entry for %s: %v", filename, err)
	}

	return nil
}
