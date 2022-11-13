package cmdpkg

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const Banner string = `___   ____    _ 
|_ _| |  _ \  | |       ___     ___ 
 | |  | |_) | | |      / _ \   / __|
 | |  |  __/  | |___  | (_) | | (__ 
|___| |_|     |_____|  \___/   \___|
`

const HelpText string = `Usage: %s [-h | -list | -all] [OPTIONS] QUERY

*Required:
	QUERY              IP/Country name

*Options:
	-ns                Don't show summery of results
	-silent            Silent output
	-offset N          Set offset for results
	-limit N           Set limit for results

*Output:
	-format FORMAT     Set output format (.BOLD json, csv, default .RESET)

*Other:
	-list             List all countries
	-all              Show all results
`

type CommandLineArguments struct {
	ListCountries bool // -list
	All           bool // -all

	Query     []string // QUERY
	NoSummery bool     // -ns
	Silent    bool     // -silent
	Format    string   // -format

	Offset, Limit int // -offset, -limit
}

func (c CommandLineArguments) Warn() []string {
	var warnings []string

	if (len(c.Query) > 0) && (c.All || c.ListCountries) {
		warnings = append(warnings,
			fmt.Sprintf(
				"You pass %squery%s with other arguments (%s-list%s, %s-all%s), those not effects.",
				COLORS.BOLD, COLORS.RESET, COLORS.BOLD, COLORS.RESET, COLORS.BOLD, COLORS.RESET,
			),
		)
	} else if c.All && c.ListCountries {
		warnings = append(warnings,
			fmt.Sprintf(
				"You use both of %s-list%s and %s-all%s, %s-all%s not effects.",
				COLORS.BOLD, COLORS.RESET, COLORS.BOLD, COLORS.RESET, COLORS.BOLD, COLORS.RESET,
			),
		)
	}

	if c.Offset < 0 {
		warnings = append(warnings,
			fmt.Sprintf("%s-offset%s is smaller than %s0%s!",
				COLORS.BOLD, COLORS.RESET, COLORS.GREEN, COLORS.RESET),
		)
	}

	if c.Limit < 0 {
		warnings = append(warnings,
			fmt.Sprintf("%s-limit%s is smaller than %s0%s!",
				COLORS.BOLD, COLORS.RESET, COLORS.GREEN, COLORS.RESET),
		)
	}

	return warnings
}

func (c CommandLineArguments) Err() error {
	if (len(c.Query) == 0) && !c.All && !c.ListCountries {
		return errors.New("use --help to see information")
	}

	return nil
}

func ParseArgs() (*CommandLineArguments, error) {
	var commands CommandLineArguments

	parser := flag.NewFlagSet("", flag.ExitOnError)
	parser.Usage = func() { fmt.Printf(ParseColors(HelpText), os.Args[0]) }

	parser.BoolVar(&commands.ListCountries, "list", false, "")
	parser.BoolVar(&commands.All, "all", false, "")
	parser.BoolVar(&commands.NoSummery, "ns", false, "")
	parser.BoolVar(&commands.Silent, "silent", false, "")
	parser.StringVar(&commands.Format, "format", "default", "")
	parser.IntVar(&commands.Offset, "offset", 0, "")
	parser.IntVar(&commands.Limit, "limit", 0, "")

	if err := (parser.Parse(os.Args[1:])); err != nil {
		return nil, err
	}

	commands.Query = parser.Args()

	return &commands, nil
}
