document.addEventListener("DOMContentLoaded", () => {
  const queryParams = new URLSearchParams(window.location.search);

  // Extract query parameters
  const bookingId = queryParams.get("booking_id");
  const newBookingDatetime = queryParams.get("newBookingDatetime");
  const newReturnDatetime = queryParams.get("newReturnDatetime");
  const originalBookingDateTime = queryParams.get("originalBookingDateTime");
  const originalReturnDateTime = queryParams.get("originalReturnDateTime");
  const additionalHours = queryParams.get("additionalHours");
  const extraAmountToPay = queryParams.get("extraAmountToPay");
  const rentalPricePerHour = queryParams.get("rentalPricePerHour");
  const location = queryParams.get("location");
  const chargeLevel = queryParams.get("chargeLevel");
  const totalDuration = queryParams.get("totalDuration");
  const totalPrice = queryParams.get("totalPrice");

  // Helper function to format date-time
  function formatDateTime(dateString) {
    if (!dateString) return "N/A";
    const date = new Date(dateString);
    return date.toLocaleString("en-US", {
      day: "2-digit",
      month: "long",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      hour12: true,
    });
  }

  // Populate the fields
  document.getElementById("bookingId").textContent = bookingId || "N/A";
  document.getElementById("newBookingDatetime").textContent =
    formatDateTime(newBookingDatetime);
  document.getElementById("newReturnDatetime").textContent =
    formatDateTime(newReturnDatetime);
  document.getElementById("originalBookingDatetime").textContent =
    formatDateTime(originalBookingDateTime);
  document.getElementById("originalReturnDatetime").textContent =
    formatDateTime(originalReturnDateTime);
  document.getElementById("additionalHours").textContent =
    additionalHours || "0";
  document.getElementById("extraAmountToPay").textContent = `$${parseFloat(
    extraAmountToPay || 0
  ).toFixed(2)}`;
  document.getElementById("rentalPricePerHour").textContent = `$${parseFloat(
    rentalPricePerHour || 0
  ).toFixed(2)}`;
  document.getElementById("location").textContent = location || "N/A";
  document.getElementById("chargeLevel").textContent = `${chargeLevel || 0}%`;
  document.getElementById("totalDuration").textContent = `${
    totalDuration || 0
  } hours`;
  document.getElementById("totalPrice").textContent = `$${parseFloat(
    totalPrice || 0
  ).toFixed(2)}`;
});
