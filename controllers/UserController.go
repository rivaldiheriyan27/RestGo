package controllers

import (
	"RestGo/helpers"
	"RestGo/models"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDB struct {
	DB *gorm.DB
}

func (db *UserDB) Register(c *gin.Context) {
	var (
		req      models.User
		findUser models.User
	)

	err := c.ShouldBindJSON(&req)
	// Maksud dari kodingan ini adalah Dalam contoh di atas, &req adalah alamat dari variabel req, yang seharusnya adalah sebuah struktur yang Anda gunakan untuk menguraikan data JSON yang diterima. Penambahan tanda & sebelum req mengubah variabel req menjadi pointer, yang berarti bahwa ShouldBindJSON akan mengisi nilai-nilai di dalam struct dengan nilai-nilai yang sesuai dari data JSON yang diterima.
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if req.Email == "" {
		c.JSON(400, gin.H{
			"message": "Email is required",
		})
	}

	if req.Username == "" {
		c.JSON(400, gin.H{
			"message": "Username is required",
		})
		return
	}

	if req.Password == "" {
		c.JSON(400, gin.H{
			"message": "Password is required",
		})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(400, gin.H{
			"message": "Minimun length of password is 6 char",
		})
		return
	}

	if req.Age < 8 {
		c.JSON(400, gin.H{
			"message": "Minimun age is 8 years",
		})
		return
	}

	_, errMailFormat := mail.ParseAddress(req.Email)
	if errMailFormat != nil {
		c.JSON(400, gin.H{
			"message": "Email format is warong",
		})
		return
	}

	db.DB.Where("email = ?", req.Email).First(&findUser)
	if findUser != (models.User{}) {
		c.JSON(400, gin.H{
			"message": "Email already used",
		})
		return
	}

	db.DB.Where("username = ?", req.Username).First(&findUser)
	if findUser != (models.User{}) {
		c.JSON(400, gin.H{
			"message": "Username already used",
		})
		return
	}

	errCreate := db.DB.Debug().Create(&req).Error
	if errCreate != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(201, gin.H{
		"age":      req.Age,
		"email":    req.Email,
		"id":       req.ID,
		"username": req.Username,
	})
}

func (db *UserDB) Login(c *gin.Context) {
	var req models.User

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	dbResult := models.User{} // Ini sudah dapat diambil database yang mau digunakan apabila dalam data sudah cocok
	errUser := db.DB.Debug().Where("email = ?", req.Email).Last(&dbResult).Error
	if errUser != nil {
		c.AbortWithError(http.StatusInternalServerError, errUser)
		return
	}

	errBcrypt := bcrypt.CompareHashAndPassword([]byte(dbResult.Password), []byte(req.Password))
	if errBcrypt != nil {
		c.AbortWithError(http.StatusBadRequest, errBcrypt)
		return
	}

	token := helpers.GenerateToken(dbResult.Username)

	// Mengambil data yang ingin ditampilkan dalam respons JSON (misalnya: id, username, email, age, dll.)
	userData := gin.H{
		"id":       dbResult.ID,
		"username": dbResult.Username,
		"email":    dbResult.Email,
		"age":      dbResult.Age,
	}

	// Menggabungkan data user dengan token dalam respons JSON
	responseData := gin.H{ // Ini digunakan sebagai pengiriman ke respon json
		"user":  userData,
		"token": token,
	}

	c.JSON(200, responseData)
}

func (db *UserDB) UserUpdate(c *gin.Context) {

}

func (db *UserDB) UserDelete(c *gin.Context) {

}
