package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	err := checkError(args)
	if err != nil {
		fmt.Println("error is ", err)
		return err
	}
	operation := args["operation"]
	fileName := args["fileName"]
	item := args["item"]
	id := args["id"]
	switch operation {
	case "list":
		data, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		_, err = writer.Write(data)
		if err != nil {
			return err
		}
	case "add":
		var users []User
		var newUser User
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		data, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(item), &newUser)
		if err != nil {
			return errors.New("")
		}
		if len(data) == 0 {
			users = append(users, newUser)
			data, err = json.Marshal(users)
			_, err = file.Write(data)
			if err != nil {
				return err
			}
			return nil
		}
		err = json.Unmarshal(data, &users)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(item), &newUser)
		if err != nil {
			return errors.New("")
		}
		for _, user := range users {
			if user.Id == newUser.Id {
				_, err = writer.Write([]byte("Item with id " + user.Id + " already exists"))
				if err != nil {
					return err
				}
				return nil
			}
		}
		users = append(users, newUser)
		data, err = json.Marshal(users)
		_, err = file.Write(data)
		if err != nil {
			return err
		}
	case "findById":
		var users []User
		data, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &users)
		if err != nil {
			return err
		}
		for _, user := range users {
			if user.Id == id {
				dataUser, err := json.Marshal(user)
				if err != nil {
					return err
				}
				_, err = writer.Write(dataUser)
				if err != nil {
					return err
				}

			}
		}
	case "remove":
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
		defer file.Close()
		date, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		var users []User
		err = json.Unmarshal(date, &users)
		if err != nil {
			return err
		}
		needRemove := true
		for i, user := range users {
			if id == user.Id && i < len(users)-1 {
				if i < len(users)-1 {
					copy(users[i:], users[i+1:])
					//users[len(users)-1] = nil // обнуляем "хвост"
					users = users[:len(users)-1]
					needRemove = false
					break
				} else if i == len(users)-1 {
					users = users[:len(users)-1]
					needRemove = false
					break
				}
			}
		}
		if needRemove {
			_, err = writer.Write([]byte("Item with id " + id + " not found"))
			if err != nil {
				return err
			}
		} else {
			data, err := json.Marshal(users)
			if err != nil {
				return err
			}
			err = file.Truncate(0)
			if err != nil {
				return err
			}
			_, err = file.Write(data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func checkError(args Arguments) error {
	//check operation
	operation := args["operation"]
	if len(operation) == 0 {
		return errors.New("-operation flag has to be specified")
	}
	switch operation {
	case "add":
	case "list":
	case "findById":
	case "remove":
		break
	default:
		return errors.New("Operation " + operation + " not allowed!")
	}
	//check add without item
	item, ok := args["item"]
	if operation == "add" {
		if !ok {
			return errors.New("-item flag is missed")
		}
		if len(item) == 0 {
			return errors.New("-item flag has to be specified")
		}
	}
	//check filename
	fileName, ok := args["fileName"]
	if !ok {
		return errors.New("-There isn't fileName flag in Arguments Map")
	}
	if len(fileName) == 0 {
		return errors.New("-fileName flag has to be specified")
	}
	//check find by id
	id := args["id"]
	if operation == "findById" && len(id) == 0 {
		return errors.New("-id flag has to be specified")
	}
	//check remove
	if operation == "remove" && len(id) == 0 {
		return errors.New("-id flag has to be specified")
	}
	return nil
}

type Arguments map[string]string

func main() {
	var operation string
	var item string
	var fileName string
	var id string
	flag.StringVar(&operation, "operation", "none", "Name of action")
	flag.StringVar(&item, "item", "none", "Date Item")
	flag.StringVar(&fileName, "filename", "none", "It is file name")
	flag.StringVar(&id, "id", "none", "it's id for find by id")
	flag.Parse()
	args := Arguments{
		"id":        "",
		"operation": operation,
		"item":      item,
		"fileName":  fileName,
	}
	fmt.Println("bellow os.Stdout")
	n, err := fmt.Fprintln(os.Stdout, args)
	fmt.Println("above os.Stdout", n, err)

	err = Perform(args, os.Stdout)
	if err != nil {
		err := errors.New("error in maim function")
		if err != nil {
			return
		}
	}
	//err := Perform(parseArgs(), os.Stdout)

	// fmt.Printf("Hello %s\n", name)

	//err := Perform(parseArgs(), os.Stdout)
	//file, err := os.OpenFile("user.json", os.O_RDWR|os.O_CREATE, 0755)

	//err := Perform([]byte{'a', 'b'}, os.Stdout)
	// if err != nil {
	// 	panic(err)
	// }
}
