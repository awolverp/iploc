package main

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aWolver/iploc/cmdpkg"
	"github.com/aWolver/iploc/geoip"
)

var (
	writer                = cmdpkg.NewPrinter()
	signum chan os.Signal = make(chan os.Signal)
)

func main() {
	CmdArgs, err := cmdpkg.ParseArgs()
	if err != nil {
		writer.Fatal(err.Error())
	}

	if err := CmdArgs.Err(); err != nil {
		writer.Fatal(err.Error())
	}

	if !CmdArgs.Silent {
		fmt.Println(cmdpkg.COLORS.BOLD, cmdpkg.Banner, cmdpkg.COLORS.RESET)

		if warnings := CmdArgs.Warn(); len(warnings) > 0 {
			for _, v := range warnings {
				writer.Warning(v)
			}
			time.Sleep(time.Second * 2)
		}
	}

	stream, err := geoip.OpenStream("geoip.csv")
	if err != nil {
		writer.Fatal(err.Error())
	}
	defer stream.Close()

	writer.SetFormat(CmdArgs.Format)

	var founded []int
	var start_time time.Time

	signal.Notify(signum, os.Interrupt)

	go func() {
		<-signum
		os.Exit(2)
	}()

	start_time = time.Now()

	if len(CmdArgs.Query) > 0 {

		for _, v := range CmdArgs.Query {
			if !CmdArgs.Silent {
				fmt.Printf("\n%s---------> %s%s%s\n", cmdpkg.COLORS.BLUE, cmdpkg.COLORS.BOLD, v, cmdpkg.COLORS.RESET)
			}

			objects, err := geoip.Resolve(stream, v, CmdArgs.Offset, CmdArgs.Limit)
			if err != nil {
				writer.Error(err.Error())
				time.Sleep(time.Second * 1)
			}

			founded = append(founded, len(objects))

			for _, obj := range objects {
				writer.Print(obj)
			}
		}
	} else if CmdArgs.ListCountries {
		objects, err := geoip.ListCountries(stream)
		if err != nil {
			writer.Error(err.Error())
			time.Sleep(time.Second * 1)
		} else {
			sort.Slice(objects, func(i, j int) bool {
				return objects[i][0] < objects[j][0]
			})

			founded = append(founded, len(objects))

			for i, obj := range objects {
				fmt.Printf("%d.\t%s%s%s [%s / %s]\n", i, cmdpkg.COLORS.BOLD, obj[2], cmdpkg.COLORS.RESET, obj[0], obj[1])
			}
		}
	} else if CmdArgs.All {
		objects, err := geoip.ResolveAll(stream, CmdArgs.Offset, CmdArgs.Limit)
		if err != nil {
			writer.Error(err.Error())
			time.Sleep(time.Second * 1)
		}

		founded = append(founded, len(objects))

		for _, obj := range objects {
			writer.Print(obj)
		}
	}

	dur := time.Since(start_time)

	if !CmdArgs.NoSummery {
		summery := "found "

		var s []string = make([]string, len(founded))
		for i, v := range founded {
			s[i] = strconv.FormatInt(int64(v), 10)
		}

		summery += strings.Join(s, " + ")

		sum := 0
		for _, v := range founded {
			sum += v
		}

		summery += " = (" + strconv.FormatInt(int64(sum), 10) + ") results / " + dur.String()
		fmt.Println("\n- " + summery)
	}
}
