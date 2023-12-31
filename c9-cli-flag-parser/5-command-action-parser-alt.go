package main

/*
5. Command Action Tanpa Menggunakan .Action()
Nilai balik statement kingpin.MustParse(), kingpin.Parse(), dan nilai balik pertama app.Parse()
ADALAH SAMA, yaitu informasi command yang ditulis pada saat pemanggilan binary.

Dari informasi command tersebut, bisa KITA KEMBANGKAN untuk MEMBUAT handler masing-masing command.
Dengan ini TAK PERLU menggunakan method .Action() untuk menulis handler command.
*/
import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
)

var app = kingpin.New("App", "Simple app")

var (
	commandAdd             = app.Command("add", "add new user")
	commandAddFlagOverride = commandAdd.Flag("override", "override existing user").Short('o').Bool()
	commandAddArgUser      = commandAdd.Arg("user", "username").Required().String()
)

var (
	commandUpdate           = app.Command("update", "update user")
	commandUpdateArgOldUser = commandUpdate.Arg("old", "old username").Required().String()
	commandUpdateArgNewUser = commandUpdate.Arg("new", "new username").Required().String()
)

var (
	commandDelete          = app.Command("delete", "delete user")
	commandDeleteFlagForce = commandDelete.Flag("force", "force deletion").Short('f').Bool()
	commandDeleteArgUser   = commandDelete.Arg("user", "username").Required().String()
)

func main() {
	commandInString := kingpin.MustParse(app.Parse(os.Args[1:]))
	/*
		jadi kita tidak perlu menggunakan .Action seperti pada sebelumnya,
		kita tinggal manfaatkan return value commandInString utk kemudian digunakan
		dalam switch
	*/
	switch commandInString {

	case commandAdd.FullCommand(): // add user
		user := *commandAddArgUser
		override := *commandAddFlagOverride
		fmt.Printf("adding user %s, override %t \n", user, override)

	case commandUpdate.FullCommand(): // update user
		oldUser := *commandUpdateArgOldUser
		newUser := *commandUpdateArgNewUser
		fmt.Printf("updating user from %s %s \n", oldUser, newUser)

	case commandDelete.FullCommand(): // delete user
		user := *commandDeleteArgUser
		force := *commandDeleteFlagForce
		fmt.Printf("deleting user %s, force %t \n", user, force)

	}
}
