package i18n

import (
	"context"

	"github.com/fleezesd/xnightwatch/pkg/i18n"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"golang.org/x/text/language"
	"google.golang.org/grpc/metadata"
)

func Translator(options ...func(*i18n.Options)) middleware.Middleware {
	i := i18n.New(options...)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (resp any, err error) {
			var lang language.Tag
			header := make(metadata.MD)
			key := "Accept-Language"

			// 从请求上下文中获取语言设置
			if tr, ok := transport.FromServerContext(ctx); ok {
				lang = language.Make(tr.RequestHeader().Get(key))
			}

			// 根据获取到的语言标签选择对应的翻译器
			ii := i.Select(lang)

			// 设置响应头中的语言信息
			header.Set(key, ii.Language().String())

			// 更新 context
			ctx = metadata.NewOutgoingContext(ctx, header) // 添加 gRPC 元数据
			ctx = i18n.NewContext(ctx, ii)                 // 将翻译器注入上下文

			return handler(ctx, req)
		}
	}
}
