package web

import (
	"bobthebuilder/logging"
	"bobthebuilder/config"
	"github.com/hoisie/web"
)


func indexMainPage(ctx *web.Context) {
  t := templates.Lookup("index")
  if t == nil {
          logging.Error("web", "No template found.")
					return
  }
  err := t.Execute(ctx.ResponseWriter, modelBasic{Config: config.All()})
  if err != nil{
          logging.Error("views-index", err)
  }
}