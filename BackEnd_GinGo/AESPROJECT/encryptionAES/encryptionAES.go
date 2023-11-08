package encryptionAES

import (
	"aesproject/createKeyAES"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func XorRound(a [][]string, b [][]string) [][]string {
	rows := len(a)
	columns := len(a[0])

	// Tạo mảng kết quả với cùng số hàng và cột như mảng a
	result := make([][]string, rows)
	for i := range result {
		result[i] = make([]string, columns)
	}

	// Thực hiện phép XOR giữa các mảng
	for column := 0; column < columns; column++ {
		for row := 0; row < rows; row++ {
			aBinary, _ := strconv.ParseInt(a[row][column], 16, 64)
			bBinary, _ := strconv.ParseInt(b[row][column], 16, 64)
			xorResult := aBinary ^ bBinary
			result[row][column] = fmt.Sprintf("%056s", strconv.FormatInt(xorResult, 2))
			hex, err := createKeyAES.ConvertBinaryToHex(result[row][column])
			if err != nil {
				return nil
			}
			result[row][column] = strings.ToUpper(hex)
			if len(result[row][column]) == 1 {
				result[row][column] = "0" + result[row][column]
			}
		}
	}

	return result
}

func shiftRow(matrix [][]byte) {
	for i := 1; i < len(matrix); i++ {
		shiftedRow := make([]byte, len(matrix[i]))
		copy(shiftedRow, matrix[i])

		shiftedRow = shiftRowLeft(shiftedRow, i)

		matrix[i] = shiftedRow
	}
}

func shiftRowLeft(row []byte, shiftAmount int) []byte {
	for i := 0; i < shiftAmount; i++ {
		firstElement := row[0]
		copy(row, row[1:])
		row[len(row)-1] = firstElement
	}

	return row
}

func hexStringToByte(stringArr [][]string) [][]byte {
	byteArr := make([][]byte, len(stringArr))

	for i, row := range stringArr {
		byteArr[i] = make([]byte, len(row))

		for j, byteVal := range row {
			intValue, _ := strconv.ParseUint(byteVal, 16, 8)
			byteValue := byte(intValue)
			byteArr[i][j] = byteValue
		}
	}
	return byteArr
}

// Hàm chuyển mảng hai chiều byte sang mảng hai chiều string
func byteArrayToString(byteArr [][]byte) [][]string {
	strArr := make([][]string, len(byteArr))

	for i, row := range byteArr {
		strArr[i] = make([]string, len(row))

		for j, byteVal := range row {
			strArr[i][j] = ByteToHexString(byteVal)
		}
	}

	return strArr
}

// Hàm chuyển đổi giá trị byte sang chuỗi hex
func ByteToHexString(byteVal byte) string {
	return fmt.Sprintf("%02X", byteVal)
}

// Hằng số ma trận MixColumns
var mixColumnsMatrix = [][]byte{
	{0x02, 0x03, 0x01, 0x01},
	{0x01, 0x02, 0x03, 0x01},
	{0x01, 0x01, 0x02, 0x03},
	{0x03, 0x01, 0x01, 0x02},
}

// Hàm MixColumn
func mixColumn(column []byte) []byte {
	result := make([]byte, 4)

	for i := 0; i < 4; i++ {
		var sum byte

		for j := 0; j < 4; j++ {
			mult := galoisMult(column[j], mixColumnsMatrix[i][j])
			sum ^= mult
		}

		result[i] = sum
	}

	return result
}

// Hàm nhân trong trường Galois (GF(2^8))
func galoisMult(a, b byte) byte {
	var result byte

	for i := 0; i < 8; i++ {
		if (b & 1) == 1 {
			result ^= a
		}

		highBit := a & 0x80
		a <<= 1

		if highBit == 0x80 {
			a ^= 0x1B
		}

		b >>= 1
	}

	return result
}

func applyMixColumns(matrix [][]byte) [][]byte {
	rows := len(matrix)
	columns := len(matrix[0])
	result := make([][]byte, rows)

	for row := 0; row < rows; row++ {
		result[row] = make([]byte, columns)
		for column := 0; column < columns; column++ {
			result[row][column] = matrix[row][column]
		}
	}

	for column := 0; column < columns; column++ {
		columnArr := make([]byte, rows)
		for row := 0; row < rows; row++ {
			columnArr[row] = matrix[row][column]
		}
		columnArr = mixColumn(columnArr)
		for row := 0; row < rows; row++ {
			result[row][column] = columnArr[row]
		}
	}

	return result
}

func HexArrayToCiphertext(array [][]string) []byte {
	ciphertext := make([]byte, 16)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			bytes, _ := hex.DecodeString(array[i][j])
			copy(ciphertext[i*4+j:], bytes)
		}
	}

	return ciphertext
}

func create2DArray(bytes []byte, rows, cols int) [][]string {
	// Tạo mảng 2 chiều có kích thước rows x cols
	twoDArray := make([][]string, rows)
	for i := range twoDArray {
		twoDArray[i] = make([]string, cols)
	}

	// Chuyển các phần tử từ chuỗi byte vào mảng 2 chiều
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			index := i*cols + j
			twoDArray[i][j] = fmt.Sprintf("%02X", bytes[index])
		}
	}

	return twoDArray
}

func EncryptionAES(key string, nonce []byte) []byte {
	// Tạo mảng 2 chiều cho nonce
	matrix0 := create2DArray(nonce, 4, 4)
	// Tạo key cho Round 0
	matrixKey := createKeyAES.CreateKeyAES(key, 0)
	matrixRound0 := XorRound(matrix0, matrixKey)

	// Chạy Round 1:
	// Sub - byte
	matrixSubBytes := createKeyAES.ConvertArray(matrixRound0, createKeyAES.SBox)
	// Chuyển mảng sub - byte  sang kiểu byte
	byteValue := hexStringToByte(matrixSubBytes)
	// Shif row
	shiftRow(byteValue)
	// Mix Colum
	result := applyMixColumns(byteValue)
	// Chuyển mảng kết quả về kiểu string
	resultString := byteArrayToString(result)
	// Tạo key cho Round
	matrixKey = createKeyAES.CreateKeyAES(key, 1)
	// Add Round Key
	matrixRound := XorRound(resultString, matrixKey)

	for i := 1; i < 10; i++ {
		// Sub - byte
		matrixSubBytes = createKeyAES.ConvertArray(matrixRound, createKeyAES.SBox)
		// Chuyển mảng sub - byte  sang kiểu byte
		byteValue = hexStringToByte(matrixSubBytes)
		// Shif row
		shiftRow(byteValue)
		// Mix Colum
		result = applyMixColumns(byteValue)
		// Chuyển mảng kết quả về kiểu string
		resultString = byteArrayToString(result)
		// Tạo key cho Round
		matrixKey = createKeyAES.CreateKeyAES(key, i)
		// Add Round Key
		matrixRound = XorRound(resultString, matrixKey)
	}

	// Sub - byte
	matrixSubBytes = createKeyAES.ConvertArray(matrixRound, createKeyAES.SBox)
	// Chuyển mảng sub - byte  sang kiểu byte
	byteValue = hexStringToByte(matrixSubBytes)
	// Shif row
	shiftRow(byteValue)
	// Chuyển mảng kết quả về kiểu string
	resultString = byteArrayToString(byteValue)
	// Tạo key cho Round
	matrixKey = createKeyAES.CreateKeyAES(key, 10)
	// Add Round Key
	matrixRound = XorRound(resultString, matrixKey)

	// Chuyển đổi mảng thành chuỗi Ciphertext 128 bit
	ciphertext := HexArrayToCiphertext(matrixRound)

	return ciphertext
}
