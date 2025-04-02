import React, { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { getChatsByUserId, createChat, deleteChat, getMessagesByChatId, addMessageToChat, generateAIResponseDialog, generateAIResponseCorrection, generateAIResponseTranslation } from '../../utils/api';
import { useNavigate } from 'react-router-dom';
import SendIcon from '@mui/icons-material/Send';
import DeleteIcon from '@mui/icons-material/Delete';
import LogoutIcon from '@mui/icons-material/Logout';

const PageHome = styled.div`
  display: flex;
  height: 100vh;
`;

const Sidebar = styled.aside`
  width: ${props => props.$isVisible ? '300px' : '0'};
  background-color: #000000;
  position: relative;
  color: #ffffff;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  transition: width 0.3s ease-in-out;

  ol {
    padding: 0;
    text-align: center;
  }

  ul {
    margin: 50px 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    list-style: none;
  }

  ul li {
    list-style: none;
    text-align: center;
    padding: 10px;
    background-color: #313131;
    border-radius: 10px;
    cursor: pointer;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin: 10px 0;
    width: 180px;

    &:hover {
      background-color: #424242;
    }

    &.active-chat {
      background-color: #4CAF50;
    }
  }

  ul li svg {
    margin-left: 10px;
    color: #fff;
    cursor: pointer;
  }

  ul li svg:hover {
    color:rgb(255, 0, 0);
  }
`;

const MainContent = styled.main`
  padding: 20px;
  background-color: #313131;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  align-items: center;
  color: #fff;

  h1 {
    margin: 25px;
  }
`;

const ToggleButton = styled.button`
  background-color: transparent;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #fff;
  position: fixed;
  top: 20px;
  left: 20px;
  z-index: 1000;
  display: ${props => props.$sidebarVisible ? 'block' : 'none'};

  ${MainContent} & {
    color: #fff;
    position: fixed;
    display: block;
    left: ${props => props.$sidebarVisible ? '320px' : '20px'};
  }
`;

const ChatContainer = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 60%;
`;

const Messages = styled.div`
  flex: 1;
  padding: 10px;
  border-radius: 10px;
  overflow-y: auto;
  overflow-x: hidden;
  max-height: 70vh;

  &::-webkit-scrollbar {
    width: 8px;
  }

  &::-webkit-scrollbar-track {
    border-radius: 4px;
  }

  &::-webkit-scrollbar-thumb {
    background: #555;
    border-radius: 4px;
  }

  &::-webkit-scrollbar-thumb:hover {
    background: #777;
  }

`;

const Message = styled.div`
  padding: 20px;
  margin: 8px;
  border-radius: 4px;
  line-height: 2;

  &.user {
    color: #e3f2fd;
    background-color: #1d1d1d;
    margin-left: auto;
    width: 50%;
  }

  &.ai {
    background-color: #242424;
    margin-right: auto;
    width: 75%;
    position: relative;
  }

  &.ai.correction {
    background-color: #1eff00;
    color: #000000;
    text-align: center;
    width: 100%;
    padding: 10px;
    margin: 10px 0;
  }

  .translation {
    margin-top: 10px;
    padding-top: 10px;
    border-top: 1px solid #555;
    font-style: italic;
    color: #aaa;
  }
`;

const ChatInput = styled.div`
  display: flex;
  gap: 10px;
  width: 100%;
  justify-content: space-between;
  align-items: center;

  input {
    flex: 1;
    padding: 20px;
    width: 100%;
    border: none;
    font-size: 1rem;
    border-radius: 4px;
    background-color: #3d3d3d;
    color: #fff;

    &:focus {
      outline: none;
    }
  }

`;

const BtCreateChat = styled.button`
  border: none;
  border: 2px solid #4CAF50;
  border-radius: 10px;
  color: #4CAF50;
  background-color: transparent;
  padding: 10px;
  font-size: 16px;
  cursor: pointer;
  font-weight: bold;
  margin: 20px 0;
  width: 200px;
  font-family: 'Karla', sans-serif;
  &:hover {
    background-color: #4CAF50;
    color: #000000;
  }
`;

const DivSpeakUp = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: row;
  margin-bottom: 25px;
`;

function Home() {

  const [isSidebarVisible, setIsSidebarVisible] = useState(false);
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');
  const [currentChatId, setCurrentChatId] = useState(null);
  const messagesEndRef = useRef(null);
  const navigate = useNavigate();

  const goToIndex = () => {
    navigate('/');
  };

  const clearCookies = () => {
    document.cookie.split(";").forEach((c) => {
      document.cookie = c
        .replace(/^ +/, "")
        .replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/");
    });
  };


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
        sender: msg.sender || 'user',
        timestamp: msg.timestamp,
        chatId: msg.chat_id,
        type: msg.type
      }));
      setMessages(formattedMessages);
    } catch (error) {
      console.error("Erro ao carregar mensagens:", error);
      setMessages([]);
    }
  };

  // Dentro do componente Home, modifique a função handleSendMessage para incluir a tradução:
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

      // Mensagem do usuário e correção da IA (sem tradução)
      const aiCorrectionResponse = await generateAIResponseCorrection(messageContent);
      const combinedMessageContent = `${messageContent}\n\nCorreção: ${aiCorrectionResponse.response}`;

      tempMessage = {
        id: Date.now(),
        text: combinedMessageContent,
        sender: 'user',
        timestamp: new Date().toISOString(),
        chatId: chatId,
        type: 'request'
      };

      setMessages(prevMessages => [...prevMessages, tempMessage]);

      // Salva a mensagem do usuário no backend
      const savedMessage = await addMessageToChat(chatId, combinedMessageContent, 'user', 'request');

      setMessages(prevMessages => [
        ...prevMessages.filter(m => m.id !== tempMessage.id),
        {
          id: savedMessage.id,
          text: savedMessage.content,
          sender: savedMessage.sender || 'user',
          timestamp: savedMessage.timestamp,
          chatId: savedMessage.chat_id,
          type: 'request'
        }
      ]);

      // Resposta da IA - Dialog
      const aiResponseDialog = await generateAIResponseDialog(aiCorrectionResponse.response);
      // Obter tradução apenas da resposta da IA
      const aiTranslation = await generateAIResponseTranslation(aiResponseDialog.response);

      const aiResponseWithTranslation = `${aiResponseDialog.response}\n\n[TRANSLATION]: ${aiTranslation.response}`;

      tempAIMessage = {
        id: Date.now() + 1,
        text: aiResponseWithTranslation,
        sender: 'ai',
        timestamp: new Date().toISOString(),
        chatId: chatId,
        type: 'response'
      };

      setMessages(prevMessages => [...prevMessages, tempAIMessage]);

      // Salva a resposta da IA com tradução no backend
      const savedAIMessage = await addMessageToChat(chatId, aiResponseWithTranslation, 'ai', 'response');

      setMessages(prevMessages => [
        ...prevMessages.filter(m => m.id !== tempAIMessage.id),
        {
          id: savedAIMessage.id,
          text: savedAIMessage.content,
          sender: savedAIMessage.sender || 'ai',
          timestamp: savedAIMessage.timestamp,
          chatId: savedAIMessage.chat_id,
          type: 'response'
        }
      ]);

    } catch (error) {
      console.error("Erro ao processar mensagem:", error);
      if (tempMessage) setMessages(prev => prev.filter(m => m.id !== tempMessage.id));
      if (tempAIMessage) setMessages(prev => prev.filter(m => m.id !== tempAIMessage.id));
      setInputMessage(messageContent);
    }
  };

  // Mantenha a mesma renderização de mensagens do exemplo anterior

  const handleDeleteChat = async (chatId) => {
    try {
      await deleteChat(chatId);
      setChats(prevChats => prevChats.filter(chat => chat.id !== chatId));
      if (currentChatId === chatId) {
        setCurrentChatId(null);
        setMessages([]);
        loadChatMessages(null);
      }
    } catch (error) {
      console.error("Erro ao excluir chat:", error);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      handleSendMessage();
    }
  };

  return (
    <PageHome>
      <Sidebar $isVisible={isSidebarVisible}>
        <ToggleButton onClick={toggleSidebar} $sidebarVisible={isSidebarVisible}>
          &times;
        </ToggleButton>
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
                <DeleteIcon
                  onClick={(e) => {
                    e.stopPropagation();
                    handleDeleteChat(chat.id);
                  }}
                ></DeleteIcon>
              </li>
            ))
          ) : (
            <ol>Crie um chat!</ol>
          )}
          <BtCreateChat onClick={() => { setCurrentChatId(null); setMessages([]); setIsSidebarVisible(false) }}>+</BtCreateChat>
        </ul>

        <LogoutIcon
          onClick={() => { goToIndex(); clearCookies(); }}
          style={{ color: "#ff0000", cursor: "pointer", margin: "20px" }}
          onMouseEnter={(e) => e.target.style.color = "#ffffff"}
          onMouseLeave={(e) => e.target.style.color = "#ff0000"}
        >
          Sair
        </LogoutIcon>

      </Sidebar>

      <MainContent $sidebarVisible={isSidebarVisible}>
        {!isSidebarVisible && (
          <ToggleButton onClick={toggleSidebar} $sidebarVisible={isSidebarVisible}>
            &#9776;
          </ToggleButton>
        )}

        <DivSpeakUp>
          <img src='./logo.png' width={75} alt="SpeakUp Logo" />
          <h1>SpeakUp</h1>
        </DivSpeakUp>
        <ChatContainer>
          <Messages>
            {messages
              .filter(msg => msg.chatId === currentChatId || msg.chat_id === currentChatId)
              .map((message) => (
                <Message key={message.id} className={`${message.sender} ${message.type}`}>
                  {(message.text || message.content).split('\n\n').map((line, index, lines) => {
                    // Verifica se é a linha de tradução
                    if (line.startsWith('[TRANSLATION]: ')) {
                      return (
                        <div key={index} className="translation">
                          <strong>Tradução:</strong> {line.replace('[TRANSLATION]: ', '')}
                        </div>
                      );
                    }

                    // Verifica se é a linha de correção
                    if (index === 1 && line.startsWith('Correção: ')) {
                      return (
                        <React.Fragment key={index}>
                          <hr />
                          <span style={{ color: '#1eff00' }}>{line.replace('Correção: ', '')}</span>
                        </React.Fragment>
                      );
                    }

                    // Linha normal
                    return (
                      <React.Fragment key={index}>
                        {line}
                        {index < lines.length - 1 && <br />}
                      </React.Fragment>
                    );
                  })}
                </Message>
              ))}
            <div ref={messagesEndRef} />
          </Messages>

          <ChatInput>
            <input
              type="text"
              placeholder="Digite sua mensagem..."
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyPress={handleKeyPress}
            />
            <SendIcon
              onKeyDown={handleKeyPress}
              onClick={handleSendMessage}
              style={{ color: "#fff", cursor: "pointer" }}
              onMouseEnter={(e) => e.target.style.color = "rgb(187, 187, 187)"}
              onMouseLeave={(e) => e.target.style.color = "#fff"}
            />
          </ChatInput>

        </ChatContainer>
      </MainContent>
    </PageHome>
  );
}

export default Home;
