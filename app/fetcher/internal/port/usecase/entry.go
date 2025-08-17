package usecase

import "github.com/samber/mo"

type ExecuteETLInput struct{}
type ExecuteETLOutput struct{}

type ExecuteETL = func(in *ExecuteETLInput) mo.Result[*ExecuteETLOutput]
