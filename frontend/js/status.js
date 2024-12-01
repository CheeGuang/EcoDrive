document.addEventListener("DOMContentLoaded", () => {
  fetch("http://localhost:5150/api/v1/vehicle/status")
    .then((response) => response.json())
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
    });
});
