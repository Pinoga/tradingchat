function submitForm(form, url, handler) {
  const formdata = new FormData(form);
  fetch(url, {
    body: formdata,
    method: "post",
    credentials: "same-origin",
  })
    .then((res) => handler(form, res))
    .catch((err) => {
      console.error(err);
    });
  return false;
}

function handleAuth(form, res) {
  let resMessageEl = document.getElementById(`${form.id}-error`);
  if (!res.ok) {
    resMessageEl.innerText =
      form.id === "login" ? "Invalid credentials!" : "Try another username!";
  } else if (res.status === 200) {
    window.location.href = "/pages/chatroom.html";
  }
}

const loginUrl = "http://localhost:8080/api/login";
const registerUrl = "http://localhost:8080/api/register";
const sendMessageUrl = "";
