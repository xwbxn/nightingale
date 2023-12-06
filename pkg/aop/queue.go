package aop

import "github.com/toolkits/pkg/container/list"

var OperationlogQueue = list.NewSafeListLimited(10000000)
