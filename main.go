package main

func main() {
	e := NewEditor()

	err := e.buffer.LoadFile("main.go")
	if err != nil {
		panic(err)
	}

	e.Run()
}
