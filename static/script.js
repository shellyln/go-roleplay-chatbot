const chatContainer = document.querySelector(".chat-container");
const myCharacterSelect = document.getElementById("my-character");
const botCharacterSelect = document.getElementById("bot-character");
const messageInput = document.getElementById("message");
let chatHistory = [];

const characters = [
    { name: "まりな", value: "marina", image: "images/1.png" },
    { name: "みさき", value: "misaki", image: "images/2.png" },
    { name: "かなこ", value: "kanako", image: "images/3.png" },
    { name: "もえか", value: "moeka", image: "images/4.png" },
];

// キャラクター選択の初期化
function initCharacterSelect(select) {
    characters.forEach((character, index) => {
        const option = document.createElement("option");
        option.value = character.value;
        option.textContent = character.name;
        select.appendChild(option);
    });
}

// キャラクター選択の変更 (自分)
myCharacterSelect.addEventListener("change", () => {
    if (myCharacterSelect.value === botCharacterSelect.value) {
        botCharacterSelect.value = myCharacterSelect.oldValue;
    }

    const myCharacter = characters.find((x) => x.value === myCharacterSelect.value) ?? characters[0];

    if (myCharacterSelect.oldValue !== myCharacter.value) {
        moveBubblesLeft();
        updateBubble("left", "right", myCharacter.value);
        myCharacterSelect.oldValue = myCharacter.value;
        notifyCharacterChange();
    }
});

// キャラクター選択の変更 (ボット)
botCharacterSelect.addEventListener("change", () => {
    if (botCharacterSelect.value === myCharacterSelect.value) {
        myCharacterSelect.value = botCharacterSelect.oldValue;
    }

    const myCharacter = characters.find((x) => x.value === myCharacterSelect.value) ?? characters[0];
    const botCharacter = characters.find((x) => x.value === botCharacterSelect.value) ?? characters[1];

    if (botCharacterSelect.oldValue !== botCharacter.value) {
        moveBubblesLeft();
        updateBubble("left", "right", myCharacter.value);
        botCharacterSelect.oldValue = botCharacter.value;
        notifyCharacterChange();
    }
});

// 発言と回答を処理する
async function sendMessage() {
    const message = messageInput.value.trim();
    if (!message) return;

    const myCharacter = characters.find((x) => x.value === myCharacterSelect.value) ?? characters[0];
    const botCharacter = characters.find((x) => x.value === botCharacterSelect.value) ?? characters[1];

    createBubble("right", message, myCharacter.image, myCharacter.value);
    messageInput.value = "";

    const botReply = await generateBotReply({
        myCharName: myCharacter.name,
        yourCharName: botCharacter.name,
        prompt: message,
    });
    createBubble("left", botReply, botCharacter.image, botCharacter.value);
}

// 発言を履歴に入れた後、リクエストして回答を履歴に入れる
async function generateBotReply({myCharName, yourCharName, prompt}) {
    chatHistory.push({
        isBot: false,
        myCharName,
        yourCharName,
        prompt,
    });

    const resp = await fetch('/chat/prompt/',{
        method: 'POST',
        body: JSON.stringify({
            history: chatHistory,
        }),
    });
    if (!resp.ok || resp.status != 200) {
        return '(ごめんなさい、会話が失敗しました)';
    }
    const respText = await resp.text()
    const respPayload = JSON.parse(respText);

    if (!respPayload.text) {
        return '(ごめんなさい、会話が失敗しました)';
    }

    chatHistory.push({
        isBot: true,
        myCharName: yourCharName,
        yourCharName: myCharName,
        prompt: respPayload.text,
    });
    if (chatHistory.length > 20) {
        chatHistory = chatHistory.slice(-20);
    }
    return respPayload.text;
}

// 吹き出しを作成する関数
function createBubble(direction, text, imgSrc, characterClass) {
    const row = document.createElement("div");
    row.classList.add("row", direction, characterClass);

    const icon = document.createElement("div");
    icon.classList.add("icon", characterClass);
    icon.style.backgroundImage = `url(${imgSrc})`;
    row.appendChild(icon);

    const bubble = document.createElement("div");
    bubble.classList.add("bubble", direction, characterClass);
    row.appendChild(bubble);

    const content = document.createElement("span");
    content.textContent = text;
    if (text.length === 1) {
        content.classList.add("emoji");
    }
    bubble.appendChild(content);

    chatContainer.appendChild(row);
    chatContainer.scrollTop = chatContainer.scrollHeight;
}

// 吹き出しの表示を切り替える関数
function updateBubble(oldDirection, newDirection, characterClass) {
    const bubbles = document.querySelectorAll(`.bubble.${characterClass}`);
    bubbles.forEach((bubble) => {
        const row = bubble.parentElement;
        row.classList.remove(oldDirection);
        row.classList.add(newDirection);
        bubble.classList.remove(oldDirection);
        bubble.classList.add(newDirection);
    });
}

// 吹き出しを左側にする (自分が担当していないキャラ)
function moveBubblesLeft() {
    const bubbles = document.querySelectorAll(`.bubble`);
    bubbles.forEach((bubble) => {
        const row = bubble.parentElement;
        row.classList.remove("right");
        row.classList.remove("left");
        bubble.classList.remove("right");
        bubble.classList.add("left");
    });
}

// キャラ変更を履歴に記録する
async function notifyCharacterChange() {
    if (chatHistory.length > 0 && chatHistory[chatHistory.length - 1].isDirective) {
        chatHistory.length--;
    }
    const myCharacter = (characters.find((x) => x.value === myCharacterSelect.value)??characters[0]);
    const botCharacter = (characters.find((x) => x.value === botCharacterSelect.value)??characters[1]);

    chatHistory.push({
        isDirective: true,
        myCharName:myCharacter.name,
        yourCharName:botCharacter.name,
        prompt: 'キャラ切り替え',
    });

    console.log(`自分のキャラクターが ${myCharacter.value} に、ボットのキャラクターが ${botCharacter.value} に変更されました。`);
}

// 発言の送信イベント (keydown)
messageInput.addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
        e.preventDefault();
        sendMessage();
    }
});

// 発言の送信イベント (ボタンクリック)
document.getElementById("send-button").addEventListener("click", () => {
    sendMessage();
});

// 初期化
initCharacterSelect(myCharacterSelect);
initCharacterSelect(botCharacterSelect);
myCharacterSelect.value = 'marina'; // 自分の初期キャラクターをキャラ1に設定
botCharacterSelect.value = 'misaki'; // 相手の初期キャラクターをキャラ2に設定
myCharacterSelect.oldValue = myCharacterSelect.value;
botCharacterSelect.oldValue = botCharacterSelect.value;
notifyCharacterChange();
