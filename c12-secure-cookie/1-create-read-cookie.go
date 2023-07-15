package main

/*
Secure Cookie (Gorilla Securecookie)
Cookie memiliki beberapa atribut, DIANTARANYA adalah secure. Dengan mengaktifkan atribut ini,
informasi cookie menjadi LEBIH AMAN karena di-enkripsi, namun kapabilitas ini HANYA AKAN AKTIF
pada kondisi aplikasi SSL/TLS enabled.

	*TL;DR; Jika atribut secure di-isi true, namun web server TIDAK menggunakan SSL/TLS,
		maka cookie disimpan seperti biasa tanpa di-enkripsi.


Lalu bagaimana cara untuk membuat cookie aman pada aplikasi yang meng-enable SSL/TLS maupun yang tidak?
caranya adalah dengan MENAMBAHKAN STEP enkripsi data sebelum disimpan dalam cookie
(dan men-decrypt data tersebut saat membaca).

Gorilla toolkit menyediakan library bernama securecookie, berguna UNTUK MEMPERMUDAH enkripsi informasi cookie,
dengan penerapan yang mudah. Pada chapter ini kita akan mempelajari penggunaannya.

*/

// 1. Create & Read Secure Cookie
/*
Penggunaan securecookie cukup mudah, buat objek secure cookie lewat securecookie.New() lalu gunakan objek tersebut untuk
operasi encode-decode data cookie. Pemanggilan fungsi .New() memerlukan 2 buah argument.
- Hash key, diperlukan untuk otentikasi data cookie menggunakan algoritma kriptografi HMAC.
- Block key, adalah opsional, diperlukan untuk enkripsi data cookie. Default algoritma enkripsi yang digunakan adalah AES.
*/

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/labstack/echo"
	"github.com/novalagung/gubrak/v2"
)

type M map[string]interface{}

var sc = securecookie.New([]byte("very-secret"), []byte("a-lot-secret-yay")) //  buat objek secure cookie lewat securecookie.New()
/*
Variabel sc adalah objek secure cookie. Objek ini kita gunakan untuk ENCODE DATA yang akan disimpan dalam cookie,
dan juga untuk decode data.
*/

// Buat fungsi setCookie(), bertugas untuk MEMPERMUDAH PEMBUATAN dan PENYIMPANAN cookie.
func setCookie(c echo.Context, name string, data M) error {
	encoded, err := sc.Encode(name, data)
	/*
		Method sc.Encode() digunakan untuk encoding data dengan identifier
		adalah isi variabel name. Variabel encoded menampung data setelah di-encode,
		lalu variabel ini dimasukan ke dalam objek cookie.
	*/
	if err != nil {
		return err
	}

	// Pembuatan cookie cukup mudah, tinggal cetak saja objek baru dari struct http.Cookie.
	// c := &http.Cookie{} --> cara create cookie
	cookie := &http.Cookie{
		Name:     name,
		Value:    encoded, // value encoded di masukan ke property Value
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(1 * time.Hour),
	}
	http.SetCookie(c.Response(), cookie) //  untuk menyimpan cookie yang baru dibuat.

	return nil
}

/*
Method sc.Encode() digunakan untuk encoding data dengan identifier adalah isi variabel name.
Variabel encoded MENAMPUNG DATA setelah di-encode, lalu variabel ini DIMASUKKAN KE dalam objek cookie.

Cara menyimpan cookie masih sama, menggunakan http.SetCookie.
*/

// Selanjutnya buat fungsi getCookie(), untuk mempermudah proses pembacaan cookie yang tersimpan
func getCookie(c echo.Context, name string) (M, error) {
	cookie, err := c.Request().Cookie(name) // ambil cookie
	if err == nil {

		data := M{}
		if err = sc.Decode(name, cookie.Value, &data); err == nil {
			return data, nil
		}
	}

	return nil, err
}

/*
Setelah cookie diambil menggunakan c.Request().Cookie(), data di dalamnya PERLU
di-decode agar bisa terbaca. Method sc.Decode() digunakan untuk decoding data.
*/

// 2. Delete Secure Cookie
func removeCookie(c echo.Context, name string) {
	cookie := &http.Cookie{}
	cookie.Name = name
	cookie.Path = "/"
	cookie.MaxAge = -1
	cookie.Expires = time.Unix(0, 0) // expired kan waktu untuk menghapus
	http.SetCookie(c.Response(), cookie)
}

func main() {
	const COOKIE_NAME = "data"
	/*
		Konstanta COOKIE_NAME disiapkan, kita gunakan sebagai identifier cookie.
	*/

	e := echo.New()

	e.GET("/index", func(c echo.Context) error {
		/*
			Dan sebuah rute juga disiapkan dengan tugas MENAMPILKAN DATA cookie jika sudah ada,
			dan membuat cookie baru jika belum ada.
		*/

		data, err := getCookie(c, COOKIE_NAME)
		if err != nil && err != http.ErrNoCookie && err != securecookie.ErrMacInvalid {
			/*
				http.ErrNoCookie adalah variabel penanda error karena cookie KOSONG,
				sedangkan securecookie.ErrMacInvalid adalah representasi dari invalid cookie.
			*/
			return err
		}

		// buat cookie baru jika belum ada
		if data == nil {
			// siapkan data baru utk cookie, data bertipe map
			data = M{"Message": "Hello", "ID": gubrak.RandomString(32)}

			err = setCookie(c, COOKIE_NAME, data)
			if err != nil {
				return err
			}
		}

		return c.JSON(http.StatusOK, data)
	})

	// 2. Delete Secure Cookie
	/*
		Securecookie perannya hanya pada bagian encode-decode data cookie, sedangkan proses simpan baca cookie masih sama seperti
		penerapan cookie biasa. Maka cara menghapus cookie PUN MASIH SAMA, yaitu dengan meng-EXPIRED-kan cookie yang sudah disimpan.
	*/
	e.GET("/delete", func(c echo.Context) error {
		removeCookie(c, COOKIE_NAME)
		return c.Redirect(http.StatusTemporaryRedirect, "/index")
	})

	e.Logger.Fatal(e.Start(":9000"))
}
