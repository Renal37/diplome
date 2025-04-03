import React, { useState, useEffect, useRef } from 'react';
import { useLocation } from 'react-router-dom';
import './AIAssistant.css';

const AIAssistant = () => {
  const location = useLocation();
  const [isOpen, setIsOpen] = useState(false);
  const [message, setMessage] = useState('');
  const [conversation, setConversation] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isMinimized, setIsMinimized] = useState(false);
  const messagesEndRef = useRef(null);
  const inputRef = useRef(null);

  // Примеры подсказок для пользователя
  const exampleQuestions = [
    "Как записаться на курс?",
    "Какие документы нужны для зачисления?",
    "Есть ли рассрочка оплаты?",
    "Когда начинается следующий поток?",
    "Какой документ выдается после окончания?"
  ];

  // Получаем контекст текущей страницы
  const getPageContext = () => {
    const pageMap = {
      '/courses': 'страница курсов',
      '/profile': 'личный кабинет',
      '/prices': 'стоимость обучения',
      '/contacts': 'контакты',
      '/about': 'о компании'
    };
    return pageMap[location.pathname] || 'главная страница';
  };

  // Автоматическая прокрутка к новым сообщениям
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [conversation]);

  // Фокус на поле ввода при открытии
  useEffect(() => {
    if (isOpen && inputRef.current) {
      inputRef.current.focus();
    }
  }, [isOpen, conversation]);

  // Обработчик отправки сообщения
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!message.trim() || isLoading) return;

    // Добавляем сообщение пользователя в чат
    const userMessage = { 
      sender: 'user', 
      text: message,
      timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    };
    setConversation(prev => [...prev, userMessage]);
    setMessage('');
    setIsLoading(true);

    try {
      const response = await fetch('http://localhost:5000/ai-assistant', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          question: message,
          context: getPageContext(),
        }),
      });

      if (!response.ok) {
        throw new Error(response.statusText);
      }

      const data = await response.json();
      const aiMessage = { 
        sender: 'ai', 
        text: data.answer,
        timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
      };
      setConversation(prev => [...prev, aiMessage]);
    } catch (error) {
      console.error('Ошибка запроса к ассистенту:', error);
      const errorMessage = { 
        sender: 'ai', 
        text: 'Произошла ошибка соединения. Пожалуйста, попробуйте позже или напишите нам на support@example.com',
        timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
      };
      setConversation(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  // Обработчик нажатия на пример вопроса
  const handleExampleClick = (example) => {
    setMessage(example);
    if (inputRef.current) {
      inputRef.current.focus();
    }
  };

  // Сброс чата
  const handleResetChat = () => {
    if (conversation.length > 0 && window.confirm('Очистить историю диалога?')) {
      setConversation([]);
    }
  };

  return (
    <div className={`ai-assistant ${isOpen ? 'open' : ''} ${isMinimized ? 'minimized' : ''}`}>
      {/* Заголовок ассистента */}
      <div className="ai-header" onClick={() => setIsOpen(!isOpen)}>
        <div className="ai-title">
          <div className="ai-icon">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z"/>
            </svg>
          </div>
          <h3>Помощник по курсам</h3>
        </div>
        <div className="ai-controls">
          <button 
            className="ai-minimize-btn"
            onClick={(e) => {
              e.stopPropagation();
              setIsMinimized(!isMinimized);
            }}
            title={isMinimized ? "Развернуть" : "Свернуть"}
          >
            {isMinimized ? '+' : '-'}
          </button>
          <button 
            className="ai-close-btn"
            onClick={(e) => {
              e.stopPropagation();
              setIsOpen(false);
              setIsMinimized(false);
            }}
            title="Закрыть"
          >
            ×
          </button>
        </div>
      </div>
      
      {/* Основное содержимое ассистента */}
      {isOpen && !isMinimized && (
        <div className="ai-content">
          {/* Сообщения чата */}
          <div className="ai-messages">
            {conversation.length === 0 ? (
              <div className="ai-welcome">
                <p>Привет! Я ваш помощник по курсам. Задайте мне вопрос, например:</p>
                <div className="ai-examples">
                  {exampleQuestions.map((question, index) => (
                    <button
                      key={index}
                      className="ai-example-btn"
                      onClick={() => handleExampleClick(question)}
                    >
                      {question}
                    </button>
                  ))}
                </div>
              </div>
            ) : (
              conversation.map((msg, index) => (
                <div key={index} className={`ai-message ${msg.sender}`}>
                  <div className="ai-message-content">
                    {msg.text.split('\n').map((paragraph, i) => (
                      <p key={i}>{paragraph}</p>
                    ))}
                  </div>
                  <div className="ai-message-time">{msg.timestamp}</div>
                </div>
              ))
            )}
            {isLoading && (
              <div className="ai-message ai">
                <div className="ai-typing-indicator">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            )}
            <div ref={messagesEndRef} />
          </div>
          
          {/* Форма ввода сообщения */}
          <form onSubmit={handleSubmit} className="ai-input-form">
            <div className="ai-input-container">
              <input
                ref={inputRef}
                type="text"
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                placeholder="Напишите ваш вопрос..."
                disabled={isLoading}
              />
              <button 
                type="submit" 
                disabled={isLoading || !message.trim()}
                className="ai-send-btn"
                title="Отправить"
              >
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"/>
                </svg>
              </button>
            </div>
            <div className="ai-actions">
              <button 
                type="button" 
                className="ai-reset-btn"
                onClick={handleResetChat}
                disabled={conversation.length === 0}
              >
                Очистить чат
              </button>
              <small className="ai-hint">
                Нажимая "Отправить", вы соглашаетесь с обработкой данных
              </small>
            </div>
          </form>
        </div>
      )}
    </div>
  );
};

export default AIAssistant;