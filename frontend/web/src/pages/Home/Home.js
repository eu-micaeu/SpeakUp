import { useState, useEffect, useRef } from 'react';
import './Home.css';
import { getChatsByUserId, createChat, getMessagesByChatId, addMessageToChat, generateAIResponse } from '../../utils/api';

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
      setChats(response.chats || []);
    }).catch(error => {
      console.error("Erro ao buscar chats:", error);
      setChats([]);
    });
  }, []);

  useEffect(() => {
    if (currentChatId) {
      loadChatMessages(currentChatId);
    } else {
      setMessages([]);
    }
  }, [currentChatId]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  const loadChatMessages = async (chatId) => {
    try {
      const response = await getMessagesByChatId(chatId);
      const formattedMessages = response.map(msg => ({
        id: msg.id,
        text: msg.content,
        sender: msg.sender || 'user', // Assume 'user' como padrão se não houver sender
        timestamp: msg.timestamp,
        chatId: msg.chat_id
      }));
      setMessages(formattedMessages);
    } catch (error) {
      console.error("Erro ao carregar mensagens:", error);
      setMessages([]);
    }
  };
  
  const handleSendMessage = async () => {
    if (inputMessage.trim() === '') return;
    
    const messageContent = inputMessage.trim();
    setInputMessage('');
    let tempMessage;
    let tempAIMessage;
  
    try {
      let chatId = currentChatId;
      
      if (!chatId) {
        const firstWord = messageContent.split(' ')[0];
        const topic = firstWord.length > 0 ? firstWord : "Novo Chat";
        
        const newChat = await createChat(topic);
        chatId = newChat.id;
        setCurrentChatId(chatId);
        setChats(prevChats => [...prevChats, newChat]);
      }
  
      // Mensagem do usuário
      tempMessage = {
        id: Date.now(),
        text: messageContent,
        sender: 'user', // Definido explicitamente como 'user'
        timestamp: new Date().toISOString(),
        chatId: chatId
      };
  
      setMessages(prevMessages => [...prevMessages, tempMessage]);
      
      // Salva a mensagem do usuário no backend
      const savedMessage = await addMessageToChat(chatId, messageContent, 'user'); // Adiciona sender 'user'
  
      setMessages(prevMessages => [
        ...prevMessages.filter(m => m.id !== tempMessage.id),
        {
          id: savedMessage.id,
          text: savedMessage.content,
          sender: savedMessage.sender || 'user', // Garante que o sender seja 'user'
          timestamp: savedMessage.timestamp,
          chatId: savedMessage.chat_id
        }
      ]);
  
      // Resposta da IA
      const aiResponse = await generateAIResponse(messageContent);
      
      tempAIMessage = {
        id: Date.now() + 1,
        text: aiResponse.response,
        sender: 'ai', // Definido explicitamente como 'ai'
        timestamp: new Date().toISOString(),
        chatId: chatId
      };
  
      setMessages(prevMessages => [...prevMessages, tempAIMessage]);
      
      // Salva a resposta da IA no backend com sender 'ai'
      const savedAIMessage = await addMessageToChat(chatId, aiResponse.response, 'ai');
  
      setMessages(prevMessages => [
        ...prevMessages.filter(m => m.id !== tempAIMessage.id),
        {
          id: savedAIMessage.id,
          text: savedAIMessage.content,
          sender: savedAIMessage.sender || 'ai', // Garante que o sender seja 'ai'
          timestamp: savedAIMessage.timestamp,
          chatId: savedAIMessage.chat_id
        }
      ]);
  
    } catch (error) {
      console.error("Erro ao processar mensagem:", error);
      if (tempMessage) setMessages(prev => prev.filter(m => m.id !== tempMessage.id));
      if (tempAIMessage) setMessages(prev => prev.filter(m => m.id !== tempAIMessage.id));
      setInputMessage(messageContent);
    }
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
                onClick={() => {
                  setCurrentChatId(chat.id);
                  loadChatMessages(chat.id);
                }}
              >
                {chat.topic || "Chat sem título"}
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
            {messages
              .filter(msg => msg.chatId === currentChatId || msg.chat_id === currentChatId)
              .map((message) => (
                <div key={message.id} className={`message ${message.sender}`}>
                  {message.text || message.content}
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