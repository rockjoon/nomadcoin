package explorer

import (
	"fmt"
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
	templateDir string = "explorer/template/"
)

var templates *template.Template

func Start(aPort int) {
	port := fmt.Sprintf(":%d", aPort)
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler := http.NewServeMux()
	handler.HandleFunc("/", handleHome)
	handler.HandleFunc("/add", handleAdd)
	log.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
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
