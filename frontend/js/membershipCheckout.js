document.addEventListener("DOMContentLoaded", () => {
  // Map for membership plans and prices
  const membershipPlans = {
    Premium: {
      price: "$150 / month",
      duration: 1, // Duration in months
      amount: 150.0,
    },
    VIP: {
      price: "$350 / month",
      duration: 12, // Duration in months
      amount: 350.0,
    },
  };

  // Function to get query parameters
  function getQueryParam(param) {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get(param);
  }

  // Helper function to format date as DD MMMM YYYY
  function formatDate(date) {
    return new Intl.DateTimeFormat("en-GB", {
      day: "2-digit",
      month: "long",
      year: "numeric",
    }).format(date);
  }

  // Function to decode JWT token
  function decodeToken(token) {
    try {
      const base64Payload = token.split(".")[1];
      const decodedPayload = atob(base64Payload);
      return JSON.parse(decodedPayload);
    } catch (error) {
      console.error("Error decoding token:", error);
      return null;
    }
  }

  // Retrieve token from localStorage and decode user_id
  const token = localStorage.getItem("token");
  if (!token) {
    showCustomAlert("User is not logged in.");
    window.location.href = "./login.html";
    return;
  }

  const decodedToken = decodeToken(token);
  const userId = decodedToken?.user_id;

  if (!userId) {
    showCustomAlert("Invalid session. Please log in again.");
    window.location.href = "./login.html";
    return;
  }

  console.log("User ID:", userId);

  const email = decodedToken?.email;
  // Retrieve membership plan from query parameters
  const membershipPlan = getQueryParam("plan");

  // Set Membership Plan, Final Price, Start Date, and End Date
  let startDate, endDate, amount;
  if (membershipPlan && membershipPlans[membershipPlan]) {
    const planDetails = membershipPlans[membershipPlan];
    document.getElementById("membershipLevel").textContent = membershipPlan;
    document.getElementById("finalPrice").textContent = planDetails.price;

    // Calculate Membership Start and End Dates
    startDate = new Date();
    endDate = new Date();
    endDate.setMonth(startDate.getMonth() + planDetails.duration);
    amount = planDetails.amount;

    document.getElementById("membershipStartDate").textContent =
      formatDate(startDate);
    document.getElementById("membershipEndDate").textContent =
      formatDate(endDate);
  } else {
    // Default values for Basic plan
    document.getElementById("membershipLevel").textContent = "Basic";
    document.getElementById("finalPrice").textContent = "$0.00";
    document.getElementById("membershipStartDate").textContent = "N/A";
    document.getElementById("membershipEndDate").textContent = "N/A";
  }

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

  function displayCardIcon(cardNumber) {
    const cardIconContainer = document.getElementById("cardIcon");

    // Remove spaces from the card number
    const cleanCardNumber = cardNumber.replace(/\s+/g, "");

    // Regex Patterns for Card Types
    const visaRegex = /^4/; // Visa card starts with 4
    const mastercardRegex = /^5[1-5]/; // MasterCard starts with 51-55
    const amexRegex = /^3[47]/; // American Express starts with 34 or 37
    const discoverRegex = /^6/; // Discover cards start with 6

    // Clear the existing icon
    cardIconContainer.innerHTML = "";

    // Determine the card type and set the image
    let cardImage = "../img/default.png"; // Default image

    if (visaRegex.test(cleanCardNumber)) {
      cardImage = "../img/visa.png";
    } else if (mastercardRegex.test(cleanCardNumber)) {
      cardImage = "../img/mastercard.png";
    } else if (amexRegex.test(cleanCardNumber)) {
      cardImage = "../img/amex.png";
    } else if (discoverRegex.test(cleanCardNumber)) {
      cardImage = "../img/discover.png";
    }

    // Display the image
    cardIconContainer.innerHTML = `<img src="${cardImage}" alt="Card Icon" style="width: 50px; height: auto;">`;
  }

  document.getElementById("cardNumber").addEventListener("input", (event) => {
    const cardNumber = event.target.value.replace(/\s+/g, "").trim();
    displayCardIcon(cardNumber);
    console.log("Card Number Entered: ", cardNumber);

    // Validate the card number using Luhn Algorithm
    if (isValidCardNumber(cardNumber)) {
      addCardNumberTick(); // Add tick icon if valid
    } else {
      removeCardNumberTick(); // Remove tick icon if invalid
    }
  });

  // Handle payment submission
  document.getElementById("payButton").addEventListener("click", () => {
    const selectedPaymentMethod = document.querySelector(
      'input[name="paymentMethod"]:checked'
    ).value;

    if (selectedPaymentMethod === "Card") {
      const cardName = document.getElementById("cardName").value.trim();
      const cardNumber = document.getElementById("cardNumber").value.trim();
      const expiryDate = document.getElementById("expiryDate").value.trim();
      const cvv = document.getElementById("cvv").value.trim();

      // Validate all card details using consolidated function
      if (!validateCreditCard(cardName, cardNumber, expiryDate, cvv)) {
        return; // Stop if validation fails
      }
    }

    // Prepare payment payload
    const payload = {
      user_id: userId,
      membership_level: membershipPlan || "Basic",
      amount: amount || 0.0,
      payment_method: selectedPaymentMethod,
      start_date: startDate.toISOString().split("T")[0],
      end_date: endDate.toISOString().split("T")[0],
      email: email,
    };

    // Send payment data to server
    fetch("http://localhost:5200/api/v1/membership/payment", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Payment processing failed");
        }
        return response.json();
      })
      .then((data) => {
        showCustomAlert(
          "Payment processed successfully!",
          "../membership.html"
        );
        console.log("Payment response:", data);
      })
      .catch((error) => {
        console.error("Payment error:", error);
        showCustomAlert("Payment failed. Please try again.");
      });
  });
});

// Validate credit card number using the Luhn algorithm
function validateCreditCard(cardName, cardNumber, expiryDate, cvv) {
  console.log(cardName, cardNumber, expiryDate, cvv);
  // Validate Cardholder Name
  if (!cardName || cardName.trim() === "") {
    showCustomAlert("Cardholder name is required.");
    return false; // Stop further validation if cardholder name is missing
  }

  // Validate Card Number using Luhn Algorithm
  const digits = cardNumber.replace(/\D/g, "").split("").reverse();
  const checksum = digits.reduce((sum, digit, idx) => {
    digit = parseInt(digit);
    if (idx % 2 === 1) {
      digit *= 2;
      if (digit > 9) digit -= 9;
    }
    return sum + digit;
  }, 0);

  if (checksum % 10 !== 0) {
    showCustomAlert("Invalid credit card number. Please try again.");
    return false; // Stop further validation if card number is invalid
  }

  // Validate Expiry Date
  const [month, year] = expiryDate.split("/");
  const currentDate = new Date();
  const currentMonth = currentDate.getMonth() + 1; // Months are 0-indexed
  const currentYear = currentDate.getFullYear() % 100; // Get last two digits of the year

  if (
    !month ||
    !year ||
    isNaN(month) ||
    isNaN(year) ||
    month < 1 ||
    month > 12 ||
    year < currentYear ||
    (year === currentYear && month < currentMonth)
  ) {
    showCustomAlert("Invalid expiry date. Please use the MM/YY format.");
    return false; // Stop further validation if expiry date is invalid
  }

  // Validate CVV
  if (!cvv || !/^\d{3,4}$/.test(cvv)) {
    showCustomAlert("Invalid CVV. Please enter a valid 3 or 4-digit CVV.");
    return false; // Stop further validation if CVV is invalid
  }

  // If all validations pass
  return true;
}

// Function to add the tick icon
function addCardNumberTick() {
  const tickIcon = document.getElementById("cardNumberTick");
  tickIcon.classList.remove("d-none"); // Show the tick icon
  tickIcon.classList.remove("fa-times-circle", "text-danger"); // Remove cross icon styling if present
  tickIcon.classList.add("fa-check-circle", "text-success"); // Add tick icon styling
}

// Function to add the cross icon
function removeCardNumberTick() {
  const tickIcon = document.getElementById("cardNumberTick");
  tickIcon.classList.remove("d-none"); // Ensure the icon is visible
  tickIcon.classList.remove("fa-check-circle", "text-success"); // Remove tick icon styling if present
  tickIcon.classList.add("fa-times-circle", "text-danger"); // Add cross icon styling
}

// Luhn Algorithm to validate the credit card number
function isValidCardNumber(cardNumber) {
  console.log("Original Card Number: ", cardNumber);

  // Check if the input contains any alphabetic characters
  if (/[a-zA-Z]/.test(cardNumber)) {
    console.log("Invalid Input: Contains alphabetic characters");
    return false; // Invalid if alphabets are present
  }

  // Remove all non-digit characters from the input
  const digits = cardNumber.replace(/\D/g, ""); // Keep only numeric characters
  console.log("Cleaned Digits: ", digits);

  // Check if the length of the card number is valid (typically 13-19 digits)
  if (digits.length < 13 || digits.length > 19) {
    console.log("Invalid Length: ", digits.length);
    return false; // Invalid length
  }
  console.log("Valid Length: ", digits.length);

  // Reverse the digits and calculate the checksum
  const reversedDigits = digits.split("").reverse();
  console.log("Reversed Digits: ", reversedDigits);

  const checksum = reversedDigits.reduce((sum, digit, idx) => {
    digit = parseInt(digit, 10); // Convert character to integer
    if (idx % 2 === 1) {
      console.log(`Index ${idx}: Doubling digit ${digit}`);
      digit *= 2; // Double every second digit
      if (digit > 9) {
        digit -= 9; // Subtract 9 if the doubled value is greater than 9
        console.log(
          `Index ${idx}: Subtracted 9 from doubled digit, new value ${digit}`
        );
      }
    }
    console.log(`Index ${idx}: Adding digit ${digit} to sum`);
    return sum + digit; // Accumulate the sum
  }, 0);

  console.log("Final Checksum: ", checksum);

  // Return true if checksum is divisible by 10, false otherwise
  const isValid = checksum % 10 === 0;
  console.log("Card is Valid: ", isValid);

  return isValid;
}
