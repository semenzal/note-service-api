package main

import (
	"context"
	"fmt"
	"log"

	desc "github.com/semenzal/note-service-api/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const address = "localhost:50051"

func main() {
	con, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("didn't connect: %s", err.Error())
	}
	defer con.Close()

	client := desc.NewNoteServiceClient(con)

	res, err := client.Create(context.Background(), &desc.CreateRequest{
		Note: &desc.NoteInfo{
			Title:  "Wow",
			Text:   "I'm surprised!",
			Author: "Semen",
		},
	})
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Id:", res.Id)

	resGetNote, err := client.Get(context.Background(), &desc.GetRequest{
		Id: 2,
	})
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Id:", resGetNote.String())

	resGetList, err := client.GetList(context.Background(), &desc.GetListRequest{
		Filter: &desc.Filter{
			Title:  nil,
			Text:   nil,
			Author: nil,
			Email:  nil,
			Limit:  wrapperspb.Int64(10),
			Offset: wrapperspb.Int64(10),
		},
	})
	if err != nil {
		log.Println(err.Error())
	}
	for _, note := range resGetList.Notes {

		fmt.Println(note.Info.Title)
		fmt.Println(note.Info.Text)
		fmt.Println(note.Info.Author)
	}
	log.Println("All Id:", resGetList.Notes)

	resUpdate, err := client.Update(context.Background(), &desc.UpdateRequest{
		Note: &desc.UpdateNoteInfo{
			Title:  nil,
			Text:   nil,
			Author: nil,
		},
	})
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("New Id:", resUpdate.String())

	resDelete, err := client.Delete(context.Background(), &desc.DeleteRequest{
		Id: 12,
	})
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Delete:", resDelete.String())
}
