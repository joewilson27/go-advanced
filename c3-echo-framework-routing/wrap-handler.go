package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// 5. Penggunaan echo.WrapHandler Untuk Routing Handler Bertipe func(http.ResponseWriter,*http.Request) atau http.HandlerFunc
/*
Echo BISA DIKOMBINASIKAN dengan handler ber-skema NON-echo-handler seperti
	func(http.ResponseWriter,*http.Request) atau http.HandlerFunc.

Caranya dengan memanfaatkan fungsi echo.WrapHandler UNTUK MENGKONVERSI handler tersebut menjadi
echo-compatible. Lebih jelasnya silakan lihat kode berikut.
*/
var ActionIndex = func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("from action index"))
}

var ActionHome = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("from action home"))
	},
)

var ActionAbout = echo.WrapHandler(
	http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("from action about"))
		},
	),
)

func main() {
	r := echo.New()

	r.GET("/index", echo.WrapHandler(http.HandlerFunc(ActionIndex)))
	r.GET("/home", echo.WrapHandler(ActionHome))
	r.GET("/about", ActionAbout)
	/*
		Untuk routing handler dengan skema func(http.ResponseWriter,*http.Request), MAKA HARUS DIBUNGKUS
		dua kali, pertama menggunakan http.HandlerFunc, lalu dengan echo.WrapHandler.

		Sedangkan untuk handler yang sudah bertipe http.HandlerFunc,
		bungkus langsung menggunakan echo.WrapHandler.

		ketiga route ini sebenarnya sama, cuma beda cara eksekusinya:
		- /index WrapHandler dan HandlerFunc dilakukan di GET semua
		- /home WrapHandler dilakukan di GET namun HandlerFunc di lakukan di object ActionHome
		- /about WrapHandler dan HandlerFunc di lakukan di var ActionAbout
	*/

	// 6. Routing Static Assets
	/*
		Cara routing static assets di echo sangatlah mudah. Gunakan method .Static(),
		- isi parameter pertama dengan prefix rute yang di-inginkan, misal /assets, /static dll
		- dan parameter ke-2 dengan path folder tujuan.

		lalu jalankan url berikut utk mengakses file nya:
		http://localhost:9000/static/layout.js
	*/
	r.Static("/static", "assets")

	r.Start(":9000")
}
