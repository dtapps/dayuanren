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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"io"
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
	RequestCostTime       int64       `json:"request_cost_time"`       // 请求消耗时长
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
	httpUri                      string           // 请求地址
	httpMethod                   string           // 请求方法
	httpHeader                   Headers          // 请求头
	httpParams                   Params           // 请求参数
	httpCookie                   string           // Cookie
	responseContent              Response         // 返回内容
	httpContentType              string           // 请求内容类型
	p12Cert                      *tls.Certificate // p12证书内容
	tlsMinVersion, tlsMaxVersion uint16           // TLS版本
	clientIP                     string           // 客户端IP
	logFunc                      LogFunc          // 日志记录函数
	trace                        bool             // OpenTelemetry链路追踪
	span                         trace.Span       // OpenTelemetry链路追踪
}

// NewHttp 实例化
func NewHttp() *App {
	c := &App{
		httpHeader: NewHeaders(),
		httpParams: NewParams(),
	}
	c.trace = true
	return c
}

// SetUri 设置请求地址
func (c *App) SetUri(uri string) {
	if uri != "" {
		c.httpUri = uri
	}
}

// SetMethod 设置请求方式
func (c *App) SetMethod(method string) {
	if method != "" {
		c.httpMethod = method
	}
}

// SetHeader 设置请求头
func (c *App) SetHeader(key, value string) {
	c.httpHeader.Set(key, value)
}

// SetHeaders 批量设置请求头
func (c *App) SetHeaders(headers Headers) {
	for key, value := range headers {
		c.httpHeader.Set(key, value)
	}
}

// SetTlsVersion 设置TLS版本
func (c *App) SetTlsVersion(minVersion, maxVersion uint16) {
	c.tlsMinVersion = minVersion
	c.tlsMaxVersion = maxVersion
}

// SetAuthToken 设置身份验证令牌
func (c *App) SetAuthToken(token string) {
	if token != "" {
		c.httpHeader.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

// SetUserAgent 设置用户代理，传空字符串就随机设置
func (c *App) SetUserAgent(ua string) {
	if ua != "" {
		c.httpHeader.Set("User-Agent", ua)
	}
}

// SetContentTypeJson 设置JSON格式
func (c *App) SetContentTypeJson() {
	c.httpContentType = httpParamsModeJson
}

// SetContentTypeForm 设置FORM格式
func (c *App) SetContentTypeForm() {
	c.httpContentType = httpParamsModeForm
}

// SetContentTypeXml 设置XML格式
func (c *App) SetContentTypeXml() {
	c.httpContentType = httpParamsModeXml
}

// SetParam 设置请求参数
func (c *App) SetParam(key string, value interface{}) {
	c.httpParams.Set(key, value)
}

// SetParams 批量设置请求参数
func (c *App) SetParams(params Params) {
	for key, value := range params {
		c.httpParams.Set(key, value)
	}
}

// SetCookie 设置Cookie
func (c *App) SetCookie(cookie string) {
	if cookie != "" {
		c.httpCookie = cookie
	}
}

// SetP12Cert 设置证书
func (c *App) SetP12Cert(content *tls.Certificate) {
	c.p12Cert = content
}

// SetClientIP 设置客户端IP
func (c *App) SetClientIP(clientIP string) {
	if clientIP != "" {
		c.clientIP = clientIP
	}
}

// Get 发起 GET 请求
func (c *App) Get(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		c.Uri = uri[0]
	}
	// 设置请求方法
	c.httpMethod = http.MethodGet
	return request(c, ctx)
}

// Head 发起 HEAD 请求
func (c *App) Head(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		c.Uri = uri[0]
	}
	// 设置请求方法
	c.httpMethod = http.MethodHead
	return request(c, ctx)
}

// Post 发起 POST 请求
func (c *App) Post(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		c.Uri = uri[0]
	}
	// 设置请求方法
	c.httpMethod = http.MethodPost
	return request(c, ctx)
}

// Put 发起 PUT 请求
func (c *App) Put(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		c.Uri = uri[0]
	}
	// 设置请求方法
	c.httpMethod = http.MethodPut
	return request(c, ctx)
}

// Patch 发起 PATCH 请求
func (c *App) Patch(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		c.Uri = uri[0]
	}
	// 设置请求方法
	c.httpMethod = http.MethodPatch
	return request(c, ctx)
}

// Delete 发起 DELETE 请求
func (c *App) Delete(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		c.Uri = uri[0]
	}
	// 设置请求方法
	c.httpMethod = http.MethodDelete
	return request(c, ctx)
}

// Connect 发起 CONNECT 请求
func (c *App) Connect(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		c.Uri = uri[0]
	}
	// 设置请求方法
	c.httpMethod = http.MethodConnect
	return request(c, ctx)
}

// Options 发起 OPTIONS 请求
func (c *App) Options(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		c.Uri = uri[0]
	}
	// 设置请求方法
	c.httpMethod = http.MethodOptions
	return request(c, ctx)
}

// Trace 发起 TRACE 请求
func (c *App) Trace(ctx context.Context, uri ...string) (httpResponse Response, err error) {
	if len(uri) == 1 {
		c.Uri = uri[0]
	}
	// 设置请求方法
	c.httpMethod = http.MethodTrace
	return request(c, ctx)
}

// Request 发起请求
func (c *App) Request(ctx context.Context) (httpResponse Response, err error) {
	return request(c, ctx)
}

// SetLogFunc 设置日志记录方法
func (c *App) SetLogFunc(logFunc LogFunc) {
	c.logFunc = logFunc
}

// 请求接口
func request(c *App, ctx context.Context) (httpResponse Response, err error) {

	// OpenTelemetry链路追踪
	ctx = c.TraceStartSpan(ctx, c.httpMethod)
	defer c.TraceEndSpan()

	// 开始时间
	start := time.Now().UTC()

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
		err = errors.New("没有请求地址")
		c.TraceRecordError(err)
		c.TraceSetStatus(codes.Error, err.Error())
		return httpResponse, err
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
			c.TraceRecordError(err)
			c.TraceSetStatus(codes.Error, err.Error())
			return httpResponse, err
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
			c.TraceRecordError(err)
			c.TraceSetStatus(codes.Error, err.Error())
			return httpResponse, err
		}
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, httpResponse.RequestMethod, httpResponse.RequestUri, reqBody)
	if err != nil {
		c.TraceRecordError(err)
		c.TraceSetStatus(codes.Error, err.Error())
		return httpResponse, err
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

	// OpenTelemetry链路追踪
	c.TraceSetAttributes(attribute.String("request.time", httpResponse.RequestTime.Format(gotime.DateTimeFormat)))
	c.TraceSetAttributes(attribute.String("request.uri", httpResponse.RequestUri))
	c.TraceSetAttributes(attribute.String("request.url", gourl.UriParse(httpResponse.RequestUri).Url))
	c.TraceSetAttributes(attribute.String("request.api", gourl.UriParse(httpResponse.RequestUri).Path))
	c.TraceSetAttributes(attribute.String("request.method", httpResponse.RequestMethod))
	c.TraceSetAttributes(attribute.String("request.cookie", httpResponse.RequestCookie))
	c.TraceSetAttributes(attribute.String("request.header", gojson.JsonEncodeNoError(httpResponse.RequestHeader)))
	c.TraceSetAttributes(attribute.String("request.params", gojson.JsonEncodeNoError(httpResponse.RequestParams)))

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		c.TraceRecordError(err)
		c.TraceSetStatus(codes.Error, err.Error())
		return httpResponse, err
	}
	defer resp.Body.Close() // 关闭连接

	// 结束时间
	end := time.Now().UTC()

	// 请求消耗时长
	httpResponse.RequestCostTime = end.Sub(start).Milliseconds()

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
		c.TraceRecordError(err)
		c.TraceSetStatus(codes.Error, err.Error())
		return httpResponse, err
	}

	// 赋值
	httpResponse.ResponseTime = gotime.Current().Time
	httpResponse.ResponseStatus = resp.Status
	httpResponse.ResponseStatusCode = resp.StatusCode
	httpResponse.ResponseHeader = resp.Header
	httpResponse.ResponseBody = body
	httpResponse.ResponseContentLength = resp.ContentLength

	// OpenTelemetry链路追踪
	c.TraceSetAttributes(attribute.Int64("request.cost_time", httpResponse.RequestCostTime))
	c.TraceSetAttributes(attribute.String("response.time", httpResponse.ResponseTime.Format(gotime.DateTimeFormat)))
	c.TraceSetAttributes(attribute.String("response.status", httpResponse.ResponseStatus))
	c.TraceSetAttributes(attribute.Int("response.status_code", httpResponse.ResponseStatusCode))
	c.TraceSetAttributes(attribute.String("response.header", gojson.JsonEncodeNoError(httpResponse.ResponseHeader)))
	if gojson.IsValidJSON(string(httpResponse.ResponseBody)) {
		c.TraceSetAttributes(attribute.String("response.body", gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(string(httpResponse.ResponseBody)))))
	} else {
		c.TraceSetAttributes(attribute.String("response.body", string(httpResponse.ResponseBody)))
	}
	c.TraceSetAttributes(attribute.Int64("response.content_length", httpResponse.ResponseContentLength))

	// 调用日志记录函数
	if c.logFunc != nil {
		c.logFunc(ctx, &LogResponse{
			SpanID:             c.TraceGetTraceID(),
			TraceID:            c.TraceGetSpanID(),
			RequestID:          httpResponse.RequestID,
			RequestTime:        httpResponse.RequestTime,
			RequestUri:         httpResponse.RequestUri,
			RequestUrl:         gourl.UriParse(httpResponse.RequestUri).Url,
			RequestApi:         gourl.UriParse(httpResponse.RequestUri).Path,
			RequestMethod:      httpResponse.RequestMethod,
			RequestParams:      gojson.JsonEncodeNoError(httpResponse.RequestParams),
			RequestHeader:      gojson.JsonEncodeNoError(httpResponse.RequestHeader),
			RequestCostTime:    httpResponse.RequestCostTime,
			RequestIP:          c.clientIP,
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
