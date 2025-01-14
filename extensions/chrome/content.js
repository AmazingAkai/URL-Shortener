chrome.runtime.sendMessage(
  { action: "shortenUrl", longUrl: window.location.href },
  (response) => {
    if (response.error) {
      alert(response.error);
    } else {
      try {
        navigator.clipboard.writeText(response.shortUrl);
        alert("Copied the short URL to clipboard!");
      } catch (err) {
        alert("Failed to copy URL: " + err);
      }
    }
  }
);
