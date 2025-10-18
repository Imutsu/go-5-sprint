package actioninfo

import "fmt"

type DataParser interface {
	Parse(dataString string) error
	ActionInfo() (string, error)
}

func Info(dataset []string, dp DataParser) {
	for _, data := range dataset {
		
		err := dp.Parse(data)
		if err != nil {
			fmt.Printf("Parse error: %v\n", err)
			continue
		}

		info, err := dp.ActionInfo()
		if err != nil {
			fmt.Printf("Information forming error: %v\n", err)
			continue
		}

		fmt.Printf("%s\n", info)
	}
}
