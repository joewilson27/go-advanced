package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"github.com/globalsign/mgo"
	"github.com/gorilla/context"
	"github.com/kidstuff/mongostore"
	"github.com/labstack/echo"
)

var sessionManager *SessionManager

const SESSION_ID = "id"

type UserModel struct {
	ID   string
	Name string
	Age  int
}

func newMongoStore() *mongostore.MongoStore {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(0)
	}

	dbCollection := session.DB("learnwebgolang").C("session")
	maxAge := 86400 * 7
	ensureTTL := true
	authKey := []byte("my-auth-key-very-secret")
	encryptionKey := []byte("my-encryption-key-very-secret123")

	store := mongostore.NewMongoStore(dbCollection, maxAge, ensureTTL, authKey, encryptionKey)
	// Statement mongostore.NewMongoStore() digunakan untuk membuat mongo db store
	return store
}

func main() {
	gob.Register(UserModel{})

	sessionManager = NewSessionManager(newMongoStore())

	e := echo.New()

	e.Use(echo.WrapMiddleware(context.ClearHandler))
	/*
		Sesuai dengan README Gorilla Session, library ini JIKA DIGABUNG dengan library lain
		selain gorilla mux, akan BERPOTENSI menyebabkan memory leak.
		Untuk mengcover isu ini maka middleware context.ClearHandler PERLU diregistrasikan.
		Middleware tersebut berada dalam library Gorilla Context.
	*/

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		report, ok := err.(*echo.HTTPError)
		if !ok {
			report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		c.Logger().Error(report)
		c.JSON(report.Code, report)
	}

	eTest := e.Group("/test") // ini preffix

	eTest.GET("/set", func(c echo.Context) error {
		user := new(UserModel)
		user.ID = "001"
		user.Name = "Joe"
		user.Age = 12

		err := sessionManager.Set(c, SESSION_ID, user)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/test/get")
	})

	eTest.GET("/get", func(c echo.Context) error {
		result, err := sessionManager.Get(c, SESSION_ID)
		if err != nil {
			return err
		}

		if result == nil {
			return c.String(http.StatusOK, "empty result")
		} else {
			user := result.(UserModel)
			return c.JSON(http.StatusOK, user)
		}
	})

	eTest.GET("/delete", func(c echo.Context) error {
		err := sessionManager.Delete(c, SESSION_ID)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/test/get")
	})

	e.Logger.Fatal(e.Start(":9000"))
}
