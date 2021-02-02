# Ably Go Terminal Chat

Basic Go chat program using Ably for networking.

## Setup

Firstly install the Ably Go library:

```bash
~ $ go get -u github.com/ably/ably-go/ably
```

Next, replace the **INSERT_YOUR_API_KEY_HERE** line in `chat.go` with your Ably API key. You can sign up for an account with [Ably](https://www.ably.com/) and access your API key from the [app dashboard](https://www.ably.com/accounts/any/apps/any/app_keys).

Finally, run the application with the following line, specify your ClientID as when asked, and you're ready to start communicating over Ably!

```bash
~ $ go run chat.go
```
