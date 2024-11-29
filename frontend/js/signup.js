document.addEventListener("DOMContentLoaded", function () {
  // Get the Send Code button
  const sendCodeButton = document.getElementById("sendCode");

  // Add a click event listener to the Send Code button
  sendCodeButton.addEventListener("click", function () {
    // Get the email input value
    const emailInput = document.getElementById("email");
    const email = emailInput.value.trim();

    // Check if the email is valid
    if (!validateEmail(email)) {
      alert("Please enter a valid email address.");
      return;
    }

    // Disable the button to prevent multiple clicks
    sendCodeButton.disabled = true;

    // Construct the endpoint dynamically using window.location.origin
    const endpoint = `http://localhost:5000/api/v1/authentication/send-verification`;

    // Send a POST request to the backend API
    fetch(endpoint, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email: email }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to send verification code.");
        }
        return response.json();
      })
      .then((data) => {
        alert(data.message || "Verification code sent successfully!");
      })
      .catch((error) => {
        alert(
          error.message ||
            "An error occurred while sending the verification code."
        );
      })
      .finally(() => {
        // Re-enable the button
        sendCodeButton.disabled = false;
      });
  });

  // Function to validate email format
  function validateEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }
});
