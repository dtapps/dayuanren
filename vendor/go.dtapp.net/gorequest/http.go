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
	RequestId             string      //【请求】编号
	RequestUri            string      //【请求】链接
	RequestParams         Params      //【请求】参数
	RequestMethod         string      //【请求】方式
	RequestHeader         Headers     //【请求】头部
	RequestCookie         string      //【请求】Cookie
	RequestTime           time.Time   //【请求】时间
	ResponseHeader        http.Header //【返回】头部
	ResponseStatus        string      //【返回】状态
	ResponseStatusCode    int         //【返回】状态码
	ResponseBody          []byte      //【返回】内容
	ResponseContentLength int64       //【返回】大小
	ResponseTime          time.Time   //【返回】时间
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
	tr                           trace.Tracer     // 链路追踪
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
func request(app *App, ctx context.Context) (httpResponse Response, err error) {

	if app.tr != nil {

		ctx, span := app.tr.Start(ctx, app.Uri)
		defer span.End()

		ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))

	}

	// 赋值
	httpResponse.RequestTime = gotime.Current().Time
	httpResponse.RequestUri = app.httpUri
	httpResponse.RequestMethod = app.httpMethod
	httpResponse.RequestParams = app.httpParams.DeepCopy()
	httpResponse.RequestHeader = app.httpHeader.DeepCopy()
	httpResponse.RequestCookie = app.httpCookie

	// 判断网址
	if httpResponse.RequestUri == "" {
		httpResponse.RequestUri = app.Uri
	}
	if httpResponse.RequestUri == "" {
		app.Error = errors.New("没有设置Uri")
		if app.debug {
			slog.DebugContext(ctx, fmt.Sprintf("{%s}------------------------\n", getIDContext(ctx)))
			slog.DebugContext(ctx, fmt.Sprintf("{%s}请求异常：%v\n", httpResponse.RequestId, app.Error))
		}
		return httpResponse, app.Error
	}

	// 创建 http 客户端
	client := &http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx)
			}),
		)}

	transportStatus := false
	transport := &http.Transport{}
	transportTls := &tls.Config{}

	if app.p12Cert != nil {
		transportStatus = true
		// 配置
		transportTls.Certificates = []tls.Certificate{*app.p12Cert}
		transport.DisableCompression = true
	}

	if app.tlsMinVersion != 0 && app.tlsMaxVersion != 0 {
		transportStatus = true
		// 配置
		transportTls.MinVersion = app.tlsMinVersion
		transportTls.MaxVersion = app.tlsMaxVersion
	}

	if transportStatus {
		transport.TLSClientConfig = transportTls
		client = &http.Client{
			Transport: transport,
		}
	}

	// 请求类型
	if app.httpContentType == "" {
		app.httpContentType = httpParamsModeJson
	}
	switch app.httpContentType {
	case httpParamsModeJson:
		httpResponse.RequestHeader.Set("Content-Type", "application/json")
	case httpParamsModeForm:
		httpResponse.RequestHeader.Set("Content-Type", "application/x-www-form-urlencoded")
	case httpParamsModeXml:
		httpResponse.RequestHeader.Set("Content-Type", "text/xml")
	}

	// 跟踪编号
	httpResponse.RequestId = getIDContext(ctx)
	if httpResponse.RequestId != "" {
		httpResponse.RequestHeader.Set("X-Request-Id", httpResponse.RequestId)
	}

	// 请求内容
	var reqBody io.Reader

	if httpResponse.RequestMethod != http.MethodGet && app.httpContentType == httpParamsModeJson {
		jsonStr, err := gojson.Marshal(httpResponse.RequestParams)
		if err != nil {
			app.Error = errors.New(fmt.Sprintf("解析出错 %s", err))
			if app.debug {
				slog.DebugContext(ctx, fmt.Sprintf("{%s}请求异常：%v\n", httpResponse.RequestId, app.Error))
			}
			return httpResponse, app.Error
		}
		// 赋值
		reqBody = bytes.NewBuffer(jsonStr)
	}

	if httpResponse.RequestMethod != http.MethodGet && app.httpContentType == httpParamsModeForm {
		// 携带 form 参数
		form := url.Values{}
		for k, v := range httpResponse.RequestParams {
			form.Add(k, GetParamsString(v))
		}
		// 赋值
		reqBody = strings.NewReader(form.Encode())
	}

	if app.httpContentType == httpParamsModeXml {
		reqBody, err = ToXml(httpResponse.RequestParams)
		if err != nil {
			app.Error = errors.New(fmt.Sprintf("解析XML出错 %s", err))
			if app.debug {
				slog.DebugContext(ctx, fmt.Sprintf("{%s}请求异常：%v\n", httpResponse.RequestId, app.Error))
			}
			return httpResponse, app.Error
		}
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, httpResponse.RequestMethod, httpResponse.RequestUri, reqBody)
	if err != nil {
		app.Error = errors.New(fmt.Sprintf("创建请求出错 %s", err))
		if app.debug {
			slog.DebugContext(ctx, fmt.Sprintf("{%s}请求异常：%v\n", httpResponse.RequestId, app.Error))
		}
		return httpResponse, app.Error
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

	if app.debug {
		slog.DebugContext(ctx, fmt.Sprintf("{%s}请求Uri：%s %s\n", httpResponse.RequestId, httpResponse.RequestMethod, httpResponse.RequestUri))
		slog.DebugContext(ctx, fmt.Sprintf("{%s}请求Params Get：%+v\n", httpResponse.RequestId, req.URL.RawQuery))
		slog.DebugContext(ctx, fmt.Sprintf("{%s}请求Params Post：%+v\n", httpResponse.RequestId, reqBody))
		slog.DebugContext(ctx, fmt.Sprintf("{%s}请求Header：%+v\n", httpResponse.RequestId, req.Header))
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		app.Error = errors.New(fmt.Sprintf("请求出错 %s", err))
		if app.debug {
			slog.DebugContext(ctx, fmt.Sprintf("{%s}请求异常：%v\n", httpResponse.RequestId, app.Error))
		}
		return httpResponse, app.Error
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
		app.Error = errors.New(fmt.Sprintf("解析内容出错 %s", err))
		if app.debug {
			slog.DebugContext(ctx, fmt.Sprintf("{%s}请求异常：%v\n", httpResponse.RequestId, app.Error))
		}
		return httpResponse, app.Error
	}

	// 赋值
	httpResponse.ResponseTime = gotime.Current().Time
	httpResponse.ResponseStatus = resp.Status
	httpResponse.ResponseStatusCode = resp.StatusCode
	httpResponse.ResponseHeader = resp.Header
	httpResponse.ResponseBody = body
	httpResponse.ResponseContentLength = resp.ContentLength

	if app.debug {
		slog.DebugContext(ctx, fmt.Sprintf("{%s}返回Status：%s\n", httpResponse.RequestId, httpResponse.ResponseStatus))
		slog.DebugContext(ctx, fmt.Sprintf("{%s}返回Header：%+v\n", httpResponse.RequestId, httpResponse.ResponseHeader))
		slog.DebugContext(ctx, fmt.Sprintf("{%s}返回Body：%s\n", httpResponse.RequestId, httpResponse.ResponseBody))
		slog.DebugContext(ctx, fmt.Sprintf("{%s}------------------------\n", getIDContext(ctx)))
	}

	// 调用日志记录函数
	if app.logFunc != nil {
		app.logFunc(ctx, &LogResponse{
			//HttpResponse:       resp,
			TraceID:            httpResponse.RequestId,
			RequestID:          httpResponse.RequestId,
			RequestTime:        httpResponse.RequestTime,
			RequestUri:         httpResponse.RequestUri,
			RequestUrl:         gourl.UriParse(httpResponse.RequestUri).Url,
			RequestApi:         gourl.UriParse(httpResponse.RequestUri).Path,
			RequestMethod:      httpResponse.RequestMethod,
			RequestParams:      gojson.JsonEncodeNoError(httpResponse.RequestParams),
			RequestHeader:      gojson.JsonEncodeNoError(httpResponse.RequestHeader),
			RequestIP:          app.clientIP,
			ResponseHeader:     gojson.JsonEncodeNoError(httpResponse.ResponseHeader),
			ResponseStatusCode: httpResponse.ResponseStatusCode,
			ResponseBody:       string(httpResponse.ResponseBody),
			ResponseBodyJson:   gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(string(httpResponse.ResponseBody))),
			ResponseBodyXml:    gojson.XmlEncodeNoError(gojson.XmlDecodeNoError(httpResponse.ResponseBody)),
			ResponseTime:       httpResponse.ResponseTime,
			GoVersion:          runtime.Version(),
			SdkVersion:         Version,
		})
	}

	return httpResponse, err
}
