package webhook

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Cache用于防重放攻击。
// 推推webhook回调每次nonce不同，可以用Cache保证nonce的唯一性。
//
// 如果Set(nonce)返回false，说明本次回调是人为重放攻击，可以直接返回错误，不执行业务逻辑。
//
// Cache的实现方式有很多种，这里示例两种简单易实现方式：
//  1. 内存模式，用于单实例执行推推webhook，可以用NewMemCache()实现；
//  2. redis模式，用于分布式执行推推webhook，可以用github.com/eachain/360-tuitui-robot/util/cache.NewRedis()实现。
type Cache interface {
	// Set将nonce写入缓存，如果成功返回true。失败返回false，表示缓存中已经存在该nonce。
	Set(nonce string) (ok bool)
}

// AuthOptions推推webhook回调安全身份验证参数，该参数有两个功能：
//  1. 确保该回调请求是由推推发起。
//  2. 功能二：解决重放安全问题。
type AuthOptions struct {
	Appid  string // 验证是哪个机器人的回调，必填参数
	Secret string // 机器人密钥，必填参数

	// 安全身份验证失败返回的http.StatusCode。
	// 默认为http.StatusUnauthorized(401)。
	FailStatusCode int

	// 默认为time.Now，可自定义。
	Now func() time.Time

	// 针对X-Tuitui-Robot-Timestamp头部，指定回调过期时间。
	// 如果now-X-Tuitui-Robot-Timestamp>Expire，表示该请求已超时，不予处理。
	// 默认为0，表示请求不会过期，即跳过该条件判断。
	Expire time.Duration

	// Cache用于解决防重放安全问题。要求Cache过期时间>=AuthOptions.Expire。
	// 默认为nil，即跳过重放安全检查逻辑。
	Cache Cache

	// Errorf用于输出错误日志，由调用方提供，可以为log.Printf。
	// 默认为nil，表示不输出日志。
	Errorf func(string, ...any)
}

// WithAuthSign安全身份验证。确保该回调请求是由推推发起。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E5%AE%89%E5%85%A8%E8%BA%AB%E4%BB%BD%E9%AA%8C%E8%AF%81
func WithAuthSign(opt *AuthOptions, handler http.Handler) http.Handler {
	if opt == nil || opt.Appid == "" || opt.Secret == "" {
		return handler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqAppid := r.Header.Get("X-Tuitui-Robot-Appid")
		if reqAppid != opt.Appid {
			authFail(opt, w, "webhook: auth sign: appid not match, expect %v, request %v",
				opt.Appid, reqAppid)
			return
		}

		reqTimestamp := r.Header.Get("X-Tuitui-Robot-Timestamp")
		reqTS, err := strconv.ParseInt(reqTimestamp, 10, 64)
		if err != nil {
			authFail(opt, w, "webhook: auth sign: parse request timestamp %q: %v",
				reqTimestamp, err)
			return
		}

		if opt.Expire > 0 {
			var now int64
			if opt.Now != nil {
				now = opt.Now().UnixMilli()
			} else {
				now = time.Now().UnixMilli()
			}

			diff := now - reqTS
			if diff < 0 {
				diff = -diff
			}
			if time.Duration(diff)*time.Millisecond > opt.Expire {
				authFail(opt, w, "webhook: auth sign: request timestamp %v, now %v, diff %v, exceeded expire duration %v",
					reqTimestamp, now, diff, opt.Expire)
				return
			}
		}

		reqNonce := r.Header.Get("X-Tuitui-Robot-Nonce")
		reqChecksum := r.Header.Get("X-Tuitui-Robot-Checksum")

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if opt.Errorf != nil {
				opt.Errorf("webhook: auth sign: read request body: %v", err)
			}
			return
		}
		r.Body.Close()
		r.Body = io.NopCloser(buf)

		hash := sha1.New()
		hash.Write([]byte(opt.Secret))
		hash.Write([]byte(reqTimestamp))
		hash.Write([]byte(reqNonce))
		hash.Write(buf.Bytes())
		checksum := hex.EncodeToString(hash.Sum(nil))

		if checksum != reqChecksum {
			authFail(opt, w, "webhook: auth sign: checksum not match, expect %v, request %v",
				checksum, reqChecksum)
			return
		}

		if opt.Cache != nil {
			if !opt.Cache.Set(reqNonce) {
				authFail(opt, w, "webhook: auth sign: duplicated nonce: %q", reqNonce)
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}

func authFail(opt *AuthOptions, w http.ResponseWriter, format string, args ...any) {
	if opt.FailStatusCode > 0 {
		w.WriteHeader(opt.FailStatusCode)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	if opt.Errorf != nil {
		opt.Errorf(format, args...)
	}
}
