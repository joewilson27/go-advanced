package main

// 5. CORS Handling Menggunakan Golang CORS Library dan Echo
/*
	Pada bagian ini kita akan mengkombinasikan library CORS golang buatan Olivier Poitrey, dan Echo,
	untuk membuat back end yang mendukung cross origin request.

	Pertama go get dulu library-nya.
	go get https://github.com/rs/cors
*/

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/rs/cors"
)

func main() {
	e := echo.New()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"https://www.apple.com", "https://www.google.com"},
		AllowedMethods: []string{"OPTIONS", "GET", "POST", "PUT"},
		AllowedHeaders: []string{"Content-Type", "X-CSRF-Token"},
		Debug:          true,
	})

	/*
		Pada kode di atas, kita meng-allow dua buah origin. Sebelumnya sudah kita bahas bahwa kebanyakan browser tidak mendukung ini.
		Dengan menggunakan CORS library, hal itu BISA TERATASI.
	*/

	e.Use(echo.WrapMiddleware(corsMiddleware.Handler))

	e.GET("/index", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	e.Logger.Fatal(e.Start(":9000"))
}
