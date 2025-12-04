package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Struktur model bioskop
type Bioskop struct {
	ID     int     `db:"id" json:"id"`
	Nama   string  `db:"nama" json:"nama"`
	Lokasi string  `db:"lokasi" json:"lokasi"`
	Rating float32 `db:"rating" json:"rating"`
}

// Global DB
var db *sqlx.DB

func main() {
	var err error

	// Koneksi ke PostgreSQL
	dsn := "host=localhost port=5432 user=sabar password=postgres dbname=bioskopdb sslmode=disable"
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	// Init Gin
	r := gin.Default()

	// Endpoint membuat data bioskop
	r.POST("/bioskop", createBioskop)

	// Jalankan server
	r.Run(":8080")
}

// Handler POST /bioskop
func createBioskop(c *gin.Context) {
	var input Bioskop

	// ambil JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	// Validasi
	if input.Nama == "" || input.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Lokasi tidak boleh kosong"})
		return
	}

	// Query insert
	query := `
		INSERT INTO bioskop (nama, lokasi, rating)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var newID int
	err := db.QueryRow(query, input.Nama, input.Lokasi, input.Rating).Scan(&newID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	input.ID = newID
	c.JSON(http.StatusCreated, gin.H{
		"message": "Bioskop berhasil ditambahkan",
		"data":    input,
	})
}
