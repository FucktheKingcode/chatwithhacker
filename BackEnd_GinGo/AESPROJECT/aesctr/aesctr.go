package aesctr

import (
	"io/ioutil"
	"os"
	"aesproject/createKeyAES"
	"aesproject/encryptionAES"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strings"
)

func ByteArrayToString(bytes []byte) string {
    return string(bytes)
}

func splitDataIntoBlocks(fileData *os.File, blockSize int) ([]string, error) {
	// Đọc dữ liệu từ file
	data, err := ioutil.ReadAll(fileData)
	if err != nil {
		return nil, fmt.Errorf("Failed to read data from file: %v", err)
	}
	
	dataBytes := ByteArrayToString(data)

	// Thêm padding cho dữ liệu nếu cần
	dataBytes = addPadding(dataBytes, blockSize)

	// Chia dữ liệu thành các khối có cùng kích thước
	var blocks []string
	for i := 0; i < len(dataBytes); i += blockSize {
		end := i + blockSize
		if end > len(dataBytes) {
			end = len(dataBytes)
		}
		block := dataBytes[i:end]
		blocks = append(blocks, string(block))
	}

	return blocks, nil
}



func addPadding(data string, blockSize int) string {
	paddingSize := blockSize - (len(data) % blockSize)
	padding := strings.Repeat(string(paddingSize), paddingSize)
	return data + padding
}

func generateIV(ivSize int) ([]byte, error) {
	iv := make([]byte, ivSize)
	_, err := rand.Read(iv)
	if err != nil {
		return nil, err
	}
	return iv, nil
}

func generateCounter(iv []byte, counter uint64) []byte {
	ctr := make([]byte, 16) // Kích thước giá trị đếm (128 bit)

	// Sao chép IV ban đầu vào ctr
	copy(ctr, iv)

	// Tăng giá trị đếm trong ctr
	binary.BigEndian.PutUint64(ctr[len(ctr)-8:], counter)

	return ctr
}

func convertBytesToStringArray(input [16]byte) [][]string {
	output := make([][]string, 4) // Tạo một mảng 2 chiều 4x4
	for i := range output {
		output[i] = make([]string, 4)
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			index := (i * 4) + j
			output[i][j] = fmt.Sprintf("%02X", input[index])
		}
	}

	return output
}

func addColumnToMatrix(matrix [][]string, column []string) [][]string {
	// Độ dài hiện tại của mảng 2 chiều
	numRows := len(matrix)
	numColumns := len(matrix[0])

	// Kiểm tra xem số hàng của mảng 2 chiều và số phần tử của cột có khớp hay không
	if numRows != len(column) {
		fmt.Println("Số hàng của mảng 2 chiều và số phần tử của cột không khớp.")
		return matrix
	}

	// Tạo mảng mới có số cột lớn hơn mảng ban đầu
	newMatrix := make([][]string, numRows)
	for i := range newMatrix {
		newMatrix[i] = make([]string, numColumns+1)
	}

	// Sao chép các giá trị từ mảng ban đầu sang mảng mới
	for i := 0; i < numRows; i++ {
		for j := 0; j < numColumns; j++ {
			newMatrix[i][j] = matrix[i][j]
		}
	}

	// Thêm cột mới vào mảng mới
	for i := 0; i < numRows; i++ {
		newMatrix[i][numColumns] = column[i]
	}

	return newMatrix
}

func byteArrayToByteString(byteArray []byte) string {
	byteString := fmt.Sprintf("%X", byteArray)

	return byteString
}

func Aesctr(plaintextFile *os.File, key string) ([]string, error) {
	keySize := 16 // Kích thước khóa (128 bit)
	iv, err := generateIV(keySize)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate IV: %v", err)
	}

	// Tạo giá trị đếm
	counter := uint64(0) // Giá trị đếm ban đầu

	ctr := generateCounter(iv, counter)
	blockCounter := 0
	result := []string{}

	// Chia File thành các chuỗi có kích thước phù hợp (128 bit)
	blocks, err := splitDataIntoBlocks(plaintextFile, keySize)
	if err != nil {
		return nil, fmt.Errorf("Failed to split data into blocks: %v", err)
	}

	for _, block := range blocks {
		if blockCounter == 0 {
			ctr = generateCounter(iv, counter)
		} else {
			// Tăng giá trị đếm và tạo lại Counter cho khối dữ liệu tiếp theo
			counter++
			ctr = generateCounter(iv, counter)
		}

		// Tạo keystream từ khóa AES và giá trị của bộ đếm
		keystream := encryptionAES.EncryptionAES(key, ctr)
		// Chuyển đổi thành mảng 4x4 kiểu string
		keystreamMatrix := convertBytesToStringArray([16]byte(keystream))
		// Tạo mảng 4x4 từ plaintext:
		blockMatrix := createKeyAES.StringToMatrix(block)
		// XOR 2 ma trận với nhau:
		ciphertextMatrix := encryptionAES.XorRound(keystreamMatrix, blockMatrix)
		// Tạo chuỗi byte ciphertext
		ciphertext := encryptionAES.HexArrayToCiphertext(ciphertextMatrix)

		result = append(result, byteArrayToByteString(ciphertext))

		blockCounter++
	}

	return result, nil
}

