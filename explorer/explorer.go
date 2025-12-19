package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/hellotect2022go/nomadcoin/blockchain"
)

var port string

const (
	templateDir string = "explorer/template/"
)

var templates *template.Template

type HomeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(w http.ResponseWriter, r *http.Request) {
	//tmpl := template.Must(template.ParseFiles("template/pages/home.gohtml"))
	data := HomeData{"Home", blockchain.GetBlockChain().Blocks()}
	//tmpl.Execute(w, data)
	templates.ExecuteTemplate(w, "home", data)
}

func add(w http.ResponseWriter, r *http.Request) {
	fmt.Println("????", r.Method)
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "add", nil)
	case "POST":

		blockchain.GetBlockChain().AddBlock()
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

}

func Start(aPort int) {
	handler := http.NewServeMux()
	port = fmt.Sprintf(":%d", aPort)
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler)) // 프로그램이 Exit(1) : {error 로 인해 종료시킴} 로그를 출력해줌
}
