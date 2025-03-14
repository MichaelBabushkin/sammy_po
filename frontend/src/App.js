import React, { useState, useEffect } from "react";
import "./App.css";

function App() {
  const [matches, setMatches] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchMatches();
  }, []);

  const fetchMatches = async () => {
    try {
      const response = await fetch("/api/matches");
      const data = await response.json();
      setMatches(data);
      setLoading(false);
    } catch (err) {
      console.error("Fetch error:", err);
      setError("Failed to fetch matches");
      setLoading(false);
    }
  };

  if (loading) return <div className="loading">Loading matches...</div>;
  if (error) return <div className="error">{error}</div>;

  return (
    <div className="App">
      <header className="header">
        <h1>Football Results</h1>
      </header>
      <main className="matches-container">
        {matches.map((match) => (
          <div key={match.id} className="match-card">
            <div className="league">{match.league}</div>
            <div className="date">{match.date}</div>
            <div className="teams">
              <div className="team home">{match.homeTeam}</div>
              <div className="score">
                {match.homeScore} - {match.awayScore}
              </div>
              <div className="team away">{match.awayTeam}</div>
            </div>
          </div>
        ))}
      </main>
      <button className="refresh-btn" onClick={fetchMatches}>
        Refresh Results
      </button>
    </div>
  );
}

export default App;
