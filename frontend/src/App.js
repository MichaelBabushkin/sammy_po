import React, { useState, useEffect } from "react";

function App() {
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchGreeting();
  }, []);

  const fetchGreeting = async () => {
    try {
      const response = await fetch("/api/greeting"); // Changed to relative path
      const data = await response.json();
      setMessage(data.text);
      setLoading(false);
    } catch (err) {
      console.error("Fetch error:", err); // Added error logging
      setError("Failed to fetch greeting");
      setLoading(false);
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="App">
      <h1>Hello from Go + React + Koyeb!</h1>
      <p>This is served from a React application.</p>
      <div>
        <h2>Message from backend:</h2>
        <p>{message}</p>
        <button onClick={fetchGreeting}>Refresh Message</button>
      </div>
    </div>
  );
}

export default App;
