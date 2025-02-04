package store

import (
	"fmt"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"

	"github.com/tilt-dev/tilt/pkg/logger"
	"github.com/tilt-dev/tilt/pkg/model"
	"github.com/tilt-dev/tilt/pkg/model/logstore"
	"github.com/tilt-dev/wmclient/pkg/analytics"
)

type ErrorAction struct {
	Error error
}

func (ErrorAction) Action() {}

func NewErrorAction(err error) ErrorAction {
	return ErrorAction{Error: err}
}

type LogAction struct {
	mn        model.ManifestName
	spanID    logstore.SpanID
	timestamp time.Time
	fields    logger.Fields
	msg       []byte
	level     logger.Level
}

func (LogAction) Action() {}

func (LogAction) Summarize(s *ChangeSummary) {
	s.Log = true
}

func (le LogAction) ManifestName() model.ManifestName {
	return le.mn
}

func (le LogAction) Level() logger.Level {
	return le.level
}

func (le LogAction) Time() time.Time {
	return le.timestamp
}

func (le LogAction) Fields() logger.Fields {
	return le.fields
}

func (le LogAction) Message() []byte {
	return le.msg
}

func (le LogAction) SpanID() logstore.SpanID {
	return le.spanID
}

func (le LogAction) String() string {
	return fmt.Sprintf("manifest: %s, spanID: %s, msg: %q", le.mn, le.spanID, le.msg)
}

func NewLogAction(mn model.ManifestName, spanID logstore.SpanID, level logger.Level, fields logger.Fields, b []byte) LogAction {
	return LogAction{
		mn:        mn,
		spanID:    spanID,
		level:     level,
		timestamp: time.Now(),
		msg:       append([]byte{}, b...),
		fields:    fields,
	}
}

func NewGlobalLogAction(level logger.Level, b []byte) LogAction {
	return LogAction{
		mn:        "",
		spanID:    "",
		level:     level,
		timestamp: time.Now(),
		msg:       append([]byte{}, b...),
	}
}

type K8sEventAction struct {
	Event        *v1.Event
	ManifestName model.ManifestName
}

func (K8sEventAction) Action() {}

func NewK8sEventAction(event *v1.Event, manifestName model.ManifestName) K8sEventAction {
	return K8sEventAction{event, manifestName}
}

func (kEvt K8sEventAction) ToLogAction(mn model.ManifestName) LogAction {
	msg := fmt.Sprintf("[event: %s] %s\n",
		objRefHumanReadable(kEvt.Event.InvolvedObject),
		strings.TrimSpace(kEvt.Event.Message))

	return LogAction{
		mn:        mn,
		spanID:    logstore.SpanID(fmt.Sprintf("events:%s", mn)),
		level:     logger.InfoLvl,
		timestamp: kEvt.Event.LastTimestamp.Time,
		msg:       []byte(msg),
	}
}

func objRefHumanReadable(obj v1.ObjectReference) string {
	kind := strings.ToLower(obj.Kind)
	if obj.Namespace == "" || obj.Namespace == "default" {
		return fmt.Sprintf("%s %s", kind, obj.Name)
	}
	return fmt.Sprintf("%s %s/%s", kind, obj.Namespace, obj.Name)
}

type AnalyticsUserOptAction struct {
	Opt analytics.Opt
}

func (AnalyticsUserOptAction) Action() {}

type AnalyticsNudgeSurfacedAction struct{}

func (AnalyticsNudgeSurfacedAction) Action() {}

type TiltCloudStatusReceivedAction struct {
	Found                    bool
	Username                 string
	TeamName                 string
	IsPostRegistrationLookup bool
	SuggestedTiltVersion     string
}

func (TiltCloudStatusReceivedAction) Action() {}

type UserStartedTiltCloudRegistrationAction struct{}

func (UserStartedTiltCloudRegistrationAction) Action() {}

type PanicAction struct {
	Err error
}

func (PanicAction) Action() {}
