package main

import (
	"fmt"
	pb "restapi/chapter_06/protofiles"

	"google.golang.org/protobuf/proto"
)

func main() {
	p := &pb.Person{
		Id: 1234,
		Name: "Shaqil Oneil",
		Email: "so@example.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "55-4312", Type: pb.Person_HOME},
		},
	}

	p1 := &pb.Person{}
	body, _ := proto.Marshal(p)
	_ = proto.Unmarshal(body, p1)
	fmt.Println("Original struct loaded from proto file: ", p, "\n")
	fmt.Println("Marshaled proto data: ", body, "\n")
	fmt.Println("Unmarshaled struct: ", p1)
}