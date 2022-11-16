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
  console.log(res.url);
  let resMessageEl = document.getElementById(`${form.id}-error`);
  if (!res.ok || (form.id === "register" && res.status === 200)) {
    resMessageEl.innerText =
      form.id === "login" ? "Invalid credentials!" : "Try another username!";
  } else if (res.status === 200 || res.status === 201) {
    window.location.href = "/pages/chatroom.html";
  }
}

const loginUrl = "http://localhost:8080/v1/login";
const registerUrl = "http://localhost:8080/v1/register";
const sendMessageUrl = "";
