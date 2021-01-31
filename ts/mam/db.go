package mam

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"github.com/gpmidi/TapeStats/ts/tsdb"
)

func getRecords(tx *pg.Tx, versionGuid string) (*tsdb.Parser, *tsdb.ParserVersion, error) {
	pv := new(tsdb.ParserVersion)
	err := tx.Model(pv).Where("guid = ?", versionGuid).Select()
	if err != nil {
		return nil, nil, err
	}
	if pv == nil || pv.Guid == "" {
		return nil, nil, errors.New("no such version")
	}
	return pv.Tool, pv, nil
}
