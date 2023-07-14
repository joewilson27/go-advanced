package main

/*
HTTP Error Handling (Validator v9, Echo)
Pada chapter ini kita akan belajar cara membuat custom error handler yang lebih readable,
SANGAT COCOK untuk web service.
*/

// 1. Error Handler
/*
Cara meng-custom default error handler milik echo, adalah dengan meng-OVERRIDE property e.HTTPErrorHandler
Langsung saja override property tersebut dengan callback berisi parameter objek error dan context.
Gunakan callback tersebut untuk bisa menampilkan error yg lebih detail.
*/
import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type User struct {
	Name  string `json:"name"  validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age"   validate:"gte=0,lte=80"`
}

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

	// error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		report, ok := err.(*echo.HTTPError)
		/*
			Pada kode di atas, objek report menampung objek error setelah di
			casting ke tipe echo.HTTPError. Error tipe ini adalah error-error yang
			BERHUBUNGAN DENGAN http, yang di-handle oleh echo. Untuk error yang BUKAN DARI
			echo, tipe nya adalah error biasa. Pada kode di atas kita standarkan,
			semua jenis error harus berbentuk echo.HTTPError.
		*/
		if !ok {
			report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		//if castedObject, ok := err.(validator.ValidationErrors); ok {
		/*
			Error yang dikembalikan oleh validator.v9 bertipe validator.ValidationErrors.
			Pada kode disini bisa kita lihat ada pengecekan apakah error tersebut
			adalah dari library validator.v9 atau bukan; jika memang iya,
			maka report.Message diubah isinya dengan kata-kata yang lebih mudah dipahami.

			Tipe validator.ValidationErrors sendiri sebenarnya merupakan slice []validator.FieldError.
			Objek tersebut di-loop, lalu diambil-lah elemen pertama sebagai nilai bailk error.
		*/
		// 	for _, err := range castedObject {
		// 		switch err.Tag() {
		// 		case "required":
		// 			report.Message = fmt.Sprintf("%s is required", err.Field())
		// 		case "email":
		// 			report.Message = fmt.Sprintf("%s is not valid email", err.Field())
		// 		case "gte":
		// 			report.Message = fmt.Sprintf("%s value must be greater than %s",
		// 				err.Field(), err.Param())
		// 		case "lte":
		// 			report.Message = fmt.Sprintf("%s value must be lower than %s",
		// 				err.Field(), err.Param())
		// 		}

		// 		break
		// 	}
		// }

		//c.Logger().Error(report)

		/*
			Selanjutnya objek error tersebut kita tampilkan ke console dan
			juga ke browser dalam bentuk JSON
		*/
		//c.JSON(report.Code, report)

		// 3. Custom Error Page
		/*
			Untuk aplikasi non-web-service, akan LEBIH BAIK jika setiap terjadi error DIMUNCULKAN error page (alihkan ke halaman error),
			atau halaman khusus yang menampilkan informasi error.
		*/
		errPage := fmt.Sprintf("%d.html", report.Code) // "error_page")
		if err := c.File(errPage); err != nil {
			c.HTML(report.Code, "Errrrooooorrrrr")
		}
	}

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
}
