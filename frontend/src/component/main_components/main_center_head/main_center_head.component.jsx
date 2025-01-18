import React from 'react';
import { useLocation } from 'react-router-dom';

const Header_center_head = () => {
    const location = useLocation();

    const getTitle = () => {
        switch (location.pathname) {
            case '/teacher':
                return 'Преподавательский состав';
            case '/professional':
                return 'Курсы профессиональной переподготовки';
            case '/promotion':
                return 'Курсы повышения квалификации';
            default:
                return 'Общая информация';
        }
    };

    return (
        <div className="main_center_header">
            <div className="line"></div>
            <h1>{getTitle()}</h1>
            <div className="line"></div>
        </div>
    );
};

export default Header_center_head;