// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.5.2
// - protoc             v3.15.6
// source: examples/singleton/api/blog/v1/blog.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationBlogListPost = "/examples.singleton.api.blog.v1.Blog/ListPost"

type BlogHTTPServer interface {
	ListPost(context.Context, *ListPostRequest) (*ListPostReply, error)
}

func RegisterBlogHTTPServer(s *http.Server, srv BlogHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/posts", _Blog_ListPost0_HTTP_Handler(srv))
}

func _Blog_ListPost0_HTTP_Handler(srv BlogHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListPostRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationBlogListPost)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListPost(ctx, req.(*ListPostRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListPostReply)
		return ctx.Result(200, reply)
	}
}

type BlogHTTPClient interface {
	ListPost(ctx context.Context, req *ListPostRequest, opts ...http.CallOption) (rsp *ListPostReply, err error)
}

type BlogHTTPClientImpl struct {
	cc *http.Client
}

func NewBlogHTTPClient(client *http.Client) BlogHTTPClient {
	return &BlogHTTPClientImpl{client}
}

func (c *BlogHTTPClientImpl) ListPost(ctx context.Context, in *ListPostRequest, opts ...http.CallOption) (*ListPostReply, error) {
	var out ListPostReply
	pattern := "/v1/posts"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationBlogListPost))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
