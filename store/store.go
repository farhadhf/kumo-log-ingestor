package store

import (
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
)

func (db *DB) InsertEvent(event *Event) error {
	metaJSON, err := json.Marshal(event.Meta)
	if err != nil {
		return fmt.Errorf("error marshaling meta: %v", err)
	}

	headersJSON, err := json.Marshal(event.Headers)
	if err != nil {
		return fmt.Errorf("error marshaling headers: %v", err)
	}

	sqlStatement := `INSERT INTO events (
        kumo_id, type, sender, recipient, queue, site, size,
        response_code, response_content, response_command, response_enhanced_code_class,
        response_enhanced_code_subject, response_enhanced_code_detail, peer_name, peer_addr,
        timestamp, created, num_attempts, bounce_classification, egress_pool, egress_source,
        source_address_address, source_address_server, source_address_protocol,
        feedback_report, meta, headers, delivery_protocol, reception_protocol, nodeid,
        tls_cipher, tls_protocol_version, tls_peer_subject_name
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
        $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33
    )`
	_, err = db.Exec(sqlStatement,
		event.ID, event.Type, event.Sender, event.Recipient, event.Queue, event.Site, event.Size,
		event.Response.Code, event.Response.Content, event.Response.Command, event.Response.EnhancedCode.Class,
		event.Response.EnhancedCode.Subject, event.Response.EnhancedCode.Detail, event.PeerAddress.Name, event.PeerAddress.Addr,
		event.Timestamp, event.Created, event.NumAttempts, event.BounceClassification, event.EgressPool, event.EgressSource,
		event.SourceAddress.Address, event.SourceAddress.Server, event.SourceAddress.Protocol,
		event.FeedbackReport, metaJSON, headersJSON, event.DeliveryProtocol, event.ReceptionProtocol, event.NodeID,
		event.TLSCipher, event.TLSProtocolVersion, pq.Array(event.TLSPeerSubjectName),
	)
	return err
}
