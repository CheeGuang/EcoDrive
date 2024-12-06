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
document.getElementById("startDate").textContent = bookingDetails.startDate;
document.getElementById("endDate").textContent = bookingDetails.endDate;
document.getElementById("rentalPricePerHour").textContent = `$${parseFloat(
  bookingDetails.rentalPricePerHour
).toFixed(2)}`;
document.getElementById(
  "totalDuration"
).textContent = `${bookingDetails.totalDuration} hours`;
document.getElementById("totalPrice").textContent = `$${parseFloat(
  bookingDetails.totalPrice
).toFixed(2)}`;

// Handle End Trip Button
document.getElementById("endBookingButton").addEventListener("click", () => {
  // Redirect to endBooking.html with all query parameters
  const redirectParams = new URLSearchParams(bookingDetails).toString();
  window.location.href = `endBooking.html?${redirectParams}`;
});
