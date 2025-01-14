const API_URL = "http://localhost:8080";

browser.browserAction.onClicked.addListener((tab) => {
  browser.tabs.executeScript({
    file: "content.js",
  });
});

browser.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === "shortenUrl") {
    const shortUrl = Math.random().toString(36).substring(2, 8);

    const form = new FormData();
    form.append("short_url", shortUrl);
    form.append("long_url", request.longUrl);

    fetch(`${API_URL}/urls`, {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: new URLSearchParams(form),
    })
      .then((response) => {
        if (response.status !== 200) {
          console.log(response.status);
          response.text().then(console.log);
          throw new Error("Failed to shorten URL");
        }
        sendResponse({ shortUrl: `${API_URL}/${shortUrl}` });
      })
      .catch((error) => {
        sendResponse({ error: error.message });
      });

    return true;
  }
});
