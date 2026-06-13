var passwordAddon = document.getElementById("password-addon");
if (passwordAddon) {
  passwordAddon.addEventListener("click", function () {
    var e = document.getElementById("password-input");
    if (e) {
      "password" === e.type ? (e.type = "text") : (e.type = "password");
    }
  });
}
