package main

/*
C.16. Secure Middleware
Pada chapter ini kita akan belajar menggunakan library secure (https://github.com/unrolled/secure)
untuk meningkatkan keamanan aplikasi web

16.1. Keamanan Web Server
Jika berbicara mengenai keamanan aplikasi web, SANGAT LUAS sebenarnya cakupannya,
ada banyak hal yang perlu diperhatian dan disiapkan. Mungkin tiga di antaranya sudah kita pelajari sebelumnya,
yaitu penerapan Secure Cookie, CORS, dan CSRF

Secure library merupakan middleware, penggunaannya sama seperti middleware pada umumnya.

*/
import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/unrolled/secure"
)

func main() {
	e := echo.New()

	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:            []string{"localhost:9000", "www.google.com"},
		FrameDeny:               true,
		CustomFrameOptionsValue: "SAMEORIGIN",
		ContentTypeNosniff:      true,
		BrowserXssFilter:        true,
	})

	e.Use(echo.WrapMiddleware(secureMiddleware.Handler))

	e.GET("/index", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")

		return c.String(http.StatusOK, "Hello")
	})

	e.Logger.Fatal(e.StartTLS(":9000", "server.crt", "server.key"))

}

/*
Perlu diketahui, aplikasi di atas di-start dengan SSL/TLS enabled
*/
