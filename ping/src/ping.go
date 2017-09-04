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
	 "strconv"
)

const FIVE_SECONDS int64 = 5000000000
const SIXTY_SECONDS int64 = 60000000000
const DEFAULT_ENDPOINT string = "www.trevorjackson.ca:80"


type Config struct {
	Endpoints []string `json:"endpoints"`
	CheckInterval int64 `json:"checkInterval"`
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
	doEvery(time.Duration(FIVE_SECONDS), checkAllEndpoints, config.Endpoints)

	glog.Info("Exiting Ping")
	glog.Flush()
	return
}

func getConfig() Config{
	var config []Config
	raw, err := ioutil.ReadFile("./config.json")
	if (err != nil){
		glog.Warning("Unable to Load the Configuration File, Using Default Configuration")
		return createDefaultConfig()
	}
	
	json.Unmarshal(raw, &config)
	glog.Info("Endpoints: " + fmt.Sprintf("%#v\n", config[0].Endpoints))
	glog.Info("CheckInterval: " +  strconv.FormatInt(config[0].CheckInterval, 10))
	return config[0]
}

func createDefaultConfig() Config {
	defaultConfig := new(Config)
	var defaultEndpoint = []string{DEFAULT_ENDPOINT}
	defaultConfig.Endpoints = defaultEndpoint
	defaultConfig.CheckInterval = FIVE_SECONDS
	glog.Info("Endpoints: " + fmt.Sprintf("%#v\n", defaultConfig.Endpoints))
	glog.Info("CheckInterval: " +  strconv.FormatInt(defaultConfig.CheckInterval, 10))
	return *defaultConfig
}

func doEvery(d time.Duration, f func([]string), endpoints []string) {
	for range time.Tick(d) {
		f(endpoints)
	}
}

func checkAllEndpoints(endpoints []string){
	for _, elment := range endpoints {
		checkEndpoint(elment)
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