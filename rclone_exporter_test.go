package main

import (
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

var coreStatsJson = `
{
	"bytes": 0,
	"checks": 20319,
	"deletes": 0,
	"elapsedTime": 0.331276801,
	"errors": 0,
	"fatalError": false,
	"retryError": false,
	"speed": 0,
	"transfers": 0
}
`

func Test_parseStats(t *testing.T) {
	type args struct {
		s CoreStats
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"successful_check_running",
			args{
				CoreStats{
					Bytes:       0,
					Checks:      20319,
					Deletes:     0,
					ElapsedTime: 0.331276801,
					Errors:      0,
					FatalError:  false,
					RetryError:  false,
					Speed:       0,
					Transfers:   0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseStats(tt.args.s)
			assert.Equal(t, tt.args.s.Bytes, testutil.ToFloat64(speedMetric))
			assert.Equal(t, tt.args.s.Bytes, testutil.ToFloat64(bytesMetric))
			assert.Equal(t, tt.args.s.Bytes, testutil.ToFloat64(deletesMetric))
			assert.Equal(t, tt.args.s.Bytes, testutil.ToFloat64(errorsMetric))
			assert.Equal(t, tt.args.s.Bytes, testutil.ToFloat64(fatalErrorMetric))
			if tt.args.s.RetryError {
				assert.Equal(t, float64(1), testutil.ToFloat64(retryErrorMetric))
			} else {
				assert.Equal(t, float64(0), testutil.ToFloat64(retryErrorMetric))
			}
			if tt.args.s.FatalError {
				assert.Equal(t, float64(1), testutil.ToFloat64(fatalErrorMetric))
			} else {
				assert.Equal(t, float64(0), testutil.ToFloat64(fatalErrorMetric))
			}
			assert.Equal(t, tt.args.s.Checks, testutil.ToFloat64(checksMetric))
			assert.Equal(t, tt.args.s.Deletes, testutil.ToFloat64(deletesMetric))
			assert.Equal(t, tt.args.s.Transfers, testutil.ToFloat64(transfersMetric))
		})
	}
}

func Test_convertJson(t *testing.T) {
	type args struct {
		reader io.Reader
		model  *CoreStats
		expected float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"unparsable_content",
			args{
				model:  &CoreStats{},
				reader: strings.NewReader("Unauthenticated"),
				expected: 0,
			},
			true,
		},
		{
			"parse_checks",
			args{
				reader: strings.NewReader(coreStatsJson),
				model:  &CoreStats{},
				expected: 20319,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := convertJson(tt.args.reader, tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			} else {
				assert.Equal(t, tt.args.expected, tt.args.model.Checks)
			}
		})
	}
}
