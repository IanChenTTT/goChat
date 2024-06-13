$(function () {
  console.log("Loaded");
  const history = $(".history");

  // expectingMessage is set to true
  // if the user has just submitted a message
  // and so we should scroll the next message into view when received.
  let expectingMessage = false;
  const conn = new WebSocket(`ws://${location.host}/sub`);

  conn.addEventListener("close", (ev) => {
    appendLog(
      `WebSocket Disconnected code: ${ev.code}, reason: ${ev.reason}`,
      true
    );
    if (ev.code !== 1001) {
      appendLog("Reconnecting in 1s", true);
      setTimeout(dial, 1000);
    }
  });
  conn.addEventListener("open", (ev) => {
    console.info("websocket connected");
  });

  // This is where we handle messages received.
  conn.addEventListener("message", (ev) => {
    if (typeof ev.data !== "string") {
      console.error("unexpected message type", typeof ev.data);
      return;
    }
    const p = appendLog(ev.data);
    if (expectingMessage) {
      p.scrollIntoView();
      expectingMessage = false;
    }
  });

  // appendLog appends the passed text to messageLog.
  function appendLog(text, error) {
    let MSG_GROUP = $(
      `<div class="MSG_GROUP" id="MSG_GROUP">\
  <img src="/assets/img/user.svg"/> \
  <p class="MSG">${new Date().toLocaleTimeString()}: ${text}</p>\
  <div></div>\
</div>`
    );
    history.append(MSG_GROUP);
    return MSG_GROUP;
  }

  const modal = document.querySelector("#nav-Modal");
  const openModal = document.querySelector(".open-btn");
  const closeModal = document.querySelector(".close-btn");
  openModal.addEventListener("click", () => {
    modal.showModal();
  });
  closeModal.addEventListener("click", () => {
    modal.close();
  });

  const modal1 = document.querySelector("#usr-Modal");
  $('.usr-img').on("click", function(){
      modal1.showModal();
  })
  $('.close-btn1').on("click", function() {
    modal1.close();
  });

});
