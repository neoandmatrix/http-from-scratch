package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	// "os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string,1)

	go func ()  {
		defer f.Close()
		defer close(out)

		str := ""
	for{
		data := make([]byte,8)
		n,err := f.Read(data)
		if err != nil {
			break
		}
		data = data[:n]
		if i := bytes.IndexByte(data,'\n'); i!=-1 {
			str += string(data[:i])
			data = data[i+1:]
			out <- str
			//		fmt.Printf("read: %s\n",str)
		str = ""
		}
		str += string(data)
	}

	if len(str) !=0{
		// fmt.Printf("read: %s\n",str)
		out <- str
	}

	
}()

return  out

}

func main(){
	listner, err := net.Listen("tcp",":42069")
	// f,err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error","error",err)
	}

	for {
		conn,err := listner.Accept()
		if err != nil {
			log.Fatal("error","error",err)
		}
			lines := getLinesChannel(conn) // reading 8 bytes form a connection instead of priviously files
	for line := range lines {
		fmt.Printf("read: %s\n",line)
	}
	}
	
	
}