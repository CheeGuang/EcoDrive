document.addEventListener("DOMContentLoaded", () => {
  // Helper function to decode a JWT token
  function decodeToken(token) {
    try {
      const base64Payload = token.split(".")[1]; // Extract the payload part
      const decodedPayload = atob(base64Payload); // Decode Base64
      return JSON.parse(decodedPayload); // Parse JSON
    } catch (error) {
      console.error("Invalid token:", error);
      return null;
    }
  }

  // Retrieve token from localStorage
  const token = localStorage.getItem("token");
  if (!token) {
    showCustomAlert("User is not logged in. Redirecting to login page.");
    window.location.href = "./login.html";
    return;
  }

  // Decode the token to get user information
  const decodedToken = decodeToken(token);
  if (!decodedToken || !decodedToken.user_id) {
    showCustomAlert("Invalid session. Please log in again.");
    window.location.href = "./login.html";
    return;
  }

  const modifyForm = document.getElementById("modifyForm");
  modifyForm.addEventListener("submit", (e) => {
    e.preventDefault();

    const bookingId = document.getElementById("bookingId").value;
    const newReturnDate = document.getElementById("newReturnDate").value;

    fetch(`http://localhost:5150/api/v1/vehicle/booking/${bookingId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ return_date: newReturnDate }),
    })
      .then((response) => response.text())
      .then((message) => showCustomAlert(message))
      .catch((error) => console.error(error));
  });
});
