package models

import "strings"

type Message struct {
	ID               string `json:"id" bson:"_id,omitempty"`
	LegacyID         string `json:"-" bson:"id,omitempty"`
	ChatID           string `json:"chat_id" bson:"chat_id"`
	Content          string `json:"content" bson:"content"`
	CorrectedContent string `json:"corrected_content,omitempty" bson:"corrected_content,omitempty"`
	Explanation      string `json:"explanation,omitempty" bson:"explanation,omitempty"`
	Translation      string `json:"translation,omitempty" bson:"translation,omitempty"`
	Sender           string `json:"sender" bson:"sender"`
	Type             string `json:"type" bson:"type"`
	CreatedAt        string `json:"created_at" bson:"created_at"`
}

// Normalize ensures backward compatibility with legacy documents where ID was stored in "id"
// and correction/translation were concatenated inside Content.
func (m *Message) Normalize() {
	if m.ID == "" && m.LegacyID != "" {
		m.ID = m.LegacyID
	}

	raw := strings.TrimSpace(m.Content)
	if raw == "" {
		return
	}

	if strings.EqualFold(m.Sender, "user") || strings.EqualFold(m.Type, "request") {
		if idx := strings.Index(raw, "\n\nCorreção:"); idx != -1 {
			userPart := strings.TrimSpace(raw[:idx])
			corrPart := strings.TrimSpace(raw[idx+len("\n\nCorreção:"):])
			m.Content = userPart
			if m.CorrectedContent == "" {
				m.CorrectedContent = corrPart
			}
		}
		return
	}

	if idx := strings.Index(raw, "\n\n[TRANSLATION]:"); idx != -1 {
		dialogPart := strings.TrimSpace(raw[:idx])
		transPart := strings.TrimSpace(raw[idx+len("\n\n[TRANSLATION]:"):])
		m.Content = dialogPart
		if m.Translation == "" {
			m.Translation = transPart
		}
	}
}

