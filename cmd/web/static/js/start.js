let user = null;
let room = null;
let socket = null;

function onClickRoom(e) {
  enterRoom(e.target.id);
}

async function auth() {
  try {
    res = await fetch("/api/authenticate", {
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
    user = data.data;
    return user;
  } catch (error) {
    throw Error("An error occurred");
  }
}

function enterRoom(r) {
  if (room === r) return;
  if (room) leaveRoom(room);

  room = r;

  const messages = document.getElementById("messages");

  socket = new WebSocket(
    "ws://localhost:8080/api/chat/enter/" + (parseInt(r) - 1),
    [],
    {
      Headers: {
        Cookie: document.cookie,
      },
    }
  );

  console.log(room);
  socket.onopen = () => {
    console.log("opened");
    document.getElementById("chatroom").removeAttribute("hidden");
    document.getElementById("chatroom-title").innerText = "Chat room #" + room;
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
    console.log(message);
    const data = JSON.parse(message.data);
    console.log("message:", data);
    const msgTag = document.createElement("div");
    msgTag.style = {
      width: "100%",
    };
    msgTag.textContent = `${data.user}: ${data.content} (${data.timestamp})`;
    messages.appendChild(msgTag);
  };
}

function leaveRoom() {
  if (!socket) return;
  socket.close();
  document.getElementById("chatroom").setAttribute("hidden", true);
  room = null;
}

function sendMessage() {
  if (!socket.OPEN) return;
  console.log("sendMessage");
  const message = document.getElementById("message").value;
  return socket.send(message);
}

auth().catch(() => (window.location.href = "unauth.html"));
