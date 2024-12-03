// Dynamically load query parameters into the confirmation page
document.addEventListener("DOMContentLoaded", () => {
  const params = new URLSearchParams(window.location.search);

  // Populate the fields with query parameter values
  document.getElementById("bookingId").textContent = params.get("booking_id");
  document.getElementById("paymentId").textContent = params.get("payment_id");
  document.getElementById("userId").textContent = params.get("user_id");
  document.getElementById("vehicleId").textContent = params.get("vehicle_id");
  document.getElementById("startDate").textContent = new Date(
    params.get("start_date")
  ).toLocaleString("en-US", {
    day: "2-digit",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  });
  document.getElementById("endDate").textContent = new Date(
    params.get("end_date")
  ).toLocaleString("en-US", {
    day: "2-digit",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  });
  document.getElementById("rentalDuration").textContent = `${params.get(
    "rental_duration"
  )} hours`;
  document.getElementById("pricePerHour").textContent = `$${parseFloat(
    params.get("price_per_hour")
  ).toFixed(2)}`;
  document.getElementById("totalPrice").textContent = `$${parseFloat(
    params.get("total_price")
  ).toFixed(2)}`;
  document.getElementById("paymentMethod").textContent =
    params.get("payment_method");
});
