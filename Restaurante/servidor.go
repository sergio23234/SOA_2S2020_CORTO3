package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"github.com/rs/cors"
)
type to_kill struct
{
 Id string `json:"Id"`
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hola\": \"mundo\"}"))
	})

	mux.HandleFunc("/idkill", func(w http.ResponseWriter, r *http.Request) {		
		var p to_kill		
		fmt.Println("entro")
		w.Header().Set("Content-Type", "application/json")
		 s, err := ioutil.ReadAll(r.Body) 
    		if err != nil {
        		panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
    		}
    		err = json.Unmarshal(s, &p)
    		if err != nil {
        		panic(err) // This would normally be a normal Error http response but I've put this here so it's easy for you to test.
    		}
		fmt.Println(p.Id)
		var ret=principal.TerminarProceso(p.Id)	
                var inf = strconv.Itoa(ret); 	
		w.Write([]byte("{\"respuesta\": \"" + inf + "\"}"))
	})
	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST).
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)
}