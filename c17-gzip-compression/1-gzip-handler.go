package main

/*
17. HTTP Gzip Compression (gziphandler)

Pada chapter ini kita akan mempelajari penerapan HTTP Compression,
dengan encoding adalah Gzip, dalam aplikasi web golang.
*/

// 17.1. Teori
/*
HTTP Compression adalah teknik KOMPRESI data pada HTTP response, agar ukuran/size output menjadi lebih kecil dan response time lebih cepat.

Pada saat sebuah endpoint diakses, di header request AKAN ADA header Accept-Encoding yang disisipkan oleh browser SECARA OTOMATIS.
example:
GET /hello HTTP/1.1
Host: localhost:9000
Accept-Encoding: gzip, deflate

Jika isinya adalah gzip atau deflate, berarti browser SIAP DAN SUPPORT untuk menerima response yang di-compress dari back end.

*Deflate adalah algoritma kompresi untuk data lossless.
*Gzip adalah salah satu teknik kompresi data yang menerapkan algoritma deflate.

Di sisi back end sendiri, jika memang output di-compress, maka response header Content-Encoding: gzip PERLU disisipkan.
example: Content-Encoding: gzip

Jika di sebuah request TIDAK ADA header Accept-Encoding: gzip, tetapi response back end tetap di-compress, maka akan muncul error di browser ERR_CONTENT_DECODING_FAILED

*/

import (
	"io"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
)

func main() {
	mux := new(http.ServeMux)

	mux.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("sample.png")
		if f != nil {
			defer f.Close()
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(w, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	server := new(http.Server)
	server.Addr = ":9000"
	server.Handler = gziphandler.GzipHandler(mux) // gzip handler

	server.ListenAndServe()
}
