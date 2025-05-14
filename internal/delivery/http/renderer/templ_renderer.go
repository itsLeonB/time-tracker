package renderer

import (
	"context"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin/render"
)

var Default = &HTMLTemplRenderer{}

type HTMLTemplRenderer struct {
	FallbackHtmlRenderer render.HTMLRender
}

func (r *HTMLTemplRenderer) Instance(s string, d any) render.Render {
	templData, ok := d.(templ.Component)
	if !ok {
		if r.FallbackHtmlRenderer != nil {
			return r.FallbackHtmlRenderer.Instance(s, d)
		}
	}

	return &Renderer{
		Ctx:       context.Background(),
		Status:    -1,
		Component: templData,
	}
}
