你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# [Ably](https://ably.io)

[![Build Status](https://travis-ci.org/ably/ably-go.png)](https://travis-ci.org/ably/ably-go)

A Go client library for [ably.io](https://ably.io), the real-time messaging service.

## Installation

```bash
go get github.com/ably/ably-go
```

## Using the Realtime API

### Subscribing to a channel

TODO

### Publishing to a channel

TODO

### Presence on a channel

TODO

## Using the REST API

### Introduction

All examples assume a client and/or channel has been created as follows:

```go
client = ably.NewRestClient(ably.ClientOptions{ApiKey: "xxxx"})
channel = client.Channel('test')
```

### Publishing a message to a channel

```go
channel.Publish("myEvent", "Hello!") # => returns an error if any
```

### Querying the History

```go
channel.History(nil) # => rest.PaginatedMessages
```

### Presence on a channel

```go
channel.Presence.Get(nil) # => rest.PaginatedPresenceMessages
```

### Querying the Presence History

```go
channel.Presence.History(nil) # => rest.PaginatedPresenceMessages
```

### Generate Token and Token Request

```go
client.Auth.RequestToken()
client.Auth.CreateTokenRequest()
```

### Fetching your application's stats

```ruby
client.Stats() #=> PaginatedResource
```

## Support and feedback

Please visit https://support.ably.io/ for access to our knowledgebase and to ask for any assistance.

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Ensure you have added suitable tests and the test suite is passing(`bundle exec rspec`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## License

Copyright (c) 2015 Ably, Licensed under an MIT license.  Refer to [LICENSE.txt](LICENSE.txt) for the license terms.
