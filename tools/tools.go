//go:build tools

package tools

import (
	_ "github.com/air-verse/air"
	_ "github.com/go-task/task/v3/cmd/task"
	_ "github.com/swaggo/swag/cmd/swag"
)

//go:generate go install github.com/go-task/task/v3/cmd/task@latest
//go:generate go install github.com/air-verse/air@latest
//go:generate go install github.com/swaggo/swag/cmd/swag@latest
