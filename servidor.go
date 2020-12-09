package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type Servidor struct {
	Usuarios       []string
	UsuariosChat   []*[]string
	UsuariosPort   []int
	NombresArchivo []string
	Archivo        [][]byte
	Chat           []string
}

var port int
var mensajes_servidor []string
var archivo_servidor []string

func (sr *Servidor) findUser(name string) int {
	for i, n := range sr.Usuarios {
		if name == n {
			return i
		}
	}
	return -1
}

func (sr *Servidor) AgregarUsuario(name string, new_port *int) error {
	if name != "" {
		port += 1
		sr.Usuarios = append(sr.Usuarios, name)
		var new_chat []string
		sr.UsuariosChat = append(sr.UsuariosChat, &new_chat)
		sr.UsuariosPort = append(sr.UsuariosPort, port)
		*new_port = port
		str := name + " se ha conectado"
		fmt.Println(str)
		return nil
	} else {
		str := name + "No te has podido conectar"
		return errors.New(str)
	}
}

func (sr *Servidor) EnviarArchivo(data [][]byte, reply *string) error {
	if string(data[0]) != "" {
		sr.NombresArchivo = append(sr.NombresArchivo, string(data[0]))
		sr.Archivo = append(sr.Archivo, data[1])
		archivo_servidor = sr.NombresArchivo
		*reply = "Recibido"
		fmt.Print("Archivo de: ")
		fmt.Println(string(data[2]))
		fmt.Println("Contenido: ")
		fmt.Println(string(data[1]))
		sr.EnviarArchivosUsuarios(data)
		return nil
	} else {
		return errors.New("Ocurrio un error al mandar tu archivo :c")
	}
}

func (sr *Servidor) Mensaje(data []string, reply *string) error {
	if data[1] != "" {
		sr.addChat(data)
		sr.ImprimirChat()
		*reply = "Enviado"
		return nil
	} else {
		return errors.New("Ocurrio un error al mandar tu mensaje, lo siento :c")
	}
}

func (sr *Servidor) EnviarArchivosUsuarios(data [][]byte) {
	i_user := sr.findUser(string(data[2]))

	for i, port := range sr.UsuariosPort {
		if i != i_user {
			port_str := "127.0.0.1:" + strconv.Itoa(port)
			c, err := rpc.Dial("tcp", port_str)
			if err != nil {
				fmt.Println(err)
				return
			}
			result := ""
			err = c.Call("ClientServidor.SetArchivo", data, &result)
			if err != nil {
				fmt.Println(err)
			}
			c.Close()
		}
	}
}

func (sr *Servidor) addChat(data []string) {
	sr.Chat = append(sr.Chat, data[0]+": "+data[1])
	mensajes_servidor = sr.Chat
	i_user := sr.findUser(data[0])

	for i, chat := range sr.UsuariosChat {
		if i == i_user {
			*chat = append(*chat, "Tu: "+data[1])
		} else {
			*chat = append(*chat, data[0]+": "+data[1])
		}
		port_str := "127.0.0.1:" + strconv.Itoa(sr.UsuariosPort[i])
		c, err := rpc.Dial("tcp", port_str)
		if err != nil {
			fmt.Println(err)
			return
		}
		result := ""
		err = c.Call("ClientServidor.SetChat", *chat, &result)
		if err != nil {
			fmt.Println(err)
		}
		c.Close()

	}
}

func (sr *Servidor) ImprimirChat() {
	for _, mssg := range sr.Chat {
		fmt.Println(mssg)
	}
}

func (sr *Servidor) End(name string, reply *string) error {
	i := sr.findUser(name)

	if i != -1 {
		copy(sr.Usuarios[i:], sr.Usuarios[i+1:])
		sr.Usuarios[len(sr.Usuarios)-1] = ""
		sr.Usuarios = sr.Usuarios[:len(sr.Usuarios)-1]
		copy(sr.UsuariosChat[i:], sr.UsuariosChat[i+1:])
		sr.UsuariosChat[len(sr.UsuariosChat)-1] = nil
		sr.UsuariosChat = sr.UsuariosChat[:len(sr.UsuariosChat)-1]
		*reply = "Se ha desconectado " + name
		fmt.Println(*reply)
		*reply = "Adios"
		return nil
	} else {
		str := name + "Ocurrio un error no te pudimos desconectar :c"
		return errors.New(str)
	}
}

func ImprimirChatLocal() {
	for _, mssg := range mensajes_servidor {
		fmt.Println(mssg)
	}
}
func ImprimirArchivos() {
	for _, mssg := range archivo_servidor {
		fmt.Println(mssg)
	}
}

func RespaldarArchivo(content, file_name string) {
	file, err := os.Create(file_name + ".txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.WriteString(content)
}

func Servidora() {
	new_Servidor := new(Servidor)
	rpc.Register(new_Servidor)
	port_str := ":" + strconv.Itoa(port)
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

func main() {
	port = 1306
	go Servidora()
	var op int64
	for {
		fmt.Println("1) Mensajes\n2) Archivos\n3) Respaldar chat\n4) Respaldar Archivos\n0) Salir")
		fmt.Scanln(&op)

		switch op {
		case 1:
			ImprimirChatLocal()
		case 2:
			ImprimirArchivos()
		case 3:
			RespaldarArchivo(strings.Join(mensajes_servidor, "\n"), "Chat")
			fmt.Println("Mensajes Respaldados")
		case 4:
			RespaldarArchivo(strings.Join(archivo_servidor, "\n"), "Archivo")
			fmt.Println("archivos Respaldados")
		case 0:
			return
		}
	}
}
