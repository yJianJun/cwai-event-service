package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestRateLimit(t *testing.T) {

	type args struct {
		r           float64
		b           int
		concurrency int
	}
	tests := []struct {
		name    string
		args    args
		ignore  int
		want    int
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "Normal",
			args: args{
				r:           1,
				b:           1,
				concurrency: 1,
			},
			ignore:  http.StatusOK,
			want:    http.StatusOK,
			wantErr: nil,
		},
		{
			name: "TooMany",
			args: args{
				r:           1,
				b:           1,
				concurrency: 2,
			},
			ignore:  http.StatusOK,
			want:    http.StatusTooManyRequests,
			wantErr: fmt.Errorf(""),
		},
		{
			name: "Block",
			args: args{
				r:           0,
				b:           0,
				concurrency: 1,
			},
			want:    http.StatusTooManyRequests,
			wantErr: fmt.Errorf(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.New()
			g.Use(RateLimit(tt.args.r, tt.args.b))
			g.GET("/", func(c *gin.Context) {
				t.Logf("receive request for client: %s", c.ClientIP())
				c.Status(http.StatusOK)
			})

			tc, cancel := context.WithTimeout(context.TODO(), time.Second*10)
			defer cancel()
			eg, ctx := errgroup.WithContext(tc)

			for i := 0; i < tt.args.concurrency; i++ {
				eg.Go(func() error {
					for {
						select {
						case <-ctx.Done():
							return nil
						default:
							req, _ := http.NewRequest(http.MethodGet, "/", nil)
							timeout, tc := context.WithTimeout(ctx, time.Second)
							req = req.WithContext(timeout)
							w := httptest.NewRecorder()
							g.ServeHTTP(w, req)
							tc()

							if tt.ignore > 0 && tt.ignore == w.Code {
								t.Logf("receive ignore response for status: %d", w.Code)
							} else if w.Code == tt.want {
								t.Logf("receive response expect for status: %d", w.Code)
								return fmt.Errorf("")
							} else {
								return fmt.Errorf("unexpect: %d", w.Code)
							}

							time.Sleep(time.Second)
						}
					}
					return nil
				})
			}

			err := eg.Wait()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
