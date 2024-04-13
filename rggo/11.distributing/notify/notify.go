package notify

import (
	"runtime" // To check the running operating system

	"golang.org/x/text/cases"    // Execute operations with string values
	"golang.org/x/text/language" // Execute operations with string values
)

const (
	SeverityLow = iota
	SeverityNormal
	SeverityUrgent
)

type Severity int

// Notification
type Notify struct {
	title    string
	message  string
	severity Severity
}

func New(title, message string, severity Severity) *Notify {
	return &Notify{
		title:    title,
		message:  message,
		severity: severity,
	}
}

func (s Severity) String() string {
	sev := "low"

	switch s {
	case SeverityLow:
		sev = "low"
	case SeverityNormal:
		sev = "normal"
	case SeverityUrgent:
		sev = "critical"
	}

	if runtime.GOOS == "darwin" {
		caser := cases.Title(language.AmericanEnglish)
		sev = caser.String(sev)
	}

	if runtime.GOOS == "windows" {
		switch s {
		case SeverityLow:
			sev = "Info"
		case SeverityNormal:
			sev = "Warning"
		case SeverityUrgent:
			sev = "Error"
		}
	}

	return sev
}
