package gorequest

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	cookiemonster "github.com/MercuryEngineering/CookieMonster"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gourl"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"runtime"
	"strings"
	"time"
)

// Response 返回内容
type Response struct {
	RequestID             string      `json:"request_id"`              // 请求编号
	RequestUri            string      `json:"request_uri"`             // 请求链接
	RequestParams         Params      `json:"request_params"`          // 请求参数
	RequestMethod         string      `json:"request_method"`          // 请求方式
	RequestHeader         Headers     `json:"request_header"`          // 请求头部
	RequestCookie         string      `json:"request_cookie"`          // 请求Cookie
	RequestTime           time.Time   `json:"request_time"`            // 请求时间
	ResponseHeader        http.Header `json:"response_header"`         // 响应头部
	ResponseStatus        string      `json:"response_status"`         // 响应状态
	ResponseStatusCode    int         `json:"response_status_code"`    // 响应状态码
	ResponseBody          []byte      `json:"response_body"`           // 响应内容
	ResponseContentLength int64       `json:"response_content_length"` // 响应大小
	ResponseTime          time.Time   `json:"response_time"`           // 响应时间
}

// LogFunc 日志函数
type LogFunc func(ctx context.Context, response *LogResponse)

// App 实例
type App struct {
	Uri                          string           // 全局请求地址，没有设置url才会使用
	Error                        error            // 错误
	httpUri                      string           // 请求地址
	httpMethod                   string           // 请求方法
	httpHeader                   Headers          // 请求头
	httpParams                   Params           // 请求参数
	httpCookie                   string           // Cookie
	responseContent              Response         // 返回内容
	httpContentType              string           // 请求内容类型
	debug                        bool             // 是否开启调试模式
	p12Cert                      *tls.Certificate // p12证书内容
	tlsMinVersion, tlsMaxVersion uint16           // TLS版本
	clientIP                     string           // 客户端IP
	logFunc                      LogFunc          // 日志记录函数
	tr                           trace.Tracer     // OpenTelemetry链路追踪
	span                         trace.Span       // OpenTelemetry追踪
}

// NewHttp 实例化
func NewHttp() *App {
	app := &App{
		httpHeader: NewHeaders(),
		httpParams: NewParams(),
	}
	return app
}

// SetDebug 设置调试模式
func (app *App) SetDebug() {
	app.debug = true
}

// SetUri 设置请求地址
func (app *App) SetUri(uri string) {
	if uri != "" {
		app.httpUri = uri
	}
}

// SetMethod 设置请求方式
func (app *App) SetMethod(method string) {
	if method != "" {
		app.httpMethod = method
	}
}

// SetHeader 设置请求头
func (app *App) SetHeader(key, value string) {
	app.httpHeader.Set(key, value)
}

// SetHeaders 批量设置请求头
func (app *App) SetHeaders(headers Headers) {
	for key, value := range headers {
		app.httpHeader.Set(key, value)
	}
}

// SetTlsVersion 设置TLS版本
func (app *App) SetTlsVersion(minVersion, maxVersion uint16) {
	app.tlsMinVersion = minVersion
	app.tlsMaxVersion = maxVersion
}

// SetAuthToken 设置身份验证令牌
func (app *App) SetAuthToken(token string) {
	if token != "" {
		app.httpHeader.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

// SetUserAgent 设置用户代理，传空字符串就随机设置
func (app *App) SetUserAgent(ua string) {
	if ua != "" {
		app.httpHeader.Set("User-Agent", ua)
	}
}

// SetContentTypeJson 设置JSON格式
func (app *App) SetContentTypeJson() {
	app.httpContentType = httpParamsModeJson
}

// SetContentTypeForm 设置FORM格式
func (app *App) SetContentTypeForm() {
	app.httpContentType = httpParamsModeForm
}

// SetContentTypeXml 设置XML格式
func (app *App) SetContentTypeXml() {
	app.httpContentType = httpParamsModeXml
}

// SetParam 设置请求参数
func (app *App) SetParam(key string, value interface{}) {
	app.httpParams.Set(key, value)
}

// SetParams 批量设置请求参数
func (app *App) SetParams(params Params) {
	for key, value := range params {
		app.httpParams.Set(key, value)
	}
}

// SetCookie 设置Cookie
func (app *App) SetCookie(cookie string) {
	if cookie != "" {
		app.httpCookie = cookie
	}
}

// SetP12Cert 设置证书
func (app *App) SetP12Cert(content *tls.Certificate) {
	app.p12Cert = content
}

// SetClientIP 设置客户端IP
func (app *App) SetClientIP(clientIP string) {
	if clientIP != "" {
		app.clientIP = clientIP
	}
}

// Get 发起 GET 请求
func (app *App) Get(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		app.Uri = uri[0]
	}
	// 设置请求方法
	app.httpMethod = http.MethodGet
	return request(app, ctx)
}

// Head 发起 HEAD 请求
func (app *App) Head(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		app.Uri = uri[0]
	}
	// 设置请求方法
	app.httpMethod = http.MethodHead
	return request(app, ctx)
}

// Post 发起 POST 请求
func (app *App) Post(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		app.Uri = uri[0]
	}
	// 设置请求方法
	app.httpMethod = http.MethodPost
	return request(app, ctx)
}

// Put 发起 PUT 请求
func (app *App) Put(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		app.Uri = uri[0]
	}
	// 设置请求方法
	app.httpMethod = http.MethodPut
	return request(app, ctx)
}

// Patch 发起 PATCH 请求
func (app *App) Patch(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		app.Uri = uri[0]
	}
	// 设置请求方法
	app.httpMethod = http.MethodPatch
	return request(app, ctx)
}

// Delete 发起 DELETE 请求
func (app *App) Delete(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		app.Uri = uri[0]
	}
	// 设置请求方法
	app.httpMethod = http.MethodDelete
	return request(app, ctx)
}

// Connect 发起 CONNECT 请求
func (app *App) Connect(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		app.Uri = uri[0]
	}
	// 设置请求方法
	app.httpMethod = http.MethodConnect
	return request(app, ctx)
}

// Options 发起 OPTIONS 请求
func (app *App) Options(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		app.Uri = uri[0]
	}
	// 设置请求方法
	app.httpMethod = http.MethodOptions
	return request(app, ctx)
}

// Trace 发起 TRACE 请求
func (app *App) Trace(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		app.Uri = uri[0]
	}
	// 设置请求方法
	app.httpMethod = http.MethodTrace
	return request(app, ctx)
}

// Request 发起请求
func (app *App) Request(ctx context.Context) (httpResponse Response, err error) {
	return request(app, ctx)
}

// SetLogFunc 设置日志记录方法
func (app *App) SetLogFunc(logFunc LogFunc) {
	app.logFunc = logFunc
}

// SetTracer 设置链路追踪
func (app *App) SetTracer(tr trace.Tracer) {
	app.tr = tr
}

// 请求接口
func request(c *App, ctx context.Context) (httpResponse Response, err error) {

	if c.tr != nil {
		ctx, span := c.tr.Start(ctx, c.Uri)
		defer span.End()
		ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
	}

	// 赋值
	httpResponse.RequestTime = gotime.Current().Time
	httpResponse.RequestUri = c.httpUri
	httpResponse.RequestMethod = c.httpMethod
	httpResponse.RequestParams = c.httpParams.DeepCopy()
	httpResponse.RequestHeader = c.httpHeader.DeepCopy()
	httpResponse.RequestCookie = c.httpCookie

	// 判断网址
	if httpResponse.RequestUri == "" {
		httpResponse.RequestUri = c.Uri
	}
	if httpResponse.RequestUri == "" {
		c.Error = errors.New("没有设置Uri")
		if c.debug {
			slog.DebugContext(ctx, fmt.Sprintf("{%s}------------------------\n", GetRequestIDContext(ctx)))
			slog.DebugContext(ctx, fmt.Sprintf("{%s}请求异常：%v\n", httpResponse.RequestID, c.Error))
		}
		return httpResponse, c.Error
	}

	// 创建 http 客户端
	client := &http.Client{
		// https://uptrace.dev/get/instrument/opentelemetry-net-http.html
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx)
			}),
		)}

	// OpenTelemetry追踪
	c.span = trace.SpanFromContext(ctx)

	transportStatus := false
	transport := &http.Transport{}
	transportTls := &tls.Config{}

	if c.p12Cert != nil {
		transportStatus = true
		// 配置
		transportTls.Certificates = []tls.Certificate{*c.p12Cert}
		transport.DisableCompression = true
	}

	if c.tlsMinVersion != 0 && c.tlsMaxVersion != 0 {
		transportStatus = true
		// 配置
		transportTls.MinVersion = c.tlsMinVersion
		transportTls.MaxVersion = c.tlsMaxVersion
	}

	if transportStatus {
		transport.TLSClientConfig = transportTls
		client = &http.Client{
			Transport: transport,
		}
	}

	// 请求类型
	if c.httpContentType == "" {
		c.httpContentType = httpParamsModeJson
	}
	switch c.httpContentType {
	case httpParamsModeJson:
		httpResponse.RequestHeader.Set("Content-Type", "application/json")
	case httpParamsModeForm:
		httpResponse.RequestHeader.Set("Content-Type", "application/x-www-form-urlencoded")
	case httpParamsModeXml:
		httpResponse.RequestHeader.Set("Content-Type", "text/xml")
	}

	// 跟踪编号
	httpResponse.RequestID = GetRequestIDContext(ctx)
	if httpResponse.RequestID != "" {
		httpResponse.RequestHeader.Set(xRequestID, httpResponse.RequestID)
	}

	// 请求内容
	var reqBody io.Reader

	if httpResponse.RequestMethod != http.MethodGet && c.httpContentType == httpParamsModeJson {
		jsonStr, err := gojson.Marshal(httpResponse.RequestParams)
		if err != nil {
			c.Error = errors.New(fmt.Sprintf("解析出错 %s", err))
			if c.debug {
				slog.DebugContext(ctx, fmt.Sprintf("{%s}请求异常：%v\n", httpResponse.RequestID, c.Error))
			}
			return httpResponse, c.Error
		}
		// 赋值
		reqBody = bytes.NewBuffer(jsonStr)
	}

	if httpResponse.RequestMethod != http.MethodGet && c.httpContentType == httpParamsModeForm {
		// 携带 form 参数
		form := url.Values{}
		for k, v := range httpResponse.RequestParams {
			form.Add(k, GetParamsString(v))
		}
		// 赋值
		reqBody = strings.NewReader(form.Encode())
	}

	if c.httpContentType == httpParamsModeXml {
		reqBody, err = ToXml(httpResponse.RequestParams)
		if err != nil {
			c.Error = errors.New(fmt.Sprintf("解析XML出错 %s", err))
			if c.debug {
				slog.DebugContext(ctx, fmt.Sprintf("{%s}请求异常：%v\n", httpResponse.RequestID, c.Error))
			}
			return httpResponse, c.Error
		}
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, httpResponse.RequestMethod, httpResponse.RequestUri, reqBody)
	if err != nil {
		c.Error = errors.New(fmt.Sprintf("创建请求出错 %s", err))
		if c.debug {
			slog.DebugContext(ctx, fmt.Sprintf("{%s}请求异常：%v\n", httpResponse.RequestID, c.Error))
		}
		return httpResponse, c.Error
	}

	// GET 请求携带查询参数
	if httpResponse.RequestMethod == http.MethodGet {
		q := req.URL.Query()
		for k, v := range httpResponse.RequestParams {
			q.Add(k, GetParamsString(v))
		}
		req.URL.RawQuery = q.Encode()
	}

	// 设置请求头
	if len(httpResponse.RequestHeader) > 0 {
		for key, value := range httpResponse.RequestHeader {
			req.Header.Set(key, fmt.Sprintf("%v", value))
		}
	}

	// 设置Cookie
	if httpResponse.RequestCookie != "" {
		cookies, _ := cookiemonster.ParseString(httpResponse.RequestCookie)
		if len(cookies) > 0 {
			for _, c := range cookies {
				req.AddCookie(c)
			}
		} else {
			req.Header.Set("Cookie", httpResponse.RequestCookie)
		}
	}

	c.span.SetAttributes(attribute.String("request.uri", httpResponse.RequestUri))
	c.span.SetAttributes(attribute.String("request.url", gourl.UriParse(httpResponse.RequestUri).Url))
	c.span.SetAttributes(attribute.String("request.api", gourl.UriParse(httpResponse.RequestUri).Path))
	c.span.SetAttributes(attribute.String("request.method", httpResponse.RequestMethod))
	c.span.SetAttributes(attribute.String("request.header", gojson.JsonEncodeNoError(httpResponse.RequestHeader)))
	c.span.SetAttributes(attribute.String("request.params", gojson.JsonEncodeNoError(httpResponse.RequestParams)))
	if c.debug {
		slog.DebugContext(ctx, fmt.Sprintf("{%s}请求Uri：%s %s\n", httpResponse.RequestID, httpResponse.RequestMethod, httpResponse.RequestUri))
		slog.DebugContext(ctx, fmt.Sprintf("{%s}请求Params Get：%+v\n", httpResponse.RequestID, req.URL.RawQuery))
		slog.DebugContext(ctx, fmt.Sprintf("{%s}请求Params Post：%+v\n", httpResponse.RequestID, reqBody))
		slog.DebugContext(ctx, fmt.Sprintf("{%s}请求Header：%+v\n", httpResponse.RequestID, req.Header))
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		c.Error = errors.New(fmt.Sprintf("请求出错 %s", err))
		if c.debug {
			slog.DebugContext(ctx, fmt.Sprintf("{%s}请求异常：%v\n", httpResponse.RequestID, c.Error))
		}
		return httpResponse, c.Error
	}

	// 最后关闭连接
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
	case "deflate":
		reader = flate.NewReader(resp.Body)
	default:
		reader = resp.Body
	}
	defer reader.Close() // nolint

	// 读取内容
	body, err := io.ReadAll(reader)
	if err != nil {
		c.Error = errors.New(fmt.Sprintf("解析内容出错 %s", err))
		if c.debug {
			slog.DebugContext(ctx, fmt.Sprintf("{%s}请求异常：%v\n", httpResponse.RequestID, c.Error))
		}
		return httpResponse, c.Error
	}

	// 赋值
	httpResponse.ResponseTime = gotime.Current().Time
	httpResponse.ResponseStatus = resp.Status
	httpResponse.ResponseStatusCode = resp.StatusCode
	httpResponse.ResponseHeader = resp.Header
	httpResponse.ResponseBody = body
	httpResponse.ResponseContentLength = resp.ContentLength

	// OpenTelemetry追踪
	c.span.SetAttributes(attribute.String("response.status", httpResponse.ResponseStatus))
	c.span.SetAttributes(attribute.Int("response.status_code", httpResponse.ResponseStatusCode))
	c.span.SetAttributes(attribute.String("response.header", gojson.JsonEncodeNoError(httpResponse.ResponseHeader)))
	if gojson.IsValidJSON(string(httpResponse.ResponseBody)) {
		c.span.SetAttributes(attribute.String("response.body", gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(string(httpResponse.ResponseBody)))))
	} else {
		c.span.SetAttributes(attribute.String("response.body", string(httpResponse.ResponseBody)))
	}
	if c.debug {
		slog.DebugContext(ctx, fmt.Sprintf("{%s}返回Status：%s\n", httpResponse.RequestID, httpResponse.ResponseStatus))
		slog.DebugContext(ctx, fmt.Sprintf("{%s}返回Header：%+v\n", httpResponse.RequestID, httpResponse.ResponseHeader))
		slog.DebugContext(ctx, fmt.Sprintf("{%s}返回Body：%s\n", httpResponse.RequestID, httpResponse.ResponseBody))
		slog.DebugContext(ctx, fmt.Sprintf("{%s}------------------------\n", GetRequestIDContext(ctx)))
	}

	// 调用日志记录函数
	if c.logFunc != nil {
		logData := LogResponse{
			RequestID:          httpResponse.RequestID,
			RequestTime:        httpResponse.RequestTime,
			RequestUri:         httpResponse.RequestUri,
			RequestUrl:         gourl.UriParse(httpResponse.RequestUri).Url,
			RequestApi:         gourl.UriParse(httpResponse.RequestUri).Path,
			RequestMethod:      httpResponse.RequestMethod,
			RequestParams:      gojson.JsonEncodeNoError(httpResponse.RequestParams),
			RequestHeader:      gojson.JsonEncodeNoError(httpResponse.RequestHeader),
			RequestIP:          c.clientIP,
			ResponseHeader:     gojson.JsonEncodeNoError(httpResponse.ResponseHeader),
			ResponseStatusCode: httpResponse.ResponseStatusCode,
			ResponseBody:       string(httpResponse.ResponseBody),
			ResponseBodyJson:   gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(string(httpResponse.ResponseBody))),
			ResponseBodyXml:    gojson.XmlEncodeNoError(gojson.XmlDecodeNoError(httpResponse.ResponseBody)),
			ResponseTime:       httpResponse.ResponseTime,
			GoVersion:          runtime.Version(),
			SdkVersion:         Version,
		}
		if c.span.SpanContext().IsValid() {
			logData.SpanID = c.span.SpanContext().SpanID().String()
			logData.TraceID = c.span.SpanContext().TraceID().String()
		}
		c.logFunc(ctx, &logData)
	}

	return httpResponse, err
}
