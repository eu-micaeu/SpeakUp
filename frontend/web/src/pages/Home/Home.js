import { useState } from 'react';
import './Home.css';

function Home() {
  const [isSidebarVisible, setIsSidebarVisible] = useState(true);

  const toggleSidebar = () => {
    setIsSidebarVisible(!isSidebarVisible);
  };

  return (
    <div className={`pageHome ${isSidebarVisible ? 'sidebarVisible' : 'sidebarHidden'}`}>
      <aside className="sidebar">
        <button className="toggleButton" onClick={toggleSidebar}>
          &times;
        </button>
      </aside>
      <main className="mainContent">
        {!isSidebarVisible && (
          <button className="toggleButton" onClick={toggleSidebar}>
            &#9776;
          </button>
        )}
        <h1>SpeakUp</h1>
        <p>O que aprender hoje?</p>
      </main>
    </div>
  );
}

export default Home;
