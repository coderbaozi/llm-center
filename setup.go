package setup

import (
	"github.com/coderbaozi/llm-center/config"
	"github.com/coderbaozi/llm-center/router"
)

func Setup() {
	config.InitDB()
	r := config.InitGin()
	router.InitRouter(r)
	r.Run()
}
