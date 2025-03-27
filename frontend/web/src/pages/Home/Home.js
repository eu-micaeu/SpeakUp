import { useState, useEffect } from 'react';
import './Home.css';

import { getChatsByUserId } from '../../utils/api';

function Home() {
  const [isSidebarVisible, setIsSidebarVisible] = useState(false);
  const [chats, setChats] = useState([]);

  const toggleSidebar = () => {
    setIsSidebarVisible(!isSidebarVisible);
  };

  useEffect(() => {
    const userId = 'c06718c9-c99c-4f18-aa99-3111e36a12c4';
    getChatsByUserId(userId).then((response) => {
      setChats(response.chats);  // Acessar a propriedade 'chats'
    }).catch(error => {
      console.error("Erro ao buscar chats:", error);
      setChats([]);  // Para evitar erro se a requisição falhar
    });
  }, []);

  return (
    <div className={`pageHome ${isSidebarVisible ? 'sidebarVisible' : 'sidebarHidden'}`}>
      <aside className="sidebar">
        <button className="toggleButton" onClick={toggleSidebar}>
          &times;
        </button>
        <ul>
          {chats.map((chat) => (
            <li key={chat.id}>{chat.topic}</li>
          ))}
        </ul>
      </aside>
      {!isSidebarVisible && (
        <button className="toggleButton" onClick={toggleSidebar}>
          &#9776;
        </button>
      )}
      <main className="mainContent">
        <h1>SpeakUp</h1>

      </main>
    </div>
  );
}

export default Home;
