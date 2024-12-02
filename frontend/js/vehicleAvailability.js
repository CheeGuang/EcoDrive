document.addEventListener("DOMContentLoaded", async () => {
  const searchButton = document.getElementById("searchButton");

  try {
    // Get the user ID from the token in localStorage
    const token = localStorage.getItem("token"); // Replace with actual token key
    if (!token) {
      showCustomAlert("User is not authenticated. Please log in.");
      return;
    }

    const userId = decodeToken(token).user_id;
    if (!userId) {
      showCustomAlert("Failed to decode user information.");
      return;
    }

    // Fetch membership status
    const membershipResponse = await fetch(
      `http://localhost:5100/api/v1/user/membership/status?user_id=${encodeURIComponent(
        userId
      )}`
    );

    if (!membershipResponse.ok) {
      throw new Error("Failed to fetch membership status.");
    }

    const { membership_level } = await membershipResponse.json();

    // Fetch active bookings for the user
    const bookingResponse = await fetch(
      `http://localhost:5150/api/v1/vehicle/booking/user/${userId}`
    );

    if (!bookingResponse.ok) {
      throw new Error("Failed to fetch active bookings.");
    }

    const activeBookingsData = await bookingResponse.json();

    if (activeBookingsData !== null) {
      // Filter active bookings: Only bookings where return_date is after today
      const activeBookings = activeBookingsData.filter((booking) => {
        const returnDate = new Date(booking.return_date);
        const today = new Date();
        return returnDate > today;
      });

      // Define booking limits based on membership level
      const bookingLimits = {
        Basic: 1,
        Premium: 3,
        VIP: 10,
      };

      const maxBookings = bookingLimits[membership_level];

      console.log(activeBookings);
      // Check if the user has exceeded their booking limit
      if (activeBookings.length >= maxBookings) {
        const messageContainer = document.getElementById("vehicleList");
        messageContainer.innerHTML = ""; // Clear any previous messages

        let message = `You have reached your limit of ${maxBookings} active booking${
          maxBookings > 1 ? "s" : ""
        } for the ${membership_level} plan.`;

        if (membership_level === "VIP") {
          message += " You cannot make additional bookings at this time.";
        } else {
          message += ` Please <a href="membership.html" class="text-primary">upgrade your membership</a> to book more.`;
        }

        messageContainer.innerHTML = `
      <div class="alert alert-warning mt-3" role="alert">
        ${message}
      </div>
    `;

        // Disable the search button and date input fields
        const searchButton = document.getElementById("searchButton");
        const startDateInput = document.getElementById("startDate");
        const endDateInput = document.getElementById("endDate");

        searchButton.disabled = true;
        startDateInput.disabled = true;
        endDateInput.disabled = true;

        return; // Stop further processing
      }
    }
  } catch (error) {
    console.error(error);
    showCustomAlert("An error occurred while loading the page.");
  }

  searchButton.addEventListener("click", async () => {
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

    try {
      // Get the user ID from the token in localStorage
      const token = localStorage.getItem("token"); // Replace with actual token key
      if (!token) {
        showCustomAlert("User is not authenticated. Please log in.");
        return;
      }

      const userId = decodeToken(token).user_id;
      if (!userId) {
        showCustomAlert("Failed to decode user information.");
        return;
      }

      // Validate the selected date range based on membership level
      const membershipResponse = await fetch(
        `http://localhost:5100/api/v1/user/membership/status?user_id=${encodeURIComponent(
          userId
        )}`
      );

      if (!membershipResponse.ok) {
        throw new Error("Failed to fetch membership status.");
      }

      const { membership_level } = await membershipResponse.json();

      const now = new Date();
      let maxRange;
      switch (membership_level) {
        case "Basic":
          maxRange = new Date(now.setMonth(now.getMonth() + 1)); // 1 month
          break;
        case "Premium":
          maxRange = new Date(now.setMonth(now.getMonth() + 6)); // 6 months
          break;
        case "VIP":
          maxRange = new Date(now.setFullYear(now.getFullYear() + 1)); // 1 year
          break;
        default:
          showCustomAlert("Invalid membership level.");
          return;
      }

      if (end > maxRange) {
        const maxRangeFormatted = maxRange.toLocaleDateString("en-GB", {
          day: "2-digit",
          month: "long",
          year: "numeric",
        });
        showCustomAlert(
          `Your membership level (${membership_level}) allows booking up to ${maxRangeFormatted}. Please adjust your date range.`
        );
        return;
      }

      // Fetch available vehicles
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
          const rentalDurationHours = Math.ceil(
            (end - start) / (1000 * 60 * 60)
          );

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
                        Cleanliness Status: $${
                          vehicle.cleanliness_status
                        } <br />
                        Rental Price per Hour: ${vehicle.rental_price_per_hour.toFixed(
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
    } catch (error) {
      console.error(error);
      showCustomAlert("An error occurred while processing your request.");
    }
  });
});

function decodeToken(token) {
  try {
    const base64Url = token.split(".")[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => `%${("00" + c.charCodeAt(0).toString(16)).slice(-2)}`)
        .join("")
    );
    return JSON.parse(jsonPayload);
  } catch (error) {
    console.error("Error decoding token:", error);
    return null;
  }
}

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
