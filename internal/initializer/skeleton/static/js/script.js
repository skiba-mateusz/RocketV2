document.addEventListener("DOMContentLoaded", function () {
  const navLinks = document.querySelectorAll(".nav__link");
  const current = document.location.pathname.split("/")[1];

  navLinks.forEach((link) => {
    if (link.pathname.split("/")[1] == current) {
      link.classList.add("nav__link--active");
    }
  });
});
