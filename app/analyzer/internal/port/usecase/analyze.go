package usecase

import "github.com/samber/mo"

type AnalyzeInput struct{}
type AnalyzeOutput struct{}

type Analyze = func(in *AnalyzeInput) mo.Result[*AnalyzeOutput]
