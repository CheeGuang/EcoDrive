document.addEventListener("DOMContentLoaded", () => {
  const queryParams = new URLSearchParams(window.location.search);

  // Extract query parameters
  const bookingId = queryParams.get("bookingId");
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

  // Populate the fields
  document.getElementById("bookingId").textContent = bookingId || "N/A";
  document.getElementById("newBookingDatetime").textContent =
    newBookingDatetime || "N/A";
  document.getElementById("newReturnDatetime").textContent =
    newReturnDatetime || "N/A";
  document.getElementById("originalBookingDatetime").textContent =
    originalBookingDateTime || "N/A";
  document.getElementById("originalReturnDatetime").textContent =
    originalReturnDateTime || "N/A";
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

  // Handle card number input
  const cardNumberInput = document.getElementById("cardNumber");
  cardNumberInput.addEventListener("input", (event) => {
    const cardNumber = event.target.value.replace(/\s+/g, "").trim();
    displayCardIcon(cardNumber);

    if (isValidCardNumber(cardNumber)) {
      addCardNumberTick();
    } else {
      removeCardNumberTick();
    }
  });

  // Display the card icon based on the card type
  function displayCardIcon(cardNumber) {
    const cardIconContainer = document.getElementById("cardIcon");
    const cleanCardNumber = cardNumber.replace(/\s+/g, "");

    const visaRegex = /^4/;
    const mastercardRegex = /^5[1-5]/;
    const amexRegex = /^3[47]/;
    const discoverRegex = /^6/;

    let cardImage = "./img/default.png"; // Default placeholder image

    if (visaRegex.test(cleanCardNumber)) {
      cardImage = "./img/visa.png";
    } else if (mastercardRegex.test(cleanCardNumber)) {
      cardImage = "./img/mastercard.png";
    } else if (amexRegex.test(cleanCardNumber)) {
      cardImage = "./img/amex.png";
    } else if (discoverRegex.test(cleanCardNumber)) {
      cardImage = "./img/discover.png";
    }

    cardIconContainer.innerHTML = `<img src="${cardImage}" alt="Card Icon" style="width: 50px; height: auto;">`;
  }

  // Function to add the tick icon
  function addCardNumberTick() {
    const tickIcon = document.getElementById("cardNumberTick");
    tickIcon.classList.remove("d-none");
    tickIcon.classList.remove("fa-times-circle", "text-danger");
    tickIcon.classList.add("fa-check-circle", "text-success");
  }

  // Function to add the cross icon
  function removeCardNumberTick() {
    const tickIcon = document.getElementById("cardNumberTick");
    tickIcon.classList.remove("d-none");
    tickIcon.classList.remove("fa-check-circle", "text-success");
    tickIcon.classList.add("fa-times-circle", "text-danger");
  }

  // Validate credit card number using the Luhn algorithm
  function isValidCardNumber(cardNumber) {
    const digits = cardNumber.replace(/\D/g, "").split("").reverse();
    const checksum = digits.reduce((sum, digit, idx) => {
      digit = parseInt(digit, 10);
      if (idx % 2 === 1) {
        digit *= 2;
        if (digit > 9) digit -= 9;
      }
      return sum + digit;
    }, 0);
    return checksum % 10 === 0;
  }

  // Payment button click
  document.getElementById("payButton").addEventListener("click", () => {
    const selectedPaymentMethod = document.querySelector(
      'input[name="paymentMethod"]:checked'
    ).value;

    if (selectedPaymentMethod === "Card") {
      const cardName = document.getElementById("cardName").value.trim();
      const cardNumber = document.getElementById("cardNumber").value.trim();
      const expiryDate = document.getElementById("expiryDate").value.trim();
      const cvv = document.getElementById("cvv").value.trim();

      if (!validateCreditCard(cardName, cardNumber, expiryDate, cvv)) {
        return; // Stop if validation fails
      }
    }

    // Construct payload for API call
    const payload = {
      start_date_time: newBookingDatetime,
      end_date_time: newReturnDatetime,
      total_price: parseFloat(extraAmountToPay || 0),
    };

    // Make API call to update booking
    fetch(`http://localhost:5150/api/v1/vehicle/booking/${bookingId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to modify booking.");
        }

        // Construct query parameters for the modifyConfirmation page
        const queryParams = new URLSearchParams({
          booking_id: bookingId,
          newBookingDatetime: newBookingDatetime,
          newReturnDatetime: newReturnDatetime,
          originalBookingDateTime: originalBookingDateTime,
          originalReturnDateTime: originalReturnDateTime,
          additionalHours: additionalHours,
          extraAmountToPay: extraAmountToPay,
          rentalPricePerHour: rentalPricePerHour,
          location,
          chargeLevel: chargeLevel,
          totalDuration: totalDuration,
          totalPrice: totalPrice,
        });

        // Show success alert and redirect to modifyConfirmation
        showCustomAlert(
          "Your payment was successful! Redirecting to the next step...",
          `../modifyConfirmation.html?${queryParams.toString()}`
        );
      })
      .catch((error) => {
        console.error("Error updating booking:", error);
        alert("Failed to update booking. Please try again.");
      });
  });

  // Validate full credit card details
  function validateCreditCard(cardName, cardNumber, expiryDate, cvv) {
    if (!cardName || cardName.trim() === "") {
      alert("Cardholder name is required.");
      return false;
    }

    if (!isValidCardNumber(cardNumber)) {
      alert("Invalid credit card number.");
      return false;
    }

    const [month, year] = expiryDate.split("/");
    const currentYear = new Date().getFullYear() % 100;
    const currentMonth = new Date().getMonth() + 1;

    if (
      isNaN(month) ||
      isNaN(year) ||
      month < 1 ||
      month > 12 ||
      year < currentYear ||
      (year == currentYear && month < currentMonth)
    ) {
      alert("Invalid expiry date.");
      return false;
    }

    if (!/^\d{3,4}$/.test(cvv)) {
      alert("Invalid CVV.");
      return false;
    }

    return true;
  }
});
