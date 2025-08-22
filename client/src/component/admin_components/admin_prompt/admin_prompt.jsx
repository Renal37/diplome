import React, { useState, useEffect } from 'react';
import './admin_prompt.css';

const AdminPrompt = () => {
    const [prompt, setPrompt] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchPrompt = async () => {
            try {
                const response = await fetch('http://localhost:5000/admin/prompt', {
                    credentials: 'include',
                });
                if (!response.ok) throw new Error('Ошибка загрузки промпта');
                const data = await response.json();
                setPrompt(data.content);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };
        fetchPrompt();
    }, []);

    const handleSave = async () => {
        try {
            setLoading(true);
            const response = await fetch('http://localhost:5000/admin/prompt/update', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({ content: prompt }),
            });
            if (!response.ok) throw new Error('Ошибка сохранения');
            alert('Промпт успешно обновлен!');
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async () => {
        if (!window.confirm('Вы уверены, что хотите удалить промпт? Будет восстановлен стандартный текст.')) return;
        
        try {
            setLoading(true);
            const response = await fetch('http://localhost:5000/admin/prompt/delete', {
                method: 'DELETE',
                credentials: 'include',
            });
            if (!response.ok) throw new Error('Ошибка удаления');
            setPrompt(''); // Очищаем поле, будет загружен дефолтный промпт
            alert('Промпт удален, будет использоваться стандартный текст');
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    if (loading) return <div>Загрузка...</div>;
    if (error) return <div className="error-message">Ошибка: {error}</div>;

    return (
        <div className="admin-prompt-container">
            <h2>Управление промптом для AI-ассистента</h2>
            
            <div className="prompt-form">
                <textarea
                    value={prompt}
                    onChange={(e) => setPrompt(e.target.value)}
                    placeholder="Введите текст промпта для AI-ассистента"
                    rows={20}
                    style={{ width: '100%', padding: '10px', fontFamily: 'monospace', height:`300px`, fontSize:`18px`}}
                />
                <div className="prompt-actions">
                    <button className='approve-btn' onClick={handleSave} disabled={loading}>
                        Сохранить промпт
                    </button>
                    <button className='reject-btn' onClick={handleDelete} disabled={loading} >
                        Удалить промпт
                    </button>
                </div>
                {error && <div className="error-message">{error}</div>}
            </div>
        </div>
    );
};

export default AdminPrompt;