import React, { useEffect, useState } from 'react';
import styles from './Chat.module.css';
import {
  Translate,
  Add,
  Delete,
  Settings,
  AccountCircle,
  Apps,
  VolumeUp,
  ExpandMore,
  Mic,
  AttachFile,
  Send,
  Dehaze
} from '@mui/icons-material';

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

export default function Chat() {
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');
  const [currentChatId, setCurrentChatId] = useState(null);
  const [isSending, setIsSending] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(true);

  useEffect(() => {
    getChatsByUserId()
      .then(res => setChats(res.chats || []))
      .catch(() => setChats([]));
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

  return (
    <div className={styles.Chat}>
      <aside className={`${styles.sidebar} ${!sidebarOpen ? styles.sidebarClosed : ''}`}>

        <div>
          <div className={styles.header}>
            <img src="/logo.png" alt="Logo" width={30} />
          </div>

          <button className={styles.newChat} onClick={() => {
            setCurrentChatId(null);
            setMessages([]);
            setInputMessage('');
          }}>
            <Add />
          </button>

          <nav>
            {chats.map(chat => (
              <a
                key={chat.id}
                className={styles.navItem}
                href="#"
                onClick={() => setCurrentChatId(chat.id)}
              >
                <span>{chat.topic}</span>
                <Delete />
              </a>
            ))}
          </nav>
        </div>

        <div className={styles.footer}>
          <button><Settings /> Configurações</button>
          <button><AccountCircle /> Perfil</button>
          <button><Apps /> Menu</button>
        </div>
      </aside>

      <main className={styles.main}>
        <header className={styles.chatHeader}>
          <Dehaze onClick={() => setSidebarOpen(prev => !prev)} className={`${styles.toggleIcon} ${!sidebarOpen ? styles.rotated : ''}`} />
        </header>

        <div className={styles.chatBody}>
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
        </div>

        <footer className={styles.chatFooter}>
          <div className={styles.inputBox}>
            <button><Mic /></button>
            <input
              type="text"
              placeholder="Digite sua mensagem aqui..."
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSend()}
            />
            <button><AttachFile /></button>
            <button onClick={handleSend} disabled={isSending}>
              <Send />
            </button>
          </div>
        </footer>
      </main>
    </div>
  );
}
