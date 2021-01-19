package ts

import (
	"github.com/gin-gonic/gin"
	"github.com/gpmidi/TapeStats/ts/mam"
	"github.com/gpmidi/TapeStats/ts/tsdb"
	"github.com/rs/zerolog"
	"strconv"
	"time"
)

func (ts *TapeStatsApp) LoadRecordHandler(c *gin.Context) {
	li, err := ts.Ctxer(c)
	if err != nil {
		ts.Log.Error().Err(c.Error(err)).Msg("Problem with getting ts")
		return
	}
	log := li.Log

	log.Error().Msg("Not implemented yet")

	c.JSON(500, gin.H{"error": "Not implemented yet"})
}

func (ts *TapeStatsApp) LoadUnparsedHandler(c *gin.Context) {
	li, err := ts.Ctxer(c)
	if err != nil {
		ts.Log.Error().Err(c.Error(err)).Msg("Problem with getting ts")
		return
	}
	l := li.Log

	accountId := c.Request.PostForm.Get("account-id")
	accountPassword := c.Request.PostForm.Get("account-password")

	l = l.With().Str("account.id", accountId).Logger()

	tx, err := ts.DB.Begin()
	if err != nil {
		l.Error().Err(c.Error(err)).Msg("Problem starting database transaction")
		return
	}
	defer func() {
		if err := tx.Close(); err != nil {
			l.Warn().Err(err).Msg("Problem closing db transaction")
		}
	}()

	l.Debug().Msg("Got unparsed submission (unverified)")

	account, err := ts.getAccount(tx, accountId)
	if err != nil {
		l.Error().Err(c.Error(err)).Msg("Problem getting account")
		return
	}

	l = l.With().Interface("account.obj", account).Logger()

	// Verify password
	// TODO: Move to wrapper
	if ok, err := account.VerifyPassword(accountPassword); err != nil || !ok {
		// Failed validation
		l.Info().Err(err).Bool("password.isok", ok).Msg("Password validation failed")
		c.JSON(403, gin.H{
			"message": "Invalid Account",
			"code":    "Forbidden",
			"request": li.Data(),
		})
		return
	}
	l.Debug().Msg("User auth'ed ok")

	data := c.Request.PostForm.Get("submission-data")
	l = l.With().Int("body.len", len(data)).Logger()
	l.Trace().Msg("Got body")

	// Parse 'em
	fields := mam.NewParser(l).ParseString(data)
	l.Trace().Msg("Got fields")

	tape, sub, err := ts.loadFields(l, account.Id, fields)
	if err != nil {
		l.Error().Err(c.Error(err)).Msg("Problem loading")
		return
	}
	l = l.With().Str("tape.id", tape.Id).Int64("sub.id", sub.Id).Logger()
	l.Trace().Msg("Loaded fields")

	c.JSON(200, gin.H{
		"message": "ok",
		"account": gin.H{
			"id":       account.Id,
			"created":  account.Created,
			"modified": account.Modified,
		},
		"tape": gin.H{
			"id":            tape.Id,
			"created":       tape.Created,
			"modified":      tape.Modified,
			"serial-number": tape.SerialNumber,
			"manufacture":   tape.Manufacture,
			"lto-version":   tape.LTOVersion,
		},
		"submission": gin.H{
			"id":                      sub.Id,
			"created":                 sub.Created,
			"modified":                sub.Modified,
			"load-count":              sub.LoadCount,
			"init-count":              sub.InitCount,
			"volume-change-ref":       sub.VolChangeRef,
			"barcode":                 sub.Barcode,
			"total-life-mbytes-write": sub.TotalMBytesLifeWrite,
			"total-life-mbytes-read":  sub.TotalMBytesLifeRead,
		},
		"request": li.Data(),
	})

	l.Trace().Msg("Done with unparsed submission")
}

func (ts *TapeStatsApp) findFieldGets(fields map[string]*mam.Field, fNames ...string) *mam.Field {
	for _, fName := range fNames {
		val, ok := fields[fName]
		if ok && val != nil && val.Value != "" {
			return val
		}
	}
	return nil
}

func (ts *TapeStatsApp) findFieldGetsValue(fields map[string]*mam.Field, fNames ...string) string {
	field := ts.findFieldGets(fields, fNames...)
	if field == nil {
		return ""
	}
	return field.Value
}

func (ts *TapeStatsApp) findFieldGetsValueInt64(fields map[string]*mam.Field, fNames ...string) int64 {
	field := ts.findFieldGets(fields, fNames...)
	if field == nil {
		return 0
	}
	val, err := strconv.ParseInt(field.Value, 16, 64)
	if err != nil {
		return 0
	}
	return val
}

func (ts *TapeStatsApp) loadFields(l zerolog.Logger, accountId string, fields map[string]*mam.Field) (*tsdb.Tape, *tsdb.Submission, error) {
	tx, err := ts.DB.Begin()
	if err != nil {
		l.Warn().Err(err).Msg("Problem starting db transaction")
		return nil, nil, err
	}
	defer func() {
		if err := tx.Close(); err != nil {
			l.Warn().Err(err).Msg("Problem closing db transaction")
		}
	}()

	// Get LTO version
	var ltoVersion int
	switch ts.findFieldGetsValue(fields, "MEDIUM DENSITY CODE", "FORMATTED DENSITY CODE") {
	case "90":
		ltoVersion = 6
	default:
		ltoVersion = 0
	}

	//// Check if tape is around
	//exists,err:=ts.tapeExists(
	//	tx,
	//	accountId,
	//	ts.findFieldGetsValue(fields,"MEDIUM MANUFACTURER"),
	//	ts.findFieldGetsValue(fields,"MEDIUM MANUFACTURE DATE"),
	//	ts.findFieldGetsValue(fields,"MEDIUM SERIAL NUMBER", "NUMERIC MEDIUM SERIAL NUMBER", "UNIQUE CARTRIDGE IDENTITY", "ALTERNATIVE UNIQUE CARTRIDGE IDENTITY"),
	//	ts.findFieldGetsValue(fields,"MEDIUM DENSITY CODE","FORMATTED DENSITY CODE"),
	//	ts.findFieldGetsValue(fields,"MEDIUM TYPE"),
	//	ltoVersion,
	//	)
	//if err != nil {
	//	l.Warn().Err(err).Msg("Problem finding if tape exists")
	//	return err
	//}

	manufactureDT, err := time.Parse("20060102", ts.findFieldGetsValue(fields, "MEDIUM MANUFACTURE DATE"))
	if err != nil {
		l.Info().Err(err).Msg("Couldn't parse manufacture date")
	}

	// Create the tape
	tape := &tsdb.Tape{
		AccountID:      accountId,
		UCI:            ts.findFieldGetsValue(fields, "UNIQUE CARTRIDGE IDENTITY"),
		AltUCI:         ts.findFieldGetsValue(fields, "ALTERNATIVE UNIQUE CARTRIDGE IDENTITY"),
		SerialNumber:   ts.findFieldGetsValue(fields, "MEDIUM SERIAL NUMBER", "NUMERIC MEDIUM SERIAL NUMBER", "UNIQUE CARTRIDGE IDENTITY", "ALTERNATIVE UNIQUE CARTRIDGE IDENTITY"),
		AssignOrg:      ts.findFieldGetsValue(fields, ""),
		Manufacture:    ts.findFieldGetsValue(fields, "MEDIUM MANUFACTURER"),
		ManufactureDT:  manufactureDT,
		DensityCode:    ts.findFieldGetsValue(fields, "MEDIUM DENSITY CODE", "FORMATTED DENSITY CODE"),
		MediumType:     ts.findFieldGetsValue(fields, "MEDIUM TYPE"),
		MediumTypeInfo: ts.findFieldGetsValue(fields, "MEDIUM TYPE INFORMATION"),
		LTOVersion:     ltoVersion,
		Submissions:    nil,
	}
	res, err := tx.Model(tape).
		OnConflict("(account_id, manufacturer, manufacture_dt, serial_number, density_code, medium_type, lto_version) DO UPDATE").
		Returning("*").
		Insert()
	if err != nil {
		l.Warn().Err(err).Msg("Problem UPSERTING tape record")
		return nil, nil, err
	}
	l = l.With().Int("tape.rows.returned", res.RowsReturned()).Int("tape.rows.affected", res.RowsAffected()).Logger()
	l.Trace().Msg("UPSERTED tape record")

	// Build the submission
	sub := &tsdb.Submission{
		TapeID:               tape.Id,
		Tape:                 tape,
		TapeAlertFlags:       ts.findFieldGetsValue(fields, "TAPEALERT FLAGS"),
		LoadCount:            ts.findFieldGetsValueInt64(fields, "LOAD COUNT"),
		MAMSpaceFree:         ts.findFieldGetsValueInt64(fields, "MAM SPACE REMAINING"),
		AssigningOrg:         ts.findFieldGetsValue(fields, "ASSIGNING ORGANIZATION"),
		FormattedDensityCode: ts.findFieldGetsValueInt64(fields, "FORMATTED DENSITY CODE"),
		InitCount:            ts.findFieldGetsValueInt64(fields, "INITIALIZATION COUNT"),
		VolChangeRef:         ts.findFieldGetsValueInt64(fields, "VOLUME CHANGE REFERENCE"),
		TotalMBytesLifeWrite: ts.findFieldGetsValueInt64(fields, "TOTAL MBYTES WRITTEN IN MEDIUM LIFE"),
		TotalMBytesLifeRead:  ts.findFieldGetsValueInt64(fields, "TOTAL MBYTES READ IN MEDIUM LIFE"),
		Barcode:              ts.findFieldGetsValue(fields, "BARCODE"),
		Raw:                  nil, // Leave empty for now - Later fill in with []Fields -> JSON
		KVS:                  map[string]string{},
	}
	// Set KVS
	for name, field := range fields {
		l := l.With().Str("field.name", name).Interface("field", field).Logger()
		l.Info().Msg("Found k:v field")
		sub.KVS[name] = field.Value
	}
	res, err = tx.Model(sub).
		Returning("*").
		Insert()
	if err != nil {
		l.Warn().Err(err).Msg("Problem INSERTING submission record")
		return nil, nil, err
	}
	l = l.With().Int("sub.rows.returned", res.RowsReturned()).Int("sub.rows.affected", res.RowsAffected()).Logger()
	l.Trace().Msg("INSERTED submission record")

	return tape, sub, nil
}
