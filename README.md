# OpenAI GPT-3.5-turbo / GPT-4 ロールプレイ チャットボット

OpenAI GPT-3.5-turbo / GPT-4 を使ったチャットボットです。  
4人のキャラクターのうち1名ずつをユーザーとAIが担当してロールプレイします。

<image src="_documents/screen.png" style="width:250px">

※ `/v1/chat/completions` エンドポイントを使用しています。 GPT-4 では動作確認していません。

シングルバイナリにビルドされます。  
アセットや設定ファイルはバイナリに埋め込まれます。


## Settings

chathandlers.go
```go
chatReqBody := CompletionRequest{
    Model:       "gpt-3.5-turbo", // "gpt-3.5-turbo", "gpt-4", "gpt-4-32k", ...
    Messages:    messages,
    MaxTokens:   100, // 生成される応答の最大トークン数を指定します。適宜調整してください
    Temperature: 0.8, // 生成されるテキストのランダム性を制御します。適宜調整してください
}
```

chathistory.go
```go
// 初期化用プロンプト
var hiddenPrompt = `・・・・・・・・・・` // 適宜作成してください
```

static/script.js - generateBotReply
```js
async function generateBotReply({myCharName, yourCharName, prompt}) {
    ...
    if (chatHistory.length > 20) {            // APIに渡す履歴の数
        chatHistory = chatHistory.slice(-20); // 長い会話でエラーが出る場合は減らしてください
    }
}
```

## Usage

```bash
# edit credential (OpenAI API key)
vi .env

# and build executable
make

# Run (port 8080)
./rpchatd

# Run (port 3000)
env PORT=3000 ./rpchatd

# If you want to use with docker, build a docker image
make docker

# and run
docker compose up
```

### Init script for OpenWrt

```bash
cp rpchatd /root/bin/.
chmod 711 /root/bin/rpchatd
cp initd/openwrt/etc/init.d/rpchatd /etc/init.d/.
chmod 755 /etc/init.d/rpchatd

service rpchatd enable
service rpchatd start
service rpchatd stop
service rpchatd disable
service
```

### Init script for Systemd

```bash
cp rpchatd /usr/local/bin/.
chmod 711 /usr/local/bin/rpchatd
cp initd/systemd/etc/default/rpchatd /etc/default/.
chmod 644 /etc/default/rpchatd
cp initd/systemd/etc/systemd/system/rpchatd.service /etc/systemd/system/.
chmod 644 /etc/systemd/system/rpchatd.service

systemctl enable rpchatd
systemctl start rpchatd
systemctl stop rpchatd
systemctl disable rpchatd
systemctl list-units --type=service
```

## License

MIT  
Copyright (c) 2023 Shellyl_N and Authors.
