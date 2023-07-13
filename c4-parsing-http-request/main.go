package main

/*
Parsing HTTP Request Payload (Echo)
Pada chapter ini kita akan belajar cara parsing beberapa variasi request payload.
Payload dalam HTTP request BISA DIKIRIMKAN dalam BERBAGAI BENTUK. Kita akan mempelajari cara untuk
handle 4 jenis payload berikut.
- Form Data (Content-Type: application/x-www-form-urlencoded)
- JSON Payload (Content-Type: application/json)
- XML Payload (Content-Type: application/xml)
- Query String (data disisipkan pada url)
*/

// 1. Parsing Request Payload
/*
Cara parsing payload request dalam echo SANGAT MUDAH, APAPUN JENIS payload nya, API
yang DIGUNAKAN untuk parsing adalah SAMA.
*/
import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type User struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

func main() {
	r := echo.New()

	/*
		Selanjutnya siapkan satu buah endpoint /user menggunakan r.Any(). Method .Any()
		MENERIMA SEGALA JENIS request dengan method GET, POST, PUT, atau lainnya.
	*/
	r.Any("/user", func(c echo.Context) (err error) {
		u := new(User)
		if err = c.Bind(u); err != nil {
			return
		}
		/*
			Bisa dilihat dalam handler, method .Bind() milik echo.Context digunakan,
			dengan DISISIPI parameter pointer objek (hasil cetakan struct User).
			Parameter tersebut nantinya AKAN MENAMPUNG payload yang dikirim,
			entah apapun jenis nya.
		*/

		return c.JSON(http.StatusOK, u)
	})

	port := fmt.Sprintf(":%d", 9000)
	fmt.Printf("server started at %s", port)
	r.Start(port)
}

// 2 Testing
/*
Jalankan aplikasi, lakukan testing. Bisa gunakan curl ataupun API
testing tools sejenis postman atau lainnya.

Di bawah ini shortcut untuk melakukan request menggunakan curl pada 4 jenis
payload yang kita telah bahas.
Response dari seluruh request adalah sama, menandakan bahwa data yang dikirim
berhasil ditampung.

• Form Data
curl -X POST http://localhost:9000/user \
     -d 'name=Joe' \
     -d 'email=nope@novalagung.com'

# output => {"name":"Nope","email":"nope@novalagung.com"}


• JSON Payload
curl -X POST http://localhost:9000/user \
     -H 'Content-Type: application/json' \
     -d '{"name":"Nope","email":"nope@novalagung.com"}'

# output => {"name":"Nope","email":"nope@novalagung.com"}

• XML Payload
curl -X POST http://localhost:9000/user \
     -H 'Content-Type: application/xml' \
     -d '<?xml version="1.0"?>\
        <Data>\
            <Name>Joe</Name>\
            <Email>nope@novalagung.com</Email>\
        </Data>'

# output => {"name":"Nope","email":"nope@novalagung.com"}


• Query String
curl -X GET http://localhost:9000/user?name=Joe&email=nope@novalagung.com

# output => {"name":"Nope","email":"nope@novalagung.com"}
*/
