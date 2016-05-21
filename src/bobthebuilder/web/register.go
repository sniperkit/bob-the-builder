package web

import (
  "bobthebuilder/logging"
  "bobthebuilder/config"
  //"golang.org/x/net/websocket"
  "github.com/hoisie/web"
)

// ### THIS FILE SHOULD CONTAIN ALL INITIALISATION CODE FOR BOTH TEMPLATES AND URL HANDLERS ###

func Initialise() {
  logging.Info("web", "Registering page handlers")
  registerCoreHandlers()
  registerViewHandlers()

  logging.Info("web", "Registering templates")
  registerCoreTemplates()
}

func registerCoreHandlers() {
  web.Get("/", indexMainPage, config.All().Web.Domain)
}

func registerViewHandlers(){

}

func registerCoreTemplates(){
  logError(registerTemplate("tailcontent.tpl", "tailcontent"), "Template load error: ")
  logError(registerTemplate("headcontent.tpl", "headcontent"), "Template load error: ")
  logError(registerTemplate("index.tpl", "index"), "Template load error: ")
  logError(registerTemplate("topnav.tpl", "topnav"), "Template load error: ")
}


func logError(e error, prefix string){
  if e != nil{
    logging.Error("web", prefix, e.Error())
  }
}