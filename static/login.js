const usernameInput = document.getElementById("username")
const passwordInput = document.getElementById("password")
const message = document.getElementsByClassName("error-message")[0]

if (localStorage.hasOwnProperty("username") && localStorage.hasOwnProperty("token")) {
	window.location.replace("/")
}

document.getElementById("login-button").addEventListener("click", _ => {
	fetch("/login", {
		method: "POST",
		headers: { "Content-Type" : "application/json" },
		body: JSON.stringify({
			username: usernameInput.value,
			password: passwordInput.value
		})
	})
	.then(response => response.json())
	.then(body => {
		if (body.error) {
			message.innerText = body.error
		} else {
			localStorage.setItem("token", body.token)
			localStorage.setItem("username",usernameInput.value)
			window.location.replace("/")
		}
	})
})
