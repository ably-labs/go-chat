> [!IMPORTANT]
> This repository uses the Ably Pub/Sub approach for building chat apps. We now offer Ably Chat—a new family of SDKs and APIs that streamline development and manage realtime chat complexity for you. For a modern, easier way to create chat experiences, visit our [Ably Chat documentation](https://ably.com/docs/chat).

# Ably Go Terminal Chat

Basic Go chat program using Ably for networking.

## Setup

Firstly you need to replace the **YOUR_API_KEY** text in `.env.example` with your Ably API key and re-name the file to `.env`. You can sign up for an account with [Ably](https://www.ably.com/) and access your API key from the [app dashboard](https://www.ably.com/accounts/any/apps/any/app_keys). We keep the API key in `.env` and ignore it in `.gitignore` to avoid accidentally sharing the API key.

Next, run the application with the following line, specifying your ClientID as an argument, and you're ready to start communicating over Ably!

```bash
~ $ go run chat.go YOUR_USERNAME
```
