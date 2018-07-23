package permission

import (
	"context"
	"testing"
)

func contextWithTokenAndService(token string, service *Service) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Authentication", token)
	ctx = context.WithValue(ctx, "permissionService", service)
	return ctx
}

func TestCheck(t *testing.T) {
	service := initService()

	type args struct {
		ctx        context.Context
		permission string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{contextWithTokenAndService("Bearer t1", service), "permission1"}, false},
		{"", args{contextWithTokenAndService("", service), "permission1"}, true},
		{"", args{contextWithTokenAndService("t1", service), "permission1"}, true},
		{"", args{contextWithTokenAndService("Bearer t1", service), "nonexistantpermission1"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Check(tt.args.ctx, tt.args.permission); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
