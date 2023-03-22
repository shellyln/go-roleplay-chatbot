package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// https://platform.openai.com/docs/api-reference/chat
// https://platform.openai.com/docs/models/gpt-3
// https://openai.com/pricing
//
// 【OpenAI リクエストパラメーターの説明】
//  1. `prompt`: これは、GPT-3に送信するテキストまたは質問です。モデルは、このプロンプトに基づいて回答を生成します。
//  2. `max_tokens`: 生成される応答の最大トークン数を指定します。トークンは、単語や句読点などの言語の基本単位です。
//  3. `n`: GPT-3に複数の応答を生成させたい場合に使用します。このパラメーターに指定された数だけ、異なる応答が生成されます。
//  4. `temperature`: 生成されるテキストのランダム性を制御します。高い値（例：1.0）は、よりランダムで多様なテキストを生成し、
//     低い値（例：0.1）は、より決定論的で一貫したテキストを生成します。
//  5. `top_p`: 生成されるテキストの多様性を制御するもう1つの方法です。
//     `top_p`は、累積確率のしきい値で、生成プロセスが選択できるトークンの範囲を制限します。通常、0.5から0.95の値が使用されます。
//
// `max_tokens`の最大値は、使用しているGPTモデルのエンジンに依存します。GPT-3の場合、最大値はモデルのトークン数上限になります。
// たとえば、`davinci-codex`エンジンの場合、最大値はモデルのトータルトークン数である2048トークンです。
// ただし、実際にAPIで使用する際には、応答の長さに注意してください。非常に長い応答は、読みやすさやレイテンシが低下する可能性があります。
// また、APIの利用制限も考慮する必要があります。
// したがって、実際の応用では、適切な応答の長さを検討し、`max_tokens`を適切な値に設定することが重要です。

// Web UIからのリクエストハンドラ
func promptHandler() {
	http.HandleFunc("/chat/prompt/", func(w http.ResponseWriter, r *http.Request) {
		var reqPayload ClientPromptReqPayload
		err := json.NewDecoder(r.Body).Decode(&reqPayload)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			log.Printf("500: %s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
			return
		}

		if len(reqPayload.History) == 0 {
			log.Printf("200: %s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
			return
		}
		messages := make([]Message, len(reqPayload.History)+1)

		messages[0].Role = "user"
		messages[0].Content = hiddenPrompt

		for i, item := range reqPayload.History {
			if item.IsBot {
				messages[i+1].Role = "assistant"
			} else {
				messages[i+1].Role = "user"
			}

			if item.IsDirective {
				messages[i+1].Content = fmt.Sprintf("【キャラ変更 私：%s あなた：%s】", item.MyCharName, item.YourCharName)
			} else {
				messages[i+1].Content = fmt.Sprint(item.Prompt)
			}
		}

		chatReqBody := CompletionRequest{
			Model:       "gpt-3.5-turbo",
			Messages:    messages,
			MaxTokens:   100,
			Temperature: 0.8,
		}

		chatRes, err := sendChatRequest(chatReqBody)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			log.Printf("500: %s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
			return
		}
		log.Println("Generated Text:", chatRes.Choices[0].Message.Content)

		resPayload := ClientPromptResPayload{
			Text: chatRes.Choices[0].Message.Content,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resPayload)

		log.Printf("200: %s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
	})
}
