package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type Proceso struct {
	Id        int
	Contador  int
	Terminado bool
}

var procesoLocal Proceso

func imprimirProceso() {
	for {
		if procesoLocal.Terminado {
			return
		}
		if procesoLocal.Id != -1 {
			fmt.Println("id ", procesoLocal.Id, ": ", procesoLocal.Contador)
			procesoLocal.Contador += 1
			time.Sleep(time.Millisecond * 500)
		}
	}
}
func getProceso() {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(procesoLocal)
	if err != nil {
		fmt.Println(err)
	}
	c.Close()
	go imprimirProceso()
}

func listenerCliente() {
	s, err := net.Listen("tcp", ":9998")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		err2 := gob.NewDecoder(c).Decode(&procesoLocal)
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		if procesoLocal.Id != -1 {
			s.Close()
			return
		}
	}
}
func cerrarCliente() {
	fmt.Println("cerrarndo cliente")
	procesoLocal.Terminado = true
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(procesoLocal)
	if err != nil {
		fmt.Println(err)
	}
	c.Close()
}

func main() {
	procesoLocal.Id = -1

	go listenerCliente()
	go getProceso()
	defer cerrarCliente()
	var aux string
	fmt.Scanln(&aux)

}
