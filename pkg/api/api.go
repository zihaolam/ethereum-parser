package api

import (
	"github.com/zihaolam/ethereum-parser/internal/api"
	"github.com/zihaolam/ethereum-parser/internal/logging"
	"github.com/zihaolam/ethereum-parser/pkg/parser"
)

func New(parser parser.Parser, logger logging.Logger) *api.Api {
	return api.New(parser, logger)
}
