package stdout

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ss49919201/keeput/app/analyzer/internal/model"
)

func PrintAnalysisReport(report *model.AnalysisReport) error {
	b, err := json.Marshal(report)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(os.Stdout, string(b))
	return err
}
