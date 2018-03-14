package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/grafana/grafana/pkg/components/renderer"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
)

func RenderToPng(c *m.ReqContext) {
	queryReader, err := util.NewUrlQueryReader(c.Req.URL)
	if err != nil {
		c.Handle(400, "Render parameters error", err)
		return
	}
	queryParams := fmt.Sprintf("?%s", c.Req.URL.RawQuery)

	//c.Logger.Info("render Api ", "Path", c.Req.URL.RawPath, "query", c.Req.URL.RawQuery)

	path := c.Params("*")

	//if strings.HasPrefix("/"+path, setting.AppSubUrl) {
	//	c.Logger.Info("pre trim", "path", path)
	//	path = strings.TrimLeft("/"+path, setting.AppSubUrl)
	//}
	//
	//c.Logger.Info("Render path", "path", path)

	renderOpts := &renderer.RenderOpts{
		Path:     path + queryParams,
		Width:    queryReader.Get("width", "800"),
		Height:   queryReader.Get("height", "400"),
		Timeout:  queryReader.Get("timeout", "60"),
		OrgId:    c.OrgId,
		UserId:   c.UserId,
		OrgRole:  c.OrgRole,
		Timezone: queryReader.Get("tz", ""),
		Encoding: queryReader.Get("encoding", ""),
	}

	c.Logger.Info("renderopts", "ops", renderOpts)

	pngPath, err := renderer.RenderToPng(renderOpts)

	if err != nil && err == renderer.ErrTimeout {
		c.Handle(500, err.Error(), err)
		return
	}

	if err != nil {
		c.Handle(500, "Rendering failed.", err)
		return
	}

	c.Resp.Header().Set("Content-Type", "image/png")
	http.ServeFile(c.Resp, c.Req.Request, pngPath)
}
