package store

import "encoding/json"

type Event struct {
	Type                 string            `json:"type"`
	ID                   string            `json:"id"`
	Sender               string            `json:"sender"`
	Recipient            string            `json:"recipient"`
	Queue                string            `json:"queue"`
	Site                 string            `json:"site"`
	Size                 int               `json:"size"`
	Response             Response          `json:"response"`
	PeerAddress          PeerAddress       `json:"peer_address"`
	Timestamp            int64             `json:"timestamp"`
	Created              int64             `json:"created"`
	NumAttempts          int               `json:"num_attempts"`
	BounceClassification string            `json:"bounce_classification"`
	EgressPool           string            `json:"egress_pool"`
	EgressSource         string            `json:"egress_source"`
	SourceAddress        SourceAddress     `json:"source_address"`
	FeedbackReport       json.RawMessage   `json:"feedback_report"` // Null or varies
	Meta                 map[string]string `json:"meta"`
	Headers              map[string]string `json:"headers"`
	DeliveryProtocol     string            `json:"delivery_protocol"`
	ReceptionProtocol    string            `json:"reception_protocol"`
	NodeID               string            `json:"nodeid"`
	TLSCipher            string            `json:"tls_cipher"`
	TLSProtocolVersion   string            `json:"tls_protocol_version"`
	TLSPeerSubjectName   []string          `json:"tls_peer_subject_name"`
}

type Response struct {
	Code         int          `json:"code"`
	EnhancedCode EnhancedCode `json:"enhanced_code"`
	Content      string       `json:"content"`
	Command      string       `json:"command"`
}

type EnhancedCode struct {
	Class   int `json:"class"`
	Subject int `json:"subject"`
	Detail  int `json:"detail"`
}

type PeerAddress struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type SourceAddress struct {
	Address  string `json:"address"`
	Server   string `json:"server"`
	Protocol string `json:"protocol"`
}
