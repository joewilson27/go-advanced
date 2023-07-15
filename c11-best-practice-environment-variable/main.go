package main

/*
Best Practice Configuration Menggunakan Environment Variable
*/

// 1. Definisi
/*
Environment variable merupakan variabel yang BERADA DI lapisan runtime sistem operasi.
Karena env var atau environment variable merupakan variabel seperti pada umumnya,
maka kita bisa melakukan operasi seperti mengubah nilainya atau mengambil nilainya.

Salah satu env var yang mungkin sering temen-temen temui adalah PATH.
PATH sendiri merupakan variabel yang digunakan oleh sistem operasi untuk
men-specify direktori tempat di mana binary atau executable berada.

Default-nya, sistem operasi PASTI MEMPUNYAI BEBERAPA env var YANG SUDAH ADA tanpa kita set,
salah satunya seperti PATH tadi, juga lainnya. Variabel-variabel tersebut digunakan oleh
sistem operasi untuk keperluan mereka. Tapi karena variabel juga bisa diakses oleh kita
(selaku developer), maka kita pun juga bisa mempergunakannya untuk kebutuhan tertentu.

Selain reserved env var, kita BISA JUGA MEMBUAT variabel baru yang hanya digunakan untuk
keperluan program secara spesifik.
*/

// 2. Penggunaan env var Sebagai Media Untuk Definisi Konfigurasi Program
/*
Pada chapter B.22. Simple Configuration dan juga C.10. Advanced Configuration: Viper, kita telah
belajar cara pendefinisian konfigurasi dengan MEMANFAATKAN FILE seperti JSON maupun YAML.

Pada chapter kali ini kita akan mendefinisikan konfigurasi yang sama TAPI TIDAK DI FILE,
melainkan di environment variable.

Definisi konfigurasi di env var banyak manfaatnya, salah satunya:
- Di support secara native oleh semua sistem operasi.
- Sudah sangat umum diterapkan di banyak aplikasi dan platform.
- Straightforward dan tidak tergantung ke file tertentu.
- Sharing konfigurasi dengan aplikasi/service lain menjadi lebih mudah.
- Mudah untuk di maintain, tidak perlu repot buka file kemudian edit lalu simpan ulang.
- ... dan banyak lagi lainnya.

Jadi bisa dibilang penulisan konfigurasi di env var MERUPAKAN BEST PRACTICE untuk banyaJadi bisa dibilang penulisan konfigurasi
di env var merupakan best practice untuk banyak jenis kasus, terutama pada microservice, pada aplikasi/service yang distributable,
maupun pada aplikasi monolith yang manajemenya ter-automatisasi.k jenis kasus, terutama pada microservice, pada aplikasi/service
yang distributable, maupun pada aplikasi monolith yang manajemenya ter-automatisasi.

Memang kalau dari sisi readability sangat kalah kalau dibandingkan dengan JSON atau YAML, tapi saya sampaikan bahwa meski effort
koding bakal lebih banyak, akan ada sangat banyak manfaat yang bisa didapat dengan menuliskan konfigurasi di env var,
terutama pada bagian devops.
*/

// 3. Praktek
import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	// ambil nilai konfigurasi nama aplikasi dari env var. Caranya kurang lebih seperti berikut.
	confAppName := os.Getenv("APP_NAME")
	if confAppName == "" {
		e.Logger.Fatal("APP_NAME config is required")
	}
	/*
		Jadi APP_NAME di situ merupakan nama env var-nya. Umumnya env var tidak dituliskan dalam bentuk camelCase,
		tapi dalam bentuk UPPERCASE dengan separator kata adalah underscore. Untuk value-nya nanti tinggal kita
		siapkan saja sebelum proses eksekusi program.

		Statement os.Getenv digunakan untuk pengambilan env var. Pada contoh di atas, terdapat pengecekan jika
		nilai APP_NAME adalah kosong, maka munculkan fatal error.
	*/
	confServerPort := os.Getenv("SERVER_PORT")
	if confServerPort == "" {
		e.Logger.Fatal("SERVER_PORT config is required")
	}

	/*
		Setelah itu, tambahkan routing untuk untuk GET /index lalu definisi objek server yang nantinya digunakan untuk
		keperluan start webserver.
	*/
	e.GET("/index", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, true)
	})

	server := new(http.Server)
	server.Addr = ":" + confServerPort // Nilai server.Addr diambil dari env var SERVER_PORT.

	/*
		Kemudian tambahkan setting untuk timeout webserver, tapi hanya ketika memang timeout didefinisikan konfigurasinya.
	*/
	if confServerReadTimeout := os.Getenv("SERVER_READ_TIMEOUT_IN_MINUTE"); confServerReadTimeout != "" {
		// var SERVER_READ_TIMEOUT_IN_MINUTE ADA NILAINYA
		duration, _ := strconv.Atoi(confServerReadTimeout)
		server.ReadTimeout = time.Duration(duration) * time.Minute // konversi ke bentuk time duration
	}
	if confServerWriteTimeout := os.Getenv("SERVER_WRITE_TIMEOUT_IN_MINUTE"); confServerWriteTimeout != "" {
		duration, _ := strconv.Atoi(confServerWriteTimeout)
		server.WriteTimeout = time.Duration(duration) * time.Minute
	}
	/*
		Bisa dilihat di atas, jika env var SERVER_READ_TIMEOUT_IN_MINUTE ADA NILAINYA, maka diambil kemudian di konversi ke
		bentuk time.Duration untuk dipergunakan pada server.ReadTimeout. Nilai balik dari os.Getenv() pasti BERUPA STRING, oleh karena itu jika
		konfigurasi dibutuhkan dalam bentuk lain, tambahkan saja statement untuk konversi datanya.
	*/
	e.Logger.Print("Starting", confAppName)
	e.Logger.Fatal(e.StartServer(server))

	// 4. Eksekusi Program
	/*
		Program sudah siap, betul, tetapi konfigurasi nya belum. Nah salah satu kelebihan dari kontrol konfigurasi lewat
		env var adalah kita bisa DEFINISIKAN SEWAKTU EKSEKUSI program (sebelum statement go run).

		Ada satu hal yang penting untuk diketahui. Cara set env var untuk Windows dibanding sistem operasi lainnya adalah berbeda.
		Untuk non-Windows, gunakan export.
			export APP_NAME=SimpleApp
			export SERVER_PORT=9000
			export SERVER_READ_TIMEOUT_IN_MINUTE=2
			export SERVER_WRITE_TIMEOUT_IN_MINUTE=2
			go run main.go

		Untuk Windows, gunakan set.
			set APP_NAME=SimpleApp
			set SERVER_PORT=9000
			set SERVER_READ_TIMEOUT_IN_MINUTE=2
			set SERVER_WRITE_TIMEOUT_IN_MINUTE=2
			go run main.go

		Agak sedikit report memang untuk bagian ini, tapi mungkin bisa DIPERINGKAS dengan membuat file .sh untuk non-Windows, dan file .bat
		untuk Windows.
		Jadi nanti bisa tinggal eksekusi file sh/bat-nya saja. Atau temen-temen bisa tulis saja dalam Makefile.
		Untuk windows bisa kok eksekusi command make caranya dengan install make lewat Chocolatey.


		jalanin run.bat lewat CMD, INI BARU BISA
	*/

}
