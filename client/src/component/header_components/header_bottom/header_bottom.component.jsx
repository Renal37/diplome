import React, { useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './header_bottom.component.css';
import { AuthContext } from '../../../content/authContext';

const Header_bottom = () => {
    const { isAuthenticated, logout } = useContext(AuthContext);
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
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
                            <li><Link to="/profile" className="a">Профиль</Link></li>
                            <li><button onClick={handleLogout} className="a">Выход</button></li>
                        </>
                    ) : (
                        <li><Link to="/registration" className="a">Вход системы</Link></li>
                    )}
                </ul>
            </nav>
        </div>
    );
};

export default Header_bottom;