import { useRef, useEffect, useState } from 'react';
import styles from './Chat.module.css';
import {
  Send,
  Dehaze,
} from '@mui/icons-material';
import CircularProgress from '@mui/material/CircularProgress'; // 1. IMPORTAÇÃO ADICIONADA

// Api functions
import {
  getChatsByUserId,
  createChat,
  getMessagesByChatId,
  addMessageToChat,
  generateAIResponseDialog,
  generateAIResponseCorrection,
  generateAIResponseTranslation,
  generateAIResponseTopic
} from '../../utils/api';

// Components
import Sidebar from '../../components/Sidebar/Sidebar';

export default function Chat() {
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');
  const [currentChatId, setCurrentChatId] = useState(null);
  const [isSending, setIsSending] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [selectedChat, setSelectedChat] = useState(null);

  const messagesEndRef = useRef(null);

  useEffect(() => {
    getChatsByUserId()
      .then(res => {
        setChats(res.chats || []);
        setSelectedChat(null);
      })
      .catch(() => {
        setChats([]);
        setSelectedChat(null);
      });
  }, []);

  useEffect(() => {
    if (currentChatId) {
      getMessagesByChatId(currentChatId)
        .then(res => {
          const msgs = res.map(m => ({
            id: m.id,
            text: m.content,
            sender: m.sender,
            type: m.type
          }));
          setMessages(msgs);
        })
        .catch(() => setMessages([]));
    }
  }, [currentChatId]);

  useEffect(() => {
    if (window.innerWidth <= 768) {
      setSidebarOpen(false);
    }
  }, []);

  useEffect(() => {

    messagesEndRef.current?.scrollIntoView({ behavior: 'auto' });
    
  }, [messages]);

  const handleSend = async () => {
    if (!inputMessage.trim() || isSending) return;

    setIsSending(true);
    const userInput = inputMessage.trim();
    setInputMessage('');

    try {
      let chatId = currentChatId;

      if (!chatId) {
        const topicRes = await generateAIResponseTopic(userInput);
        const newChat = await createChat(topicRes.response);
        chatId = newChat.id;
        setCurrentChatId(chatId);
        setChats(prev => [...prev, newChat]);
      }

      const correction = await generateAIResponseCorrection(userInput);
      const fullUserMsg = `${userInput}\n\nCorreção: ${correction.response}`;

      const savedUserMsg = await addMessageToChat(chatId, fullUserMsg, 'user', 'request');
      setMessages(prev => [...prev, {
        id: savedUserMsg.id,
        text: fullUserMsg,
        sender: 'user',
        type: 'request'
      }]);

      const dialogRes = await generateAIResponseDialog(correction.response, chatId);
      const translation = await generateAIResponseTranslation(dialogRes.response);
      const fullAIResponse = `${dialogRes.response}\n\n[TRANSLATION]: ${translation.response}`;

      const savedBotMsg = await addMessageToChat(chatId, fullAIResponse, 'ai', 'response');
      setMessages(prev => [...prev, {
        id: savedBotMsg.id,
        text: fullAIResponse,
        sender: 'ai',
        type: 'response'
      }]);

    } catch (err) {
      console.error(err);
      setInputMessage(userInput);
    } finally {
      setIsSending(false);
    }
  };

  const handleChatSelect = (chatId) => {
    setCurrentChatId(chatId);
    const chat = chats.find(c => c.id === chatId);
    setSelectedChat(chat);
  };

  return (
    <div className={styles.Chat}>

      <Sidebar
        chats={chats}
        setChats={setChats}
        setCurrentChatId={handleChatSelect}
        setMessages={setMessages}
        setInputMessage={setInputMessage}
        sidebarOpen={sidebarOpen}
        setSidebarOpen={setSidebarOpen}
        selectedChat={selectedChat}
        setSelectedChat={setSelectedChat}
      />

      <main className={styles.main}>
        <header className={styles.chatHeader}>
          <Dehaze onClick={() => setSidebarOpen(prev => !prev)} className={`${styles.toggleIcon} ${!sidebarOpen ? styles.rotated : ''}`} />
        </header>

        <div className={styles.chatBody}>

          {messages.length === 0 && (
            <div className={styles.welcomeBox}>
              <h2>Bem-vindo ao Chat de Prática!</h2>
              <p>
                Aqui você pode praticar o idioma conversando com nossa IA.
              </p>
              <div>
                <strong>Exemplos:</strong><br /><br />
                <em>Entrada:</em> "I ned a car"<br />
                <em>Saída:</em> "I need a car"<br /><br />
                <em>Entrada:</em> "How are you doin?"<br />
                <em>Saída:</em> "How are you doing?"<br /><br />
                <em>Entrada:</em> "Let's go beach tomorrow?"<br />
                <em>Saída:</em> "Let's go to the beach tomorrow?"<br /><br />
                <em>Entrada:</em> "I don't know how say this."<br />
                <em>Saída:</em> "I don't know how to say this."
              </div>
            </div>
          )}


          {messages.map(msg => (
            <div
              key={msg.id}
              className={msg.sender === 'user' ? styles.userMessage : styles.botMessage}
            >
              <div>
                <p>{msg.text.split('\n\n')[0]}</p>
                {msg.type === 'request' && (
                  <div className={styles.responseExtras}>
                    <p><strong>Correção:</strong> {msg.text.split('\n\n')[1]?.replace('Correção: ', '')}</p>
                  </div>
                )}
                {msg.type === 'response' && (
                  <div className={styles.responseExtras}>
                    <p><strong>Tradução:</strong> {msg.text.split('[TRANSLATION]: ')[1]}</p>
                  </div>
                )}
              </div>
            </div>
          ))}

          <div ref={messagesEndRef} />

        </div>

        <footer className={styles.chatFooter}>
          <div className={styles.inputBox}>
            <input
              type="text"
              placeholder="Digite sua mensagem aqui..."
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && !isSending && handleSend()} // 2. MODIFICADO: Adicionado !isSending
              disabled={isSending} // 3. MODIFICADO: Adicionado disabled
            />
            <button onClick={handleSend} disabled={isSending}>
              {/* 4. MODIFICAÇÃO PRINCIPAL: Renderização condicional do ícone ou spinner */}
              {isSending ? (
                <CircularProgress size={24} color="inherit" />
              ) : (
                <Send />
              )}
            </button>
          </div>
        </footer>
      </main>
    </div>
  );
}