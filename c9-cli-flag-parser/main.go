package main

/*
CLI Flag Parser (Kingpin v2)
Tidak jarang, sebuah aplikasi DALAM EKSEKUSINYA MEMBUTUHKAN ARGUMEN untuk DISISIPKAN, entah itu mandatory atau
tidak. Contohnya seperti berikut.

	$ ./main --port=3000

Pada chapter ini kita akan belajar cara PARSING ARGUMEN EKSEKUSI aplikasi. Parsing sebenarnya
BISA DILAKUKAN dengan CUKUP MEMAINKAN property os.Args, akan tetapi pada pembelajaran kali ini kita
AKAN MENGGUNAKAN 3rd party library github.com/alecthomas/kingpin untuk mempermudah pelaksanaannya
*/

// 1. Parsing Argument
/*
Kita akan buat aplikasi yang bisa menerima bentuk argument seperti berikut.

	# $ ./main ArgAppName <ArgPort>
	$ ./main "My Application" 4000

Argument ArgAppName mandatory, HARUS DIISI, sedangkan argument ArgPort adalah
opsional (ada nilai default-nya).

*/

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
)

var (
	// argAppName = kingpin.Arg("name", "Application name").Required().String()
	// argPort    = kingpin.Arg("port", "Web server port").Default("9000").Int()

	// 2. Penggunaan Kingpin Application Instance
	// app        = kingpin.New("App", "Simple app")
	// argAppName = app.Arg("name", "Application name").Required().String()
	// argPort    = app.Arg("port", "Web server port").Default("9000").Int()
	/*
		Statement kingpin.Arg() diganti dengan app.Arg().
	*/

	// 3. Parsing Flag
	// app         = kingpin.New("App", "Simple App")
	// flagAppName = app.Flag("name", "Application name").Required().String()
	// flagPort    = app.Flag("port", "Web server port").Short('p').Default("9000").Int()
	/*
		Method .Short() digunakan untuk mendefinisikan short flag. Pada kode di atas,
		flag port bisa ditulis dalam bentuk --port=value ataupun -p=value.
	*/

	// 4. Parsing Command
	app = kingpin.New("App", "Simple app")
	/*
		Method .Command() digunakan UNTUK MEMBUAT command. Pembuatan argument dan flag masing-masing command
		BISA DILAKUKAN SEPERTI BIASANYA, dengan mengakses method .Arg() atau .Flag(),
		hanya saja pengaksesannya LEWAT OBJEK command masing-masing.

		var commandSomething = app.Command("something", "do something")
		var commandSomethingArgX = commandSomething.Flag("x", "arg x").String()
		var commandSomethingFlagY = commandSomething.Flag("y", "flag y").String()

	*/
	// Command add, beserta flag dan argument-nya.
	commandAdd             = app.Command("add", "add new userx")
	commandAddFlagOverride = commandAdd.Flag("override", "override existing user").Short('o').Bool()
	commandAddArgUser      = commandAdd.Arg("user", "username").Required().String()
	// Command update, beserta argument-nya.
	commandUpdate           = app.Command("update", "update user")
	commandUpdateArgOldUser = commandUpdate.Arg("old", "old username").Required().String()
	commandUpdateArgNewUser = commandUpdate.Arg("new", "new username").Required().String()
	// Command delete, beserta flag dan argument-nya.
	commandDelete          = app.Command("delete", "delete user")
	commandDeleteFlagForce = commandDelete.Flag("force", "force deletion").Short('f').Bool()
	commandDeleteArgUser   = commandDelete.Arg("user", "username").Required().String()
)

/*
Statement kingpin.Arg() DIGUNAKAN UNTUK MENYIAPKAN objek penampung argument.
Tulis nama argument sebagai parameter pertama, dan deskripsi argument sebagai
parameter kedua. 2 Informasi tersebut nantinya akan muncul ketika flag --help digunakan.

Untuk aplikasi YANG MEME banyak argument, deklarasi variabel penampungnya HARUS DITULISKAN
BERURUTAN. Seperti contoh di atas argAppName merupakan argument pertama,
dan argPort adalah argument kedua.

Chain statement kingpin.Arg() dengan beberapa method yang tersedia SESUAI dengan KEBUTUHAN.
Berikut adalah penjelasan dari 4 method yang digunakan di atas.
- Method .Required() membuat argument yang ditulis menjadi mandatory.
	Jika tidak disisipkan maka muncul error.
- Method .String() menandakan bahwa argument ditampung dalam tipe string.
- Method .Default() digunakan untuk menge-set default value dari argument.
	Method ini adalah kebalikan dari .Required(). Jika default value di-set maka
	argument BOLEH TIDAK DIISI. Objek penampung akan berisi default value.
- Method .Int() menandakan bahwa argument ditampung dalam tipe int.

Perlu diketahui, dalam pendefinisian argument, penulisan statement-nya HARUS DIAKHIRI
dengan pemanggilan method .String(), .Int(), .Bool(), atau method tipe lainnya
yang di-support oleh kingpin. Lebih jelasnya silakan cek laman dokumentasi
https://godoc.org/github.com/alecthomas/kingpin#ArgClause


*/

func main() {
	// kingpin.Parse()

	// 2. Penggunaan Kingpin Application Instance
	//command, err := app.Parse(os.Args[1:]) // Juga, kingpin.Parse() diganti dengan app.Parse(), dengan pemanggilannya harus disisipkan os.Args[1:].
	/*
		Manfaatkan objek err kembalian app.Parse() untuk MEMBUAT CUSTOM ERROR handling.
		Atau bisa tetap gunakan default custom error handling milik kingpin,
		caranya dengan membungkus statement app.Parse() ke dalam kingpin.MustParse().

		kingpin.MustParse(app.Parse(os.Args[1:]))
	*/
	// if err != nil {
	// 	// handler error here ...
	// }

	// 2. Penggunaan Kingpin Application Instance
	/*
		Dari yang sudah kita praktekan di atas, fungsi-fungsi diakses langsung dari package kingpin (kingpin.Arg(), kingpin.Parse()).
		Kingpin menyediakan fasilitas untuk membuat -what-so-called- objek kingpin application. Lewat objek ini, semua fungsi yang
		biasa digunakan (seperti .Arg() atau .Parse()) BISA DIAKSES SEBAGAI METHOD.

		Kelebihan menggunakan kingpin application, kita bisa buat CUSTOM HANDLER untuk antisipasi error. Pada aplikasi yg sudah dibuat
		di atas, jika argument yang required tidak disisipkan dalam eksekusi binary, maka aplikasi langsung exit dan error muncul. Error sejenis ini bisa kita override jika menggunakan kingpin application.
	*/
	// kingpin.Parse()

	// appName := *argAppName
	// port := fmt.Sprintf(":%d", *argPort)

	// fmt.Printf("Starting %s at %s", appName, port)

	// e := echo.New()
	// e.GET("/index", func(c echo.Context) (err error) {
	// 	return c.JSON(http.StatusOK, true)
	// })

	// e.Logger.Fatal(e.Start(port))

	// 3. Parsing Flag
	/*
		Flag adalah argument yang lebih terstruktur. Golang SEBENARNYA SUDAH MENYEDIAKAN package flag, isinya API untuk parsing flag.

		Contoh argument:
				$ ./executable "My application" 4000

		Contoh flag:
				$ ./executable --name="My application" --port=4000
				$ ./executable --name "My application" --port 4000
				$ ./executable --name "My application" -p 4000

		Kita tetap menggunakan kingpin pada bagian ini. Pembuatan flag di kingpin tidak sulit, CUKUP GUNAKAN .Flag() (tidak menggunakan .Arg()).
		Contohnya seperti berikut.

			app         = kingpin.New("App", "Simple app")
			flagAppName = app.Flag("name", "Application name").Required().String()
			flagPort    = app.Flag("port", "Web server port").Short('p').Default("9000").Int()

			kingpin.MustParse(app.Parse(os.Args[1:]))

			Method .Short() digunakan untuk mendefinisikan short flag. Pada kode di atas,
			flag port bisa ditulis dalam bentuk --port=value ataupun -p=value.

			Penggunaan flag --help akan memunculkan keterangan mendetail tiap-tiap flag.


			...bersambung ke C.9.4. Parsing Command


	*/
	// kingpin.MustParse(app.Parse(os.Args[1:]))

	// appName := *flagAppName
	// port := fmt.Sprintf(":%d", *flagPort)

	// fmt.Printf("Starting %s at %s", appName, port)

	// e := echo.New()
	// e.GET("/index", func(c echo.Context) (err error) {
	// 	return c.JSON(http.StatusOK, true)
	// })

	// 4. Parsing Command
	/*
		Command adalah BENTUK YANG LEBIH ADVANCE dari argument. Banyak command bisa dibuat, pendefinisian flag ataupun argument
		bisa dilakukan lebih spesifik, untuk masing-masing command.

		contoh disini, akan dibuat aplikasi simulasi manajemen user, 3 buah command dibuat dengan skema berikut:
		Command add
				Flag --override
				Argument user
		Command update
				Argument old user
				Argument new user
		Command delete
				Flag --force
				Argument user


		Buat fungsi main, lalu di dalamnya siapkan action untuk masing-masing command.
		Gunakan method .Action() dengan parameter adalah fungsi ber-skema
		func(*kingpin.ParseContext)error untuk menambahkan action.
	*/

	commandAdd.Action(func(ctx *kingpin.ParseContext) error {
		user := *commandAddArgUser
		override := *commandAddFlagOverride
		fmt.Printf("adding user %s, override %t \n", user, override)

		return nil
	})

	commandUpdate.Action(func(ctx *kingpin.ParseContext) error {
		oldUser := *commandUpdateArgOldUser
		newUser := *commandUpdateArgNewUser
		fmt.Printf("updating user from %s %s \n", oldUser, newUser)

		return nil
	})

	commandDelete.Action(func(ctx *kingpin.ParseContext) error {
		user := *commandDeleteArgUser
		force := *commandDeleteFlagForce
		fmt.Printf("deleting user %s, force %t \n", user, force)

		return nil
	})

	kingpin.MustParse(app.Parse(os.Args[1:]))
	/*
		// 1. Parsing Argument
			run in cli utk cek argumen:
				go build main.go (akan membuat file .exe)
				./main --help

			nanti setelah file .exe nya jadi, jalankan file tersebut dengan argument
				./main "My Application" 3000


		// 3. Parsing Flag
		go build main.go
		./main --help (akan keluar argument yg bisa dipake dan flag akan muncul)

		// 4. Parsing Command
		go build main.go
		./main add --help
		Atau gunakan --help-long dalam eksekusi binary, untuk menampilkan help yang mendetail (argument dan flag tiap command juga dimunculkan).
		./main --help-long

	*/
}
