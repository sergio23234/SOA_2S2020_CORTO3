package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

type pedido struct {
	Desc string `json:"Desc"`
	Id   string `json:"Id"`
}
type resultado struct {
	Resultado string `json:"respuesta"`
}

var p pedido

func main() {
	servidor()
}
func servidor() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if (pedido{}) == p {
			w.Write([]byte("{\"respuesta\": \"entregado o aun no recibido\"}"))
		} else {
			w.Write([]byte("{\"respuesta\": \"en camino\"}"))
		}
	})
	mux.HandleFunc("/cli-res", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var recibidod string = recibir_signal1()
		w.Write([]byte("{\"respuesta\": \"" + recibidod + "\"}"))

	})
	mux.HandleFunc("/recibido-rep", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		enviar_signal()
		w.Write([]byte("{\"respuesta\": \"llego al orquestador repartidor\"}"))

	})
	mux.HandleFunc("/entregado-rep", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		enviar_signal_entregado()
		w.Write([]byte("{\"respuesta\": \"llego al orquestador repartidor\"}"))

	})
	mux.HandleFunc("/cli-rep", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var recibidod string = recibir_signal2()
		w.Write([]byte("{\"respuesta\": \"" + recibidod + "\"}"))
	})
	mux.HandleFunc("/idpedido-cli", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var ped pedido
		s, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
		}
		err = json.Unmarshal(s, &ped)
		if err != nil {
			panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
		}
		enviar_signal_cli(ped)
		w.Write([]byte("{\"respuesta\": \"enviado al restaurante\"}"))
	})
	mux.HandleFunc("/idpedido-res", func(w http.ResponseWriter, r *http.Request) {
		var ped pedido
		w.Header().Set("Content-Type", "application/json")
		s, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
		}
		err = json.Unmarshal(s, &ped)
		if err != nil {
			panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
		}
		enviar_signal_res(ped)
		w.Write([]byte("{\"respuesta\": \"enviado al repartidor\"}"))
	})
	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST).
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":9093", handler)
}

/*Notificar al cliente*/
func enviar_signal() {
	jjson := `{"respuesta":"En camino"}`
	b := strings.NewReader(jjson)
	resp, err := http.Post("http://localhost:9092/recibido", "application/json", b)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var r resultado
	err = json.Unmarshal(body, &r)
	if err != nil {
		panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
	}
	fmt.Println(r.Resultado)
}

func enviar_signal_entregado() {
	jjson := `{"respuesta":"Entregado"}`
	b := strings.NewReader(jjson)
	resp, err := http.Post("http://localhost:9092/entregado", "application/json", b)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var r resultado
	err = json.Unmarshal(body, &r)
	if err != nil {
		panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
	}
	p = pedido{}
	fmt.Println(r.Resultado)
}

/*idpedido-cli: enviar pedido del cliente al restaurante.*/
func enviar_signal_cli(buf pedido) { //enviar pedido al restaurante
	var r resultado
	jjson := `{"Id":"` + buf.Id + `","Desc":"` + buf.Desc + `"}`
	b := strings.NewReader(jjson)
	resp, err := http.Post("http://localhost:9091/idpedido", "application/json", b)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
	}
	fmt.Println(r.Resultado)
}

/*idpedido-res: enviar pedido del restaurante al repartidor.*/
func enviar_signal_res(buf pedido) { //enviar pedido al restaurante
	var r resultado
	jjson := `{"Id":"` + buf.Id + `","Desc":"` + buf.Desc + `"}`
	b := strings.NewReader(jjson)
	resp, err := http.Post("http://localhost:9090/idpedido", "application/json", b)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
	}
	fmt.Println(r.Resultado)
}
func recibir_signal1() string { //restaurante
	var r resultado
	resp, err := http.Get("http://localhost:9091/")
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
	}
	return r.Resultado
}
func recibir_signal2() string { //repartidor
	var r resultado
	resp, err := http.Get("http://localhost:9090/")
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
	}
	return r.Resultado
}

/*
idpedido-res: enviar pedido del restaurante al repartidor.
idpedido-cli: enviar pedido del cliente al restaurante.
cli-res: verificar pedido al restaurante
cli-rep: verificar pedido al repartidor
recibido-rep: enviar notificacion al cliente que esta con el repartidor.
*/
