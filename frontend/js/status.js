document.addEventListener("DOMContentLoaded", () => {
  // Helper function to add token to fetch requests
  async function authenticatedFetch(url, options = {}) {
    const token = localStorage.getItem("token");
    if (!token) {
      showCustomAlert("User is not authenticated. Please log in.");
      window.location.href = "./login.html";
      throw new Error("User is not authenticated");
    }

    const headers = {
      ...options.headers,
      Authorization: `Bearer ${token}`,
    };

    return fetch(url, { ...options, headers });
  }

  // Fetch vehicle status with authentication
  authenticatedFetch("http://localhost:5150/api/v1/vehicle/status")
    .then((response) => {
      if (!response.ok) {
        throw new Error("Failed to fetch vehicle status");
      }
      return response.json();
    })
    .then((vehicles) => {
      const vehicleStatusList = document.getElementById("vehicleStatusList");
      vehicles.forEach((vehicle) => {
        const statusCard = `
            <div class="card mb-3">
              <div class="card-body">
                <h5 class="card-title">${vehicle.model}</h5>
                <p class="card-text">
                  Location: ${vehicle.location} <br />
                  Charge Level: ${vehicle.charge_level}% <br />
                  Cleanliness: ${vehicle.cleanliness_status}
                </p>
              </div>
            </div>`;
        vehicleStatusList.innerHTML += statusCard;
      });
    })
    .catch((error) => {
      console.error("Error fetching vehicle status:", error);
      showCustomAlert("An error occurred while fetching vehicle status.");
    });
});
