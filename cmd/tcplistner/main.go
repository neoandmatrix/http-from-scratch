package main

import (
	"fmt"
	"httpformscratch/internal/request"
	"log"
	"net"
)

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	out := make(chan string,1)

// 	go func ()  {
// 		defer f.Close()
// 		defer close(out)

// 		str := ""
// 	for{
// 		data := make([]byte,8)
// 		n,err := f.Read(data)
// 		if err != nil {
// 			break
// 		}
// 		data = data[:n]
// 		if i := bytes.IndexByte(data,'\n'); i!=-1 {
// 			str += string(data[:i])
// 			data = data[i+1:]
// 			out <- str
// 			//		fmt.Printf("read: %s\n",str)
// 		str = ""
// 		}
// 		str += string(data)
// 	}

// 	if len(str) !=0{
// 		// fmt.Printf("read: %s\n",str)
// 		out <- str
// 	}

	
// }()

// return  out

// }

// TODO: headers and body not workig currntly have to fix them

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

		r,err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error",err)
		}

		fmt.Printf("Request Line:\n")
		fmt.Printf("- Method: %s\n",r.RequestLine.Method)
		fmt.Printf("- Target: %s\n",r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n",r.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		r.Headers.ForEach(func (n,v string)  {
		 	fmt.Printf("- %s: %s\n",n,v)
		})		
		fmt.Printf("Body:\n")
		fmt.Printf("%s\n",r.Body)		
	// 		lines := getLinesChannel(conn) // reading 8 bytes form a connection instead of priviously files
	// for line := range lines {
	// 	fmt.Printf("read: %s\n",line)
	// }
	}
	
	
}