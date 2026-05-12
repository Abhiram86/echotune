package manual

var argsMap = map[string]string{
	"--repeat":  "repeat",
	"-r":        "repeat",
	"--shuffle": "shuffle",
	"-sh":       "shuffle",
	"--limit":   "limit",
	"-l":        "limit",
}

func OrderedArgParse(args []string) []string {
	var result []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--" {
			break
		}

		if name, isFlag := argsMap[arg]; isFlag {
			result = append(result, name)

			if name != "shuffle" && i+1 < len(args) {
				i++
			}
		}
	}

	return result
}
