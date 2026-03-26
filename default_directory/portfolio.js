full = new fullpage("#fullpage", {
  autoScrolling: true,
  scrollHorizontally: false,
  navigation: false,
});
document.getElementById("step").addEventListener("click", () => {
  fullpage_api.moveTo(2);
});
document.getElementById("About").addEventListener("click", () => {
  fullpage_api.moveTo(2);
});
document.getElementById("Demos").addEventListener("click", () => {
  fullpage_api.moveTo(3);
});
document.getElementById("Home").addEventListener("click", () => {
  fullpage_api.moveTo(1);
});
var active = fullpage_api.getActiveSection();

if (active !== 1) {
  document.getElementById("Home").style.display = "unset;";
}

const imageContainerEl = document.querySelector(".image-container");

let x = 0;
let timer;

function updateGallery() {
  imageContainerEl.style.transform = `perspective(1000px) rotateY(${x}deg)`;
  timer = setTimeout(() => {
    x -= 90;
    updateGallery();
  }, 3000);
}

updateGallery();
