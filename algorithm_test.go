package cose

import (
	"crypto"
	"crypto/sha256"
	"hash"
	"io"
	"reflect"
	"testing"
)

func TestAlgorithm_String(t *testing.T) {
	// run tests
	tests := []struct {
		alg  Algorithm
		want string
	}{
		{AlgorithmPS256, "PS256"},
		{AlgorithmPS384, "PS384"},
		{AlgorithmPS512, "PS512"},
		{AlgorithmRS256, "RS256"},
		{AlgorithmRS384, "RS384"},
		{AlgorithmRS512, "RS512"},
		{AlgorithmES256, "ES256"},
		{AlgorithmES384, "ES384"},
		{AlgorithmES512, "ES512"},
		{AlgorithmEdDSA, "EdDSA"},
		{AlgorithmReserved, "Reserved"},
		{AlgorithmHMAC256_256, "HMAC256/256"},
		{AlgorithmHMAC384_384, "HMAC384/384"},
		{AlgorithmHMAC512_512, "HMAC512/512"},
		{123, "Algorithm(123)"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.alg.String(); got != tt.want {
				t.Errorf("Algorithm.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_computeHash(t *testing.T) {
	// run tests
	data := []byte("hello world")
	tests := []struct {
		name    string
		alg     Algorithm
		want    []byte
		wantErr error
	}{
		{
			name: "PS256",
			alg:  AlgorithmPS256,
			want: []byte{
				0xb9, 0x4d, 0x27, 0xb9, 0x93, 0x4d, 0x3e, 0x08, 0xa5, 0x2e, 0x52, 0xd7, 0xda, 0x7d, 0xab, 0xfa,
				0xc4, 0x84, 0xef, 0xe3, 0x7a, 0x53, 0x80, 0xee, 0x90, 0x88, 0xf7, 0xac, 0xe2, 0xef, 0xcd, 0xe9,
			},
		},
		{
			name: "PS384",
			alg:  AlgorithmPS384,
			want: []byte{
				0xfd, 0xbd, 0x8e, 0x75, 0xa6, 0x7f, 0x29, 0xf7, 0x01, 0xa4, 0xe0, 0x40, 0x38, 0x5e, 0x2e, 0x23,
				0x98, 0x63, 0x03, 0xea, 0x10, 0x23, 0x92, 0x11, 0xaf, 0x90, 0x7f, 0xcb, 0xb8, 0x35, 0x78, 0xb3,
				0xe4, 0x17, 0xcb, 0x71, 0xce, 0x64, 0x6e, 0xfd, 0x08, 0x19, 0xdd, 0x8c, 0x08, 0x8d, 0xe1, 0xbd,
			},
		},
		{
			name: "PS512",
			alg:  AlgorithmPS512,
			want: []byte{
				0x30, 0x9e, 0xcc, 0x48, 0x9c, 0x12, 0xd6, 0xeb, 0x4c, 0xc4, 0x0f, 0x50, 0xc9, 0x02, 0xf2, 0xb4,
				0xd0, 0xed, 0x77, 0xee, 0x51, 0x1a, 0x7c, 0x7a, 0x9b, 0xcd, 0x3c, 0xa8, 0x6d, 0x4c, 0xd8, 0x6f,
				0x98, 0x9d, 0xd3, 0x5b, 0xc5, 0xff, 0x49, 0x96, 0x70, 0xda, 0x34, 0x25, 0x5b, 0x45, 0xb0, 0xcf,
				0xd8, 0x30, 0xe8, 0x1f, 0x60, 0x5d, 0xcf, 0x7d, 0xc5, 0x54, 0x2e, 0x93, 0xae, 0x9c, 0xd7, 0x6f,
			},
		},
		{
			name: "ES256",
			alg:  AlgorithmES256,
			want: []byte{
				0xb9, 0x4d, 0x27, 0xb9, 0x93, 0x4d, 0x3e, 0x08, 0xa5, 0x2e, 0x52, 0xd7, 0xda, 0x7d, 0xab, 0xfa,
				0xc4, 0x84, 0xef, 0xe3, 0x7a, 0x53, 0x80, 0xee, 0x90, 0x88, 0xf7, 0xac, 0xe2, 0xef, 0xcd, 0xe9,
			},
		},
		{
			name: "ES384",
			alg:  AlgorithmES384,
			want: []byte{
				0xfd, 0xbd, 0x8e, 0x75, 0xa6, 0x7f, 0x29, 0xf7, 0x01, 0xa4, 0xe0, 0x40, 0x38, 0x5e, 0x2e, 0x23,
				0x98, 0x63, 0x03, 0xea, 0x10, 0x23, 0x92, 0x11, 0xaf, 0x90, 0x7f, 0xcb, 0xb8, 0x35, 0x78, 0xb3,
				0xe4, 0x17, 0xcb, 0x71, 0xce, 0x64, 0x6e, 0xfd, 0x08, 0x19, 0xdd, 0x8c, 0x08, 0x8d, 0xe1, 0xbd,
			},
		},
		{
			name: "ES512",
			alg:  AlgorithmES512,
			want: []byte{
				0x30, 0x9e, 0xcc, 0x48, 0x9c, 0x12, 0xd6, 0xeb, 0x4c, 0xc4, 0x0f, 0x50, 0xc9, 0x02, 0xf2, 0xb4,
				0xd0, 0xed, 0x77, 0xee, 0x51, 0x1a, 0x7c, 0x7a, 0x9b, 0xcd, 0x3c, 0xa8, 0x6d, 0x4c, 0xd8, 0x6f,
				0x98, 0x9d, 0xd3, 0x5b, 0xc5, 0xff, 0x49, 0x96, 0x70, 0xda, 0x34, 0x25, 0x5b, 0x45, 0xb0, 0xcf,
				0xd8, 0x30, 0xe8, 0x1f, 0x60, 0x5d, 0xcf, 0x7d, 0xc5, 0x54, 0x2e, 0x93, 0xae, 0x9c, 0xd7, 0x6f,
			},
		},
		{
			name:    "Ed25519",
			alg:     AlgorithmEdDSA,
			wantErr: ErrUnavailableHashFunc,
		},
		{
			name: "HMAC256/256",
			alg:  AlgorithmHMAC256_256,
			want: []byte{
				0xb9, 0x4d, 0x27, 0xb9, 0x93, 0x4d, 0x3e, 0x08, 0xa5, 0x2e, 0x52, 0xd7, 0xda, 0x7d, 0xab, 0xfa,
				0xc4, 0x84, 0xef, 0xe3, 0x7a, 0x53, 0x80, 0xee, 0x90, 0x88, 0xf7, 0xac, 0xe2, 0xef, 0xcd, 0xe9,
			},
		},
		{
			name: "HMAC384/384",
			alg:  AlgorithmHMAC384_384,
			want: []byte{
				0xfd, 0xbd, 0x8e, 0x75, 0xa6, 0x7f, 0x29, 0xf7, 0x01, 0xa4, 0xe0, 0x40, 0x38, 0x5e, 0x2e, 0x23,
				0x98, 0x63, 0x03, 0xea, 0x10, 0x23, 0x92, 0x11, 0xaf, 0x90, 0x7f, 0xcb, 0xb8, 0x35, 0x78, 0xb3,
				0xe4, 0x17, 0xcb, 0x71, 0xce, 0x64, 0x6e, 0xfd, 0x08, 0x19, 0xdd, 0x8c, 0x08, 0x8d, 0xe1, 0xbd,
			},
		},
		{
			name: "HMAC512/512",
			alg:  AlgorithmHMAC512_512,
			want: []byte{
				0x30, 0x9e, 0xcc, 0x48, 0x9c, 0x12, 0xd6, 0xeb, 0x4c, 0xc4, 0x0f, 0x50, 0xc9, 0x02, 0xf2, 0xb4,
				0xd0, 0xed, 0x77, 0xee, 0x51, 0x1a, 0x7c, 0x7a, 0x9b, 0xcd, 0x3c, 0xa8, 0x6d, 0x4c, 0xd8, 0x6f,
				0x98, 0x9d, 0xd3, 0x5b, 0xc5, 0xff, 0x49, 0x96, 0x70, 0xda, 0x34, 0x25, 0x5b, 0x45, 0xb0, 0xcf,
				0xd8, 0x30, 0xe8, 0x1f, 0x60, 0x5d, 0xcf, 0x7d, 0xc5, 0x54, 0x2e, 0x93, 0xae, 0x9c, 0xd7, 0x6f,
			},
		},
		{
			name:    "unknown algorithm",
			alg:     0,
			wantErr: ErrUnavailableHashFunc,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.alg.computeHash(data)
			if err != tt.wantErr {
				t.Errorf("Algorithm.computeHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Algorithm.computeHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

type badHash struct{}

func badHashNew() hash.Hash {
	return &badHash{}
}

func (*badHash) Write(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (*badHash) Sum(b []byte) []byte {
	return b
}

func (*badHash) Reset() {}

func (*badHash) Size() int {
	return 0
}
func (*badHash) BlockSize() int {
	return 0
}

func Test_computeHash(t *testing.T) {
	crypto.RegisterHash(crypto.SHA256, badHashNew)
	defer crypto.RegisterHash(crypto.SHA256, sha256.New)

	_, err := computeHash(crypto.SHA256, nil)
	if err != io.EOF {
		t.Fatalf("computeHash() error = %v, wantErr %v", err, io.EOF)
	}
}
