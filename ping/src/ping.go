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

const FIVE_SECONDS int64 = 5000000000
const SIXTY_SECONDS int64 = 60000000000

type Config struct {
	Endpoint string `json:"endpoint"`
}

func init(){
	flag.Usage = usage
	flag.Parse()
	flag.Set("stderrthreshold", "INFO")
	//fmt.Println(flag.Lookup("stderrthreshold").Value)
}

func usage(){
	fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n", )
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	glog.Info("Starting Ping\n")

	config := getConfig()
	doEvery(time.Duration(FIVE_SECONDS), checkAllEndpoints, config)

	glog.Info("Exiting Ping")
	glog.Flush()
	return
}

func getConfig() []Config{
	var config []Config
	raw, err := ioutil.ReadFile("./config.json")
	if (err != nil){
		glog.Warning("Unable to Load the Configuration File, Using Default Endpoint: www.trevorjackson.ca:80")
		var defaultEndpoint Config
		defaultEndpoint.Endpoint = "www.trevorjackson.ca:80"
		config = make([]Config, 1)
		config[0] = defaultEndpoint
		return config
	}
	
	json.Unmarshal(raw, &config)

	return config
}

func doEvery(d time.Duration, f func([]Config), endpoints []Config) {
	for range time.Tick(d) {
		f(endpoints)
	}
}

func checkAllEndpoints(endpoints []Config){
	for _, elment := range endpoints {
		checkEndpoint(elment.Endpoint)
	}
}

func checkEndpoint(endpoint string){
	conn, err := net.DialTimeout("tcp",endpoint, time.Duration(FIVE_SECONDS))
	if err != nil {
		glog.Info(endpoint + " is unreachable")
	}else{
		glog.Info(endpoint + " is reachable");
		conn.Close();
	}
}