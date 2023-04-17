package rbac

import (
	"bytes"
	stderrors "errors"
	"io"
	stdhttp "net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type Context struct {
	http.Context
}

func (c *Context) Bind(v interface{}) error {
	r := c.Context.Request()
	//if contentSubtype(r.Header.Get("Content-Type")) == "json" {
	data, err := io.ReadAll(r.Body)

	// reset body.
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	if err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	if len(data) == 0 {
		return nil
	}

	a := accessorOrVisitor(c.Value("author"))
	return ((codec)(append([]string{a.GetRoleType()}, a.GetRoles()...))).Unmarshal(data, v)
	//}
	//return ErrUnsupportContentType
	//return c.Context.Bind(v)
}

func (c *Context) Result(code int, v interface{}) error {
	if rd, ok := v.(http.Redirector); ok {
		c.Response().WriteHeader(code)
		url, code := rd.Redirect()
		stdhttp.Redirect(c.Response(), c.Request(), url, code)
		return nil
	}
	return c.JSON(code, v)
}

func (c *Context) JSON(code int, v interface{}) error {
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(code)
	a := accessorOrVisitor(c.Value("author"))
	data, err := ((codec)(append([]string{a.GetRoleType()}, a.GetRoles()...))).Marshal(v)
	if err != nil {
		return err
	}
	_, err = c.Response().Write(data)
	return err
}

func contentSubtype(contentType string) string {
	left := strings.Index(contentType, "/")
	if left == -1 {
		return ""
	}
	right := strings.Index(contentType, ";")
	if right == -1 {
		right = len(contentType)
	}
	if right < left {
		return ""
	}
	return contentType[left+1 : right]
}

var ErrUnsupportContentType = stderrors.New("unsupported content type")
