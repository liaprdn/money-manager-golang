package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

// deklarasi variable dengan struct
type Catatan struct {
	Id_transaksi     int
	Tgl_pengeluaran  string
	Nominal          int
	Jenis_kebutuhan  string
	Sumber_dana      string
	Nama_pengeluaran string
}

// koneksi ke database
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_project_ap3")
	if err != nil {
		panic(err.Error())
	}
	return db
}

// jalanin semua file html di folder form
var tmpl = template.Must(template.ParseGlob("display/*"))

// index halaman
func Index(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Index", nil)
}

// fungsi menampilkan halaman catatan transaksi
func Tampil(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM catatan ORDER BY id_transaksi DESC")
	if err != nil {
		panic(err.Error())
	}
	cttn := Catatan{}
	res := []Catatan{}
	for selDB.Next() {
		var id_transaksi int
		var tgl_pengeluaran string
		var nominal int
		var jenis_kebutuhan, sumber_dana, nama_pengeluaran string
		err = selDB.Scan(&id_transaksi, &tgl_pengeluaran, &nominal, &jenis_kebutuhan, &sumber_dana, &nama_pengeluaran)
		if err != nil {
			panic(err.Error())
		}
		cttn.Id_transaksi = id_transaksi
		cttn.Tgl_pengeluaran = tgl_pengeluaran
		cttn.Nominal = nominal
		cttn.Jenis_kebutuhan = jenis_kebutuhan
		cttn.Sumber_dana = sumber_dana
		cttn.Nama_pengeluaran = nama_pengeluaran
		res = append(res, cttn)
	}
	tmpl.ExecuteTemplate(w, "Tampil", res)
	defer db.Close()
}

// menampilkan halaman tambah
func Tambah(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Tambah", nil)
}

// insert data ke db lewat halaman tambah pake form action insert
func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		tgl_pengeluaran := r.FormValue("tgl_pengeluaran")
		nominal := r.FormValue("nominal")
		jenis_kebutuhan := r.FormValue("jenis_kebutuhan")
		sumber_dana := r.FormValue("sumber_dana")
		nama_pengeluaran := r.FormValue("nama_pengeluaran")
		insForm, err := db.Prepare("INSERT INTO catatan(tgl_pengeluaran, nominal, jenis_kebutuhan, sumber_dana, nama_pengeluaran) VALUES(?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(tgl_pengeluaran, nominal, jenis_kebutuhan, sumber_dana, nama_pengeluaran)
		log.Println("INSERT: Tanggal Pengeluaran: " + tgl_pengeluaran + " | Nominal: " + nominal + " | Jenis Kebutuhan: " + jenis_kebutuhan + " | Sumber Dana: " + sumber_dana + " | Nama Pengeluaran: " + nama_pengeluaran)
	}
	defer db.Close()
	http.Redirect(w, r, "/tampil", 301)
}

// edit data berdasarkan id untuk halaman edit
func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id_transaksi")
	selDB, err := db.Query("SELECT * FROM Catatan WHERE id_transaksi=?", nId)
	if err != nil {
		panic(err.Error())
	}
	cttn := Catatan{}
	for selDB.Next() {
		var id_transaksi int
		var tgl_pengeluaran string
		var nominal int
		var jenis_kebutuhan, sumber_dana, nama_pengeluaran string
		err = selDB.Scan(&id_transaksi, &tgl_pengeluaran, &nominal, &jenis_kebutuhan, &sumber_dana, &nama_pengeluaran)
		if err != nil {
			panic(err.Error())
		}
		cttn.Id_transaksi = id_transaksi
		cttn.Tgl_pengeluaran = tgl_pengeluaran
		cttn.Nominal = nominal
		cttn.Jenis_kebutuhan = jenis_kebutuhan
		cttn.Sumber_dana = sumber_dana
		cttn.Nama_pengeluaran = nama_pengeluaran
	}
	tmpl.ExecuteTemplate(w, "Edit", cttn)
	defer db.Close()
}

// update data berdasarkan id untuk halaman edit di form action update
func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		tgl_pengeluaran := r.FormValue("tgl_pengeluaran")
		nominal := r.FormValue("nominal")
		jenis_kebutuhan := r.FormValue("jenis_kebutuhan")
		sumber_dana := r.FormValue("sumber_dana")
		nama_pengeluaran := r.FormValue("nama_pengeluaran")
		id_transaksi := r.FormValue("id_transaksi")

		insForm, err := db.Prepare("UPDATE catatan SET tgl_pengeluaran=?, nominal=?, jenis_kebutuhan=?, sumber_dana=?, nama_pengeluaran=? WHERE id_transaksi=?")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(tgl_pengeluaran, nominal, jenis_kebutuhan, sumber_dana, nama_pengeluaran, id_transaksi)

		log.Println("UPDATE: Tanggal Pengeluaran: " + tgl_pengeluaran + " | Nominal: " + nominal + " | Jenis Kebutuhan : " + jenis_kebutuhan + " | Sumber Dana: " + sumber_dana + " | Nama Pengeluaran: " + nama_pengeluaran)
	}
	defer db.Close()
	http.Redirect(w, r, "/tampil", 301)
}

// hapus data berdasarkan id
func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	cttn := r.URL.Query().Get("id_transaksi")
	delForm, err := db.Prepare("DELETE FROM Catatan WHERE id_transaksi=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(cttn)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// program utama
func main() {
	log.Println("Server started on: http://localhost:1393")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/", Index)
	http.HandleFunc("/tampil", Tampil)
	http.HandleFunc("/tambah", Tambah)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":1393", nil)
}
