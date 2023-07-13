package main

/*
Echo Framework & Routing
Pada chapter ini kita akan belajar cara mudah routing menggunakan Echo Framework.

Mulai chapter C1 hingga C6 kita akan mempelajari banyak aspek dalam framework Echo dan mengkombinasikannya
dengan beberapa library lain.
*/

// 1 Echo Framework1 Echo Framework
/*
Echo adalah framework bahasa golang untuk PENGEMBANGAN aplikasi web. Framework ini cukup terkenal di komunitas.
Echo merupakan framework besar, di dalamnya terdapat BANYAK SEKALI dependensi.

Salah satu dependensi yang ada di dalamnya adalah router, dan pada chapter ini kita akan mempelajarinya.

Dari banyak routing library yang sudah penulis gunakan, hampir seluruhnya mempunyai kemiripan dalam hal penggunaannya,
cukup panggil fungsi/method yang dipilih (biasanya namanya sama dengan HTTP Method), lalu sisipkan rute pada parameter
pertama dan handler pada parameter kedua.

Berikut contoh sederhana penggunaan echo framework.

r := echo.New()
r.GET("/", handler)
r.Start(":9000")

Sebuah objek router r dicetak lewat echo.New(). Lalu lewat objek router tersebut,
dilakukan registrasi rute untuk / dengan method GET dan handler adalah closure handler. Terakhir,
dari objek router di-start-lah sebuah web server pada port 9000.
*/

// 2. Praktek
import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

type M map[string]interface{}

func main() {

	r := echo.New() // Statement echo.New() mengembalikan objek mux/router

	r.GET("/", func(ctx echo.Context) error {
		/*
			Handler dari method routing milik echo membutuhkan satu argument saja, dengan tipe adalah echo.Context

			Dari argumen tersebut (ctx) objek http.ResponseWriter dan http.Request BISA di-akses. Namun kedua objek
			tersebut AKAN JARANG kita gunakan karena echo.Context MEMILIKI BANYAK METHOD yang beberapa tugasnya
			SUDAH MENG-COVER operasi umum yang biasanya kita lakukan lewat objek request dan response, di antara seperti:
			- Render output (dalam bentuk html, plain text, json, atau lainnya).
			- Parsing request data (json payload, form data, query string).
			- URL Redirection.
			- ... dan lainnya.

				*Untuk mengakses objek http.Request gunakan ctx.Request().
				 Sedang untuk objek http.ResponseWriter gunakan ctx.Response()

			Salah satu alasan lain kenapa penulis memilih framework ini, adalah karena desain route-handler-nya MENARIK.
			Dalam handler cukup KEMBALIKAN objek error KETIKA MEMANG ADA KESALAHAN terjadi,
			sedangkan jika tidak ada error maka kembalikan nilai nil.

			Ketika terjadi error pada saat mengakses endpoint, IDEALNYA HTTP Status error DIKEMBALIKAN sesuai dengan
			jenis errornya. Tapi terkadang juga ada kebutuhan dalam KONDISI TERTENTU http.StatusOK atau status 200
			dikembalikan dengan disisipi informasi error dalam response body-nya.
			Kasus sejenis ini menjadikan standar error reporting menjadi kurang bagus. Pada konteks ini echo unggul
			menurut penulis, karena default-nya SEMUA ERROR DIKEMBALIKAN sebagai response dalam bentuk yang sama.


			Method ctx.String() dari objek context milik handler digunakan untuk MEMPERMUDAH rendering data string
			sebagai OUTPUT. Method ini mengembalikan objek error, jadi bisa digunakan langsung sebagai nilai balik handler.
			Argumen pertama adalah http status dan argumen ke-2 adalah data yang dijadikan output.

		*/
		data := "Hello from /"
		return ctx.String(http.StatusOK, data)
	})

	// Method .String()
	// Digunakan untuk render plain text sebagai output
	r.GET("/index", func(ctx echo.Context) error {
		data := "Hello from /index"
		return ctx.String(http.StatusOK, data)
	})

	// Method .HTML()
	// Digunakan untuk render html sebagai output. Isi response header Content-Type adalah text/html.
	r.GET("/html", func(ctx echo.Context) error {
		data := "Hello from /html"
		return ctx.HTML(http.StatusOK, data)
	})

	// Method .Redirect()
	// Digunakan untuk redirect, pengganti http.Redirect().
	r.GET("/index", func(ctx echo.Context) error {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	})

	// Method .JSON()
	// Digunakan untuk render data JSON sebagai output. Isi response header Content-Type adalah application/json.
	r.GET("/json", func(ctx echo.Context) error {
		data := M{"Message": "Hello", "Counter": 2}
		return ctx.JSON(http.StatusOK, data)
	})

	/*
		Routing dengan memanfaatkan package net/http dalam penerapannya adalah menggunakan http.HandleFunc()
		atau http.Handle(). BERBEDA dengan Echo, routingnya adalah method-based, tidak hanya endpoint dan handler
		yang di-registrasi, method juga.

		Statement echo.New() mengembalikan objek mux/router. Pada kode di atas rute / dengan method GET di-daftarkan.
		Selain r.GET() ada banyak lagi method lainnya, semua method dalam spesifikasi REST seperti PUT, POST, dan lainnya
		bisa digunakan.
	*/

	// 4. Parsing Request
	/*
		Echo juga menyediakan beberapa method untuk keperluan parsing request, di antaranya:
	*/
	// Parsing Query String
	// Method .QueryParam() digunakan untuk MENGAMBIL data pada query string request, sesuai dengan key yang diinginkan.
	r.GET("/page1", func(ctx echo.Context) error {
		name := ctx.QueryParam("name")
		data := fmt.Sprintf("Hello %s", name)

		return ctx.String(http.StatusOK, data)
	})

	// Parsing URL Path Param
	// Method .Param() digunakan untuk MENGAMBIL data PATH PARAMETER sesuai skema rute.
	r.GET("/page2/:name", func(ctx echo.Context) error {
		name := ctx.Param("name")
		/*
			Bisa dilihat, terdapat :name pada pendeklarasian rute (/page2/:name). Nantinya url apapun yang ditulis sesuai
			skema di-atas akan bisa diambil path parameter-nya. Misalkan /page2/halo maka ctx.Param("name") mengembalikan string halo.
			atau /page2/lorem akan mengambalikan string lorem
		*/
		data := fmt.Sprintf("Hello %s", name)

		return ctx.String(http.StatusOK, data)
	})

	// Parsing URL Path Param dan Setelahnya
	/*
		Selain mengambil parameter sesuai spesifik path, kita juga bisa mengambil data parameter path dan setelahnya.

		Statement ctx.Param("*") MENGEMBALIKAN SEMUA PATH sesuai dengan skema url-nya.
		Misal url adalah /page3/tim/a/b/c/d/e/f/g/h maka yang dikembalikan adalah a/b/c/d/e/f/g/h.

		test dengan:
		curl -X GET http://localhost:9000/page3/tim/need/some/sleep
		will printed "Hello tim, I have message for you: need/some/sleep"
	*/
	r.GET("/page3/:name/*", func(ctx echo.Context) error {
		name := ctx.Param("name")
		message := ctx.Param("*") // url path *

		data := fmt.Sprintf("Hello %s, I have message for you: %s", name, message)

		return ctx.String(http.StatusOK, data)
	})

	// Parsing Form Data
	/*
		Data yang dikirim sebagai request body dengan jenis adalah Form Data bisa di-ambil dengan mudah menggunakan ctx.FormValue().

		test dengan postman atau curl berikut:
		curl -X POST -F name=damian -F message=angry http://localhost:9000/page4
	*/
	r.POST("/page4", func(ctx echo.Context) error {
		name := ctx.FormValue("name")
		message := ctx.FormValue("message")

		data := fmt.Sprintf(
			"Hello %s, I have message for you: %s",
			name,
			strings.Replace(message, "/", "", 1),
		)

		return ctx.String(http.StatusOK, data)
	})

	var port = fmt.Sprintf(":%d", 9000)
	r.Start(port)

}
