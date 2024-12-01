document.addEventListener("DOMContentLoaded", () => {
  const modifyForm = document.getElementById("modifyForm");
  modifyForm.addEventListener("submit", (e) => {
    e.preventDefault();

    const bookingId = document.getElementById("bookingId").value;
    const newReturnDate = document.getElementById("newReturnDate").value;

    fetch(`http://localhost:5150/api/v1/vehicle/booking/${bookingId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ return_date: newReturnDate }),
    })
      .then((response) => response.text())
      .then((message) => alert(message))
      .catch((error) => console.error(error));
  });
});
