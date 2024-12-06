// Extract query parameters from the URL
const urlParams = new URLSearchParams(window.location.search);
const bookingDetails = {
  bookingId: urlParams.get("bookingId"),
  vehicleId: urlParams.get("vehicleId"),
  startDate: urlParams.get("startDate"),
  endDate: urlParams.get("endDate"),
  location: urlParams.get("location"),
  chargeLevel: urlParams.get("chargeLevel"),
  cleanlinessStatus: urlParams.get("cleanlinessStatus"),
  rentalPricePerHour: urlParams.get("rentalPricePerHour"),
  totalDuration: urlParams.get("totalDuration"),
  totalPrice: urlParams.get("totalPrice"),
};

// Populate booking details in the form
document.getElementById("bookingId").textContent = bookingDetails.bookingId;
document.getElementById("vehicleModel").textContent = bookingDetails.vehicleId;
document.getElementById("vehicleLocation").textContent =
  bookingDetails.location;
document.getElementById(
  "chargeLevel"
).textContent = `${bookingDetails.chargeLevel}%`;
document.getElementById("cleanlinessStatus").textContent =
  bookingDetails.cleanlinessStatus;

// Handle form submission
document.getElementById("updateForm").addEventListener("submit", (event) => {
  event.preventDefault();

  // Gather new cleanliness and charging level inputs
  const newCleanlinessStatus = document.getElementById(
    "newCleanlinessStatus"
  ).value;
  const newChargeLevel = document.getElementById("newChargeLevel").value;

  console.log(parseInt(bookingDetails.bookingId));

  // Make an API call to update the vehicle details
  fetch("http://localhost:5150/api/v1/vehicle/update", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${localStorage.getItem("token")}`,
    },
    body: JSON.stringify({
      vehicle_id: parseInt(bookingDetails.vehicleId), // Convert to integer
      cleanliness_status: newCleanlinessStatus,
      charge_level: parseInt(newChargeLevel), // Convert to integer
      location: bookingDetails.location,
      booking_id: parseInt(bookingDetails.bookingId, 10),
    }),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("Failed to update vehicle details.");
      }
      return response.json();
    })
    .then(() => {
      showCustomAlert("Have a safe trip!");

      // Redirect to activeBooking.html with updated query parameters
      const redirectParams = new URLSearchParams({
        ...bookingDetails,
        cleanlinessStatus: newCleanlinessStatus,
        chargeLevel: newChargeLevel,
      });
      window.location.href = `activeBooking.html?${redirectParams.toString()}`;
    })
    .catch((error) => {
      console.error(error);
    });
});
