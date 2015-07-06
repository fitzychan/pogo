package main

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

//Functions
func ParseArgs(r *http.Request) string {
	var argsBuffer bytes.Buffer

	pUrl, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Error("Error Parsing URL")
	}

	for key, value := range pUrl.Query() {
		argsBuffer.WriteString(" " + key)
		argsBuffer.WriteString(" " + value[0])
	}

	return argsBuffer.String()
}

func exec_script(sc string) string {
	sc = string(strings.TrimSpace(sc))

	cmd := exec.Command("powershell.exe", "-NoLogo", "-NonInteractive", "-Command", "&{", sc, "}", "| ConvertTo-Json")
	out, err := cmd.Output()

	if err != nil {
		log.WithFields(log.Fields{
			"command": string(sc),
		}).Error(err.Error())
	}

	return string(out)
}

//Handlers
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Service Running"))
}

func ExitHandler(w http.ResponseWriter, r *http.Request) {
	defer os.Exit(0)

	log.Info("Shutting Down")
	w.Write([]byte(fmt.Sprintf("Shutting Down")))
	time.Sleep(3000 * time.Millisecond)
}

func RunShell(w http.ResponseWriter, r *http.Request) {
	var commbuffer bytes.Buffer
	commbuffer.WriteString(mux.Vars(r)["name"])
	commbuffer.WriteString(ParseArgs(r))

	w.Write([]byte(fmt.Sprintf(exec_script(commbuffer.String()))))
}

func RunScript(w http.ResponseWriter, r *http.Request) {
	var commbuffer bytes.Buffer

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error(err)
	}

	commbuffer.WriteString("&\"")
	commbuffer.WriteString(filepath.Join(dir, "\\scripts\\", mux.Vars(r)["name"]))
	commbuffer.WriteString(".ps1\"")
	commbuffer.WriteString(ParseArgs(r))

	w.Write([]byte(fmt.Sprintf(exec_script(commbuffer.String()))))
}

//Runtime
func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".\\")
	viper.ReadInConfig()

	viper.SetDefault("Binding", ":8080")
}

func main() {
	mx := mux.NewRouter()

	mx.HandleFunc("/", IndexHandler)
	mx.HandleFunc("/exit", ExitHandler)
	mx.HandleFunc("/command/{name:\\S+}", RunShell)
	mx.HandleFunc("/script/{name:\\S+}", RunScript)

	log.Info("Listening at " + viper.GetString("Binding"))
	http.ListenAndServe(viper.GetString("Binding"), mx)
}
