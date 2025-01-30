import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import './header_bottom.component.css';

const Header_bottom = ({ isAuthenticated, setIsAuthenticated }) => {
    const navigate = useNavigate();

    const handleLogout = () => {
        // Удаление куки с токеном вручную
        document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    
        setIsAuthenticated(false);
        navigate('/');
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