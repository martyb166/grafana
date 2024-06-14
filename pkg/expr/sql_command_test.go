package expr

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"go.opentelemetry.io/otel/trace"
)

func TestNewCommand(t *testing.T) {
	cmd, err := NewSQLCommand("a", "select a from foo, bar", &testEngine{})
	if err != nil && strings.Contains(err.Error(), "feature is not enabled") {
		return
	}

	if err != nil {
		t.Fail()
		return
	}

	for _, v := range cmd.varsToQuery {
		if strings.Contains("foo bar", v) {
			continue
		}
		t.Fail()
		return
	}
}

func TestExecute(t *testing.T) {
	cmd, err := NewSQLCommand("a", "select a from foo, bar", &testEngine{})
	if err != nil {
		t.Fail()
		return
	}

	ctx := context.Background()
	_, err = cmd.Execute(ctx, time.Now(), nil, &testTracer{})
	if err != nil {
		t.Fail()
		return
	}
}

func TestExecuteShouldTruncate(t *testing.T) {
	cmd, err := NewSQLCommand("a", "select a from foo, bar", &testEngine{})
	if err != nil {
		t.Fail()
		return
	}

	cmd.limit = 1

	ctx := context.Background()
	res, err := cmd.Execute(ctx, time.Now(), nil, &testTracer{})
	if err != nil {
		t.Fail()
		return
	}

	f := res.Values.AsDataFrames("a")
	if len(f) != 1 {
		t.Fail()
		return
	}
}

type testEngine struct{}

func (e *testEngine) QueryFrames(name string, query string, frames []*data.Frame) (string, error) {
	return "", nil
}
func (e *testEngine) QueryFramesInto(name string, query string, frames []*data.Frame, f *data.Frame) error {
	frame := &data.Frame{}
	fld := data.NewField("a", nil, []string{"1", "2", "3"})
	frame.Fields = append(frame.Fields, fld)
	f.Fields = frame.Fields
	return nil
}
func (e *testEngine) RunCommands(cmds []string) (string, error) {
	return `[{"table_name": "foo"}]`, nil
}
func (e *testEngine) Destroy() error {
	return nil
}

type testTracer struct {
	trace.Tracer
}

func (t *testTracer) Start(ctx context.Context, name string, s ...trace.SpanStartOption) (context.Context, trace.Span) {
	return ctx, &testSpan{}
}
func (t *testTracer) Inject(context.Context, http.Header, trace.Span) {

}

type testSpan struct {
	trace.Span
}

func (ts *testSpan) End(opt ...trace.SpanEndOption) {
}
