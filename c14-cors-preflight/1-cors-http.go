package main

/*
C.14. CORS & Preflight Request
*/

// 1. Teori & Penerapan
/*
CORS adalah mekanisme untuk MEMBERI TAHU browser, apakah sebuah request yang di-dispatch dari aplikasi web domain lain atau origin lain,
ke aplikasi web kita itu DIPERBOLEHKAN ATAU TIDAK. Jika aplikasi kita tidak mengijinkan maka akan muncul error, dan request pasti digagalkan oleh browser.

CORS HANYA BERLAKU pada request-request yang DILAKUKAN LEWAT BROWSER, dari javascript; dan tujuan request-nya berbeda domain/origin.
Jadi request yang dilakukan dari curl maupun dari back end, tidak terkena dampak aturan CORS.
	*Request jenis ini biasa disebut dengan istilah cross-origin HTTP request.


Konfigurasi CORS dilakukan di RESPONSE HEADER aplikasi web. Penerapannya DI SEMUA BAHASA pemrograman yang web-based adalah SAMA, yaitu dengan memanipulasi
response header-nya. Berikut merupakan list header yang bisa digunakan untuk konfigurasi CORS.

	Access-Control-Allow-Origin
	Access-Control-Allow-Methods
	Access-Control-Allow-Headers
	Access-Control-Allow-Credentials
	Access-Control-Max-Age

Konfigurasi CORS BERADA di SISI SERVER, di aplikasi web tujuan request.
*/

// 2. Aplikasi dengan konfigurasi CORS sederhana
/*
Buat project baru, lalu isi fungsi main() dengan kode berikut. Aplikasi sederhana ini akan kita jalankan pada domain atau origin http://localhost:3000/,
lalu akan kita coba AKSES dari DOMAIN BERBEDA.
*/

import (
	"log"
	"net/http"
)

func mainHold() {
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://www.google.com")
		/*
			Kode di atas artinya request yang di-dispatch dari https://www.google.com DIIJINJAN untuk masuk (melakukan request ke server kita);
		*/
		// multiple allow origin
		//w.Header().Set("Access-Control-Allow-Origin", "https://www.google.com, https://apple.com")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT")
		/*
			Header Access-Control-Allow-Methods menentukan HTTP Method mana saja yang diperbolehkan masuk (penulisannya dengan pembatas koma).
		*/
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
		/*
			Header Access-Control-Allow-Headers menentukan key header mana saja yang diperbolehkan di-dalam request.
		*/

		/*
			Seperti yang sudah dijelaskan, bahwa konfigurasi CORS BERADA di HEADER RESPONSE. Pada kode di atas 3 buah property header untuk
			keperluan CORS digunakan.

			Header Access-Control-Allow-Origin digunakan untuk menentukan domain mana saja yang DIPERBOLEHKAN MENGAKSES aplikasi ini.
			Kita BISA SET value-nya dengan BANYAK ORIGIN, hal ini diperbolehkan dalam spesifikasi CORS namun sayangnya BANYAK BROWSER YANG TIDAK SUPPORT.

			Simulasi pada chapter ini adalah aplikasi web localhost:9000 DIAKSES dari google.com (eksekusi request sendiri kita lakukan dari browser
			dengan memanfaatkan developer tools milik chrome). BUKAN google.com diakses dari aplikasi web localhost:9000, jangan sampai dipahami terbalik.

			Khusus untuk beberapa header seperti Accept, Origin, Referer, dan User-Agent tidak terkena efek CORS, karena header-header tersebut secara otomatis di-set di setiap request.
		*/

		if r.Method == "OPTIONS" {
			w.Write([]byte("allowed"))
			return
		}

		w.Write([]byte("hello"))
	})

	log.Println("Starting app at :9000")
	http.ListenAndServe(":9000", nil)

	// 3. Testing CORS
	/*
		Lalu install extension jQuery Injector. Buka https://www.google.com lalu inject jQuery dengan tools extension tadi. Dengan melakukan inject
		jQuery secara paksa maka dari situs google kita bisa menggunakan jQuery.
		Buka chrome developer tools, klik tab console. Lalu jalankan perintah jQuery AJAX berikut.

		lalu lakukan hal yang sama dari web berbeda, misalnya buka web apple.com, maka akan ada error access-control-allow-origin cors


		Masih tetap error, tapi berbeda dengan error sebelumnya (ketika menggunakan multiple allow origin).
		Sebenarnya sudah kita singgung juga di atas, bahwa di spesifikasi adalah DIPERBOLEHKAN isi header Access-Control-Allow-Origin lebih dari satu website.
		Namun, KEBANYAKKAN browser TIDAK MENDUKUNG bagian ini. Oleh karena itu error di atas muncul. Konfigurasi ini termasuk tidak valid,
		hasilnya kedua website tersebut tidak punya ijin masuk.


		Allow All
		Gunakan tanda asteriks (*) sebagai nilai ketiga CORS header untuk memberi ijin ke semua.
		// semua origin mendapat ijin akses
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// semua method diperbolehkan masuk
		w.Header().Set("Access-Control-Allow-Methods", "*")

		// semua header diperbolehkan untuk disisipkan
		w.Header().Set("Access-Control-Allow-Headers", "*")
	*/

	// 4. Preflight Request
	/*
		Teori
		Dalam konteks CORS, request dikategorikan menjadi 2 yaitu, SIMPLE REQUEST dan Preflighted Request.
		Beberapa contoh request yang sudah kita pelajari di atas termasuk SIMPLE REQUEST.

		Ketika melakukan cross origin request dengan payload adalah JSON, atau request jenis lainnya, biasanya di developer tools -> network log MUNCUL 2 kali request,
		request pertama method-nya OPTIONS dan request ke-2 adalah actual request.

		Request ber-method OPTIONS tersebut DISEBUT dengan Preflight Request. Request ini AKAN OTOMATIS muncul ketika http request yang kita dispatch MEMENUHI KRITERIA
		preflighted request.

		Tujuan dari preflight request adalah UNTUK MENGECEK APAKAH DESTINASI mendukung CORS.
		Tiga buah informasi dikirimkan Access-Control-Request-Method, Access-Control-Request-Headers, dan Origin, dengan method adalah OPTIONS.


		Berikut merupakan kriteria preflighted request.

		Method yang digunakan adalah salah satu dari method berikut:
		PUT
		DELETE
		CONNECT
		OPTIONS
		TRACE
		PATCH

		Lebih detailnya mengenai simple dan preflighted request silakan baca https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS.

		Praktek
		lakukan test seperti sebelumnya dengan jQuery Injector dan request dengan console pada developer tools, lalu gunakan headers Content-Type.
		Nanti pada networks akan ada 2 request. yang pertama preflight request dan kedua actual request

	*/
}
