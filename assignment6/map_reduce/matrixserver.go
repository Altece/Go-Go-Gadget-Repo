package main

import(
	"./imatrix"
	"rpc"
	"net"
	"log"
	"http"
	"os"
)

/**
 * The Matrix database. Contains a map of string names associated with  
 * matrix structs
 */
type Matrix_database struct {
	Matrices map[string]*matrix.Matrix_struct
}

//Creates a new matrix database
func NewMatrixDatabase() *Matrix_database{
	news := new(Matrix_database)
	news.Matrices = make(map[string]*matrix.Matrix_struct)
	return news
}

//finds the dimensions of a given matrix
func (m *Matrix_database) Dim(name string, matrix *matrix.Matrix_struct) os.Error{
	mat, ok := m.Matrices[name]
	if(!ok){
		return os.NewError("Matrix of that name is not in the database.")
	}
	*matrix = *mat
	return nil
}

//makes of a matrix of the given name and dimensions
func (m *Matrix_database) Make(mak matrix.Make, a *bool) os.Error{
	_, ok := m.Matrices[mak.Name]
	if ( ok ){
		return os.NewError("Matrix of that name is already in the database.")
	}
	m.Matrices[mak.Name] = matrix.NewMatrix(mak.Name, mak.Rows, mak.Cols)
	return nil
}

//removes a matrix of a given name from the database
func (m *Matrix_database) Remove(name string, a *bool ) os.Error{
	_, ok:= m.Matrices[name]
	if (!ok){
		return os.NewError("Matrix of that name is not in the database.")
	}
	m.Matrices[name] = matrix.NewMatrix("", 0, 0), false
	return nil
}

//Gets the value at a given location of a given matrix
func (m *Matrix_database) Get(get matrix.Get, val *float64) os.Error{
	mat, ok := m.Matrices[get.Name]
	if (!ok){
		return os.NewError("Matrix of that name is not in the database")
	}else if(get.I < 0 || get.I > mat.Rows ||get.J < 0 || get.J > mat.Cols){
		return os.NewError("I or J is out of bounds of the matrix.")
	}
	*val = mat.Matrix[get.I][get.J]
	return nil
}

//Sets the value at a given location of a given matrix
func (m *Matrix_database) Set(set matrix.Set, a *bool) os.Error{
	mat, ok := m.Matrices[set.Name]
	if (!ok){
		return os.NewError("Matrix of that name is not in the database")
	}else if(set.I < 0 || set.I > mat.Rows ||set.J < 0 || set.J > mat.Cols){
		return os.NewError("I or J is out of bounds of the matrix.")
	}
	mat.Matrix[set.I][set.J] = set.Value
	return nil
}

//Starts the server
func main(){
	md := NewMatrixDatabase()
	rpc.Register(md)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}