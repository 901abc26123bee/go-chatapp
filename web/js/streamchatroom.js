const initializeWebSocketURL = realtimeHostName + `/api/realtime/v1/chatroom/stream`
const statusDiv = document.getElementById('status');
let socket;
let currentRoomId;
const token = sessionStorage.getItem('access-token');
const encodedToken = encodeURIComponent(token);

function initializeWebSocket() {
    if (token == null) {
        displayError('WebSocket error: user not login');
        return;
    }
    const url = initializeWebSocketURL + `?access_token=${encodedToken}`

    const timeout = setTimeout(() => {
        if (socket.readyState !== WebSocket.OPEN) {
            socket.close();
            console.error("WebSocket connection timed out");
        }
    }, 36000000); // 10min
    socket = new WebSocket(url);

    socket.onopen = function () {
        clearTimeout(timeout); // Clear timeout once connected
        console.log('WebSocket connection opened.');
        updateStatus('Connected');
    };

    socket.onmessage = function (event) {
        const message = JSON.parse(event.data);
        displayMessage(`${message.sender_id}: ${message.chat}`);
    };

    socket.onerror = function (error) {
        displayError('WebSocket error: ' + error.message);
    };

    socket.onclose = function () {
        displayError('WebSocket connection closed.');
    };
}

function joinChatRoom() {
    const chatRoomId = document.getElementById('chatRoomId').value;
    if (!chatRoomId) {
        displayError('Please enter a chat room ID.');
        return;
    }
    if (socket && socket.readyState === WebSocket.OPEN) {
        const joinMessage = JSON.stringify({
            type: "CHAT_ROOM_ACTION",
            payload: {
                action: "JOIN_CHAT_ROOM",
                room_id: chatRoomId,
            }
        });
        try {
            socket.send(joinMessage);
            console.log("WebSocket readyState:", socket.readyState);
            console.log("joinChatRoom", joinMessage)

            document.getElementById('chatBox').innerHTML = '';
            displayError('');
            currentRoomId = chatRoomId
            displayCurChatroomID(chatRoomId)
        } catch (error) {
            displayError('failed to join chat room');
            console.error("failed to join chat room:", error);
        }
    } else {
        displayError('WebSocket connection is not open.');
    }
}

function leaveChatRoom() {
    const chatRoomId = document.getElementById('chatRoomId').value;
    if (!chatRoomId) {
        displayError('Please enter a chat room ID.');
        return;
    }
    if (socket && socket.readyState === WebSocket.OPEN) {
        const leaveMessage = JSON.stringify({
            type: "CHAT_ROOM_ACTION",
            payload: {
                action: "LEAVE_CHAT_ROOM",
                room_id: chatRoomId,
            }
        });
        try {
            socket.send(leaveMessage);
            console.log("WebSocket readyState:", socket.readyState);
            console.log("leaveChatRoom", leaveMessage)

            document.getElementById('chatBox').innerHTML = '';
            displayError('');
            currentRoomId = null
            displayCurChatroomID('')
        } catch (error) {
            displayError('failed to leave room');
            console.error("failed to leave room:", error);
        }
    } else {
        displayError('WebSocket connection is not open.');
    }
}

function sendMessage() {
    const message = document.getElementById('message').value;
    if (!message || !currentRoomId) {
        displayError('Please join a chat room and enter a message.');
        return;
    }
    console.log("currentRoomId",currentRoomId)

    if (socket && socket.readyState === WebSocket.OPEN) {
        const chatMessage = JSON.stringify({
            type: "CHAT",
            payload: {
                room_id: currentRoomId,
                action: "CHAT_ROOM_MESSAGE",
                chat: message
            }
        });
        console.log("sendMessage", chatMessage)
        try {
            socket.send(chatMessage);
            document.getElementById('message').value = '';
        } catch (error) {
            displayError('failed to send message');
            console.error("failed to send message:", error);
        }
    } else {
        displayError('WebSocket connection is not open.');
    }
}

function displayMessage(message) {
    const data = JSON.parse(event.data);
    console.log("data", data)
    if (data.type == "CHAT_STREAM") {
        const chatBox = document.getElementById('chatBox');
        chatBox.innerHTML += `<p><strong>${data.payload.sender_id}:</strong> ${data.payload.chat}</p>`;
        chatBox.scrollTop = chatBox.scrollHeight;
    }
}

function displayError(error) {
    document.getElementById('error').innerText = error;
}

function displayCurChatroomID(currentRoomId) {
    document.getElementById('currentRoomId').innerText = currentRoomId;
}

function updateStatus(status) {
    statusDiv.textContent = status;
    statusDiv.style.color = status === 'Connected' ? 'green' : 'red';
}

window.addEventListener('load', function(event) {
    closeConnection();
    initializeWebSocket();
});

window.addEventListener('beforeunload', function(event) {
    closeConnection();
});

window.addEventListener('unload', function(event) {
    closeConnection();
});
function closeConnection() {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
        console.log('WebSocket connection closed manually.');
    }
}