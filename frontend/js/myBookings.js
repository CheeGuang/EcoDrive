document.addEventListener("DOMContentLoaded", () => {
  const activeBookingList = document.getElementById("activeBookingList");
  const pastBookingList = document.getElementById("pastBookingList");

  // Retrieve and decode token from localStorage
  const token = localStorage.getItem("token");
  if (!token) {
    showCustomAlert("You must be logged in to view your bookings.");
    window.location.href = "./login.html"; // Redirect to login page if no token is found
    return;
  }

  // Decode the token (assuming it's a base64 encoded JWT)
  function decodeToken(token) {
    try {
      const base64Payload = token.split(".")[1]; // Extract payload
      const decodedPayload = atob(base64Payload); // Decode Base64
      return JSON.parse(decodedPayload); // Parse JSON
    } catch (error) {
      console.error("Invalid token:", error);
      return null;
    }
  }

  const decodedToken = decodeToken(token);
  if (!decodedToken || !decodedToken.user_id) {
    showCustomAlert("Invalid session. Please log in again.");
    window.location.href = "./login.html";
    return;
  }

  const userId = decodedToken.user_id;

  // Fetch bookings for the user
  fetch(`http://localhost:5150/api/v1/vehicle/booking/user/${userId}`)
    .then((response) => {
      if (!response.ok) {
        throw new Error("Failed to fetch bookings.");
      }
      return response.json();
    })
    .then((bookings) => {
      if (!bookings || bookings.length === 0) {
        activeBookingList.innerHTML =
          "<p>You have no active bookings at the moment. Start renting today!</p>";
        pastBookingList.innerHTML =
          "<p>You have no past bookings at the moment.</p>";
        return;
      }

      const now = new Date();

      // Separate active and past bookings
      const activeBookings = bookings.filter(
        (booking) => new Date(booking.return_date) > now
      );
      const pastBookings = bookings.filter(
        (booking) => new Date(booking.return_date) <= now
      );

      // Render Active Bookings
      if (activeBookings.length === 0) {
        activeBookingList.innerHTML =
          "<p>You have no active bookings at the moment. Start renting today!</p>";
      } else {
        activeBookings.forEach((booking) => {
          const startDate = new Date(booking.booking_date);
          const endDate = new Date(booking.return_date);

          const durationInHours = Math.ceil(
            (endDate - startDate) / (1000 * 60 * 60)
          );
          const totalPrice = Math.ceil(
            durationInHours * booking.rental_price_per_hour
          );

          const bookingCard = `
              <div class="card mb-3">
                <div class="card-body">
                  <h5 class="card-title">${booking.model}</h5>
                  <p class="card-text">
                    Booking ID: ${booking.booking_id} <br />
                    Start Date: ${startDate.toLocaleString()} <br />
                    End Date: ${endDate.toLocaleString()} <br />
                    Location: ${booking.location} <br />
                    Charge Level: ${booking.charge_level}% <br />
                    Rental Price per Hour: $${booking.rental_price_per_hour.toFixed(
                      2
                    )} <br />
                    Total Duration: ${durationInHours} hour(s) <br />
                    Total Price: $${totalPrice.toFixed(2)}
                  </p>
                  <button class="btn btn-warning" onclick="openModifyModal(${
                    booking.booking_id
                  }, '${booking.booking_date}', '${booking.return_date}', ${
            booking.rental_price_per_hour
          })">
                    Modify Booking
                  </button>
                  <button class="btn btn-danger" onclick="cancelBooking(${
                    booking.booking_id
                  })">
                    Cancel Booking
                  </button>
                </div>
              </div>`;
          activeBookingList.innerHTML += bookingCard;
        });
      }

      // Render Past Bookings
      if (pastBookings.length === 0) {
        pastBookingList.innerHTML =
          "<p>You have no past bookings at the moment.</p>";
      } else {
        pastBookings.forEach((booking) => {
          const startDate = new Date(booking.booking_date);
          const endDate = new Date(booking.return_date);

          const durationInHours = Math.ceil(
            (endDate - startDate) / (1000 * 60 * 60)
          );
          const totalPrice = Math.ceil(
            durationInHours * booking.rental_price_per_hour
          );

          const bookingCard = `
              <div class="card mb-3">
                <div class="card-body">
                  <h5 class="card-title">${booking.model}</h5>
                  <p class="card-text">
                    Booking ID: ${booking.booking_id} <br />
                    Start Date: ${startDate.toLocaleString()} <br />
                    End Date: ${endDate.toLocaleString()} <br />
                    Location: ${booking.location} <br />
                    Charge Level: ${booking.charge_level}% <br />
                    Rental Price per Hour: $${booking.rental_price_per_hour.toFixed(
                      2
                    )} <br />
                    Total Duration: ${durationInHours} hour(s) <br />
                    Total Price: $${totalPrice.toFixed(2)}
                  </p>
                </div>
              </div>`;
          pastBookingList.innerHTML += bookingCard;
        });
      }
    })
    .catch((error) => {
      console.error(error);
      activeBookingList.innerHTML =
        "<p>An error occurred while fetching your bookings. Please try again later.</p>";
      pastBookingList.innerHTML =
        "<p>An error occurred while fetching your bookings. Please try again later.</p>";
    });

  // Open the modal for modifying a booking
  window.openModifyModal = (
    bookingId,
    startDate,
    endDate,
    rentalPricePerHour
  ) => {
    document.getElementById("bookingId").value = bookingId;
    document.getElementById("startDate").value = startDate.slice(0, 16); // Format for datetime-local input
    document.getElementById("endDate").value = endDate.slice(0, 16); // Format for datetime-local input
    document.getElementById("rentalPricePerHour").value = rentalPricePerHour;
    const modifyModal = new bootstrap.Modal(
      document.getElementById("modifyBookingModal")
    );
    modifyModal.show();
  };

  // Save changes and call the API
  document.getElementById("saveChangesButton").addEventListener("click", () => {
    const bookingId = document.getElementById("bookingId").value;
    const startDate = document.getElementById("startDate").value;
    const endDate = document.getElementById("endDate").value;
    const rentalPricePerHour = parseFloat(
      document.getElementById("rentalPricePerHour").value
    );

    if (!startDate || !endDate || isNaN(rentalPricePerHour)) {
      showCustomAlert("Please fill in all fields.");
      return;
    }

    const startDateTime = new Date(startDate);
    const endDateTime = new Date(endDate);

    if (startDateTime >= endDateTime) {
      showCustomAlert("End date must be later than start date.");
      return;
    }

    // Calculate the duration in hours and total price
    const durationInHours = Math.ceil(
      (endDateTime - startDateTime) / (1000 * 60 * 60)
    );
    const totalPrice = durationInHours * rentalPricePerHour;

    fetch(`http://localhost:5150/api/v1/vehicle/booking/${bookingId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        start_date_time: startDate,
        end_date_time: endDate,
        total_price: totalPrice,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to modify booking.");
        }
        return response.text();
      })
      .then((message) => {
        showCustomAlert(message);
        window.location.reload(); // Refresh the page
      })
      .catch((error) => {
        console.error(error);
        showCustomAlert("An error occurred while modifying the booking.");
      });
  });
});

// Cancel booking function
function cancelBooking(bookingId) {
  if (!confirm("Are you sure you want to cancel this booking?")) {
    return;
  }

  fetch(`http://localhost:5150/api/v1/vehicle/booking/${bookingId}`, {
    method: "DELETE",
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("Failed to cancel the booking.");
      }
      showCustomAlert("Booking cancelled successfully.");
      window.location.reload(); // Refresh the page to update the booking list
    })
    .catch((error) => {
      console.error(error);
      showCustomAlert("An error occurred while cancelling the booking.");
    });
}
