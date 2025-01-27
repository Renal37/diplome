import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './profile_component.css';

const Profile = () => {
  const [residence, setResidence] = useState('');
  const [education, setEducation] = useState('');
  const [birthPlace, setBirthPlace] = useState('');
  const [homeAddress, setHomeAddress] = useState('');
  const [passportData, setPassportData] = useState('');
  const [snils, setSnils] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('/auth/profile');
    }
  }, [navigate]);
  const handleSubmit = async (e) => {
    e.preventDefault();
    const userId = 'user-id'; // Замените на реальный идентификатор пользователя
    const response = await fetch(`http://localhost:5000/update-profile/${userId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        residence,
        education,
        birthPlace,
        homeAddress,
        passportData,
        snils,
      }),
    });

    if (response.ok) {
      alert('Профиль успешно обновлен!');
    } else {
      alert('Ошибка при обновлении профиля');
    }
  };

  return (
    <div className="profile-component">
      <h1>Профиль</h1>
      <form onSubmit={handleSubmit} className="profile-form">
        <input
          type="text"
          placeholder="Место жительства"
          value={residence}
          onChange={(e) => setResidence(e.target.value)}
        />
        <input
          type="text"
          placeholder="Образование"
          value={education}
          onChange={(e) => setEducation(e.target.value)}
        />
        <input
          type="text"
          placeholder="Место рождения"
          value={birthPlace}
          onChange={(e) => setBirthPlace(e.target.value)}
        />
        <input
          type="text"
          placeholder="Домашний адрес (прописка)"
          value={homeAddress}
          onChange={(e) => setHomeAddress(e.target.value)}
        />
        <input
          type="text"
          placeholder="Данные паспорта (серия и номер)"
          value={passportData}
          onChange={(e) => setPassportData(e.target.value)}
        />
        <input
          type="text"
          placeholder="СНИЛС"
          value={snils}
          onChange={(e) => setSnils(e.target.value)}
        />
        <button type="submit">Сохранить</button>
      </form>
    </div>
  );
};

export default Profile;