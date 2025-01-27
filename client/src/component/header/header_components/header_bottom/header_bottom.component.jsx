import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './header_bottom.component.css';

const Header_bottom = ({ isAuthenticated, setIsAuthenticated }) => {
    const navigate = useNavigate();

    const handleLogout = () => {
        localStorage.removeItem('token');
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
                            <li><button onClick={handleLogout} className="a">Выход</button></li>
                        </>
                    ) : (
                        <li><Link to="/auth" className="a">Вход системы</Link></li>
                    )}
                </ul>
            </nav>
        </div>
    );
};

export default Header_bottom;
