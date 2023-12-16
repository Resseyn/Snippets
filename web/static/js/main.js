var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

	document.addEventListener("DOMContentLoaded", function() {
	var toggleButton = document.getElementById("toggleButton");
	var popupForm = document.getElementById("popupForm");

	toggleButton.addEventListener("click", function() {
		console.log('Привет от JavaScript!');
	if (popupForm.style.display === "none") {
		popupForm.style.display = "flex";
} else {
		popupForm.style.display = "none";
}
});
});