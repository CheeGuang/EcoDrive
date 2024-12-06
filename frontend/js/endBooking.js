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

console.debug("Extracted booking details from URL parameters:", bookingDetails);

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

console.debug("Populated booking details in the form.");

// Handle form submission
document
  .getElementById("endBookingForm")
  .addEventListener("submit", (event) => {
    event.preventDefault();

    console.debug("Form submission event triggered.");

    // Gather new cleanliness and charging level inputs
    const newCleanlinessStatus = document.getElementById(
      "newCleanlinessStatus"
    ).value;
    const newChargeLevel = document.getElementById("newChargeLevel").value;

    console.debug("New cleanliness status:", newCleanlinessStatus);
    console.debug("New charge level:", newChargeLevel);

    // Make an API call to end the booking and update vehicle details
    fetch(
      `http://localhost:5150/api/v1/vehicle/booking/end/${bookingDetails.bookingId}`,
      {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({
          vehicle_id: parseInt(bookingDetails.vehicleId, 10), // Convert to integer
          cleanliness_status: newCleanlinessStatus,
          charge_level: parseInt(newChargeLevel, 10), // Convert to integer
          location: bookingDetails.location,
        }),
      }
    )
      .then((response) => {
        console.debug("API response status:", response.status);
        if (!response.ok) {
          throw new Error("Failed to end the booking.");
        }

        showCustomAlert("Thank you for using EcoDrive!");

        // Redirect to memberHome.html
        console.debug("Redirecting to memberHome.html...");
        window.location.href = "memberHome.html";
      })
      .catch((error) => {
        console.error("Error occurred while ending booking:", error);
        showCustomAlert("Error ending the booking. Please try again.");
      });
  });
