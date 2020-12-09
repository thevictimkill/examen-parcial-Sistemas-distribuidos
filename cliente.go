package main

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

var name string
var Port int
var Chat []string

type ClientServidor struct {
	Chat []string
}

func client() {
	c, err := rpc.Dial("tcp", "127.0.0.1:1306")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Usuario: ")
	fmt.Scanln(&name)
	err = c.Call("Servidor.AgregarUsuario", name, &Port)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Bienvenido")
	}
	go Servidor()

	var op int64
	for {
		fmt.Println("1) Enviar mensaje\n2) Enviar archivo\n3) Mostrar chat\n0) Salir")
		fmt.Scanln(&op)

		switch op {
		case 1:
			var result string
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Mensaje: ")
			scanner.Scan()
			mssg := scanner.Text()
			err = c.Call("Servidor.Mensaje", []string{name, mssg}, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
		case 2:
			var result string
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Ruta: ")
			scanner.Scan()
			file_name := scanner.Text()
			file, err := os.Open(file_name)

			if err != nil {
				fmt.Println(err)
				return
			}

			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				fmt.Println(err)
				return
			}

			total := stat.Size()

			bs := make([]byte, total)
			count, err := file.Read(bs)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("El archivo, contine: ", string(bs), "Tamaño: ", count)
			err = c.Call("Servidor.EnviarArchivo", [][]byte{[]byte(file_name), bs, []byte(name)}, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
		case 3:
			ImprimirChatlocal()
		case 0:
			var result string
			err = c.Call("Servidor.End", name, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
				return
			}
		}
	}
}

func Servidor() {
	new_Servidor := new(ClientServidor)
	rpc.Register(new_Servidor)
	port_str := ":" + strconv.Itoa(Port)
	ln, err := net.Listen("tcp", port_str)
	if err != nil {
		fmt.Println(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func (sr *ClientServidor) SetChat(chat []string, reply *string) error {
	sr.Chat = chat
	Chat = chat
	sr.ImprimirChat()
	return nil
}
func (sr *ClientServidor) SetArchivo(data [][]byte, reply *string) error {
	*reply = "Archivo recibido"
	fmt.Print("Archivo de: ")
	fmt.Println(string(data[2]))
	fmt.Println("Contenido: ")
	fmt.Println(string(data[1]))
	return nil
}
func (sr *ClientServidor) ImprimirChat() {
	for _, mssg := range sr.Chat {
		fmt.Println(mssg)
	}
}

func ImprimirChatlocal() {
	for _, mssg := range Chat {
		fmt.Println(mssg)
	}
}

func main() {
	client()
}
