const usernameInput = document.getElementById("username")
const passwordInput = document.getElementById("password")
const verifyPasswordInput = document.getElementById("verify-password")
const message = document.getElementById("error-message")

if (localStorage.hasOwnProperty("username") && localStorage.hasOwnProperty("token")) {
	window.location.replace("/")
}

document.getElementById("signin-button").addEventListener("click", _ => {
	if (verifyPasswordInput.value != passwordInput.value) {
		return
	}

	fetch("/signin", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({
			username: usernameInput.value,
			password: passwordInput.value
		})
	})
	.then(response => response.json())
	.then(body => {
		if (body.error) {
			message.textContent = body.error
		} else {
			localStorage.setItem("username", usernameInput.value)
			localStorage.setItem("token", body.token)
			window.location.replace("/")
		}
	})
})
