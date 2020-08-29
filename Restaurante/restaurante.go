package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strings"
	"os"
	"github.com/rs/cors"
)
type pedido struct
{
 Desc 	string 	`json:"Desc"`
 Id 	string	`json:"Id"`
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
		fmt.Println("1.Enviar pedido al repartidor")
		fmt.Println("2.Salir")
		fmt.Scanf("%d\n", &numero)
		if(numero==1){
			enviar_signal(p)
		}else if(numero ==2){
			os.Exit(1)
		}else{
			fmt.Println("error")
		}
	}
}
func servidor(){
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if ((pedido{})==p){
			w.Write([]byte("{\"respuesta\": \"ya en el repartidor\"}"))
		}else{
			w.Write([]byte("{\"respuesta\": \"aun en restaurante\"}"))
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
		w.Write([]byte("{\"respuesta\": \"recibido\"}"))
	})
	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST).
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":9091", handler)
}
func enviar_signal(buf pedido){
	jjson := `{"Id":"`+buf.Id+`","Desc":"`+buf.Desc+`"}`
	b := strings.NewReader(jjson)
	resp, err := http.Post("http://localhost:9090/idpedido","application/json",b)
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
	p = pedido{}
}