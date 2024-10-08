package main

import (
	"1/internal/pkg/storage"
	"fmt"
)

func main() {
	my_storage := storage.NewStorage()
	my_storage.Set("1", "321")
	my_storage.Set("3", "qwer")
	res1 := my_storage.Get("1")
	res2 := my_storage.Get("3")
	res3 := my_storage.Get("2")
	if res1 != nil {
		fmt.Printf("Value for key '1': %s\n", *res1)
	} else {
		fmt.Println("Key '1' not found.")
	}

	if res2 != nil {
		fmt.Printf("Value for key '3': %s\n", *res2)
	} else {
		fmt.Println("Key '3' not found.")
	}

	if res3 != nil {
		fmt.Printf("Value for key '2': %s\n", *res3)
	} else {
		fmt.Println("Key '2' not found.")
	}
	kind1 := my_storage.GetKind("1")
	kind2 := my_storage.GetKind("3")
	kind3 := my_storage.GetKind("2")
	fmt.Printf("Kind of value for key '1': %s\n", kind1)
	fmt.Printf("Kind of value for key '3': %s\n", kind2)
	fmt.Printf("Kind of value for key '2': %s\n", kind3)
}
