package explorer

import (
	"github.com/rockjoon/nomadcoin/blockchain"
	"log"
	"net/http"
	"text/template"
)

type homeData struct {
	HomeStory string
	Blocks    []*blockchain.Block
}

const (
	port        string = ":4000"
	templateDir string = "explorer/template/"
)

var templates *template.Template

func Start() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/add", handleAdd)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleHome(rw http.ResponseWriter, r *http.Request) {
	homeData := homeData{"home story", blockchain.GetBlockChain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", homeData)
}

func handleAdd(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templates.ExecuteTemplate(rw, "add", nil)
	case http.MethodPost:
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockChain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}
