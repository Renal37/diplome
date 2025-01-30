import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './header_bottom.component.css';
import Cookies from 'js-cookie';

const Header_bottom = ({ isAuthenticated, setIsAuthenticated }) => {
    const navigate = useNavigate();

    const handleLogout = async () => {
        try {
            // Отправляем запрос на сервер для удаления куки
            const response = await fetch('http://localhost:5000/logout', {
                method: 'POST',
                credentials: 'include', // Отправляем куки
            });

            if (response.ok) {
                // Удаляем куки на клиенте, если они не HttpOnly
                Cookies.remove("token", { path: "/" });
                localStorage.clear();
                setIsAuthenticated(false);
                navigate('/');
            } else {
                console.error('Ошибка при выходе из системы');
            }
        } catch (error) {
            console.error('Ошибка при выходе из системы:', error);
        }
    };

    return (
        <div className="header_bottom">
            <nav>
                <ul>
                    <li><Link to="/professional" className="a">Профессиональная переподготовка</Link></li>
                    <li><Link to="/promotion" className="a">Повышение квалификации</Link></li>
                    {isAuthenticated ? (
                        <>
                            <li><Link to="/auth/profile" className="a">Профиль</Link></li>
                            <li><button onClick={handleLogout} className="a">Выйти</button></li>
                        </>
                    ) : (
                        <li><Link to="/auth" className="a">Войти</Link></li>
                    )}
                </ul>
            </nav>
        </div>
    );
};

export default Header_bottom;