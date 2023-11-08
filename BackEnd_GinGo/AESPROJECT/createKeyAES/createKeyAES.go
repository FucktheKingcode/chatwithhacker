package createKeyAES

import (
	"fmt"
	"strconv"
	"strings"
)

func charToHex(c byte) string {
	hex := fmt.Sprintf("%X", c)
	return hex
}

func StringToMatrix(input string) [][]string {
	matrix := make([][]string, 4) // Tạo một mảng 2 chiều 4x4
	for i := range matrix {
		matrix[i] = make([]string, 4)
	}

	row := 0
	col := 0
	for i := 0; i < len(input); i++ {
		matrix[row][col] = charToHex(input[i])
		row++
		if row == 4 {
			row = 0
			col++
		}
	}

	return matrix
}

func rotwork(matrix [][]string) [][]string {
	rows := len(matrix)
	columns := len(matrix[0])
	result := make([][]string, rows)
	for row := 0; row < rows; row++ {
		result[row] = make([]string, columns)
		for column := 0; column < columns; column++ {
			result[row][column] = matrix[row][column]
		}
	}

	temp := result[0][columns-1]
	for i := 0; i < 4; i++ {
		if i == 3 {
			result[i][columns-1] = temp
		} else {
			result[i][columns-1] = result[i+1][columns-1]
		}
	}
	return result
}

var SBox = [256]byte{
	0x63, 0x7C, 0x77, 0x7B, 0xF2, 0x6B, 0x6F, 0xC5, 0x30, 0x01, 0x67, 0x2B, 0xFE, 0xD7, 0xAB, 0x76,
	0xCA, 0x82, 0xC9, 0x7D, 0xFA, 0x59, 0x47, 0xF0, 0xAD, 0xD4, 0xA2, 0xAF, 0x9C, 0xA4, 0x72, 0xC0,
	0xB7, 0xFD, 0x93, 0x26, 0x36, 0x3F, 0xF7, 0xCC, 0x34, 0xA5, 0xE5, 0xF1, 0x71, 0xD8, 0x31, 0x15,
	0x04, 0xC7, 0x23, 0xC3, 0x18, 0x96, 0x05, 0x9A, 0x07, 0x12, 0x80, 0xE2, 0xEB, 0x27, 0xB2, 0x75,
	0x09, 0x83, 0x2C, 0x1A, 0x1B, 0x6E, 0x5A, 0xA0, 0x52, 0x3B, 0xD6, 0xB3, 0x29, 0xE3, 0x2F, 0x84,
	0x53, 0xD1, 0x00, 0xED, 0x20, 0xFC, 0xB1, 0x5B, 0x6A, 0xCB, 0xBe, 0x39, 0x4A, 0x4C, 0x58, 0xCF,
	0xD0, 0xEF, 0xAA, 0xFB, 0x43, 0x4D, 0x33, 0x85, 0x45, 0xF9, 0x02, 0x7F, 0x50, 0x3C, 0x9F, 0xA8,
	0x51, 0xA3, 0x40, 0x8F, 0x92, 0x9D, 0x38, 0xF5, 0xBC, 0xB6, 0xDA, 0x21, 0x10, 0xFF, 0xF3, 0xD2,
	0xCD, 0x0C, 0x13, 0xEC, 0x5F, 0x97, 0x44, 0x17, 0xC4, 0xA7, 0x7E, 0x3D, 0x64, 0x5D, 0x19, 0x73,
	0x60, 0x81, 0x4F, 0xDC, 0x22, 0x2A, 0x90, 0x88, 0x46, 0xEE, 0xB8, 0x14, 0xDE, 0x5E, 0x0B, 0xDB,
	0xE0, 0x32, 0x3A, 0x0A, 0x49, 0x06, 0x24, 0x5C, 0xC2, 0xD3, 0xAC, 0x62, 0x91, 0x95, 0xE4, 0x79,
	0xE7, 0xC8, 0x37, 0x6D, 0x8D, 0xD5, 0x4E, 0xA9, 0x6C, 0x56, 0xF4, 0xEA, 0x65, 0x7A, 0xAE, 0x08,
	0xBA, 0x78, 0x25, 0x2E, 0x1C, 0xA6, 0xB4, 0xC6, 0xE8, 0xDD, 0x74, 0x1F, 0x4B, 0xBD, 0x8B, 0x8A,
	0x70, 0x3E, 0xB5, 0x66, 0x48, 0x03, 0xF6, 0x0E, 0x61, 0x35, 0x57, 0xB9, 0x86, 0xC1, 0x1D, 0x9E,
	0xE1, 0xF8, 0x98, 0x11, 0x69, 0xD9, 0x8E, 0x94, 0x9B, 0x1E, 0x87, 0xE9, 0xCE, 0x55, 0x28, 0xDF,
	0x8C, 0xA1, 0x89, 0x0D, 0xBF, 0xE6, 0x42, 0x68, 0x41, 0x99, 0x2D, 0x0F, 0xB0, 0x54, 0xBB, 0x16,
}

func ConvertArray(input [][]string, SBox [256]byte) [][]string {
	output := make([][]string, len(input))
	for i, row := range input {
		output[i] = make([]string, len(row))
		for j, val := range row {
			// Lấy giá trị tương ứng từ bảng SBox

			var rowIndex, colIndex int
			if val[0] >= '0' && val[0] <= '9' {
				rowIndex = int(val[0] - '0')
			} else {
				switch val[0] {
				case 'A':
					rowIndex = 10
				case 'B':
					rowIndex = 11
				case 'C':
					rowIndex = 12
				case 'D':
					rowIndex = 13
				case 'E':
					rowIndex = 14
				case 'F':
					rowIndex = 15
				}
			}
			if val[1] >= '0' && val[1] <= '9' {
				colIndex = int(val[1] - '0')
			} else {
				switch val[1] {
				case 'A':
					colIndex = 10
				case 'B':
					colIndex = 11
				case 'C':
					colIndex = 12
				case 'D':
					colIndex = 13
				case 'E':
					colIndex = 14
				case 'F':
					colIndex = 15
				}
			}

			SBoxIndex := rowIndex*16 + colIndex
			output[i][j] = strings.ToUpper(strconv.FormatUint(uint64(SBox[SBoxIndex]), 16))
		}
	}
	return output
}

var matrixRCON = [][]string{
	{"01", "02", "04", "08", "10", "20", "40", "80", "18", "36"},
	{"00", "00", "00", "00", "00", "00", "00", "00", "00", "00"},
	{"00", "00", "00", "00", "00", "00", "00", "00", "00", "00"},
	{"00", "00", "00", "00", "00", "00", "00", "00", "00", "00"},
}

func ConvertWordsToBinary(words []string) []string {
	binaryWords := make([]string, len(words))
	for i, word := range words {
		value, err := strconv.ParseInt(word, 16, 64)
		if err != nil {
			// Xử lý lỗi nếu có
			// Ví dụ: return nil hoặc gán giá trị mặc định
		}
		binary := strconv.FormatInt(value, 2)
		// Kiểm tra độ dài của chuỗi nhị phân
		if len(binary) < 4 {
			// Thêm các kí tự "0" vào đầu chuỗi để có đúng 4 kí tự
			binary = strings.Repeat("0", 4-len(binary)) + binary
		}
		binaryWords[i] = binary
	}
	return binaryWords
}

func printArray(arr []string) {
	for _, value := range arr {
		fmt.Println(value)
	}
}

func ConvertBinaryToHex(binary string) (string, error) {
	value, err := strconv.ParseInt(binary, 2, 64)
	if err != nil {
		return "", err
	}
	hex := strconv.FormatInt(value, 16)
	return hex, nil
}

func createKey(a [][]string, b [][]string, c [][]string, k int) [][]string {
	rows := len(a)
	columns := len(a[0])

	// Tạo mảng kết quả với cùng số hàng và cột như mảng a
	result := make([][]string, rows)
	for i := range result {
		result[i] = make([]string, columns)
	}

	// Lấy cột 1 của mảng a và chuyển từng phần tử về mã nhị phân
	aColumn1x := make([]string, rows)
	for row := 0; row < rows; row++ {
		aColumn1x[row] = a[row][0]
	}
	aColumn1 := ConvertWordsToBinary(aColumn1x)

	// Lấy cột 4 của mảng b và chuyển từng phần tử về mã nhị phân
	bColumn4x := make([]string, rows)
	for row := 0; row < rows; row++ {
		bColumn4x[row] = b[row][columns-1]
	}
	bColumn4 := ConvertWordsToBinary(bColumn4x)

	// Lấy cột 1 của mảng c và chuyển từng phần tử về mã nhị phân
	cColumn1x := make([]string, rows)
	for row := 0; row < rows; row++ {
		cColumn1x[row] = c[row][k-1]
	}
	cColumn1 := ConvertWordsToBinary(cColumn1x)

	// Thực hiện phép XOR giữa các chuỗi nhị phân
	for row := 0; row < rows; row++ {
		aBinary, _ := strconv.ParseInt(aColumn1[row], 2, 64)
		bBinary, _ := strconv.ParseInt(bColumn4[row], 2, 64)
		cBinary, _ := strconv.ParseInt(cColumn1[row], 2, 64)
		xorResult := aBinary ^ bBinary ^ cBinary
		result[row][0] = fmt.Sprintf("%08s", strconv.FormatInt(xorResult, 2))
		hex, err := ConvertBinaryToHex(result[row][0])
		if err != nil {
			return nil
		}
		result[row][0] = strings.ToUpper(hex)
		if len(result[row][0]) == 1 {
			result[row][0] = "0" + result[row][0]
		}

	}

	for column := 1; column < columns; column++ {

		// Lấy cột column hiện tại của mảng result và chuyển từng phần tử về mã nhị phân
		ColumnNx := make([]string, rows)
		for row := 0; row < rows; row++ {
			ColumnNx[row] = result[row][column-1]
		}
		ColumnN := ConvertWordsToBinary(ColumnNx)

		// Lấy cột n của mảng b và chuyển từng phần tử về mã nhị phân
		ColumnN1x := make([]string, rows)
		for row := 0; row < rows; row++ {
			ColumnN1x[row] = a[row][column]
		}
		ColumnN1 := ConvertWordsToBinary(ColumnN1x)

		// Thực hiện phép XOR giữa các chuỗi nhị phân
		for row := 0; row < rows; row++ {

			NBinary, _ := strconv.ParseInt(ColumnN[row], 2, 64)
			N1Binary, _ := strconv.ParseInt(ColumnN1[row], 2, 64)
			xorResult := NBinary ^ N1Binary
			result[row][column] = fmt.Sprintf("%08s", strconv.FormatInt(xorResult, 2))
			hex, err := ConvertBinaryToHex(result[row][column])
			if err != nil {
				return nil
			}
			result[row][column] = strings.ToUpper(hex)
			if len(result[row][column]) == 1 {
				result[row][column] = "0" + result[row][column]
			}
		}
	}

	// Trả về mảng kết quả kiểu string
	return result
}

func CreateKeyAES(input string, k int) [][]string {
	matrix := StringToMatrix(input)
	switch k {
	case 0:
		matrixKey := matrix
		return matrixKey

	case 1:
		matrixrot := rotwork(matrix)
		matrixSBox := ConvertArray(matrixrot, SBox)
		matrixKey := createKey(matrix, matrixSBox, matrixRCON, k)
		return matrixKey
	default:
		matrixrot := rotwork(matrix)
		matrixSBox := ConvertArray(matrixrot, SBox)
		matrixKey := createKey(matrix, matrixSBox, matrixRCON, 1)
		for i := 2; i <= k; i++ {
			matrixSBox = ConvertArray(matrixKey, SBox)
			matrixKey = createKey(matrixKey, matrixSBox, matrixRCON, i)
		}
		return matrixKey
	}

}
