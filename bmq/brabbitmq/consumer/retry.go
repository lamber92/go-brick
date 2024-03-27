package consumer

import (
	"errors"
	"go-brick/berror"
	"go-brick/blog/logger"
	"net"
	"strings"
	"time"
)

const _defaultMaxRetryTimes = 5

type RetryHandler struct {
	enable        bool // whether to enable retry
	retryTimes    uint // current retry times
	maxRetryTimes uint // maximum retry times (default is 5, infinite retries when set to 0)

	timer            *time.Timer              // retry wait timer
	timeIntervalFunc func(uint) time.Duration // retry interval strategy method

	consumer *Consumer
}

func newRetryHandler(consumer *Consumer) *RetryHandler {
	res := &RetryHandler{
		consumer:      consumer,
		enable:        true,
		timer:         time.NewTimer(0),
		maxRetryTimes: _defaultMaxRetryTimes,
	}
	res.timeIntervalFunc = defaultRetryTimeInterval
	return res
}

func (hdr *RetryHandler) Enable() bool {
	return hdr.enable
}

// InfiniteRetry whether to DO-NOT limit the number of retries
func (hdr *RetryHandler) InfiniteRetry(err error) bool {
	var (
		infinitely = false
		opError    *net.OpError
		innError   berror.Error
	)
	switch {
	// there are two trigger scenarios for infinite retry:
	// 1. controlled by the business side
	// 2. for errors related to network connection types
	case errors.As(err, &opError):
		infinitely = true
		logger.Infra.WithError(err).Warn(hdr.buildLogPrefix() + "internal network error. triggers continuous retry...")
	case errors.As(err, &innError):
		if innError.Status().Code() == EventCodeRetryInfinitely {
			infinitely = true
		} else {
			// try to get the original error
			innerErr := innError.Cause()
			if innerErr == nil {
				// if the original error type cannot be obtained,
				// give up and retry infinitely.
				break
			}
			var opError *net.OpError
			if errors.As(innerErr, &opError) {
				infinitely = true
			}
		}
	default:
		// this is a last resort.
		// when the error content contains the following content, retry infinitely:
		// - bad connection
		// - connection refused
		// - invalid connection
		if strings.Contains(err.Error(), "connection") {
			infinitely = true
		}
	}
	return infinitely
}

// ExceededLimit check whether the max retried times has been exceeded
func (hdr *RetryHandler) ExceededLimit() bool {
	return hdr.maxRetryTimes > 0 && hdr.retryTimes >= hdr.maxRetryTimes
}

func (hdr *RetryHandler) ClearRetriedTimes() {
	hdr.retryTimes = 0
}

func (hdr *RetryHandler) CloneConfig(source *RetryHandler) {
	// release orig timer
	hdr.timer.Stop()

	hdr.enable = source.enable
	hdr.retryTimes = source.retryTimes
	hdr.maxRetryTimes = source.maxRetryTimes
	hdr.timer = source.timer
	hdr.timeIntervalFunc = source.timeIntervalFunc
}

func (hdr *RetryHandler) Close() {
	hdr.timer.Stop()
}

func (hdr *RetryHandler) waitForNextRetry(err error, event string) error {
	// the retry times of retries increases automatically
	hdr.retryTimes++
	// get the time interval for the next retry
	interval := hdr.timeIntervalFunc(hdr.retryTimes)
	logger.Infra.WithError(err).
		Infof(hdr.buildLogPrefix()+"[%s] retry consumption will be executed in %.3f seconds. current retry times: %d",
			event, interval.Seconds(), hdr.retryTimes)

	hdr.timer.Reset(interval)
	defer hdr.timer.Stop()
	select {
	case _, ok := <-hdr.consumer.exitWorker:
		if ok {
			close(hdr.consumer.exitWorker)
		}
		return berror.NewClientClose(err, hdr.buildLogPrefix()+"trigger shutdown or reconnection mechanism")
	case <-hdr.timer.C:
		// wait to return
	}
	return nil
}

func defaultRetryTimeInterval(retryTimes uint) time.Duration {
	switch retryTimes {
	case 1:
		return time.Second
	case 2:
		return 5 * time.Second
	case 3:
		return 10 * time.Second
	case 4:
		return 30 * time.Second
	default:
		return time.Minute
	}
}

func (hdr *RetryHandler) buildLogPrefix() string {
	return hdr.consumer.buildLogPrefix()
}
