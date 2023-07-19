package main

// 6. Mongo DB Store
/*
go get dulu
"github.com/gorilla/context"
"github.com/kidstuff/mongostore"
"github.com/labstack/echo"
"gopkg.in/mgo.v2" /////// not compatible to mongostore
"github.com/globalsign/mgo"

-- done install --

*/

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/globalsign/mgo"
	"github.com/gorilla/context"
	"github.com/kidstuff/mongostore"
	"github.com/labstack/echo"
)

const SESSION_ID = "id"

func newMongoStore() *mongostore.MongoStore {
	mgoSession, err := mgo.Dial("localhost:27017")
	/*
		Statement mgo.Dial() digunakan untuk TERHUBUNG dengan mongo db server.
		Method dial mengembalikan dua objek, salah satunya adalah mgo session.
	*/
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(0)
	}

	dbCollection := mgoSession.DB("learnwebgolang").C("session")
	/*
		Dari mgo session akses database lewat method .DB(), lalu akses collection
		yang ingin DIGUNAKAN sebagai MEDIA PENYIMPANAN data asli session lewat method .C().

		Database : learnwebgolang
		collection: session
	*/
	maxAge := 86400 * 7
	ensureTTL := true
	authKey := []byte("my-auth-key-very-secret")
	encryptionKey := []byte("my-encryption-key-very-secret123")

	store := mongostore.NewMongoStore(
		dbCollection,
		maxAge,
		ensureTTL,
		authKey,
		encryptionKey,
	)
	/*
		Statement mongostore.NewMongoStore() digunakan untuk MEMBUAT mongo db store.
		Ada beberapa parameter yang diperlukan: objek collection mongo di atas,
		dan dua lagi lainnya adalah authentication key dan encryption key.
	*/

	return store
}

func mainHold() {
	store := newMongoStore() // buat object Store

	e := echo.New()

	e.Use(echo.WrapMiddleware(context.ClearHandler))
	/*
		Sesuai dengan README Gorilla Session, library ini JIKA DIGABUNG dengan library lain
		selain gorilla mux, akan BERPOTENSI menyebabkan memory leak.
		Untuk mengcover isu ini maka middleware context.ClearHandler PERLU diregistrasikan.
		Middleware tersebut berada dalam library Gorilla Context.
	*/

	e.GET("/set", func(c echo.Context) error {
		session, _ := store.Get(c.Request(), SESSION_ID)
		/*
			Dari objek store store, dalam handler,
			kita bisa mengakses objek session dengan menyisipkan context http request.
		*/
		session.Values["message1"] = "hello"
		session.Values["message2"] = "world"
		session.Save(c.Request(), c.Response()) // save data to datbase

		return c.Redirect(http.StatusTemporaryRedirect, "/get")
	})

	e.GET("/get", func(c echo.Context) error {
		session, _ := store.Get(c.Request(), SESSION_ID)

		if len(session.Values) == 0 {
			return c.String(http.StatusOK, "empty result")
		}

		return c.String(http.StatusOK, fmt.Sprintf(
			"%s %s",
			session.Values["message1"],
			session.Values["message2"],
		))
	})

	e.GET("/delete", func(c echo.Context) error {
		session, _ := store.Get(c.Request(), SESSION_ID)
		session.Options.MaxAge = -1             // expired kan session
		session.Save(c.Request(), c.Response()) // simpan kembali sessionnya

		return c.Redirect(http.StatusTemporaryRedirect, "/get")
	})

	e.Logger.Fatal(e.Start(":9000"))
}
