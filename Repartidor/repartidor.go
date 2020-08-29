package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"github.com/rs/cors"
)
type pedido struct
{
 Desc 	string 	`json:"Desc"`
 Id 	string 	`json:"Id"`
}
type resultado struct
{
	Resultado 	string 	`json:"respuesta"`
}
var p pedido
func main() {
	go menu()
	servidor()
}
func menu(){
	var numero int = 0
	for(numero!=4){
		fmt.Println("1.Informar al cliente")
		fmt.Println("2.Marcar como entregado")
		fmt.Println("3.Salir")
		fmt.Scanf("%d\n", &numero)
		if(numero==1){
			enviar_signal()
		}else if(numero ==2){
			enviar_signal_entregado()
		}else if(numero ==3){
			os.Exit(1)
		}else{
			fmt.Println("error")
		}
	}
}
func servidor() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if ((pedido{})==p){
			w.Write([]byte("{\"respuesta\": \"entregado o aun no recibido\"}"))
		}else{
			w.Write([]byte("{\"respuesta\": \"en camino\"}"))
		}
	})

	mux.HandleFunc("/idpedido", func(w http.ResponseWriter, r *http.Request) {				
		w.Header().Set("Content-Type", "application/json")
		 s, err := ioutil.ReadAll(r.Body) 
    		if err != nil {
        		panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
    		}
    		err = json.Unmarshal(s, &p)
    		if err != nil {
        		panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
    		}
		w.Write([]byte("{\"respuesta\": \"repartidor recibido\"}"))
	})
	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST).
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":9090", handler)
}

func enviar_signal(){
	jjson := `{"respuesta":"En camino"}`
	b := strings.NewReader(jjson)
	resp, err := http.Post("http://localhost:9092/recibido","application/json",b)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var r resultado
	err = json.Unmarshal(body,&r)
    		if err != nil {
        		panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
    		}
	fmt.Println(r.Resultado)
}
func enviar_signal_entregado(){
	jjson := `{"respuesta":"Entregado"}`
	b := strings.NewReader(jjson)
	resp, err := http.Post("http://localhost:9092/entregado","application/json",b)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var r resultado
	err = json.Unmarshal(body,&r)
    		if err != nil {
        		panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
			}
	p = pedido{}
	fmt.Println(r.Resultado)
}