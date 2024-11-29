// navbarFooter.js
document.addEventListener("DOMContentLoaded", () => {
  const loadComponent = (url, placeholderId) => {
    fetch(url)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`Failed to load ${url}.`);
        }
        return response.text();
      })
      .then((data) => {
        document.getElementById(placeholderId).innerHTML = data;
      })
      .catch((error) => console.error("Error loading component:", error));
  };

  // Load Navbar
  loadComponent("./navbar.html", "navbar-container");

  // Load Footer
  loadComponent("./footer.html", "footer-container");
});
