package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var GlobalMaterias map[string]map[string]float64
var GlobalAlumno map[string]map[string]float64
var PromTodos map[string][]float64

type Registro struct {
	Alumno       string
	Materia      string
	Calificacion float64
}

func cargarHTML(a string) string {
	html, _ := ioutil.ReadFile(a)

	return string(html)
}

func agregarAGlobal(datos Registro) {
	materiasDelAlumno, existeElAlumno := GlobalAlumno[datos.Alumno]
	materia, existeLaMateria := GlobalMaterias[datos.Materia]

	if existeElAlumno {
		materiasDelAlumno[datos.Materia] = datos.Calificacion
	} else {
		materiasDelAlumno = make(map[string]float64)
		materiasDelAlumno[datos.Materia] = datos.Calificacion
		GlobalAlumno[datos.Alumno] = materiasDelAlumno
	}

	if existeLaMateria {
		materia[datos.Alumno] = datos.Calificacion
	} else {
		materia = make(map[string]float64)
		materia[datos.Alumno] = datos.Calificacion
		GlobalMaterias[datos.Materia] = materia
	}

	for key, element := range GlobalMaterias {
		fmt.Println("Materia:", key)
		for keyAl, Calificacion := range element {
			fmt.Println("	Alumno:", keyAl, "=>", "Calificacion:", Calificacion)
		}
	}
	fmt.Println("*****************")
	fmt.Println("Calificacion Agregada")
}

func MostrarTodosCalificaciones() string {
	var html string
	for materia, element := range GlobalMaterias {
		for keyAl, Calificacion := range element {
			html += "<tr>" +
				"<td>" + keyAl + "</td>" +
				"<td>" + materia + "</td>" +
				"<td>" + fmt.Sprintf("%g", Calificacion) + "</td>" +
				"</tr>"
		}
	}
	return html
}

func MostrarTodosPromedios() string {
	PromTodos = make(map[string][]float64)
	for _, element := range GlobalMaterias {
		for keyAl, Calificacion := range element {
			_, existeElPromedio := PromTodos[keyAl]
			if existeElPromedio {
				PromTodos[keyAl][0] += Calificacion
				PromTodos[keyAl][1]++
			} else {
				PromTodos[keyAl] = append(PromTodos[keyAl], Calificacion)
				PromTodos[keyAl] = append(PromTodos[keyAl], 1)
			}

		}
	}
	var html string
	for key, element := range PromTodos {
		html += "<tr>" +
			"<td>" + key + "</td>" +
			"<td>" + fmt.Sprintf("%.2f", element[0]/element[1]) + "</td>" +
			"<td>" + fmt.Sprintf("%g", element[1]) + "</td>" +
			"</tr>"
	}
	return html
}

func Registros(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		f, _ := strconv.ParseFloat(req.FormValue("calificacion"), 64)
		datos := Registro{Alumno: req.FormValue("alumno"), Materia: req.FormValue("materia"), Calificacion: f}
		agregarAGlobal(datos)
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHTML("paginasWeb/tablaRegistros.html"),
			MostrarTodosCalificaciones(),
		)
		// case "GET":
		// 	res.Header().Set(
		// 		"Content-Type",
		// 		"text/html",
		// 	)
		// 	fmt.Fprintf(
		// 		res,
		// 		cargarHTML("paginasWeb/tabla.html"),
		// 		MostrarTodosCalificaciones(),
		// 	)
	}
}

func PromAlum(alumno string) string {
	total := 0.0
	cantMaterias := 0.0
	for key, element := range GlobalMaterias {
		fmt.Println("Materia:", key)
		for keyAl, Calificacion := range element {
			if keyAl == alumno {
				fmt.Println("	Alumno:", keyAl, "=>", "Calificacion:", Calificacion)
				cantMaterias++
				total += Calificacion
			}
		}
	}
	promedio := total / cantMaterias
	fmt.Println("*****************")
	promedioString := alumno + ": " + fmt.Sprintf("%.2f", promedio)
	return promedioString
}

func PromAlumno(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		alumno := req.FormValue("alumno")
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHTML("paginasWeb/promAlumno.html"),
			PromAlum(alumno),
		)
	}
}

func PromGeneralMateria(materia string) string {
	total := 0.0
	cantAlumnos := 0.0
	for key, element := range GlobalMaterias {
		if key == materia {
			for keyAl, Calificacion := range element {
				fmt.Println("	Alumno:", keyAl, "=>", "Calificacion:", Calificacion)
				cantAlumnos++
				total += Calificacion

			}
		}
	}
	promedio := total / cantAlumnos
	fmt.Println("*****************")
	promedioString := materia + ": " + fmt.Sprintf("%.2f", promedio)
	return promedioString
}

func PromMateria(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		materia := req.FormValue("materia")
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHTML("paginasWeb/promAlumno.html"),
			PromGeneralMateria(materia),
		)
	}
}

func PromGeneral(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHTML("paginasWeb/tablaTodosPromedios.html"),
			MostrarTodosPromedios(),
		)
	}
}

func root(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHTML("paginasWeb/form.html"),
	)
}

func main() {
	GlobalAlumno = make(map[string]map[string]float64)
	GlobalMaterias = make(map[string]map[string]float64)
	http.HandleFunc("/", root)
	http.HandleFunc("/registros", Registros)
	http.HandleFunc("/promedioAlumno", PromAlumno)
	http.HandleFunc("/promGeneral", PromGeneral)
	http.HandleFunc("/promMateria", PromMateria)
	fmt.Println("Arrancando el servidor...")
	http.ListenAndServe(":9000", nil)
}
