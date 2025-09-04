package printer

import "github.com/ss49919201/keeput/app/analyzer/internal/model"

type PrintAnalysisReport = func(*model.AnalysisReport) error
