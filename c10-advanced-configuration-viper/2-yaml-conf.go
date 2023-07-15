package main

/*
Utk penggunaan configuration yaml sama seperti json, tinggal sesuain tipe nya
*/

import (
	"net/http"

	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

func main() {
	e := echo.New()

	viper.SetConfigType("yaml") // dengan tipe yaml
	viper.AddConfigPath(".")
	viper.SetConfigName("app.config")

	// 3. Watcher Configuration
	/*
		Viper memiliki banyak fitur, satu di antaranya adalah MENGAKTIFKAN WATCHER pada file konfigurasi.
		Dengan adanya watcher, maka kita bisa membuat callback yang AKAN DIPANGGIL setiap kali ADA PERUBAHAN
		konfigurasi.

		JANGAN LUPA IMPORT https://github.com/fsnotify/fsnotify
	*/
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	err := viper.ReadInConfig()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/index", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, true)
	})

	e.Logger.Print("Starting", viper.GetString("appName"))
	e.Logger.Fatal(e.Start(":" + viper.GetString("server.port")))
}
