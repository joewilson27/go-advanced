package main

// 15. CSFR

// 1. Teori
/*
Cross-Site Request Forgery atau CSRF adalah SALAH SATU TEKNIK HACKING yang dilakukan dengan cara
mengeksekusi perintah yang seharusnya tidak diizinkan, tetapi output yang dihasilkan sesuai dengan
yang seharusnya. Contoh serangan jenis ini: mencoba untuk login lewat media selain web browser,
seperti menggunakan CURL, menembak langsung endpoint login.
Masih banyak contoh lainnya yang lebih ekstrim.

Ada beberapa cara untuk mencegah serangan ini, salah satunya adalah dengan memanfaatkan csrf token.
Di setiap halaman yang ada form nya, csrf token di-generate. Pada saat submit form, csrf disisipkan di request,
lalu di sisi back end dilakukan pengecekan apakah csrf yang dikirim valid atau tidak.

Csrf token sendiri merupakan sebuah random string yang di-generate setiap kali halaman form muncul. Biasanya di
tiap POST request, token tersebut disisipkan sebagai header, atau form data, atau query string.

Lebih detailnya silakan merujuk ke https://en.wikipedia.org/wiki/Cross-site_request_forgery.
*/

// 2. Praktek: Back End
/*
Di golang, pencegahan CSRF bisa dilakukan dengan membuat middleware untuk pengecekan setiap request POST yang masuk.
Cukup mudah sebenarnya, namun agar lebih mudah lagi kita akan gunakan salah satu middleware milik echo framework untuk
belajar.

Di setiap halaman, jika di dalam html nya terdapat form, maka harus disisipkan token csrf. Token tersebut di-generate
oleh middleware.

Di tiap POST request hasil dari form submit, token tersebut harus ikut dikirimkan. Proses validasi token sendiri di-handle oleh middleware.

Mari kita praktekkan, siapkan project baru. Buat file main.go, isi dengan kode berikut.
*/

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type M map[string]interface{}

const CSRF_TOKEN_HEADER = "X-Csrf-Token"
const CSRF_KEY = "csrf_token"

func main() {
	port := 9000
	tmpl := template.Must(template.ParseGlob("./*.html"))

	e := echo.New()

	// Objek middleware CSRF dibuat lewat statement middleware.CSRF(), konfigurasi default digunakan.
	// Atau bisa juga dibuat dengan disertakan konfigurasi custom, lewat middleware.CSRFWithConfig() seperti pada kode di atas.
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:" + CSRF_TOKEN_HEADER,
		ContextKey:  CSRF_KEY,
	}))
	/*
		Property TokenLookup adalah acuan di bagian mana informasi csrf disisipkan dalam objek request, apakah dari header, query string, atau form data.
		Ini penting karena dibutuhkan oleh middleware yang bersangkutan untuk memvalidasi token tersebut.
		header:X-CSRF-Token, artinya csrf token dalam request akan disisipkan dalam HEADER dengan key adalah X-CSRF-Token
		*Isi value TokenLookup dengan "form:<name>" jika token disispkan dalam form data request, dan "query:<name>" jika token disisipkan dalam query string.


		Property ContextKey digunakan untuk MENGAKSES token csrf yang tersimpan di echo.Context, pembuatan token sendiri terjadi pada saat ada http
		request GET masuk.
		Property tersebut kita isi dengan konstanta CSRF_KEY, maka dalam pengambilan token cukup panggil c.Get(CSRFKey).
	*/

	e.GET("/index", func(c echo.Context) error {
		data := make(M)
		data[CSRF_KEY] = c.Get(CSRF_KEY) // maka dalam pengambilan token cukup panggil c.Get(CSRFKey), Token kemudian disisipkan sebagai data pada saat rendering view.html.
		data["paramtest"] = "testttt"
		return tmpl.Execute(c.Response(), data)
	})

	e.POST("/sayhello", func(c echo.Context) error {
		/*
			Pada handler endpoint /sayhello TIDAK ADA pengecekan token csrf, karena sudah ditangani secara implisit
			oleh middleware.
		*/
		data := make(M)
		if err := c.Bind(&data); err != nil {
			return err
		}

		message := fmt.Sprintf("Hello %s", data["name"])
		return c.JSON(http.StatusOK, message)
	})

	addr := fmt.Sprintf(":%d", port)
	e.Logger.Fatal(e.Start(addr))
}

/*
Percobaan:

Coba tembak langsung endpoint nya lewat CURL.

curl -X POST http://localhost:9000/sayhello \
     -H 'Content-Type: application/json' \
     -d '{"name":"noval","gender":"male"}'


Hasilnya error, karena token csrf tidak di-sisipkan.

Lewat teknik pencegahan ini, BUKAN BERARTI serangan CSRF tidak bisa dilakukan, si hacker masih bisa menembak
endpoint secara paksa lewat CURL, hanya saja membutuhkan usaha ekstra jika ingin sukses
*/
