package testdata

import (
	"errors"
	ex "github.com/jackskj/protoc-gen-map/examples"
)

var CallbackPool = map[string]int{}

func BlogBeforeQueryCallback(queryString string, req *EmptyRequest) error {
	incrementPool("BlogBeforeQueryCallback")
	return nil
}
func BlogsBeforeQueryCallback(queryString string, req *EmptyRequest) error {
	incrementPool("BlogsBeforeQueryCallback")
	return nil
}

func BlogAfterQueryCallback(queryString string, req *EmptyRequest, resp *ex.BlogResponse) error {
	incrementPool("BlogAfterQueryCallback")
	return nil
}
func BlogsAfterQueryCallback(queryString string, req *EmptyRequest, resp []*ex.BlogResponse) error {
	incrementPool("BlogsAfterQueryCallback")
	return nil
}

func BlogCache(queryString string, req *EmptyRequest) (*ex.BlogResponse, error) {
	resp := ex.BlogResponse{
		Title: "cached result",
	}
	incrementPool("BlogCache")
	return &resp, nil
}
func BlogsCache(queryString string, req *EmptyRequest) ([]*ex.BlogResponse, error) {
	resp := []*ex.BlogResponse{
		&ex.BlogResponse{
			Title: "cached result 1",
		},
		&ex.BlogResponse{
			Title: "cached result 2",
		},
		&ex.BlogResponse{
			Title: "cached result 3",
		},
	}
	incrementPool("BlogsCache")
	return resp, nil
}

func FailedBlogBeforeQueryCallback(queryString string, req *EmptyRequest) error {
	incrementPool("FailedBlogBeforeQueryCallback")
	return errors.New("FailedBlogBeforeQueryCallback")
}
func FailedBlogsBeforeQueryCallback(queryString string, req *EmptyRequest) error {
	incrementPool("FailedBlogsBeforeQueryCallback")
	return errors.New("FailedBlogsBeforeQueryCallback")
}

func FailedBlogAfterQueryCallback(queryString string, req *EmptyRequest, resp *ex.BlogResponse) error {
	incrementPool("FailedBlogAfterQueryCallback")
	return errors.New("FailedBlogAfterQueryCallback")
}
func FailedBlogsAfterQueryCallback(queryString string, req *EmptyRequest, resp []*ex.BlogResponse) error {
	incrementPool("FailedBlogsAfterQueryCallback")
	return errors.New("FailedBlogsAfterQueryCallback")
}

func FailedBlogCache(queryString string, req *EmptyRequest) (*ex.BlogResponse, error) {
	incrementPool("FailedBlogCache")
	return nil, errors.New("FailedBlogCache")
}
func FailedBlogsCache(queryString string, req *EmptyRequest) ([]*ex.BlogResponse, error) {
	incrementPool("FailedBlogsCache")
	return nil, errors.New("FailedBlogsCache")
}

func incrementPool(msg string) {
	CallbackPool[msg] = CallbackPool[msg] + 1
}
