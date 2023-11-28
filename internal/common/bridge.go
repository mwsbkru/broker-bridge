package common

import "context"

type Bridge interface {
	Run(ctx context.Context) error
}
