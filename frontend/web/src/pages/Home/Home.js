import { useState, useEffect, useRef } from 'react';
import './Home.css';
import { getChatsByUserId, createChat, getMessagesByChatId, addMessageToChat } from '../../utils/api';

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
        sender: 'user',
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
    let tempMessage; // Declaramos a variável aqui para que seja acessível no bloco catch
  
    try {
      let chatId = currentChatId;
      
      // Se não há chat atual, cria um novo
      if (!chatId) {
        const firstWord = messageContent.split(' ')[0];
        const topic = firstWord.length > 0 ? firstWord : "Novo Chat";
        
        const newChat = await createChat(topic);
        chatId = newChat.id;
        setCurrentChatId(chatId);
        setChats(prevChats => [...prevChats, newChat]);
      }
  
      // Cria a mensagem temporária localmente
      tempMessage = {
        id: Date.now(), // ID temporário
        text: messageContent,
        sender: 'user',
        timestamp: new Date().toISOString(),
        chatId: chatId
      };
  
      setMessages(prevMessages => [...prevMessages, tempMessage]);
  
      // Envia a mensagem para o backend
      const savedMessage = await addMessageToChat(chatId, messageContent);
  
      // Atualiza a mensagem local com os dados do servidor
      setMessages(prevMessages => [
        ...prevMessages.filter(m => m.id !== tempMessage.id), // Remove a temporária
        {
          id: savedMessage.id,
          text: savedMessage.content,
          sender: 'user',
          timestamp: savedMessage.timestamp,
          chatId: savedMessage.chat_id
        }
      ]);
  
    } catch (error) {
      console.error("Erro ao processar mensagem:", error);
      // Remove a mensagem temporária em caso de erro
      if (tempMessage) {
        setMessages(prevMessages => prevMessages.filter(m => m.id !== tempMessage.id));
      }
      setInputMessage(messageContent); // Devolve a mensagem para o input
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