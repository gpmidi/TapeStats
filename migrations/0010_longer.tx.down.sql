ALTER TABLE submissions
    ALTER COLUMN tape_alert_flags TYPE VARCHAR(8),
    ALTER COLUMN assigning_org TYPE VARCHAR(8),
    ALTER COLUMN barcode TYPE VARCHAR(32);

ALTER TABLE tapes
    ALTER COLUMN density_code TYPE CHAR(1),
    ALTER COLUMN medium_type TYPE CHAR(1),
    ALTER COLUMN manufacturer TYPE VARCHAR(8),
    ALTER COLUMN assigning_org TYPE VARCHAR(8),
    ALTER COLUMN medium_type_info TYPE CHAR(2);
