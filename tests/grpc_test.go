package tests

import (
	api "Mongo/api/protoc"
	"Mongo/internal/pkg/constErr"
	"Mongo/internal/pkg/data"
	"Mongo/internal/pkg/storage"
	s "Mongo/internal/pkg/storage/grpc"
	mem "Mongo/internal/pkg/storage/in-memory"
	"log"
	"net"
	"os"
	"testing"

	"google.golang.org/grpc"
)

type TestGRPC struct {
	Desc           string
	ExcpectedError error
	testFunc       func(s storage.InterfaceStorage) error
}

func TestCommandStorageGRPC(t *testing.T) {
	tests := []TestGRPC{
		{
			Desc:           "Add ad valid",
			ExcpectedError: nil,
			testFunc: func(s storage.InterfaceStorage) error {
				ad := data.Ad{
					Brand: "Mazda",
					Model: "CX-5",
					Color: "Red",
					Price: 25000,
				}
				return s.Add(&ad)
			},
		},

		{
			Desc:           "Add ad invalid",
			ExcpectedError: constErr.AdIsNil,
			testFunc: func(s storage.InterfaceStorage) error {
				return s.Add(nil)
			},
		},

		{
			Desc:           "Get ad valid",
			ExcpectedError: nil,
			testFunc: func(s storage.InterfaceStorage) error {
				_, err := s.Get(1)
				return err
			},
		},

		{
			Desc:           "Get ad invalid",
			ExcpectedError: constErr.InvalidID,
			testFunc: func(s storage.InterfaceStorage) error {
				_, err := s.Get(20)
				return err
			},
		},

		{
			Desc:           "Update ad valid",
			ExcpectedError: nil,
			testFunc: func(s storage.InterfaceStorage) error {
				ad := data.Ad{
					Brand: "Subaru",
					Model: "Forester",
					Color: "Blue",
					Price: 50000,
				}
				return s.Update(&ad, 1)
			},
		},

		{
			Desc:           "Update ad invalid",
			ExcpectedError: constErr.AdIsNil,
			testFunc: func(s storage.InterfaceStorage) error {
				return s.Update(nil, 1)
			},
		},

		{
			Desc:           "Delete ad invalid",
			ExcpectedError: constErr.InvalidID,
			testFunc: func(s storage.InterfaceStorage) error {
				return s.Delete(2)
			},
		},

		{
			Desc:           "Delete ad valid",
			ExcpectedError: nil,
			testFunc: func(s storage.InterfaceStorage) error {
				return s.Delete(1)
			},
		},
	}
	go RunServerGRPCServer()
	m := mem.NewStorage()
	testCommands(t, tests, m)
}

func testCommands(t *testing.T, tests []TestGRPC, s storage.InterfaceStorage) {
	for _, test := range tests {
		err := test.testFunc(s)

		if test.ExcpectedError != err {
			t.Errorf("Desc: %v\nError commands - returned %v, want %v", test.Desc, err, test.ExcpectedError)
		}
	}

}

func RunServerGRPCServer() {
	ls, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()
	api.RegisterServiceProtobufServer(srv, s.NewStorageGrpcServer(mem.NewStorage()))
	log.Fatal(srv.Serve(ls))
}
