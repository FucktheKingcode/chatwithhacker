package main

import (
	"path/filepath"
	"io"
	"archive/zip"
	"io/ioutil"
	"os"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"math/rand"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type EncryptionResult struct {
	Hash          string   `json:"hash"`
	EncryptedData string   `json:"encrypted_data"`
	PublicKey     []string `json:"public_key"`
	PrivateKey    []string `json:"private_key"`
	EncryptionKey string   `json:"encryption_key"`
}

type DecryptRequest struct {
	Hash       string `json:"hash"`
	AesM       string `json:"aesM"`
	PrivateKey string `json:"privateKey"`
	CipherKey  string `json:"cipherKey"`
}

func HashString(text string) string {
	hash := sha256.Sum256([]byte(text))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

func generateRandomKey(length int) string {
	rand.Seed(time.Now().UnixNano())

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func main() {
	r := gin.Default()

	// Áp dụng middleware Cors
	r.Use(cors.Default())

	r.POST("/api/mahoadulieu", func(c *gin.Context) {
		// Lấy fileData và filePublicKey từ request
		fileData, err := c.FormFile("fileData")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filePublicKey, err := c.FormFile("filePublicKey")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Tạo thư mục tạm để lưu trữ file
		tempDir, err := ioutil.TempDir("", "temp")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer os.RemoveAll(tempDir)

		// Lưu trữ file mã hóa vào thư mục tạm
		fileEncryptionPath := filepath.Join(tempDir, fileData.Filename)
		if err := c.SaveUploadedFile(fileData, fileEncryptionPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Lưu trữ file 2 vào thư mục tạm
		fileHashPath := filepath.Join(tempDir, filePublicKey.Filename)
		if err := c.SaveUploadedFile(filePublicKey, fileHashPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Tạo file zip
		zipFilePath := filepath.Join(tempDir, "output.zip")
		err = createZipFile(zipFilePath, fileEncryptionPath, fileHashPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Trả về file zip
		c.Header("Content-Disposition", "attachment; filename=output.zip")
		c.Header("Content-Type", "application/zip")

		zipFile, err := os.Open(zipFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer zipFile.Close()

		if _, err := io.Copy(c.Writer, zipFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	})

	r.POST("/api/mahoadulieu/giaima", func(c *gin.Context) {
		// Lấy file zip từ request
		fileZip, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		// Lấy file privateKey từ request
		filePrivateKey, err := c.FormFile("privateKey")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		// Tạo thư mục tạm để lưu trữ file
		tempDir, err := ioutil.TempDir("", "temp")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer os.RemoveAll(tempDir)
	
		// Lưu trữ file zip vào thư mục tạm
		zipFilePath := filepath.Join(tempDir, fileZip.Filename)
		if err := c.SaveUploadedFile(fileZip, zipFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		// Lưu trữ file privateKey vào thư mục tạm
		privateKeyFilePath := filepath.Join(tempDir, filePrivateKey.Filename)
		if err := c.SaveUploadedFile(filePrivateKey, privateKeyFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		// Đường dẫn tới file output
		outputFilePath := filepath.Join(tempDir, "filedata.txt")
	
		// Giải mã và giải nén file zip
		err = unzipAndDecryptFile(zipFilePath, outputFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		// Trả về file data giải mã
		c.File(outputFilePath)
	})
	
	

	// Run the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// Hàm giải mã và giải nén file zip
func unzipAndDecryptFile(zipFilePath string, outputFilePath string) error {
	// Mở file zip
	zipFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Duyệt qua các file trong zip
	for _, file := range zipFile.File {
		if file.Name == "fileEncryption.txt" {
			// Mở file fileData.txt
			fileEncryption, err := file.Open()
			if err != nil {
				return err
			}
			defer fileEncryption.Close()

			// Đọc nội dung của fileEncryption.txt
			fileDataEncryption, err := ioutil.ReadAll(fileEncryption)
			if err != nil {
				return err
			}

			// Giải mã dữ liệu
			decryptedData, err := giaiMaData(string(fileDataEncryption))
			if err != nil {
				return err
			}

			// Lưu dữ liệu giải mã vào file output
			err = ioutil.WriteFile(outputFilePath, []byte(decryptedData), 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


// Hàm thêm file vào zip
func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Lấy thông tin file
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Tạo header cho file trong zip
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Đặt tên file trong zip
	header.Name = filepath.Base(filePath)

	// Tạo entry mới trong zip
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Copy nội dung file vào entry trong zip
	if _, err := io.Copy(writer, file); err != nil {
		return err
	}

	return nil
}

// Hàm giải mã file zip và trả về nội dung file txt
func giaiMaData(fileZipPath string) (string, error) {
	// Giải mã cipherKey bằng privateKey
	// ...

	// Giải nén file zip
	zipFile, err := zip.OpenReader(fileZipPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	// Lặp qua từng file trong file zip
	for _, file := range zipFile.File {
		// Mở file trong file zip
		zippedFile, err := file.Open()
		if err != nil {
			return "", err
		}
		defer zippedFile.Close()

		// Đọc nội dung file
		content, err := ioutil.ReadAll(zippedFile)
		if err != nil {
			return "", err
		}
		

		// Lưu nội dung file vào file txt
		txtFilePath := "path/to/output.txt" // Thay đổi đường dẫn theo ý muốn
		err = ioutil.WriteFile(txtFilePath, content, 0644)
		if err != nil {
			return "", err
		}

		return txtFilePath, nil
	}

	return "", nil
}


func checkNewData(hash, result string) bool {
	if HashString(result) == hash {
		return true
	}
	return false
}



func createZipFile(zipFilePath string, file1Path string, file2Path string) error {
	// Tạo file zip
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Tạo writer cho file zip
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Thêm file 1 vào zip
	if err := addFileToZip(zipWriter, file1Path); err != nil {
		return err
	}

	// Thêm file 2 vào zip
	if err := addFileToZip(zipWriter, file2Path); err != nil {
		return err
	}

	return nil
}

