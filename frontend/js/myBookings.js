document.addEventListener("DOMContentLoaded", () => {
  const bookingList = document.getElementById("bookingList");

  // Retrieve and decode token from localStorage
  const token = localStorage.getItem("token");
  if (!token) {
    alert("You must be logged in to view your bookings.");
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
    alert("Invalid session. Please log in again.");
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
      if (bookings.length === 0) {
        bookingList.innerHTML =
          "<p>You have no bookings at the moment. Start renting today!</p>";
        return;
      }

      bookings.forEach((booking) => {
        const startDate = new Date(booking.booking_date);
        const endDate = new Date(booking.return_date);

        // Calculate the duration in hours
        const durationInHours = Math.ceil(
          (endDate - startDate) / (1000 * 60 * 60)
        );

        // Calculate the total price based on price per hour
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
        bookingList.innerHTML += bookingCard;
      });
    })
    .catch((error) => {
      console.error(error);
      bookingList.innerHTML =
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
      alert("Please fill in all fields.");
      return;
    }

    const startDateTime = new Date(startDate);
    const endDateTime = new Date(endDate);

    if (startDateTime >= endDateTime) {
      alert("End date must be later than start date.");
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
        alert(message);
        window.location.reload(); // Refresh the page
      })
      .catch((error) => {
        console.error(error);
        alert("An error occurred while modifying the booking.");
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
      alert("Booking cancelled successfully.");
      window.location.reload(); // Refresh the page to update the booking list
    })
    .catch((error) => {
      console.error(error);
      alert("An error occurred while cancelling the booking.");
    });
}
