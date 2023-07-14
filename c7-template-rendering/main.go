package main

/*
Template Rendering in Echo
Pada dasarnya proses parsing dan rendering template TIDAK DI-HANDLE oleh echo sendiri,
melainkan oleh API dari package html/template. Jadi bisa dibilang cara render template di echo adalah
SAMA SEPERTI pada aplikasi yang murni menggunakan golang biasa, seperti yang sudah dibahas pada chapter sebelumnya.

Echo MENYEDIAKAN SATU fasilitas yang bisa kita manfaatkan UNTUK STANDARISASI rendering template.
Cara penggunaannya, dengan meng-OVERRIDE default .Renderer property milik echo menggunakan objek cetakan struct,
yang di mana pada struct tersebut harus ada method bernama .Render() dengan skema sesuai dengan kebutuhan echo.
Nah, di dalam method .Render() inilah kode untuk parsing dan rendering template ditulis.
*/

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type M map[string]interface{}

type Renderer struct {
	template *template.Template
	debug    bool
	location string
}

/*
Berikut adalah tugas dan penjelasan mengenai ketiga property di atas.
- Property .template bertanggung jawab untuk parsing dan rendering template.
- Property .location MENGARAH KE PATH folder di mana file template berada.
- Property .debug menampung nilai bertipe bool.
		- Jika false, maka parsing template HANYA DILAKUKAN SEKALI saja pada saat
			aplikasi di start. Mode ini sangat cocok untuk diaktifkan pada stage production.
		- Sedangkan jika nilai adalah true, maka parsing template dilakukan TIAP PENGAKSESAN
			rute. Mode ini cocok diaktifkan untuk stage development,
			karena perubahan kode pada file html sering pada stage ini.
*/

// Selanjutnya buat fungsi NewRenderer() untuk MEMPERMUDAH inisialisasi objek renderer.
func NewRenderer(location string, debug bool) *Renderer {
	tpl := new(Renderer)
	tpl.location = location
	tpl.debug = debug

	tpl.ReloadTemplates()

	return tpl
}

func (t *Renderer) ReloadTemplates() {
	t.template = template.Must(template.ParseGlob(t.location))
}

/*
Method .ReloadTemplates() bertugas untuk parsing template.
Method ini WAJIB dipanggil pada saat inisialisasi objek renderer.
Jika .debug == true, maka method ini HARUS DIPANGGIL setiap kali rute diakses
(jika tidak, maka perubahan pada view tidak akan muncul).
*/

func (t *Renderer) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {
	if t.debug {
		t.ReloadTemplates()
	}

	return t.template.ExecuteTemplate(w, name, data)
}

/*
Method .Render() berguna untuk RENDER TEMPLATE YANG SUDAH DIPARSING sebagai output.
Method ini HARUS DIBUAT dalam skema berikut:
// skema method Render()
func (io.Writer, string, interface{}, echo.Context) error
*/

func main() {
	e := echo.New()

	// override property renderer nya, dan siapkan sebuah rute.
	e.Renderer = NewRenderer("./*.html", true)
	/*
		Saat pemanggilan NewRenderer() sisipkan path folder tempat file template html berada. Gunakan ./*.html
		agar mengarah ke semua file html pada current folder.
	*/

	e.GET("/index", func(c echo.Context) error {
		data := M{"message": "Hello World!"}
		return c.Render(http.StatusOK, "index.html", data)
	})

	e.Logger.Fatal(e.Start(":9000"))
}

// 2. Render Parsial dan Spesifik Template
/*
Proses parsing dan rendering TIDAK DI-HANDLE oleh echo, melainkan menggunakan API
dari html/template. Echo hanya MENYEDIAKAN TEMPAT untuk mempermudah pemanggilan fungsi
rendernya. Nah dari sini berarti untuk render parsial, render spesifik template,
maupun operasi template lainnya dilakukan seperti biasa, menggunakan html/template.
*/
