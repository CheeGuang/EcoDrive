document.addEventListener("DOMContentLoaded", () => {
  // Helper function to add token to fetch requests
  async function authenticatedFetch(url, options = {}) {
    const token = localStorage.getItem("token");
    if (!token) {
      showCustomAlert("User is not authenticated. Redirecting to login page.");
      window.location.href = "./login.html";
      throw new Error("User is not authenticated");
    }

    const headers = {
      ...options.headers,
      Authorization: `Bearer ${token}`,
    };

    return fetch(url, { ...options, headers });
  }

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
  authenticatedFetch(
    `http://localhost:5150/api/v1/vehicle/booking/user/${userId}`
  )
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

      console.log(bookings);

      const now = new Date();

      // Separate active and past bookings
      const activeBookings = bookings.filter(
        (booking) =>
          new Date(booking.return_date) > now && booking.status != "Completed"
      );

      const pastBookings = bookings.filter(
        (booking) =>
          new Date(booking.return_date) <= now || booking.status == "Completed"
      );

      // Render Active Bookings
      if (activeBookings.length === 0) {
        activeBookingList.innerHTML =
          "<p>You have no active bookings at the moment. Start renting today!</p>";
      } else {
        // Render Active Bookings
        if (activeBookings.length === 0) {
          activeBookingList.innerHTML =
            "<p>You have no active bookings at the moment. Start renting today!</p>";
        } else {
          // Sort bookings by start date in ascending order
          activeBookings.sort(
            (a, b) => new Date(a.booking_date) - new Date(b.booking_date)
          );

          activeBookings.forEach((booking) => {
            console.debug("Processing active booking:", booking);

            const startDate = new Date(booking.booking_date);
            const endDate = new Date(booking.return_date);

            const durationInHours = Math.ceil(
              (endDate - startDate) / (1000 * 60 * 60)
            );
            const totalPrice = Math.ceil(
              durationInHours * booking.rental_price_per_hour
            );

            const now = new Date();
            const timeDiffInMinutes = Math.ceil(
              (startDate - now) / (1000 * 60)
            );

            console.debug("Current time:", now);
            console.debug("Booking start time:", startDate);
            console.debug("Time difference in minutes:", timeDiffInMinutes);

            const bookingCard = `
            <div class="card mb-3 booking-card shadow">
  <div class="card-body">
    <div class="row">
      <!-- Column 1: Booking Details -->
      <div class="col-md-6">
        <h5 class="card-title text-center">
          <i class="fas fa-car"></i> ${booking.model}
        </h5>
        <p class="card-text">
          <i class="fas fa-id-badge"></i> Booking ID: ${
            booking.booking_id
          } <br />
          <i class="fas fa-car"></i> Vehicle ID: ${booking.vehicle_id} <br />
          <i class="fas fa-calendar-alt"></i> Start Date: ${startDate.toLocaleString(
            "en-US",
            {
              day: "2-digit",
              month: "long",
              year: "numeric",
              hour: "2-digit",
              minute: "2-digit",
              hour12: true,
            }
          )} <br />
          <i class="fas fa-calendar-alt"></i> End Date: ${endDate.toLocaleString(
            "en-US",
            {
              day: "2-digit",
              month: "long",
              year: "numeric",
              hour: "2-digit",
              minute: "2-digit",
              hour12: true,
            }
          )} <br />
          <i class="fas fa-map-marker-alt"></i> Location: ${
            booking.location
          } <br />
          <i class="fas fa-battery-half"></i> Charge Level: ${
            booking.charge_level
          }% <br />
          <i class="fas fa-shower"></i> Cleanliness Status: ${
            booking.cleanliness_status
          } <br />
          <i class="fas fa-dollar-sign"></i> Rental Price per Hour: $${booking.rental_price_per_hour.toFixed(
            2
          )} <br />
          <i class="fas fa-clock"></i> Total Duration: ${durationInHours} hour(s) <br />
          <i class="fas fa-money-check-alt"></i> Total Price: $${totalPrice.toFixed(
            2
          )}
        </p>
        <div class="button-group text-center mt-3">
          ${
            booking.status === "Pending" && timeDiffInMinutes <= 60
              ? `<button class="btn btn-success me-2" onclick="startTrip(${
                  booking.booking_id
                }, ${booking.vehicle_id}, '${startDate.toLocaleString("en-US", {
                  day: "2-digit",
                  month: "long",
                  year: "numeric",
                  hour: "2-digit",
                  minute: "2-digit",
                  hour12: true,
                })}', '${endDate.toLocaleString("en-US", {
                  day: "2-digit",
                  month: "long",
                  year: "numeric",
                  hour: "2-digit",
                  minute: "2-digit",
                  hour12: true,
                })}', '${booking.location}', ${booking.charge_level}, '${
                  booking.cleanliness_status
                }', ${
                  booking.rental_price_per_hour
                }, ${durationInHours}, ${totalPrice})">
                <i class="fas fa-play"></i> Start Trip
              </button>`
              : booking.status === "Active"
              ? `<button class="btn btn-primary me-2" onclick="viewTripDetails(
                  ${booking.booking_id}, 
                  '${booking.vehicle_id}', 
                  '${startDate.toLocaleString("en-US", {
                    day: "2-digit",
                    month: "long",
                    year: "numeric",
                    hour: "2-digit",
                    minute: "2-digit",
                    hour12: true,
                  })}', 
                  '${endDate.toLocaleString("en-US", {
                    day: "2-digit",
                    month: "long",
                    year: "numeric",
                    hour: "2-digit",
                    minute: "2-digit",
                    hour12: true,
                  })}', 
                  '${booking.location}', 
                  ${booking.charge_level}, 
                  '${booking.cleanliness_status}', 
                  ${booking.rental_price_per_hour}, 
                  ${durationInHours}, 
                  ${totalPrice}
                )">
                <i class="fas fa-info-circle"></i> Trip Details
              </button>`
              : ""
          }
          <button class="btn btn-warning me-2" onclick="openModifyModal(${
            booking.booking_id
          }, '${booking.booking_date}', '${booking.return_date}', ${
              booking.rental_price_per_hour
            }, '${booking.vehicle_id}')">
            <i class="fas fa-edit"></i> Modify
          </button>
          <button class="btn btn-danger" onclick="cancelBooking(${
            booking.booking_id
          })">
            <i class="fas fa-trash-alt"></i> Cancel
          </button>
        </div>
      </div>

      <div class="col-md-6">
        <div class="google-map" style="height:100%">
          <iframe
            width="100%"
            height="100%"
            frameborder="0"
            style="border:0"
            src="https://www.google.com/maps?q=${encodeURIComponent(
              booking.location
            )}&output=embed"
            allowfullscreen
          ></iframe>
        </div>
      </div>
    </div>
  </div>
</div>`;
            activeBookingList.innerHTML += bookingCard;
          });
        }

        // Function to handle Start Trip redirection
        window.startTrip = (
          bookingId,
          vehicleId,
          startDate,
          endDate,
          location,
          chargeLevel,
          cleanlinessStatus,
          rentalPricePerHour,
          totalDuration,
          totalPrice
        ) => {
          console.debug("Start Trip button clicked.");
          console.debug("Booking Details:", {
            bookingId,
            vehicleId,
            startDate,
            endDate,
            location,
            chargeLevel,
            cleanlinessStatus,
            rentalPricePerHour,
            totalDuration,
            totalPrice,
          });

          const queryParams = new URLSearchParams({
            bookingId,
            vehicleId,
            startDate,
            endDate,
            location,
            chargeLevel,
            cleanlinessStatus,
            rentalPricePerHour: rentalPricePerHour.toFixed(2),
            totalDuration,
            totalPrice: totalPrice.toFixed(2),
          }).toString();

          console.debug(
            "Redirecting to startBooking.html with query params:",
            queryParams
          );
          window.location.href = `startBooking.html?${queryParams}`;
        };

        // Updated Function to view trip details
        window.viewTripDetails = (
          bookingId,
          vehicleId,
          startDate,
          endDate,
          location,
          chargeLevel,
          cleanlinessStatus,
          rentalPricePerHour,
          totalDuration,
          totalPrice
        ) => {
          console.debug(
            "View Trip Details button clicked for booking ID:",
            bookingId
          );

          const queryParams = new URLSearchParams({
            bookingId,
            vehicleId,
            startDate,
            endDate,
            location,
            chargeLevel,
            cleanlinessStatus,
            rentalPricePerHour: rentalPricePerHour.toFixed(2),
            totalDuration,
            totalPrice: totalPrice.toFixed(2),
          }).toString();

          console.debug(
            "Redirecting to activeBooking.html with query params:",
            queryParams
          );
          window.location.href = `activeBooking.html?${queryParams}`;
        };
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
            <div class="card mb-3 booking-card shadow">
              <div class="card-body">
                <h5 class="card-title text-center">
                  <i class="fas fa-car"></i> ${booking.model}
                </h5>
                <p class="card-text">
                  <i class="fas fa-id-badge"></i> Booking ID: ${
                    booking.booking_id
                  } <br />
                  <i class="fas fa-car"></i> Vehicle ID: ${
                    booking.vehicle_id
                  } <br />
                  <i class="fas fa-calendar-alt"></i> Start Date: ${startDate.toLocaleString(
                    "en-US",
                    {
                      day: "2-digit",
                      month: "long",
                      year: "numeric",
                      hour: "2-digit",
                      minute: "2-digit",
                      hour12: true,
                    }
                  )} <br />
                  <i class="fas fa-calendar-alt"></i> End Date: ${endDate.toLocaleString(
                    "en-US",
                    {
                      day: "2-digit",
                      month: "long",
                      year: "numeric",
                      hour: "2-digit",
                      minute: "2-digit",
                      hour12: true,
                    }
                  )} <br />
                  <i class="fas fa-map-marker-alt"></i> Location: ${
                    booking.location
                  } <br />
                  <i class="fas fa-battery-half"></i> Charge Level: ${
                    booking.charge_level
                  }% <br />
                  <i class="fas fa-shower"></i> Cleanliness Status: ${
                    booking.cleanliness_status
                  } <br />
                  <i class="fas fa-dollar-sign"></i> Rental Price per Hour: $${booking.rental_price_per_hour.toFixed(
                    2
                  )} <br />
                  <i class="fas fa-clock"></i> Total Duration: ${durationInHours} hour(s) <br />
                  <i class="fas fa-money-check-alt"></i> Total Price: $${totalPrice.toFixed(
                    2
                  )}
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

  const unavailableTimeslots = document.getElementById("unavailableTimeslots");

  // Function to fetch unavailable timeslots for a specific vehicle
  function fetchUnavailableTimeslots(
    vehicleId,
    originalBookingDatetime,
    originalReturnDatetime
  ) {
    console.log("Fetching unavailable timeslots for vehicle ID:", vehicleId);

    authenticatedFetch(
      `http://localhost:5150/api/v1/vehicle/booking/vehicle/${vehicleId}`
    )
      .then((response) => {
        console.log("Received response status:", response.status);
        if (!response.ok) {
          throw new Error("Failed to fetch unavailable timeslots.");
        }
        return response.json();
      })
      .then((timeslots) => {
        console.log("Fetched timeslots:", timeslots);

        const now = new Date();
        const originalBookingDateTimeObj = new Date(originalBookingDatetime);
        const originalReturnDateTimeObj = new Date(originalReturnDatetime);

        // Filter timeslots: exclude original booking timeslot and past timeslots
        const filteredTimeslots = timeslots.filter((slot) => {
          const existingStart = new Date(slot.booking_date);
          const existingEnd = new Date(slot.return_date);

          // Exclude original booking timeslot
          const isOriginalTimeslot =
            existingStart.getTime() === originalBookingDateTimeObj.getTime() &&
            existingEnd.getTime() === originalReturnDateTimeObj.getTime();

          // Include only future timeslots that are not the original
          return !isOriginalTimeslot && existingEnd > now;
        });

        console.log("Filtered timeslots:", filteredTimeslots);

        if (!filteredTimeslots || filteredTimeslots.length === 0) {
          console.log("No unavailable timeslots after now.");
          unavailableTimeslots.innerHTML =
            "<p>All timeslots are currently available.</p>";
          return;
        }

        // Helper function to format dates
        function formatDate(date) {
          const options = { day: "2-digit", month: "long", year: "numeric" };
          return new Date(date).toLocaleDateString("en-GB", options);
        }

        // Helper function to format time in 12-hour clock
        function formatTime(date) {
          const options = { hour: "2-digit", minute: "2-digit", hour12: true };
          return new Date(date).toLocaleTimeString("en-GB", options);
        }

        // Render unavailable timeslots
        unavailableTimeslots.innerHTML = filteredTimeslots
          .map(
            (slot) =>
              `<p>Start: ${formatDate(slot.booking_date)} ${formatTime(
                slot.booking_date
              )}, 
       End: ${formatDate(slot.return_date)} ${formatTime(slot.return_date)}</p>`
          )
          .join("");

        console.log("Rendered unavailable timeslots to the DOM.");
      })
      .catch((error) => {
        console.error("Error fetching unavailable timeslots:", error);
        unavailableTimeslots.innerHTML =
          "<p>An error occurred while fetching unavailable timeslots.</p>";
      });
  }

  // Open the modal for modifying a booking
  window.openModifyModal = (
    bookingId,
    startDate,
    endDate,
    rentalPricePerHour,
    vehicleId
  ) => {
    console.log(
      "openModifyModal called with:",
      bookingId,
      startDate,
      endDate,
      rentalPricePerHour,
      vehicleId
    );

    // Parse the booking card details dynamically
    const bookingCard = document
      .querySelector(`button[onclick*="openModifyModal(${bookingId}"]`)
      .closest(".card-body");

    if (!bookingCard) {
      console.error(`Booking card for ID ${bookingId} not found.`);
      return;
    }

    console.log("Booking card found:", bookingCard);

    const location =
      bookingCard
        .querySelector(".card-text")
        .textContent.match(/Location: (.+)/)?.[1]
        .trim() || "N/A";
    console.log("Location:", location);

    const chargeLevel =
      bookingCard
        .querySelector(".card-text")
        .textContent.match(/Charge Level: (.+)%/)?.[1]
        .trim() || "0";
    console.log("Charge Level:", chargeLevel);

    const totalDuration =
      bookingCard
        .querySelector(".card-text")
        .textContent.match(/Total Duration: (.+) hour\(s\)/)?.[1]
        .trim() || "0";
    console.log("Total Duration:", totalDuration);

    const totalPrice =
      bookingCard
        .querySelector(".card-text")
        .textContent.match(/Total Price: \$([0-9.]+)/)?.[1]
        .trim() || "0.00";
    console.log("Total Price:", totalPrice);

    // Populate modal fields (these could be hidden inputs for passing to the next screen)
    document.getElementById("bookingId").value = bookingId;
    document.getElementById("vehicleId").value = vehicleId;
    document.getElementById("newBookingDatetime").value = startDate.slice(
      0,
      16
    ); // Format for datetime-local input
    document.getElementById("newReturnDatetime").value = endDate.slice(0, 16); // Format for datetime-local input
    document.getElementById("originalBookingDateTime").value = startDate.slice(
      0,
      16
    );
    document.getElementById("originalReturnDateTime").value = endDate.slice(
      0,
      16
    );
    document.getElementById("rentalPricePerHour").value = rentalPricePerHour;
    document.getElementById("location").value = location;
    document.getElementById("chargeLevel").value = chargeLevel;
    document.getElementById("totalDuration").value = totalDuration;
    document.getElementById("totalPrice").value = totalPrice;

    console.log("Modal fields populated. Booking details prepared.");

    // Fetch unavailable timeslots for the selected vehicle
    console.log("Fetching unavailable timeslots for vehicleId:", vehicleId);
    fetchUnavailableTimeslots(
      vehicleId,
      document.getElementById("originalBookingDateTime").value,
      document.getElementById("originalReturnDateTime").value
    ); // Use vehicle_id directly

    const modifyModal = new bootstrap.Modal(
      document.getElementById("modifyBookingModal")
    );
    modifyModal.show();
    console.log("Modify modal opened.");
  };

  // Save changes and validate against unavailable timeslots
  document.getElementById("saveChangesButton").addEventListener("click", () => {
    const bookingId = document.getElementById("bookingId").value;
    const vehicleId = document.getElementById("vehicleId").value;
    const newBookingDatetime =
      document.getElementById("newBookingDatetime").value;
    const newReturnDatetime =
      document.getElementById("newReturnDatetime").value;
    const rentalPricePerHour = parseFloat(
      document.getElementById("rentalPricePerHour").value
    );
    const location = document.getElementById("location").value;
    const chargeLevel = document.getElementById("chargeLevel").value;
    const totalDuration = document.getElementById("totalDuration").value;
    const totalPrice = document.getElementById("totalPrice").value;

    if (
      !newBookingDatetime ||
      !newReturnDatetime ||
      isNaN(rentalPricePerHour)
    ) {
      showCustomAlert("Please fill in all fields.");
      return;
    }

    const newBookingDateTimeObj = new Date(newBookingDatetime);
    const newReturnDateTimeObj = new Date(newReturnDatetime);
    const today = new Date();

    // Validation 1: Check if date-times are in the future
    if (newBookingDateTimeObj < today || newReturnDateTimeObj < today) {
      showCustomAlert("Start and end date-times must be in the future.");
      return;
    }

    // Validation 2: Check if end date is after start date
    if (newBookingDateTimeObj >= newReturnDateTimeObj) {
      showCustomAlert("End date-time must be later than start date-time.");
      return;
    }

    const originalBookingDateTime = document.getElementById(
      "originalBookingDateTime"
    ).value;
    const originalReturnDateTime = document.getElementById(
      "originalReturnDateTime"
    ).value;

    const originalBookingDateTimeObj = new Date(originalBookingDateTime);
    const originalReturnDateTimeObj = new Date(originalReturnDateTime);

    // Calculate additional hours and total amount extra to pay
    const originalDuration = Math.ceil(
      (originalReturnDateTimeObj - originalBookingDateTimeObj) /
        (1000 * 60 * 60)
    );
    const newDuration = Math.ceil(
      (newReturnDateTimeObj - newBookingDateTimeObj) / (1000 * 60 * 60)
    );

    const additionalHours = newDuration - originalDuration;
    const extraAmountToPay =
      additionalHours > 0 ? additionalHours * rentalPricePerHour : 0;

    // Validation 3: Check for time slot overlap
    authenticatedFetch(
      `http://localhost:5150/api/v1/vehicle/booking/vehicle/${vehicleId}`
    )
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to validate timeslots.");
        }
        return response.json();
      })
      .then((timeslots) => {
        const isOverlap = timeslots.some((slot) => {
          const existingStart = new Date(slot.booking_date);
          const existingEnd = new Date(slot.return_date);

          // Skip the original booking timeslot from the overlap check
          if (
            existingStart.getTime() === originalBookingDateTimeObj.getTime() &&
            existingEnd.getTime() === originalReturnDateTimeObj.getTime()
          ) {
            return false;
          }

          // Check for overlap with other bookings
          return (
            (newBookingDateTimeObj >= existingStart &&
              newBookingDateTimeObj < existingEnd) ||
            (newReturnDateTimeObj > existingStart &&
              newReturnDateTimeObj <= existingEnd) ||
            (newBookingDateTimeObj <= existingStart &&
              newReturnDateTimeObj >= existingEnd)
          );
        });

        if (isOverlap) {
          showCustomAlert(
            "The selected timeslot overlaps with an existing booking. Please choose a different time."
          );
          return;
        }

        // All validations passed, proceed to redirect
        const queryParams = new URLSearchParams({
          bookingId,
          newBookingDatetime,
          newReturnDatetime,
          originalBookingDateTime: originalBookingDateTime || "",
          originalReturnDateTime: originalReturnDateTime || "",
          additionalHours: additionalHours.toFixed(2),
          extraAmountToPay: extraAmountToPay.toFixed(2),
          rentalPricePerHour,
          location,
          chargeLevel,
          totalDuration,
          totalPrice,
        }).toString();

        window.location.href = `modifyCheckout.html?${queryParams}`;
      })
      .catch((error) => {
        console.error(error);
        showCustomAlert("An error occurred while validating the timeslot.");
      });
  });
});

// Cancel booking function
function cancelBooking(bookingId) {
  if (!confirm("Are you sure you want to cancel this booking?")) {
    return;
  }

  // Get the token from localStorage
  const token = localStorage.getItem("token");
  if (!token) {
    showCustomAlert("Authentication token is missing.");
    return;
  }

  // Decode the token to extract the email
  const tokenPayload = JSON.parse(atob(token.split(".")[1]));
  const email = tokenPayload.email; // Assuming the email field exists in the token payload
  if (!email) {
    showCustomAlert("Email information is missing in the token.");
    return;
  }

  // Make the DELETE request with the email in the body
  fetch(`http://localhost:5150/api/v1/vehicle/booking/${bookingId}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ email }), // Pass the email in the request body
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("Failed to cancel the booking.");
      }
      showCustomAlert(
        "Booking cancelled successfully. A confirmation email has been sent, and the refund will be processed within 3 working days.",
        "../myBookings.html"
      );
    })
    .catch((error) => {
      console.error(error);
      showCustomAlert("An error occurred while cancelling the booking.");
    });
}
