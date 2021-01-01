package fibonacci

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

type mockRestoreRdb struct {
	value string
	err   error
}

func (mrrh mockRestoreRdb) Get(ctx context.Context, key string) *redis.StringCmd {
	return redis.NewStringResult(mrrh.value, mrrh.err)
}

func (mrrh mockRestoreRdb) Set(
	ctx context.Context, key string, value interface{}, expiration time.Duration,
) *redis.StatusCmd {
	return redis.NewStatusResult(mrrh.value, mrrh.err)
}

func Test_restoreFibonacci(t *testing.T) {
	type args struct {
		rdb RedisClient
	}

	tests := []struct {
		name    string
		args    args
		want    *Fibonacci
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				rdb: mockRestoreRdb{
					value: "5",
					err:   nil,
				},
			},
			want: &Fibonacci{
				current:  5,
				next:     8,
				previous: 3,
				rwMutex:  &sync.RWMutex{},
			},
			wantErr: false,
		},
		{
			name: "redis get error",
			args: args{
				rdb: mockRestoreRdb{
					value: "",
					err:   errors.New("mock error"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "redis bad number",
			args: args{
				rdb: mockRestoreRdb{
					value: "five",
					err:   nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := restoreFibonacci(tt.args.rdb)
			if (err != nil) != tt.wantErr {
				t.Errorf("restoreFibonacci() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("restoreFibonacci() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeFibonacci(t *testing.T) {
	type args struct {
		rdb RedisClient
	}
	tests := []struct {
		name string
		args args
		want *Fibonacci
	}{
		{
			name: "no restore",
			args: args{
				rdb: mockRestoreRdb{
					value: "",
					err:   errors.New("mock error"),
				},
			},
			want: &Fibonacci{
				current:  0,
				next:     1,
				previous: 0,
				rwMutex:  &sync.RWMutex{},
			},
		},
		{
			name: "with restore",
			args: args{
				rdb: mockRestoreRdb{
					value: "5",
					err:   nil,
				},
			},
			want: &Fibonacci{
				current:  5,
				next:     8,
				previous: 3,
				rwMutex:  &sync.RWMutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitializeFibonacci(tt.args.rdb); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitializeFibonacci() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFibonacci_GetCurrent(t *testing.T) {
	type fields struct {
		current  uint32
		next     uint32
		previous uint32
		rwMutex  *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "simple",
			fields: fields{
				current:  5,
				next:     8,
				previous: 3,
				rwMutex:  &sync.RWMutex{},
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fibonacci{
				current:  tt.fields.current,
				next:     tt.fields.next,
				previous: tt.fields.previous,
				rwMutex:  tt.fields.rwMutex,
			}
			if got := f.GetCurrent(); got != tt.want {
				t.Errorf("Fibonacci.GetCurrent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFibonacci_GetNext(t *testing.T) {
	type fields struct {
		current  uint32
		next     uint32
		previous uint32
		rwMutex  *sync.RWMutex
	}
	type args struct {
		rdb RedisClient
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		{
			name: "simple",
			fields: fields{
				current:  5,
				next:     8,
				previous: 3,
				rwMutex:  &sync.RWMutex{},
			},
			args: args{
				rdb: mockRestoreRdb{
					value: "",
					err:   nil,
				},
			},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fibonacci{
				current:  tt.fields.current,
				next:     tt.fields.next,
				previous: tt.fields.previous,
				rwMutex:  tt.fields.rwMutex,
			}
			if got := f.GetNext(tt.args.rdb); got != tt.want {
				t.Errorf("Fibonacci.GetNext() = %v, want %v", got, tt.want)
			}

			updated := &Fibonacci{
				current:  tt.fields.next,
				next:     tt.fields.current + tt.fields.next,
				previous: tt.fields.current,
				rwMutex:  &sync.RWMutex{},
			}

			if !reflect.DeepEqual(f, updated) {
				t.Errorf("Failed to update state properly to: %v, got: %v instead", updated, f)
			}
		})
	}
}

func TestFibonacci_GetPrevious(t *testing.T) {
	type fields struct {
		current  uint32
		next     uint32
		previous uint32
		rwMutex  *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "simple",
			fields: fields{
				current:  5,
				next:     8,
				previous: 3,
				rwMutex:  &sync.RWMutex{},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fibonacci{
				current:  tt.fields.current,
				next:     tt.fields.next,
				previous: tt.fields.previous,
				rwMutex:  tt.fields.rwMutex,
			}
			if got := f.GetPrevious(); got != tt.want {
				t.Errorf("Fibonacci.GetPrevious() = %v, want %v", got, tt.want)
			}
		})
	}
}
