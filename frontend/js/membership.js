document.addEventListener("DOMContentLoaded", () => {
  // Function to decode a JWT
  function decodeToken(token) {
    try {
      const base64Url = token.split(".")[1];
      const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split("")
          .map((c) => `%${("00" + c.charCodeAt(0).toString(16)).slice(-2)}`)
          .join("")
      );
      return JSON.parse(jsonPayload);
    } catch (error) {
      console.error("Error decoding token:", error);
      return null;
    }
  }

  // Function to get the user's membership status
  async function fetchMembershipStatus() {
    const token = localStorage.getItem("token"); // Replace 'authToken' with the actual key if needed
    if (!token) {
      console.error("No token found in localStorage");
      return;
    }

    const decoded = decodeToken(token);
    if (!decoded || !decoded.user_id) {
      console.error("Invalid token or user_id not found");
      return;
    }

    const userID = decoded.user_id;
    const apiURL = `http://localhost:5100/api/v1/user/membership/status?user_id=${encodeURIComponent(
      userID
    )}`;

    try {
      const response = await fetch(apiURL);
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      const data = await response.json();
      updateMembershipUI(data.membership_level);
    } catch (error) {
      console.error("Error fetching membership status:", error);
    }
  }

  function updateMembershipUI(currentMembership) {
    const tiers = ["Basic", "Premium", "VIP"];
    tiers.forEach((tier) => {
      // Find the card by looking for the header text
      const cards = document.querySelectorAll(".card");
      let card = null;

      cards.forEach((currentCard) => {
        const header = currentCard.querySelector("h3");
        if (header && header.textContent.trim() === tier) {
          card = currentCard;
        }
      });

      if (!card) {
        console.error(`Card for ${tier} not found`);
        return;
      }

      const cardBody = card.querySelector(".card-body");

      // Handle the current membership tier
      if (tier === currentMembership) {
        // Add "This is your current plan" message
        if (cardBody) {
          cardBody.insertAdjacentHTML(
            "beforeend",
            '<div class="current-plan fw-bold mt-3">This is your current plan</div>'
          );
        }
        const button = card.querySelector(".btn");
        if (button) {
          button.style.display = "none";
        }
      } else if (tiers.indexOf(tier) < tiers.indexOf(currentMembership)) {
        // For lower tiers, remove "This is your current plan" and add a friendly message
        if (cardBody) {
          const currentPlanMessage = cardBody.querySelector(".current-plan");
          if (currentPlanMessage) {
            currentPlanMessage.remove();
          }

          cardBody.insertAdjacentHTML(
            "beforeend",
            '<div class="text-muted fw-bold" style="margin-top: 32px">You’ve already unlocked this and more with your current plan!</div>'
          );
        }
        const button = card.querySelector(".btn");
        if (button) {
          button.style.display = "none";
        }
      }
    });
  }

  // Function to redirect to membershipCheckout.html with the selected plan
  function redirectToCheckout(plan) {
    const url = `membershipCheckout.html?plan=${encodeURIComponent(plan)}`;
    window.location.href = url;
  }

  // Add event listeners to the "Upgrade now" buttons
  document.querySelector(".btn-warning").addEventListener("click", () => {
    redirectToCheckout("Premium");
  });

  document.querySelector(".btn-danger").addEventListener("click", () => {
    redirectToCheckout("VIP");
  });

  // Fetch membership status on page load
  fetchMembershipStatus();
});
