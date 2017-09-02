package main

import (
	"os"
	"fmt"
	"flag"
	 "net"
	 "encoding/json"
	 "time"
	 "io/ioutil"
	 "github.com/golang/glog"
)

type Config struct {
	Endpoint string `json:"endpoint"`
}

func usage(){
	
	flag.Lookup("stderrthreshold").Value.Set("INFO")
	fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n", )
	flag.PrintDefaults()
	os.Exit(2)
}

func init(){
	flag.Usage = usage
	flag.Parse()
}

func main() {
	glog.Info("Hello World\n")

	config := getConfig()

	for _, elment := range config {
		conn, err := net.DialTimeout("tcp",elment.Endpoint, time.Duration(5000000000))
		if err != nil {
			glog.Info(elment.Endpoint + " is unreachable")
		}else{
			glog.Info(elment.Endpoint + " is reachable");
			conn.Close();
		}
	}

	glog.Flush()
}

func getConfig() []Config{
	raw, err := ioutil.ReadFile("./config.json")
	if (err != nil){
		glog.Error(err.Error())
	}

	var config []Config
	json.Unmarshal(raw, &config)

	return config
}
