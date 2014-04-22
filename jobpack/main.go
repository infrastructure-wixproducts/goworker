package main

import (
	"flag"
	"fmt"

	"github.com/discoproject/goworker/jobutil"

	"errors"
	"os"
	"strings"
)

type Inputs []string

func (i *Inputs) String() string {
	return fmt.Sprint(*i)
}
func (i *Inputs) Set(value string) error {
	if len(*i) > 0 {
		return errors.New("Inputs already set")
	}
	for _, input := range strings.Split(value, ",") {
		*i = append(*i, input)
	}
	return nil
}

func main() {
	var master string
	var confFile string
	var inputs Inputs
	var workerDir string

	const (
		defaultMaster = "localhost"
		masterUsage   = "The master node."
		defaultConf   = "/etc/disco/settings.py"
		confUsage     = "The setting file which contains disco settings"

		// TODO also accept a single go file or an executable.
		defaultWorker = ""
		workerUsage   = "The worker directory"

		defaultInputs = ""
		inputUsage    = "The comma separated list of inputs to the job."
	)
	flag.StringVar(&master, "Master", defaultMaster, masterUsage)
	flag.StringVar(&master, "M", defaultMaster, masterUsage)
	flag.StringVar(&confFile, "Conf", defaultConf, confUsage)
	flag.StringVar(&confFile, "C", defaultConf, confUsage)
	flag.StringVar(&workerDir, "Worker", defaultWorker, workerUsage)
	flag.StringVar(&workerDir, "W", defaultWorker, workerUsage)

	flag.Var(&inputs, "Inputs", inputUsage)
	flag.Var(&inputs, "I", inputUsage)

	flag.Parse()

	if workerDir == "" || len(inputs) == 0 {
		fmt.Println("Usage: jobpack -W worker_dir -I input(s)")
		os.Exit(1)
	}

	jobutil.AddFile(confFile)
	jobutil.SetKeyValue("DISCO_MASTER", master)

	CreateJobPack(inputs, workerDir)
	Post()
	Cleanup()
}