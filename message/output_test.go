package message

import "testing"

func TestOutput(t *testing.T) {
	t.Log("--------------- Print ------------------")
	Print("test 1\n")
	t.Log("--------------- Printf ------------------")
	Printf("test %v\n", 1)
	t.Log("--------------- Println ------------------")
	Println("test", 1, 2)
	t.Log("--------------- Print End ------------------")

}
