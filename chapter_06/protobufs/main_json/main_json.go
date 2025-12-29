package main

import (
	"encoding/json"
	"fmt"
	pb "restapi/chapter_06/protofiles"
)

func main() {
	p := &pb.Person{
		Id: 4321,
		Name: "Michael Jordan",
		Email: "mj@example.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "555-7654", Type: pb.Person_HOME},
		},
	}

	body, _ := json.Marshal(p)
	fmt.Println(string(body))
}