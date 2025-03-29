import React, { useState, useEffect } from "react";
import "./App.css";

function App() {
  const [matches, setMatches] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [stadiumInfo, setStadiumInfo] = useState(null);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);

      // Fetch stadium info
      const stadiumResponse = await fetch("/api/stadium/sammyofer");
      if (!stadiumResponse.ok) {
        throw new Error(`HTTP error! Status: ${stadiumResponse.status}`);
      }
      const stadiumData = await stadiumResponse.json();
      setStadiumInfo(stadiumData);

      // Fetch Sammy Ofer matches
      const matchesResponse = await fetch("/api/fotmob/sammyofer");
      if (!matchesResponse.ok) {
        throw new Error(`HTTP error! Status: ${matchesResponse.status}`);
      }

      const matchesData = await matchesResponse.json();
      console.log("Matches received:", matchesData?.length || 0);

      // Make sure we have an array of matches
      if (!Array.isArray(matchesData)) {
        throw new Error("Expected an array of matches");
      }

      // Filter for only upcoming matches
      const upcomingMatches = getUpcomingMatches(matchesData);
      console.log("Upcoming matches:", upcomingMatches.length);

      setMatches(upcomingMatches);
      setLoading(false);
    } catch (err) {
      console.error("Error fetching data:", err);
      setError(`Failed to fetch data: ${err.message}`);
      setLoading(false);
    }
  };

  // Helper function to get upcoming matches
  const getUpcomingMatches = (allMatches) => {
    const now = new Date();
    return allMatches.filter((match) => {
      if (match.status && match.status.utcTime) {
        const matchDate = new Date(match.status.utcTime);
        return matchDate > now;
      }
      return false;
    });
  };

  const formatDate = (timestamp) => {
    if (!timestamp) return "";
    const date = new Date(timestamp * 1000);
    return date.toLocaleDateString(undefined, {
      weekday: "short",
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const formatTime = (timestamp) => {
    if (!timestamp) return "";
    const date = new Date(timestamp * 1000);
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  };

  const getMatchStatus = (match) => {
    if (!match.status) return "scheduled";
    const isFinished = match.status.finished;
    const isStarted = match.status.started;

    if (isFinished) return "finished";
    if (isStarted) return "live";
    return "scheduled";
  };

  if (loading) return <div className="loading">Loading match data...</div>;
  if (error) return <div className="error">{error}</div>;

  return (
    <div className="App">
      <header className="header">
        <h1>Sammy Ofer Stadium</h1>
        <p className="subtitle">Home of Maccabi Haifa and Hapoel Haifa</p>
      </header>

      {stadiumInfo && (
        <div className="stadium-showcase">
          <div className="stadium-image-container">
            <img
              src={stadiumInfo.imageUrl}
              alt="Sammy Ofer Stadium"
              className="stadium-image"
              loading="eager" // Prioritize loading this image
              onError={(e) => {
                // Fallback if the primary image fails to load
                e.target.src =
                  "https://upload.wikimedia.org/wikipedia/commons/5/50/Sammy_Ofer_Stadium.jpg";
                e.target.onerror = null; // Prevent infinite loop
              }}
            />
          </div>
          <div className="stadium-details">
            <h2>{stadiumInfo.name}</h2>
            <div className="stadium-info-grid">
              <div className="info-item">
                <span className="info-label">Location:</span>
                <span className="info-value">
                  {stadiumInfo.city}, {stadiumInfo.country}
                </span>
              </div>
              <div className="info-item">
                <span className="info-label">Capacity:</span>
                <span className="info-value">
                  {stadiumInfo.capacity.toLocaleString()} spectators
                </span>
              </div>
              <div className="info-item">
                <span className="info-label">Address:</span>
                <span className="info-value">{stadiumInfo.address}</span>
              </div>
              <div className="info-item">
                <span className="info-label">Home Teams:</span>
                <span className="info-value">
                  {stadiumInfo.teams.join(", ")}
                </span>
              </div>
            </div>
            <p className="stadium-description">{stadiumInfo.description}</p>
          </div>
        </div>
      )}

      <h2 className="section-title">Upcoming Matches at Sammy Ofer Stadium</h2>

      {matches.length > 0 ? (
        <div className="matches-grid">
          {matches.map((match) => {
            const status = getMatchStatus(match);

            // Format date and time from the utcTime
            let matchDate = "";
            let matchTime = "";

            if (match.status && match.status.utcTime) {
              const date = new Date(match.status.utcTime);
              matchDate = date.toLocaleDateString(undefined, {
                weekday: "short",
                year: "numeric",
                month: "short",
                day: "numeric",
              });
              matchTime = date.toLocaleTimeString([], {
                hour: "2-digit",
                minute: "2-digit",
              });
            }

            return (
              <div key={match.id} className="match-card upcoming">
                <div className="match-header">
                  <div className="match-date">{matchDate}</div>
                  <div className={`match-status ${status}`}>
                    {status === "live" ? (
                      <span className="live-indicator">LIVE</span>
                    ) : status === "finished" ? (
                      "FINAL"
                    ) : (
                      matchTime
                    )}
                  </div>
                </div>

                {/* Move the upcoming badge after the header */}
                <div className="upcoming-badge">Upcoming</div>

                <div className="match-teams">
                  <div className="team home">
                    <span className="team-name">
                      {match.home?.name || "Home Team"}
                    </span>
                    <span className="team-score">
                      {status !== "scheduled" ? match.home?.score || 0 : ""}
                    </span>
                  </div>

                  <div className="match-separator">
                    {status === "scheduled" ? "vs" : "-"}
                  </div>

                  <div className="team away">
                    <span className="team-score">
                      {status !== "scheduled" ? match.away?.score || 0 : ""}
                    </span>
                    <span className="team-name">
                      {match.away?.name || "Away Team"}
                    </span>
                  </div>
                </div>

                <div className="match-footer">
                  <div className="match-competition">
                    {match.tournament?.name || "Israeli League"}
                  </div>
                  {match.round && (
                    <div className="match-round">Round {match.round}</div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      ) : (
        <div className="no-matches">
          No upcoming matches scheduled at Sammy Ofer Stadium
        </div>
      )}

      <button className="refresh-btn" onClick={fetchData}>
        Refresh Data
      </button>

      <footer className="footer">
        <p>Data provided by Fotmob API</p>
        <p>
          © {new Date().getFullYear()} Sammy Ofer Stadium - Home of Haifa
          Football
        </p>
      </footer>
    </div>
  );
}

export default App;
