document.addEventListener("DOMContentLoaded", () => {
  const searchButton = document.getElementById("searchButton");

  searchButton.addEventListener("click", () => {
    const startDate = document.getElementById("startDate").value;
    const endDate = document.getElementById("endDate").value;

    if (!startDate || !endDate) {
      showCustomAlert("Please select both start and end date-time.");
      return;
    }

    const start = new Date(startDate);
    const end = new Date(endDate);

    if (end <= start) {
      showCustomAlert("End date and time must be after start date and time.");
      return;
    }

    fetch(
      `http://localhost:5150/api/v1/vehicle/availability?start_date=${encodeURIComponent(
        startDate
      )}&end_date=${encodeURIComponent(endDate)}`
    )
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to fetch available vehicles");
        }
        return response.json();
      })
      .then((vehicles) => {
        const vehicleList = document.getElementById("vehicleList");
        vehicleList.innerHTML = ""; // Clear previous results

        if (vehicles.length === 0) {
          vehicleList.innerHTML =
            "<p>No vehicles available for the selected date range.</p>";
          return;
        }

        // Calculate rental duration in hours (rounded up)
        const rentalDurationHours = Math.ceil((end - start) / (1000 * 60 * 60));

        vehicles.forEach((vehicle) => {
          const totalPrice =
            rentalDurationHours * vehicle.rental_price_per_hour;

          // Generate a Google Maps query URL for the vehicle's location
          const googleMapsUrl = `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(
            vehicle.location
          )}`;

          const vehicleCard = `
                  <div class="card mb-3">
                    <div class="card-body">
                      <h5 class="card-title">${vehicle.model}</h5>
                      <p class="card-text">
                        Location: ${
                          vehicle.location
                        } <a href="${googleMapsUrl}" target="_blank" class="text-primary">View in Google Maps</a> <br />
                        Rental Price per Hour: $${vehicle.rental_price_per_hour.toFixed(
                          2
                        )} <br />
                        Total Rental Price (for ${rentalDurationHours} ${
            rentalDurationHours > 1 ? "hours" : "hour"
          }): $${totalPrice.toFixed(2)}
                      </p>
                      <button class="btn btn-primary" onclick="makeBooking(${
                        vehicle.vehicle_id
                      }, '${startDate}', '${endDate}', ${vehicle.rental_price_per_hour.toFixed(
            2
          )})">Book Now</button>
                    </div>
                  </div>`;
          vehicleList.innerHTML += vehicleCard;
        });
      })
      .catch((error) => {
        console.error(error);
        showCustomAlert(
          "An error occurred while fetching vehicle availability."
        );
      });
  });
});

function makeBooking(vehicleId, startDate, endDate, pricePerHour) {
  // Calculate rental duration in hours
  const start = new Date(startDate);
  const end = new Date(endDate);
  const rentalDurationHours = Math.ceil((end - start) / (1000 * 60 * 60));

  const totalPrice = (rentalDurationHours * parseFloat(pricePerHour)).toFixed(
    2
  );

  // Encode query parameters
  const queryParams = new URLSearchParams({
    vehicleId,
    start_date: startDate,
    end_date: endDate,
    rentalDuration: rentalDurationHours,
    pricePerHour,
    totalPrice,
  }).toString();

  // Redirect to the checkout page with query parameters
  window.location.href = `./checkout.html?${queryParams}`;
}
