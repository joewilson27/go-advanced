package main

// 7. Postgres SQL Store
/*
Pembuatan postgres store caranya kurang lebih SAMA DENGAN mongo store. Library yang dipakai
adalah github.com/antonlindstrom/pgstore.

Gunakan pgstore.NewPGStore() untuk MEMBUAT store. Isi parameter pertama dengan connection string postgres server,
lalu authentication key dan encryption key.
*/

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/context"
	"github.com/labstack/echo"
)

const SESSION_ID = "id"

func newPostgresStore() *pgstore.PGStore {
	url := "postgres://postgres:1234@127.0.0.1:5432/go-advanced?sslmode=disable"

	/*
		Let's break down the components:

		postgres://: This indicates that the connection is intended for a PostgreSQL database.
		username: The username to authenticate with the database.
		password: The password associated with the provided username. If there's no password, it can be left empty.
		@host: The IP address or hostname of the PostgreSQL server.
		:port: The port number on which the PostgreSQL server is listening.
		/database_name: The name of the PostgreSQL database to which you want to connect.
		?connection_options: Additional options can be specified using query parameters. In your case, sslmode=disable is used to disable SSL for the connection.
	*/

	authKey := []byte("my-auth-key-very-secret")
	encryptionKey := []byte("my-encryption-key-very-secret123")

	store, err := pgstore.NewPGStore(url, authKey, encryptionKey)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(0)
	}

	return store
}

func mainHold2() {
	store := newPostgresStore()

	e := echo.New()

	e.Use(echo.WrapMiddleware(context.ClearHandler))

	e.GET("/set", func(c echo.Context) error {
		session, _ := store.Get(c.Request(), SESSION_ID)
		session.Values["message1"] = "hello"
		session.Values["message2"] = "world"
		session.Save(c.Request(), c.Response()) // save session

		return c.Redirect(http.StatusTemporaryRedirect, "/get")
	})

	e.GET("/get", func(c echo.Context) error {
		session, _ := store.Get(c.Request(), SESSION_ID)

		if len(session.Values) == 0 {
			// jika session kosong / atau expired, balikan response ini
			return c.String(http.StatusOK, "empty result")
		}

		/*
			value session sebagai berikut:

			&{CBUKU4QSYCDDAGZYGGYA6LNIWWUPGB3TQUJYPIO335HNY56NUTTQ map[message1:hello message2:world] 0xc000098140 false 0xc000053880 id}

		*/

		// artinya session ada dan aktif, tampilkan session
		return c.String(http.StatusOK, fmt.Sprintf(
			"%s %s",
			session.Values["message1"],
			session.Values["message2"],
		))
	})

	e.GET("/delete", func(c echo.Context) error {
		session, _ := store.Get(c.Request(), SESSION_ID)
		session.Options.MaxAge = -1 // set session to expired for deleting session
		session.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusTemporaryRedirect, "/get")
	})

	e.Logger.Fatal(e.Start(":9000"))
}
