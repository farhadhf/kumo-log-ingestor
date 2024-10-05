package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func Connect(dataSourceName string) (*DB, error) {
	// dataSourceName should be in this format: postgresql://username:password@host:port/db_name?sslmode=disable
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
	//return db, nil
}

func Close(db *DB) error {
	err := db.Close()
	return err
}

func (db *DB) InitDatabase() error {
	var tableOID sql.NullString
	err := db.QueryRow(`SELECT to_regclass('public.events');`).Scan(&tableOID)
	if err != nil {
		return err
	}

	if tableOID.Valid && tableOID.String != "" {
		return nil
	}

	// @TODO: You might want to add partitions on `created` if you're processing a large number of messages.
	_, err = db.Exec(`CREATE TABLE events (
		id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		kumo_id text NOT NULL,
		type text NOT NULL,
		sender text,
		sender_domain text GENERATED ALWAYS AS (substring(sender from '@(.*)$')) STORED,
		recipient text,
		recipient_domain text GENERATED ALWAYS AS (substring(recipient from '@(.*)$')) STORED,
		queue text,
		site text,
		size int,
		response_code int,
		response_content text,
		response_command text,
		response_enhanced_code_class int,
		response_enhanced_code_subject int,
		response_enhanced_code_detail int,
		peer_name text,
		peer_addr text,
		timestamp bigint,
		created bigint,
		num_attempts int,
		bounce_classification text,
		egress_pool text,
		egress_source text,
		source_address_address text,
		source_address_server text,
		source_address_protocol text,
		feedback_report jsonb,
		meta jsonb,
		headers jsonb,
		delivery_protocol text,
		reception_protocol text,
		nodeid text,
		tls_cipher text,
		tls_protocol_version text,
		tls_peer_subject_name text[]
	);`)
	if err != nil {
		return err
	}

	// @TODO Review you queries, you'll want to drop unnecessary indexes to improve write performance.
	// @TODO You'll most probably need to create additional composite indexes depending on your actual queries.
	_, err = db.Exec(`CREATE EXTENSION pg_trgm;
		CREATE INDEX ON events (response_code);
		CREATE INDEX ON events USING GIN (response_content gin_trgm_ops);
		CREATE INDEX ON events (response_enhanced_code_class);
		CREATE INDEX ON events (response_enhanced_code_subject);
		CREATE INDEX ON events (response_enhanced_code_detail);
		CREATE INDEX ON events (delivery_protocol);
		CREATE INDEX ON events (reception_protocol);
		CREATE INDEX ON events (nodeid);
		CREATE INDEX ON events (recipient);
		CREATE INDEX ON events (recipient_domain);
		CREATE INDEX ON events (sender);
		CREATE INDEX ON events (sender_domain);
		CREATE INDEX ON events (type);
		CREATE INDEX ON events (timestamp);
		CREATE INDEX ON events (created);
		CREATE INDEX ON events (bounce_classification);
		CREATE INDEX ON events (egress_pool);
		CREATE INDEX ON events (egress_source);`)
	if err != nil {
		return err
	}
	return nil
}
