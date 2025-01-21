import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './header_bottom.component.css';

const Header_bottom = () => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    const navigate = useNavigate();

    useEffect(() => {
        // Проверка состояния аутентификации (например, проверка наличия токена)
        const token = localStorage.getItem('token');
        if (token) {
            checkToken(token);
        }
    }, []);

    const checkToken = async (token) => {
        const response = await fetch('http://localhost:5000/check-token', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
        });

        if (response.ok) {
            console.log('Токен действителен');
            setIsAuthenticated(true);
        } else {
            console.log('Неверный токен');
            setIsAuthenticated(false);
        }
    };

    const handleLogout = () => {
        // Удаление токена и обновление состояния аутентификации
        localStorage.removeItem('token');
        setIsAuthenticated(false);
        navigate('/');
    };
    return (
        <div className="header_bottom">
            <nav>
                <ul>
                    <li><Link to="/professional" className="a">Профессиональная переподготовка</Link></li>
                    <li><Link to="/promotion" className="a">Повышение квалификации  </Link></li>
                    {isAuthenticated ? (
                        <>
                            <li><Link to="/profile" className="a">Профиль</Link></li>
                            <li><button onClick={handleLogout} className="a">Выход</button></li>
                        </>
                    ) : (
                        <li><Link to="/registration" className="a">Вход системы</Link></li>
                    )}
                </ul>
            </nav>
        </div>
    )

};


export default Header_bottom;