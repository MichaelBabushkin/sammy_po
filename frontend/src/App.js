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

  // Helper function to generate iCalendar string for a match
  const createCalendarEvent = (match) => {
    if (!match.status || !match.status.utcTime) {
      return null;
    }

    try {
      // Parse match time
      const startTime = new Date(match.status.utcTime);

      // End time (assume matches are 2 hours)
      const endTime = new Date(startTime);
      endTime.setHours(endTime.getHours() + 2);

      // Format dates for iCalendar
      const formatDate = (date) => {
        return date.toISOString().replace(/-|:|\.\d+/g, "");
      };

      const start = formatDate(startTime);
      const end = formatDate(endTime);

      // Create event title and description
      const title = `${match.home?.name || "Home"} vs ${
        match.away?.name || "Away"
      }`;
      const description = `${
        match.tournament?.name || "Israeli League"
      } match at Sammy Ofer Stadium`;
      const location = "Sammy Ofer Stadium, Haifa, Israel";

      // Generate iCalendar format
      const icsData = [
        "BEGIN:VCALENDAR",
        "VERSION:2.0",
        "BEGIN:VEVENT",
        `DTSTART:${start}`,
        `DTEND:${end}`,
        `SUMMARY:${title}`,
        `DESCRIPTION:${description}`,
        `LOCATION:${location}`,
        "END:VEVENT",
        "END:VCALENDAR",
      ].join("\n");

      // Convert to data URI for download
      const dataUri = `data:text/calendar;charset=utf-8,${encodeURIComponent(
        icsData
      )}`;

      return {
        icsData,
        dataUri,
        googleCalendarUrl: createGoogleCalendarUrl(
          match,
          startTime,
          endTime,
          title,
          description,
          location
        ),
      };
    } catch (error) {
      console.error("Error creating calendar event:", error);
      return null;
    }
  };

  // Create Google Calendar URL
  const createGoogleCalendarUrl = (
    match,
    startTime,
    endTime,
    title,
    description,
    location
  ) => {
    const formatGoogleDate = (date) => {
      return date.toISOString().replace(/-|:|\.\d+/g, "");
    };

    const googleParams = new URLSearchParams({
      action: "TEMPLATE",
      text: title,
      dates: `${formatGoogleDate(startTime)}/${formatGoogleDate(endTime)}`,
      details: description,
      location: location,
      sf: true,
      output: "xml",
    });

    return `https://calendar.google.com/calendar/render?${googleParams.toString()}`;
  };

  // Add to Calendar functionality for mobile devices
  const addToCalendar = (match) => {
    const calendarEvent = createCalendarEvent(match);
    if (!calendarEvent) {
      alert("Sorry, calendar data couldn't be created for this match.");
      return;
    }

    // Different handling for mobile vs desktop
    const isMobile = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);

    if (isMobile) {
      // For mobile - open Google Calendar (works on most devices)
      window.open(calendarEvent.googleCalendarUrl, "_blank");
    } else {
      // For desktop - offer download of .ics file
      const link = document.createElement("a");
      link.href = calendarEvent.dataUri;
      link.download = `match-${match.id}.ics`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  };

  // Add functionality to add all matches to calendar
  const addAllMatchesToCalendar = () => {
    if (matches.length === 0) {
      alert("No upcoming matches to add to your calendar");
      return;
    }

    // For mobile devices, we need to handle differently than desktop
    const isMobile = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);

    if (isMobile) {
      // On mobile, we'll create a calendar event for the first match
      // and then let the user repeat the process for other matches
      const calendarEvent = createCalendarEvent(matches[0]);
      if (calendarEvent) {
        // Open the Google Calendar URL for the first match
        window.open(calendarEvent.googleCalendarUrl, "_blank");

        if (matches.length > 1) {
          setTimeout(() => {
            alert(
              `Added first match to calendar. There are ${
                matches.length - 1
              } more matches. Please repeat for each match.`
            );
          }, 1000);
        }
      }
    } else {
      // On desktop, create a combined iCalendar file with all events
      const combinedICS = createCombinedCalendarEvents(matches);
      if (combinedICS) {
        // Download the combined .ics file
        const link = document.createElement("a");
        link.href = combinedICS;
        link.download = `sammy-ofer-matches.ics`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      } else {
        alert("Failed to create calendar events");
      }
    }
  };

  // Create a combined ICS file with all matches
  const createCombinedCalendarEvents = (matchList) => {
    if (!matchList || matchList.length === 0) return null;

    try {
      // Start building the iCalendar content
      let icsContent = [
        "BEGIN:VCALENDAR",
        "VERSION:2.0",
        "PRODID:-//SammyOfer//UpcomingMatches//EN",
      ];

      // Add each match as an event
      for (const match of matchList) {
        if (!match.status || !match.status.utcTime) continue;

        // Parse match time
        const startTime = new Date(match.status.utcTime);

        // End time (assume matches are 2 hours)
        const endTime = new Date(startTime);
        endTime.setHours(endTime.getHours() + 2);

        // Format dates for iCalendar
        const formatDate = (date) =>
          date.toISOString().replace(/-|:|\.\d+/g, "");
        const start = formatDate(startTime);
        const end = formatDate(endTime);

        // Create event title and description
        const title = `${match.home?.name || "Home"} vs ${
          match.away?.name || "Away"
        }`;
        const description = `${
          match.tournament?.name || "Israeli League"
        } match at Sammy Ofer Stadium`;
        const location = "Sammy Ofer Stadium, Haifa, Israel";

        // Add this event to the calendar
        icsContent.push(
          "BEGIN:VEVENT",
          `DTSTART:${start}`,
          `DTEND:${end}`,
          `SUMMARY:${title}`,
          `DESCRIPTION:${description}`,
          `LOCATION:${location}`,
          `UID:match-${match.id}@sammyofer.com`,
          "END:VEVENT"
        );
      }

      // Close the calendar
      icsContent.push("END:VCALENDAR");

      // Convert to data URI for download
      const dataUri = `data:text/calendar;charset=utf-8,${encodeURIComponent(
        icsContent.join("\n")
      )}`;
      return dataUri;
    } catch (error) {
      console.error("Error creating combined calendar events:", error);
      return null;
    }
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
        <>
          <div className="calendar-actions">
            <button className="add-all-btn" onClick={addAllMatchesToCalendar}>
              Add All Matches to Calendar
            </button>
          </div>
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

                  {/* Add calendar button */}
                  <button
                    className="calendar-btn"
                    onClick={(e) => {
                      e.preventDefault();
                      addToCalendar(match);
                    }}
                  >
                    Add to Calendar
                  </button>
                </div>
              );
            })}
          </div>
        </>
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
