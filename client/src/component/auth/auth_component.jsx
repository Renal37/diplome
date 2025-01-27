import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Registration from './registration/registration_component.jsx';
import Authorization from './authorization/authorization_component.jsx';
import Profile from './profile/profile_component.jsx';
import './auth_component.css'

const Auth = () => {
    return (
        <div className='auth'>
            <Routes>
                <Route index element={<Registration />} />
                <Route path="authorization" element={<Authorization />} />
                <Route path="profile" element={<Profile />} />
            </Routes>
        </div>

    );
};

export default Auth;