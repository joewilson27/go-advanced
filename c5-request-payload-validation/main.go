package main

/*
HTTP Request Payload Validation (Validator v9, Echo)

Pada chapter ini kita akan belajar CARA VALIDASI payload request di sisi back end.
Library yang kita gunakan adalah github.com/go-playground/validator/v10,
library ini sangat berguna untuk KEPERLUAN VALIDASI data.
*/

// 1. Payload Validation
/*
Penggunaan validator cukup mudah, di struct penampung payload, TAMBAHKAN TAG BARU pada masing-masing
property dengan skema validate:"<rules>".

*/
import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type User struct {
	Name  string `json:"name"  validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age"   validate:"gte=0,lte=80"`
}

/*
Struct User memiliki 3 field, berisikan aturan/rule validasi, yang berbeda satu sama lain (bisa dilihat pada tag validate).
Kita bahas validasi per-field agar lebih mudah untuk dipahami.
- Field Name, tidak boleh kosong.
- Field Email, tidak boleh kosong, dan isinya harus dalam format email.
- Field Age, tidak harus di-isi; namun jika ada isinya, maka harus berupa numerik dalam kisaran angka 0 hingga 80.

Kurang lebih berikut adalah penjelasan singkat mengenai beberapa rule yang kita gunakan di atas.
- Rule required, menandakan bahwa field harus di isi.
- Rule email, menandakan bahwa value pada field harus dalam bentuk email.
- Rule gte=n, artinya isi harus numerik dan harus di atas n atau sama dengan n.
- Rule lte=n, berarti isi juga harus numerik, dengan nilai di bawah n atau sama dengan n.

Jika sebuah field MEMBUTUHKAN DUA ATAU LEBIH rule, maka tulis seluruhnya dengan delimiter tanda koma (,).
*/

/*
OK, selanjutnya buat struct baru CustomValidator dengan isi sebuah property bertipe *validator.Validate dan satu buah method ber-skema
Validate(interface{})error. Objek cetakan struct ini akan kita gunakan SEBAGAI PENGGANTI default validator milik echo.
*/
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error { // method Validate() belongs to struct CustomValidator
	/*
		Method .Struct() milik *validator.Validate, digunakan untuk mem-validasi data objek dari struct.

			*Library validator menyediakan banyak sekali cakupan data yang bisa divalidasi, tidak hanya struct,
			 lebih jelasnya silakan lihat di laman github https://github.com/go-playground/validator.
	*/
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.POST("/users", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		/*
			Dalam endpoint ini method Validate milik CustomValidator dipanggil.
		*/
		if err := c.Validate(u); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, true)
	})

	e.Logger.Fatal(e.Start(":9000"))

	/*
		run project, lakukan testing di postman. Lakukan beberapa test dengan mengisi form lengkap dan juga mengosongkan salah satu inputan
		yang akan menyebabkan error validasi

		pada field age, jika kita isi null, tidak akan error, karena kita tidak men-set required pada struct User

		Dari testing di atas bisa kita simpulkan bahwa fungsi validasi berjalan sesuai harapan. Namun MASIH ADA YANG KURANG,
		ketika ada yang tidak valid, error yang dikembalikan SELALU SAMA, yaitu message Internal server error.

		Sebenarnya error 500 ini SUDAH SESUAI jika muncul pada page yang sifatnya MENAMPILKAN CONTENT.
		Pengguna TIDAK PERLU TAHU secara mendetail mengenai detail error yang sedang terjadi.
		Mungkin dibuat saja halaman custom error agar lebih menarik.

		Tapi untuk web service (RESTful API?), AKAN LEBIH BAIK jika errornya detail (terutama pada fase development),
		agar aplikasi consumer bisa lebih bagus dalam meng-handle error tersebut.

		Nah, pada chapter selanjutnya kita akan belajar cara membuat custom error handler untuk
		meningkatkan kualitas error reporting.
	*/
}
