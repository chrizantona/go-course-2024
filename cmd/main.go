package main

import (
	"GO-COURSE-2024/internal/pkg/storage"
	"fmt"
	"os"
)

func main() {
	myStorage := storage.NewStorage()
	filename := "storage.json"
	if _, err := os.Stat(filename); err == nil {
		if err := myStorage.LoadFromFile(filename); err != nil {
			fmt.Printf("Error loading storage from file: %v\n", err)
		} else {
			fmt.Println("Storage loaded from file.")
		}
	} else {
		fmt.Println("No previous storage found. Starting fresh.")
	}
	myStorage.Set("1", "321")
	myStorage.Set("3", "qwer")

	res1 := myStorage.Get("1")
	res2 := myStorage.Get("3")
	res3 := myStorage.Get("2")

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

	kind1 := myStorage.GetKind("1")
	kind2 := myStorage.GetKind("3")
	kind3 := myStorage.GetKind("2")
	fmt.Printf("Kind of value for key '1': %s\n", kind1)
	fmt.Printf("Kind of value for key '3': %s\n", kind2)
	fmt.Printf("Kind of value for key '2': %s\n", kind3)



	fmt.Println("LPUSH list1 1 2 3")
	count := myStorage.LPUSH("list1", "1", "2", "3")
	fmt.Printf("(integer) %d\n", count)

	fmt.Println("LPOP list1 0 -1")
	popped := myStorage.LPOP("list1", 0, -1)
	for i, v := range popped {
	fmt.Printf("%d) %s\n", i+1, v)
	}

	fmt.Println("RPUSH list1 4 5 6")
	count = myStorage.RPUSH("list1", "4", "5", "6")
	fmt.Printf("(integer) %d\n", count)

	fmt.Println("RPOP list1 2")
	popped = myStorage.RPOP("list1", 2)
	for i, v := range popped {
	fmt.Printf("%d) %s\n", i+1, v)
	}

	fmt.Println("RADDTOSET list1 3 5 8 4 8")
	count = myStorage.RADDTOSET("list1", "3", "5", "8", "4", "8")
	fmt.Printf("(integer) %d\n", count)

	fmt.Println("LPOP list1 0 -1")
	popped = myStorage.LPOP("list1", 0, -1)
	for i, v := range popped {
	fmt.Printf("%d) %s\n", i+1, v)
	}

	fmt.Println("LSET list1 1 30")
	if err := myStorage.LSET("list1", 1, "30"); err != nil {
	fmt.Println(err)
	} else {
	fmt.Println("OK")
	}

	fmt.Println("LGET list1 1")
	value, err := myStorage.LGET("list1", 1)
	if err != nil {
	fmt.Println(err)
	} else {
	fmt.Println(value)
	}


	if err := myStorage.SaveToFile(filename); err != nil {
		fmt.Printf("Error saving storage to file: %v\n", err)
	} else {
		fmt.Println("Storage saved to file.")
	}
}
