document.addEventListener("DOMContentLoaded", () => {
  const urlParams = new URLSearchParams(window.location.search);
  const vehicleId = urlParams.get("vehicleId");

  const token = localStorage.getItem("token"); // Assuming the token is stored in localStorage
  if (!token) {
    alert("You must be logged in to book a vehicle.");
    window.location.href = "/login.html"; // Redirect to login page if no token is found
    return;
  }

  // Decode the base64 payload of the JWT
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
    alert("Invalid or expired session. Please log in again.");
    window.location.href = "/login.html";
    return;
  }

  const userId = decodedToken.user_id;

  const bookingForm = document.getElementById("bookingForm");
  bookingForm.addEventListener("submit", (e) => {
    e.preventDefault();

    const bookingDate = document.getElementById("bookingDate").value;
    const returnDate = document.getElementById("returnDate").value;

    fetch("http://localhost:5150/api/v1/vehicle/booking", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`, // Optionally send the token for verification
      },
      body: JSON.stringify({
        vehicle_id: parseInt(vehicleId),
        user_id: userId,
        booking_date: bookingDate,
        return_date: returnDate,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to create booking.");
        }
        return response.text();
      })
      .then((message) => alert(message))
      .catch((error) => console.error(error));
  });
});
