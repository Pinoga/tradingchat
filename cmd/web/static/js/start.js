let user;
let room;
let socket;

function onClickRoom(e) {
  if (room === undefined) {
    room = e.target.id;
    enterRoom(room);
  } else if (room !== e.target.id) {
    leaveRoom(room).then(() => enterRoom(e.target.id));
  }
}

async function auth() {
  try {
    res = await fetch("/api/chat/authenticate", {
      body: JSON.stringify({}),
      method: "post",
      credentials: "same-origin",
      headers: new Headers({
        Accept: "application/json, text/plain, */*",
        "Content-Type": "application/json",
      }),
    });
    if (!res.ok || res.status !== 200) {
      throw Error("user not authenticated");
    }
    const data = await res.json();
    user = data.body;
    return user;
  } catch (error) {
    throw Error("An error occurred");
  }
}

function enterRoom(r) {
  const messages = document.getElementById("messages");

  socket = new WebSocket(
    "ws://localhost:8080/api/chat/enter/" + (parseInt(r) - 1)
  );
  socket.onopen = () => {
    console.log("opened");
    messages.replaceChildren();
  };
  socket.onclose = () => {
    console.log("closed");
  };
  socket.onerror = () => {
    leaveRoom();
    console.log("error");
  };
  socket.onmessage = (message) => {
    const content = JSON.parse(message);
    console.log("message:", content);
    const msgTag = document.createElement("div");
    msgTag.textContent = content;
    messages.appendChild(msgTag);
  };
}

function leaveRoom() {
  if (!socket) return;
  return socket.close();
}

function sendMessage(e) {}

auth().catch((err) => (window.location.href = "unauth.html"));
