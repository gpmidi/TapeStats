package mam

import (
	"github.com/rs/zerolog"
	"regexp"
	"strconv"
)

var ReParseStupid = regexp.MustCompile(`^(\w+)\s+(.+)\s+\((\w+),\s+(\d+)\s+bytes,\s+(read-(?:write|only))\):(.*)$`)

type Parser struct {
	Log zerolog.Logger
}

func NewParser(log zerolog.Logger) *Parser {
	return &Parser{
		Log: log,
	}
}

// FIXME: Add versioning!
func (p *Parser) ParseString(data string) map[string]*Field {
	log := p.Log
	ret := make(map[string]*Field)
	// Will fix up later
	for i, match := range ReParseStupid.FindAllStringSubmatch(data, 0) {
		log := log.With().Int("idx", i).Strs("matches", match).Logger()

		log.Info().Msg("Found match") // FIXME: To Trace

		id, err := strconv.ParseInt(match[1], 16, 64)
		if err != nil {
			id = 0
			log.Info().Err(err).Msg("Failed to parse int for id")
			continue
		}

		bytes, err := strconv.ParseInt(match[4], 10, 64)
		if err != nil {
			bytes = 0
			log.Info().Err(err).Msg("Failed to parse int for bytes")
			continue
		}

		var rw bool
		if match[5] == "read-only" {
			rw = false
		} else if match[5] == "read-write" {
			rw = true
		} else {
			log.Info().Str("mode", match[5]).Msg("Invalid access mode value")
			continue
		}

		if _, ok := ret[match[2]]; ok {
			log.Info().Str("name.existing", match[2]).Msg("Duplicate name found - Not overwriting")
			continue
		}

		ret[match[2]] = &Field{
			Id:          id,
			Name:        match[2],
			Type:        match[3],
			Bytes:       (int)(bytes),
			IsReadWrite: rw,
			Value:       match[6],
		}

		log.Info().Interface("field", ret[match[2]]).Msg("Made field") // FIXME: To Trace
	}

	return ret
}
