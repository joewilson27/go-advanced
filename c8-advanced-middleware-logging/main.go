package main

/*
Advanced Middleware & Logging (Logrus, Echo Logger)
Middleware adalah sebuah blok kode yang DIPANGGIL SEBELUM ataupun SESUDAH http request di-proses.
Middleware biasanya dibuat per-fungsi-nya, contohnya: middleware autentikasi, middleware untuk logging,
middleware untuk gzip compression, dan lainnya.
*/

// 1. Custom Middleware
/*
Pembuatan middleware pada echo sangat mudah, CUKUP GUNAKAN method .Use() milik objek echo untuk
registrasi middleware. Method ini BISA DIPANGGIL BERKALI-KALI, dan eksekusi middleware-nya sendiri
adalah berurutan sesuai dengan urutan registrasi.
*/

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func middlewareOne(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("from middleware one")
		return next(c)
	}
}

func middlewareTwo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("from middleware two")
		return next(c)
	}
}

// 2. Integrasi Middleware ber-skema Non-Echo-Middleware
func middlewareSomething(next http.Handler) http.Handler {
	/*
		Bisa dilihat, fungsi middlewareSomething TIDAK MENGGUNAKAN skema middleware milik echo
		melainkan skema  func(http.Handler)http.Handler
	*/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("from middleware something")
		next.ServeHTTP(w, r)
	})
}

// 4. 3rd Party Logging Middleware: Logrus
func makeLogEntry(c echo.Context) *log.Entry {
	if c == nil {
		return log.WithFields(log.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return log.WithFields(log.Fields{
		"at":     time.Now().Format("2006-01-02 15:04:05"),
		"method": c.Request().Method,
		"uri":    c.Request().URL.String(),
		"ip":     c.Request().RemoteAddr,
	})
}

/*
Fungsi makeLogEntry() bertugas MEMBUAT BASIS LOG objek yang akan ditampilkan.
Informasi standar seperti waktu, dibentuk di dalam fungsi ini. Khusus untuk
log yang berhubungan dengan http request, maka informasi yang lebih detail
dimunculkan (http method, url, dan IP).
*/

func middlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		makeLogEntry(c).Info("incoming request")
		return next(c)
	}
}

/*
Fungsi middlewareLogging() bertugas untuk menampilkan log setiap ada http request
masuk. Dari objek *log.Entry -yang-dicetak-lewat-fungsi-makeLogEntry()- (return
dari fungsi makeLogEntry), panggil method Info() untuk menampilkan pesan
log dengan level adalah INFO.
*/

func errorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if ok {
		report.Message = fmt.Sprintf("http error %d - %v", report.Code, report.Message)
	} else {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	makeLogEntry(c).Error(report.Message)
	c.HTML(report.Code, report.Message.(string))
}

/*
Sedang fungsi errorHandler akan digunakan untuk meng-OVERRIDE default http error handler milik echo.
Dalam fungsi ini log dengan level ERROR dimunculkan lewat pemanggilan method Error() milik *log.Entry.
*/

func main() {
	e := echo.New()

	// middleware here
	// e.Use(middlewareOne)
	// e.Use(middlewareTwo)

	// 2. Integrasi Middleware ber-skema Non-Echo-Middleware
	/*
		Di echo, fungsi middleware HARUS MEMILIKI skema func(echo.HandlerFunc)echo.HandlerFunc
		(bisa dilihat pada method middlewareOne & middlewarTwo). Untuk 3rd party middleware,
		TETAP BISA DIKOMBINASIKAN dengan echo, namun MEMBUTUHKAN sedikit PENYESUAIAN tentunya.

		Echo menyediakan solusi mudah untuk membantu integrasi 3rd party middleware,
		yaitu dengan menggunakan fungsi echo.WrapMiddleware() untuk mengkonversi
		middleware menjadi echo-compatible-middleware, DENGAN SYARAT skema harus
		dalam bentuk func(http.Handler)http.Handler.
	*/
	// e.Use(echo.WrapMiddleware(middlewareSomething))
	/*
		Bisa dilihat, fungsi middlewareSomething TIDAK MENGGUNAKAN skema middleware milik echo,
		namun tetap bisa digunakan dalam .Use() dengan cara dibungkus fungsi echo.WrapMiddleware().
	*/

	// 3. Echo Middleware: Logger
	/*
		Echo merupakan framework besar, di dalamnya terdapat BANYAK dependency dan library, salah satunya adalah
		logging middleware.

		Cara menggunakan logging middleware (ataupun middleware lainnya milik echo) adalah dengan meng-IMPORT
		package github.com/labstack/echo/middleware, lalu panggil nama middleware nya.
		Lebih detailnya silakan baca dokumentasi echo mengenai middleware di https://echo.labstack.com/middleware
	*/
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "method=${method}, uri=${uri}, status=${status}\n", // tulis konfigurasinya sebagai property objek cetakan  middleware.LoggerConfig
	// }))
	/*
		Cara menggunakan echo logging middleware adalah dengan membuat objek logging baru lewat statement middleware.Logger(),
		lalu membungkusnya dengan e.Use(). Atau bisa juga menggunakan middleware.LoggerWithConfig() JIKA logger yang dibuat MEMERLUKAN BEBERAPA
		konfigurasi (tulis konfigurasinya sebagai property objek cetakan middleware.LoggerConfig, lalu tempatkan sebagai
		parameter method pemanggilan .LoggerWithConfig()).
	*/
	// e.GET("/index", func(c echo.Context) (err error) {
	// 	fmt.Println("threeeeee!")

	// 	return c.JSON(http.StatusOK, true)
	// })

	// e.Logger.Fatal(e.Start(":9000"))

	// 4. 3rd Party Logging Middleware: Logrus
	/*
		Selain dengan membuat middleware sendiri, ataupun menggunakan echo middleware, kita juga bisa menggunakan 3rd party middleware lain.
		Tinggal sesuaikan sedikit agar sesuai dengan skema fungsi middleware milik echo untuk bisa digunakan.

		Next, kita akan coba untuk meng-implementasi salah satu golang library terkenal untuk keperluan logging, yaitu logrus https://github.com/sirupsen/logrus.
	*/
	e.Use(middlewareLogging)
	e.HTTPErrorHandler = errorHandler

	e.GET("/index", func(c echo.Context) error {
		return c.JSON(http.StatusOK, true)
	})

	lock := make(chan error)
	/*
		Web server di start dalam sebuah goroutine. Karena method .Start() milik echo
		adalah blocking, kita manfaatkan nilai baliknya untuk di kirim ke channel lock.
	*/
	go func(lock chan error) { lock <- e.Start(":9000") }(lock)

	time.Sleep(1 * time.Millisecond)
	makeLogEntry(nil).Warning("application started without ssl/tls enabled") // warning simulasi, bukan error/warning sebenarnya
	/*
		Selanjutnya dengan delay waktu 1 milidetik, log dengan level WARNING dimunculkan.
		Ini hanya simulasi saja, karena memang aplikasi tidak di start menggunakan
		ssl/tls. Dengan memberi delay 1 milidetik, maka log WARNING bisa muncul setelah
		log default dari echo muncul.
	*/

	err := <-lock
	if err != nil {
		makeLogEntry(nil).Panic("failed to start application")
	}
	/*
		Nah pada bagian penerimaan channel, jika nilai baliknya tidak nil MAKA PASTI
		terjadi error pada saat start web server, dan pada saat itu juga munculkan log
		dengan level PANIC.
	*/

}
