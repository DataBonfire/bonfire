package ac

import (
	stdhttp "net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/transport/http"
)

func ReadHTTPTransporter(t http.Transporter) (action, res string) {
	req := t.Request()
	chips := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	res = chips[0]
	switch req.Method {
	case stdhttp.MethodGet:
		if len(chips) == 1 {
			// GET /posts
			action = ActionBrowse
		} else {
			// GET /posts/1/comments
			action = ActionShow
		}
	case stdhttp.MethodPost, stdhttp.MethodPut, stdhttp.MethodPatch:
		if len(chips) == 1 {
			// [POST|PUT|PATCH] /posts
			action = ActionCreate
		} else {
			// [POST|PUT|PATCH] /posts/1
			// [POST|PUT|PATCH] /posts/1/{action}
			// [POST|PUT|PATCH] /posts/1/comments
			// [POST|PUT|PATCH] /posts/1/comments/1
			action = ActionEdit
		}
	case stdhttp.MethodDelete:
		if len(chips) == 1 {
			// DELETE /posts/1
			action = ActionDelete
		} else {
			// DELETE /posts/1/comments/1
			action = ActionEdit
		}
	}
	return
}
