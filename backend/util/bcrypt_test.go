package util

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func Test_HashPassword(t *testing.T) {
	type args struct {
		password string
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "success",
			args:    args{password: "password"},
			wantErr: nil,
		},
		{
			name:    "failed by password is empty",
			args:    args{password: ""},
			wantErr: ErrPasswordEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.args.password)
			if err != nil {
				if err.Error() == tt.wantErr.Error() {
					return
				} else {
					t.Errorf("HashPassword() error = %v, wantErr = %v", err, tt.wantErr)
					return
				}
			}
			if err := bcrypt.CompareHashAndPassword(
				[]byte(got), []byte(tt.args.password),
			); err != nil {
				t.Errorf("HashPassword() = %v, want %v", got, tt.wantErr)
				return
			}
		})
	}
}
