package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const privateKey string = "faef81589d8420c5cf179ba55536ae9a5d"

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Upload File",
		})
	})

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		start := time.Now()
		file, _ := c.FormFile("file")
		if err := c.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
			})
		}
		key, _ := hex.DecodeString("554fa21d43a6f294fe34746ff1481f" + privateKey)
		data, err := ioutil.ReadFile("./uploads/" + file.Filename)
		if err != nil {
			log.Print("read error")
			log.Fatal(err)
		}
		block, err := aes.NewCipher(key)
		if err != nil {
			log.Print("read block")
			panic(err.Error())
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			log.Panic(err)
		}
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			log.Fatal(err)
		}
		ciphertext := gcm.Seal(nonce, nonce, data, nil)
		err = ioutil.WriteFile("./uploads/"+file.Filename, ciphertext, 0777)
		if err != nil {
			log.Print("read write")
			panic(err.Error())
		}
		end := time.Now()
		elapsed := end.Sub(start)
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"time":   elapsed,
		})
	})

	router.GET("/private", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"key": privateKey,
		})
	})

	router.GET("/download/:filename", func(c *gin.Context) {
		// start := time.Now()
		filename := c.Param("filename")
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		fmt.Sprintf("attachment; filename=%s", filename)
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.File("./uploads/" + filename)
		// key, _ := hex.DecodeString("554fa21d43a6f294fe34746ff1481ffaef81589d8420c5cf179ba55536ae9a5d")
		// ciphertext, err := ioutil.ReadFile("./uploads/" + filename)
		// block, err := aes.NewCipher(key)
		// if err != nil {
		// 	log.Panic(err)
		// }
		// gcm, err := cipher.NewGCM(block)
		// if err != nil {
		// 	log.Panic(err)
		// }
		// nonce := ciphertext[:gcm.NonceSize()]
		// ciphertext = ciphertext[gcm.NonceSize():]
		// plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		// if err != nil {
		// 	log.Panic(err)
		// }
		// err = ioutil.WriteFile("./uploads/"+filename, plaintext, 0777)
		// if err != nil {
		// 	log.Print("read write")
		// 	panic(err.Error())
		// }
		// end := time.Now()
		// elapsed := end.Sub(start)
		// c.JSON(http.StatusOK, gin.H{
		// 	"status": http.StatusOK,
		// 	"key":    privateKey,
		// })
	})
	router.Run(":8080")
}
