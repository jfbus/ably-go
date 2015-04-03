package config_test

import (
	"net/url"
	"time"

	. "github.com/ably/ably-go/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/ably/ably-go/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/ably/ably-go/config"
)

func tsToTime(ts int) time.Time {
	return time.Unix(0, int64(ts)*int64(time.Millisecond)/int64(time.Nanosecond))
}

var _ = Describe("PaginationParams", func() {
	var (
		params config.PaginateParams
		values *url.Values
	)

	BeforeEach(func() {
		values = &url.Values{}
		params = config.PaginateParams{}
	})

	Describe("Values", func() {
		It("returns nil with no values", func() {
			params = config.PaginateParams{}
			err := params.EncodeValues(values)
			Expect(err).NotTo(HaveOccurred())
			Expect(values.Encode()).To(Equal(""))
		})

		It("returns the full params encoded", func() {
			params = config.PaginateParams{
				Limit:     1,
				Direction: "backwards",
				ScopeParams: config.ScopeParams{
					Start: tsToTime(123),
					End:   tsToTime(124),
					Unit:  "hello",
				},
			}
			err := params.EncodeValues(values)
			Expect(err).NotTo(HaveOccurred())
			Expect(values.Encode()).To(Equal("direction=backwards&end=124&limit=1&start=123&unit=hello"))
		})

		Context("with values", func() {
			BeforeEach(func() {
				params = config.PaginateParams{
					Limit:     10,
					Direction: "backwards",
				}
			})

			It("does not return an error", func() {
				err := params.EncodeValues(values)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with a value for ScopeParams", func() {
			BeforeEach(func() {
				params.Start = tsToTime(123)
			})

			It("creates a valid url", func() {
				err := params.EncodeValues(values)
				Expect(err).NotTo(HaveOccurred())
				Expect(values.Get("start")).To(Equal("123"))
			})
		})

		Context("with invalid value for direction", func() {
			BeforeEach(func() {
				params.Direction = "unknown"
			})

			It("resets the value to the default", func() {
				err := params.EncodeValues(values)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with invalid value for limit", func() {
			BeforeEach(func() {
				params.Limit = -1
			})

			It("resets the value to the default", func() {
				err := params.EncodeValues(values)
				Expect(err).NotTo(HaveOccurred())
				Expect(values.Get("limit")).To(Equal("100"))
			})
		})
	})
})
