package i18n

import (
	"testing"

	"github.com/nicksnyder/go-i18n/i18n"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func loadTranslation() {
	enTranslation := `[{
    "id": "hello",
    "translation": "Hello world"
  }]`
	jaTranslation := `[{
    "id": "hello",
    "translation": "こんにちは"
  }]`
	_ = i18n.ParseTranslationFileBytes("en-us.all.json", []byte(enTranslation))
	_ = i18n.ParseTranslationFileBytes("ja-jp.all.json", []byte(jaTranslation))
}

func newMetadataContext(ctx context.Context, val string) context.Context {
	md := metadata.Pairs("accept-language", val)
	return metadata.NewContext(ctx, md)
}

func TestDefaultLanguage(t *testing.T) {
	loadTranslation()
	req := "request"
	info := &grpc.UnaryServerInfo{FullMethod: "/test/test"}
	_, err := UnaryI18nHandler(context.Background(), req, info, func(ctx context.Context, _ interface{}) (interface{}, error) {

		T := MustTfunc(ctx)
		if got, want := T("hello"), "Hello world"; got != want {
			t.Errorf("expect T() = %q, but got %q", want, got)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRespectAcceptLanguage(t *testing.T) {
	loadTranslation()
	req := "request"
	info := &grpc.UnaryServerInfo{FullMethod: "/test/test"}
	ctx := newMetadataContext(context.Background(), "ja")
	_, err := UnaryI18nHandler(ctx, req, info, func(ctx context.Context, _ interface{}) (interface{}, error) {

		T := MustTfunc(ctx)
		if got, want := T("hello"), "こんにちは"; got != want {
			t.Errorf("expect T() = %q, but got %q", want, got)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFallbackDefaultLanguage(t *testing.T) {
	loadTranslation()
	req := "request"
	info := &grpc.UnaryServerInfo{FullMethod: "/test/test"}
	ctx := newMetadataContext(context.Background(), "da")
	_, err := UnaryI18nHandler(ctx, req, info, func(ctx context.Context, _ interface{}) (interface{}, error) {

		T := MustTfunc(ctx)
		if got, want := T("hello"), "Hello world"; got != want {
			t.Errorf("expect T() = %q, but got %q", want, got)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
