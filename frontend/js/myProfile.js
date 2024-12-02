document.addEventListener("DOMContentLoaded", function () {
  const profileForm = document.getElementById("profileForm");
  const apiUrl = "http://localhost:5100/api/v1/user"; // Base API URL
  const editButton = document.getElementById("editButton");
  const saveButton = document.getElementById("saveButton");

  // Decode the JWT to extract user information
  function getUserIdFromToken() {
    const token = localStorage.getItem("token");
    if (!token) {
      return null;
    }

    try {
      const payload = JSON.parse(atob(token.split(".")[1])); // Decode the token payload
      return payload.user_id; // Assuming the token contains a user_id field
    } catch (error) {
      console.error("Error decoding token:", error);
      return null;
    }
  }

  // Fetch and populate user profile data
  async function fetchProfile() {
    const token = localStorage.getItem("token");
    const userId = getUserIdFromToken();

    if (!token || !userId) {
      showCustomAlert("You are not logged in. Redirecting to login page.");
      window.location.href = "./login.html";
      return;
    }

    try {
      const response = await fetch(`${apiUrl}/profile?user_id=${userId}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to fetch profile data.");
      }

      const data = await response.json();
      document.getElementById("membership").value = data.membership_level;
      document.getElementById("name").value = data.name;
      document.getElementById("email").value = data.email;
      document.getElementById("contact").value = data.contact_number;
      document.getElementById("address").value = data.address;
    } catch (error) {
      console.error("Error fetching profile:", error);
      showCustomAlert(
        "An error occurred while loading your profile. Please try again."
      );
    }
  }

  // Toggle between edit and save mode
  function toggleEditMode(enable) {
    const inputs = ["name", "contact", "address"];
    inputs.forEach((id) => {
      document.getElementById(id).readOnly = !enable;
    });
    editButton.classList.toggle("d-none", enable);
    saveButton.classList.toggle("d-none", !enable);
  }

  // Enable edit mode on button click
  editButton.addEventListener("click", () => {
    toggleEditMode(true);
  });

  // Save updated profile data
  profileForm.addEventListener("submit", async function (e) {
    e.preventDefault();

    const updatedProfile = {
      user_id: 2, // Example user ID
      name: document.getElementById("name").value.trim(),
      contact_number: document.getElementById("contact").value.trim(),
      address: document.getElementById("address").value.trim(),
    };

    try {
      const response = await fetch(`${apiUrl}/profile/update`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify(updatedProfile),
      });

      if (!response.ok) {
        throw new Error("Failed to update profile.");
      }

      showCustomAlert("Profile updated successfully!");
      toggleEditMode(false);
    } catch (error) {
      console.error("Error updating profile:", error);
      showCustomAlert(
        "An error occurred while saving your profile. Please try again."
      );
    }
  });

  // Initialize the page
  fetchProfile();
});
