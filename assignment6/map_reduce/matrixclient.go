package main

import(
	"rpc"
	"./imatrix"
	"./parser"
	"strconv"
	"os"
	"fmt"
	"bufio"
	"log"
)

//Wrapper for an rpc client
type Client_wrapper struct{
	Client *rpc.Client
}

//Create a new client wrapper
func NewClientWrapper(client *rpc.Client) *Client_wrapper{
	news := new(Client_wrapper)
	news.Client = client;
	return news
}

// return the dimensions of an existing matrix.
func (client *Client_wrapper) Dim(name string) (rows, cols int, err os.Error){
	var reply matrix.Matrix_struct
	err = client.Client.Call("Matrix_database.Dim", name, &reply)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Println("Dimensions error: ", err)
		return 0, 0, err
	}
	return reply.Rows, reply.Cols, nil
}

// create a (unique) matrix with (positive) dimensions.
func (client *Client_wrapper) Make(name string, rows, cols int) os.Error{
	var reply bool
	mk := matrix.NewMake(name, rows, cols)
	err := client.Client.Call("Matrix_database.Make", *mk, &reply)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Println("Make error: ", err)
		return err
	}
	return nil
}

// delete an existing matrix.
func (client *Client_wrapper) Remove(name string) os.Error{
	var reply bool
	err := client.Client.Call("Matrix_database.Remove", name, &reply)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Println("Remove error: ", err)
		return err
	}
	return nil
}

// get an element value of an existing matrix.
func (client *Client_wrapper) Get(name string, i, j int) (value float64, err os.Error){
	var reply float64
	get := matrix.NewGet(name, i, j)
	err = client.Client.Call("Matrix_database.Get", get, &reply)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Println("Get error: ", err)
		return 0, err
	}
	return reply, nil
}

// set an element value of an existing matrix.
func (client *Client_wrapper) Set(name string, i, j int, value float64) os.Error{
	var reply bool
	set := matrix.NewSet(name, i, j, value)
	err := client.Client.Call("Matrix_database.Set", set, &reply)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Println("Set error: ", err)
		return err
	}
	return nil
}
	
// disconnect from the service.
func (client *Client_wrapper) Close() os.Error{
	err := client.Client.Close();
	if err != nil{
		log.SetOutput(os.Stderr)
		log.Println("Close error: ", err)
		return err
	}
	return os.NewError("Closing")
}

//Make the commands to be used by the parser
func make_commands(cw *Client_wrapper) *parser.Commands{
	commands := parser.NewCommands()
	dim := func(input []string) os.Error{
		if len(input) != 1{
			log.SetOutput(os.Stderr)
			log.Println(os.NewError("Invalid arguments"))
			return nil
		}
		r, c, err := cw.Dim(input[0])
		if err == nil{
			fmt.Printf("R: %d C: %d \n", r, c)
			return nil
		}
		return nil
	}
	mak := func(input []string) os.Error{
		if (len(input) != 3){
			log.SetOutput(os.Stderr)
			log.Println(os.NewError("Invalid arguments"))
			return nil
		}
		rows, err := strconv.Atoi(input[1])
		if (err != nil){
			log.SetOutput(os.Stderr)
			log.Println(os.NewError("Expected an int for number of rows"))
			return nil
		}
		cols, e := strconv.Atoi(input[2])
		if (e != nil){
			log.SetOutput(os.Stderr)
			log.Println(os.NewError("Expected an int for number of cols"))
			return nil
		}
		return cw.Make(input[0], rows, cols)
	}
	rm := func (input []string) os.Error{
		if (len(input) != 1){
			log.SetOutput(os.Stderr)
			log.Println(os.NewError("Invalid arguments"))
			return nil
		}
		return cw.Remove(input[0])
	}

	get := func(input []string) os.Error{
		if len(input) != 3{
			log.SetOutput(os.Stderr)
			log.Println(os.NewError("Invalid arguments"))
			return nil
		}
		i, e1 := strconv.Atoi(input[1])
		if (e1 != nil){
			log.SetOutput(os.Stderr)
			log.Println(os.NewError("Expected an int for i"))
			return nil
		}
		j, e2 := strconv.Atoi(input[2])
		if (e2 != nil){
			log.SetOutput(os.Stderr)
			log.Println(os.NewError("Expected an int for j"))
			return nil
		}
		value, err := cw.Get(input[0], i, j)
		if err == nil{
			fmt.Printf("I: %d J: %d Val: %3.3f \n", i, j ,value )
		}
		return nil
	}

	set := func(input []string) os.Error{
		if len(input) != 4{
			log.SetOutput(os.Stderr)
			log.Println(os.NewError("Invalid arguments"))
			return nil
		}
		i, e1 := strconv.Atoi(input[1])
		if (e1 != nil){
			log.SetOutput(os.Stderr)
			log.Println(os.NewError("Expected an int for i"))
			return nil
		}
		j, e2 := strconv.Atoi(input[2])
		if (e2 != nil){
			log.SetOutput(os.Stderr)
			log.Println(os.NewError("Expected an int for j"))
			return nil
		}
		v, e3 := strconv.Atof64(input[3])
		if (e3 != nil){
			log.SetOutput(os.Stderr)
			log.Println(os.NewError("Expected a float64 for value"))
			return nil
		}
		return cw.Set(input[0], i, j, v)
		
	}

	clo := func() os.Error{
		return cw.Close()
	}

	commands.AddInputCommand("Make", mak)
	commands.AddInputCommand("Dim", dim)
	commands.AddInputCommand("Remove", rm)
	commands.AddInputCommand("Get", get)
	commands.AddInputCommand("Set", set)
	commands.AddCommand("Close", clo)
	return commands
}

//main for the matrix client
func main (){
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	cw := NewClientWrapper(client)
	commands := make_commands(cw)
	sin := bufio.NewReader(os.Stdin)
	var l string
	for err != os.EOF{
		l,err = sin.ReadString('\n')
		e := commands.Parseln(l)
		if(e.String() == "Closing"){
			break
		}
		if (e != nil){
			fmt.Fprintln(os.Stderr, e)
		}
	}
}
