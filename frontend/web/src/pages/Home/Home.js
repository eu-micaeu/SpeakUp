import { useState, useEffect, useRef } from 'react';
import './Home.css';

import { getChatsByUserId, createChat } from '../../utils/api';

function Home() {
  const [isSidebarVisible, setIsSidebarVisible] = useState(false);
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');
  const [currentChatId, setCurrentChatId] = useState(null);
  const messagesEndRef = useRef(null);

  const toggleSidebar = () => {
    setIsSidebarVisible(!isSidebarVisible);
  };

  useEffect(() => {
    getChatsByUserId().then((response) => {
      setChats(response.chats);
    }).catch(error => {
      console.error("Erro ao buscar chats:", error);
      setChats([]);
    });
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  const handleSendMessage = async () => {
    if (inputMessage.trim() === '') return;
    
    let chatId = currentChatId;
    
    // Se não há chat atual, cria um novo
    if (!chatId) {
      const firstWord = inputMessage.trim().split(' ')[0]; // Pega a primeira palavra
      const topic = firstWord.length > 0 ? firstWord : "Novo Chat"; // Usa como tópico
      
      try {
        const newChat = await createChat(topic);
        chatId = newChat.id;
        setCurrentChatId(chatId);
        setChats(prevChats => [...prevChats, newChat]);
      } catch (error) {
        console.error("Erro ao criar chat:", error);
        return;
      }
    }
    
    // Adiciona a mensagem ao chat atual (recém-criado ou existente)
    const newMessage = {
      id: Date.now(),
      text: inputMessage,
      sender: 'user',
      timestamp: new Date().toISOString(),
      chatId: chatId
    };
    
    setMessages(prevMessages => [...prevMessages, newMessage]);
    setInputMessage('');
    
    // Aqui você pode adicionar a lógica para enviar a mensagem para o backend
    // usando o chatId como referência
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      handleSendMessage();
    }
  };

  return (
    <div className={`pageHome ${isSidebarVisible ? 'sidebarVisible' : 'sidebarHidden'}`}>
      <aside className="sidebar">
        <button className="toggleButton" onClick={toggleSidebar}>
          &times;
        </button>
        <ul>
          {chats?.length > 0 ? (
            chats.map((chat) => (
              <li 
                key={chat.id} 
                className={chat.id === currentChatId ? 'active-chat' : ''}
                onClick={() => setCurrentChatId(chat.id)}
              >
                {chat.topic}
              </li>
            ))
          ) : (
            <ol>Sem chats, por ora</ol>
          )}
        </ul>
      </aside>
      {!isSidebarVisible && (
        <button className="toggleButton" onClick={toggleSidebar}>
          &#9776;
        </button>
      )}
      <main className="mainContent">
        <h1>SpeakUp</h1>
        <div className="chat-container">
          <div className="messages">
            {messages.filter(msg => msg.chatId === currentChatId).map((message) => (
              <div key={message.id} className={`message ${message.sender}`}>
                {message.text}
              </div>
            ))}
            <div ref={messagesEndRef} />
          </div>
          
          <div className='chat-input'>
            <input 
              type="text" 
              placeholder="Digite sua mensagem..." 
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyPress={handleKeyPress}
            />
            <button onClick={handleSendMessage}>Enviar</button>
          </div>
        </div>
      </main>
    </div>
  );
} 

export default Home;