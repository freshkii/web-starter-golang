const connButton = document.getElementById("head-conn-button")
const usernameLabel = document.getElementById("head-username-label")

if (localStorage.hasOwnProperty("username") && localStorage.hasOwnProperty("token")) {
	const username = localStorage.getItem("username")
	const token = localStorage.getItem("token")

	connButton.textContent = "Log out"
	usernameLabel.innerText = username

	connButton.addEventListener("click", _ => {
		fetch("/logout", {
			method: "DELETE",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				username: username,
				token: token
			})
		})
		.then(_ => {
			localStorage.clear()
			window.location.reload()
		})
	})
} else {
	connButton.textContent = "Login"
	connButton.addEventListener("click", _ => {
		window.location.replace("/login")
	})
}
