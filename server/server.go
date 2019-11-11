package server

import (
	"books/server/proto/library"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
)

var books = []*library.Book{
	&library.Book{
		Isbn:     60929871,
		Title:    "Brave New World",
		Author:   "Aldous Huxley",
		BookType: library.BookType_HARDCOVER,
		PublishingMethod: &library.Book_Publisher{
			Publisher: &library.Publisher{
				Name: "Chatto & Windus",
			},
		},
		PublicationDate: &timestamp.Timestamp{
			Seconds: time.Date(1932, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
		},
	},
	&library.Book{
		Isbn:     140009728,
		Title:    "Nineteen Eighty-Four",
		Author:   "George Orwell",
		BookType: library.BookType_PAPERBACK,
		PublishingMethod: &library.Book_Publisher{
			Publisher: &library.Publisher{
				Name: "Secker & Warburg",
			},
		},
		PublicationDate: &timestamp.Timestamp{
			Seconds: time.Date(1949, time.June, 8, 0, 0, 0, 0, time.UTC).Unix(),
		},
	},
	&library.Book{
		Isbn:     9780140301694,
		Title:    "Alice's Adventures in Wonderland",
		Author:   "Lewis Carroll",
		BookType: library.BookType_AUDIOBOOK,
		PublishingMethod: &library.Book_Publisher{
			Publisher: &library.Publisher{
				Name: "Macmillan",
			},
		},
		PublicationDate: &timestamp.Timestamp{
			Seconds: time.Date(1865, time.November, 26, 0, 0, 0, 0, time.UTC).Unix(),
		},
	},
	&library.Book{
		Isbn:     140008381,
		Title:    "Animal Farm",
		Author:   "George Orwell",
		BookType: library.BookType_HARDCOVER,
		PublishingMethod: &library.Book_Publisher{
			Publisher: &library.Publisher{
				Name: "Secker & Warburg",
			},
		},
		PublicationDate: &timestamp.Timestamp{
			Seconds: time.Date(1945, time.August, 17, 0, 0, 0, 0, time.UTC).Unix(),
		},
	},
	&library.Book{
		Isbn:     1501107739,
		Title:    "Still Alice",
		Author:   "Lisa Genova",
		BookType: library.BookType_PAPERBACK,
		PublishingMethod: &library.Book_SelfPublished{
			SelfPublished: true,
		},
		PublicationDate: &timestamp.Timestamp{
			Seconds: time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC).Unix(),
		},
	},
}

type BookService struct {
	b broadcaster
}
type broadcaster struct {
	listenerMu sync.Mutex
	listeners  map[string]chan<- string
}

func (s *BookService) GetBook(ctx context.Context, bookquery *library.GetBookRequest) (*library.Book, error) {
	fmt.Println("Get Book Request Happended")
	for _, bk := range books {
		if bk.Isbn == bookquery.Isbn {
			return bk, nil
		}
	}
	return nil, grpc.Errorf(codes.NotFound, "Book could not found")
}
func (s *BookService) QueryBooks(bookQuery *library.QueryBooksRequest, stream library.BookService_QueryBooksServer) error {
	for _, book := range books {
		select {
		case <-stream.Context().Done():
			return nil
		default:
		}

		if strings.HasPrefix(book.Author, bookQuery.AuthorPrefix) {
			err := stream.Send(book)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *BookService) MakeCollection(srv library.BookService_MakeCollectionServer) error {
	collection := &library.Collection{}
	for {
		bk, err := srv.Recv()
		if err == io.EOF {
			return srv.SendAndClose(collection)
		}
		if err != nil {
			return err
		}

		collection.Books = append(collection.Books, bk)
	}
}
func (s *BookService) GetAllBooks(ctx context.Context, req *library.EmptyRequest) (*library.Collection, error) {
	collection := &library.Collection{}
	for _, bk := range books {
		collection.Books = append(collection.Books, bk)

	}
	return collection, nil

}

