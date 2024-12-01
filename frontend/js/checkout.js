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

  const userId = decodedToken.user_id;
  const email = decodedToken.email;

  // Extract query parameters
  const urlParams = new URLSearchParams(window.location.search);
  const vehicleId = urlParams.get("vehicleId");
  const startDate = urlParams.get("start_date");
  const endDate = urlParams.get("end_date");
  const rentalDuration = urlParams.get("rentalDuration");
  const pricePerHour = urlParams.get("pricePerHour");
  const totalPrice = urlParams.get("totalPrice");

  // Populate static fields
  document.getElementById("startDate").textContent = new Date(
    startDate
  ).toLocaleString();
  document.getElementById("endDate").textContent = new Date(
    endDate
  ).toLocaleString();
  document.getElementById(
    "rentalDuration"
  ).textContent = `${rentalDuration} hours`;
  document.getElementById("pricePerHour").textContent = `$${parseFloat(
    pricePerHour
  ).toFixed(2)}`;
  document.getElementById("totalPrice").textContent = `$${parseFloat(
    totalPrice
  ).toFixed(2)}`;

  // Fetch membership level
  fetch(`http://localhost:5100/api/v1/user/profile?user_id=${userId}`)
    .then((response) => {
      if (!response.ok) {
        throw new Error("Failed to fetch user details");
      }
      return response.json();
    })
    .then((userData) => {
      const membershipLevel = userData.membership_level;
      document.getElementById("membershipLevel").textContent = membershipLevel;

      // Fetch discount using the membership level
      fetch(
        `http://localhost:5200/api/v1/payment/real-time-bill?membership_level=${membershipLevel}&duration_hours=${rentalDuration}&price_per_hour=${pricePerHour}`
      )
        .then((response) => {
          if (!response.ok) {
            throw new Error("Failed to fetch discount");
          }
          return response.json();
        })
        .then((data) => {
          document.getElementById("discount").textContent = `$${parseFloat(
            data.discount
          ).toFixed(2)}`;
          document.getElementById("finalPrice").textContent = `$${parseFloat(
            data.final_price
          ).toFixed(2)}`;
        })
        .catch((error) => {
          console.error("Error fetching discount:", error);
        });
    })
    .catch((error) => {
      console.error("Error fetching user details:", error);
    });

  // Toggle payment method content
  const cardDetails = document.getElementById("cardDetails");
  const paynowQRCode = document.getElementById("paynowQRCode");
  const paymentMethodInputs = document.querySelectorAll(
    'input[name="paymentMethod"]'
  );

  paymentMethodInputs.forEach((input) => {
    input.addEventListener("change", (event) => {
      if (event.target.value === "Card") {
        cardDetails.classList.remove("d-none");
        paynowQRCode.classList.add("d-none");
      } else if (event.target.value === "PayNow") {
        cardDetails.classList.add("d-none");
        paynowQRCode.classList.remove("d-none");
      }
    });
  });

  // Validate credit card number using the Luhn algorithm
  function validateCreditCard(cardNumber) {
    const digits = cardNumber.replace(/\D/g, "").split("").reverse();
    const checksum = digits.reduce((sum, digit, idx) => {
      digit = parseInt(digit);
      if (idx % 2 === 1) {
        digit *= 2;
        if (digit > 9) digit -= 9;
      }
      return sum + digit;
    }, 0);
    return checksum % 10 === 0;
  }

  // Validate expiry date (MM/YY format)
  function validateExpiryDate(expiry) {
    const [month, year] = expiry.split("/").map((val) => parseInt(val, 10));
    if (!month || !year || month < 1 || month > 12) return false;

    const now = new Date();
    const currentYear = now.getFullYear() % 100; // Last two digits of the current year
    const currentMonth = now.getMonth() + 1;

    return (
      year > currentYear || (year === currentYear && month >= currentMonth)
    );
  }

  // Handle payment submission
  document.getElementById("payButton").addEventListener("click", () => {
    const paymentMethod = document.querySelector(
      'input[name="paymentMethod"]:checked'
    ).value;

    if (paymentMethod === "Card") {
      const cardNumber = document.getElementById("cardNumber").value;
      const expiryDate = document.getElementById("expiryDate").value;

      if (!validateCreditCard(cardNumber)) {
        showCustomAlert("Invalid credit card number.");
        return;
      }

      if (!validateExpiryDate(expiryDate)) {
        showCustomAlert("Invalid expiry date.");
        return;
      }
    }

    const payload = {
      user_id: userId,
      vehicle_id: vehicleId,
      start_date: startDate,
      end_date: endDate,
      rental_duration: rentalDuration,
      price_per_hour: pricePerHour,
      total_price: totalPrice,
      payment_method: paymentMethod,
      email: email,
    };

    fetch("http://localhost:5200/api/v1/payment/process", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Payment processing failed");
        }
        return response.json();
      })
      .then((data) => {
        const bookingId = data.booking_id;
        const paymentId = data.payment_id;

        // Construct query parameters to include booking and payment details
        const queryParams = new URLSearchParams({
          booking_id: bookingId,
          payment_id: paymentId,
          user_id: userId,
          vehicle_id: vehicleId,
          start_date: startDate,
          end_date: endDate,
          rental_duration: rentalDuration,
          price_per_hour: pricePerHour,
          total_price: totalPrice,
          payment_method: paymentMethod,
        }).toString();

        // Redirect to confirmation.html with query parameters
        showCustomAlert("Payment successful!");
        window.location.href = `./confirmation.html?${queryParams}`;
      })
      .catch((error) => {
        console.error("Payment error:", error);
        showCustomAlert("Payment failed. Please try again.");
      });
  });
});
