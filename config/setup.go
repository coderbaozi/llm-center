package config

func Setup() {
	r := InitGin()
	InitRouter(r)
	InitDB()
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()
}
