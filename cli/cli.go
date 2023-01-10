package cli

//go sai lenh compile thi no se in ra printUsage

//client
import (
	"flag"
	"fmt"
	"gitlabapi/database"
	"gitlabapi/util"
	"log"
	"os"
	"runtime"
)

// cli.go de goi db.go
type CommandLine struct {
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Please run go run . import")
	fmt.Println("import -u update -d day -h hour -m min -s sec")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) Run() {
	cli.validateArgs()
	// init db flag
	importCmd := flag.NewFlagSet("import", flag.ContinueOnError)
	isUpdate := importCmd.Bool("u", false, "is updated")
	import_day := importCmd.Int64("d", 0, "Last day to import")
	import_hour := importCmd.Int64("h", 0, "Last hour to import")
	import_min := importCmd.Int64("m", 0, "Last min to import")
	import_sec := importCmd.Int64("s", 0, "Last sec to import")
	switch os.Args[1] {
	case "import":
		err := importCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	}
	if importCmd.Parsed() {
		db := database.Database{}
		db.HandleDB(*isUpdate, util.GetDate(*import_day*24+*import_hour, *import_min, *import_sec))
	}
}
