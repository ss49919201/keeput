package persister

import (
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
)

type SaveAnalysisReport = func(*model.AnalysisReport) error
