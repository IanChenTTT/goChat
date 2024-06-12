$(function () {
  console.log("Loaded");
  $(".host_set").on("mouseenter mouseleave", function (e, data) {
    switch (e.type) {
      case 'mouseenter':
        console.log("mouseEnter");
        break;
      case 'mouseleave':
        console.log("mouseLeave");
        break;
    }
  })
const modal = document.querySelector("#nav-Modal");
const openModal = document.querySelector(".open-btn");
const closeModal = document.querySelector(".close-btn");
openModal.addEventListener("click", () => {
  modal.showModal();
});
closeModal.addEventListener("click", () => {
  modal.close();
});
 
});
