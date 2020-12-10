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

var listaProcesos []*Proceso

func imprimirProcesos() {
	for {
		fmt.Println("")
		fmt.Println("---------------------")
		for _, p := range listaProcesos {
			fmt.Println("id ", p.Id, ": ", p.Contador)
			p.Contador += 1
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func servidor() {
	go imprimirProcesos()
	s, err := net.Listen("tcp", ":9999")
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
		go handleClient(c)
	}
}

func handleClient(c net.Conn) {
	var p Proceso
	err := gob.NewDecoder(c).Decode(&p)
	if err != nil {
		fmt.Println(err)
		return
	} else if p.Terminado == false {
		c2, err2 := net.Dial("tcp", ":9998")
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		err3 := gob.NewEncoder(c2).Encode(listaProcesos[0])
		if err3 != nil {
			fmt.Println(err3)
		}
		c2.Close()
		listaProcesos = append(listaProcesos[:0], listaProcesos[1:]...)
	} else {
		p.Terminado = false
		fmt.Println("Procesos en servidor")
		listaProcesos = append(listaProcesos, &p)
	}
}

func main() {
	for i := 0; i <= 5; i++ {
		listaProcesos = append(listaProcesos, &Proceso{
			Id:        i,
			Contador:  0,
			Terminado: false,
		})
	}
	go servidor()
	var aux string
	fmt.Scan(&aux)
}
