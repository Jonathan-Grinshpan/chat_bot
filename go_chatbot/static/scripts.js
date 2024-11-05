async function sendMessage() {
  startTypingIndicator();
  const messageInput = document.getElementById("messageInput");
  const message = messageInput.value.trim();
  const chatbox = document.getElementById("chatbox");

  if (message) {
    // Display the user message in the chatbox
    chatbox.value += `You: ${message}\n`;
    chatbox.scrollTop = chatbox.scrollHeight;

    // Send the message to the backend API for chatbot response
    try {
      const response = await fetch("/chatbot", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ question: message }),
      });

      if (!response.ok) {
        throw new Error("Network response was not ok");
      }

      const data = await response.json();

      // Display the bot's response in the chatbox
      chatbox.value += `Bot: ${data.answer}\n`;
      chatbox.scrollTop = chatbox.scrollHeight;
    } catch (error) {
      console.error("Error fetching chatbot response:", error);
      chatbox.value += `Bot: Sorry, there was an error fetching the response.\n`;
      chatbox.scrollTop = chatbox.scrollHeight;
    } finally {
      stopTypingIndicator(); // Stop typing indicator when done
    }

    // Clear the input field
    messageInput.value = "";
  }
}

// Allow sending messages with the Enter key
const messageInput = document.getElementById("messageInput");
messageInput.addEventListener("keypress", function (event) {
  if (event.key === "Enter") {
    event.preventDefault();
    sendMessage();
  }
});

let typingIndicator = "";
let typingInterval;

function startTypingIndicator() {
  typingIndicator = ".";
  typingInterval = setInterval(() => {
    typingIndicator += ".";
    if (typingIndicator.length > 3) {
      typingIndicator = ".";
    }
    document.getElementById("typing-indicator").innerText = typingIndicator;
  }, 500); // Update every 500ms
}

function stopTypingIndicator() {
  clearInterval(typingInterval);
  document.getElementById("typing-indicator").innerText = ""; // Clear the indicator
}
