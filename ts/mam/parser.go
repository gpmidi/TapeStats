package mam

import (
	"github.com/go-pg/pg/v10"
	"github.com/gpmidi/TapeStats/ts/tsdb"
	"github.com/rs/zerolog"
	"regexp"
	"strconv"
	"strings"
)

var ReParseStupid = regexp.MustCompile(`^\s*(\w+)\s+(.+)\s+\((\w+),\s+(\d+)\s+bytes,\s+(read-(?:write|only))\):(.*)\s*$`)

type Parser struct {
	Log           zerolog.Logger
	ParserRecord  *tsdb.Parser
	VersionRecord *tsdb.ParserVersion
}

func NewParser(log zerolog.Logger, tx *pg.Tx) (*Parser, error) {
	p, pv, err := getRecords(tx, versionGuid)
	if err != nil {
		return nil, err
	}
	return &Parser{
		Log:           log,
		ParserRecord:  p,
		VersionRecord: pv,
	}, nil
}

// FIXME: Add versioning!
func (p *Parser) ParseString(data string) map[string]*Field {
	log := p.Log.With().Logger()
	ret := make(map[string]*Field)

	log.Trace().Str("data", data).Msg("Starting parsing")

	// Will fix up later
	for lineNo, line := range strings.Split(data, "\n") {
		log := log.With().Str("data.line", line).Int("data.lineno", lineNo).Logger()
		log.Trace().Msg("Working on line")
		for i, match := range ReParseStupid.FindAllStringSubmatch(line, -1) {
			log := log.With().Int("match.idx", i).Strs("matches", match).Logger()

			log.Trace().Msg("Found match") // FIXME: To Trace

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

			ret[strings.TrimSpace(match[2])] = &Field{
				Id:          id,
				Name:        strings.TrimSpace(match[2]),
				Type:        strings.TrimSpace(match[3]),
				Bytes:       (int)(bytes),
				IsReadWrite: rw,
				Value:       strings.TrimSpace(match[6]),
			}

			log.Trace().Interface("field", ret[match[2]]).Msg("Made field") // FIXME: To Trace
		}
	}

	log.Trace().Int("fields.len", len(ret)).Msg("Ending parsing")

	return ret
}
