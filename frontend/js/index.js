document.addEventListener("DOMContentLoaded", () => {
  fetch(window.location.origin, { method: "HEAD" })
    .then((response) => {
      const servedBy = response.headers.get("X-Served-By");
      const servedByElement = document.getElementById("servedBy");
      if (servedBy) {
        servedByElement.textContent = `Served by: ${servedBy}`;
      } else {
        servedByElement.textContent = `Served by: Unknown`;
      }
    })
    .catch((error) => {
      console.error("Error fetching endpoint info:", error);
    });
});
