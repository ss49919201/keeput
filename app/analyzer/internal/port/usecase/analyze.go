package usecase

import (
	"context"

	"github.com/samber/mo"
)

type AnalyzeInput struct{}
type AnalyzeOutput struct{}

type Analyze = func(context.Context, *AnalyzeInput) mo.Result[*AnalyzeOutput]
